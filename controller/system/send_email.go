package system

import (
	"IM-Server/config/confDetail"
	"IM-Server/global"
	"IM-Server/res"
	"IM-Server/service/user"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"time"
)

func (s SysPublicController) SendEmail(c *gin.Context) {
	Account := c.PostForm("account")
	// 检查是否存在相同账号（邮箱）的用户
	if user.IsUserExists(Account) {
		c.JSON(http.StatusOK, gin.H{"error": "已存在相同邮箱的用户"})
		global.Logger.Warn("已存在相同邮箱的用户")
		return
	}

	// 生成并发送验证邮件
	verificationCode := generateVerificationCode()
	if err := sendVerificationEmail(Account, verificationCode); err != nil {
		global.Logger.Warn("发送邮箱验证码失败")
		res.FailWithCode(http.StatusOK, c)
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "发送验证邮件失败"})
		return
	}

	// 在 Redis 中存储用户信息和验证码
	err := user.StoreUserAndCodeInRedis(Account, verificationCode)
	if err != nil {
		global.Logger.Warn("在 Redis 中存储用户和验证码失败")
		res.FailWithCode(http.StatusOK, c)
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "在 Redis 中存储用户和验证码失败"})
		return
	}

	// 返回注册成功的信息
	c.JSON(http.StatusOK, gin.H{"message": "邮箱验证码发送成功，请查看您的邮箱进行验证。"})
}

func sendVerificationEmail(email, code string) error {
	// 从 Viper 中获取邮件配置信息
	emailConfig := getEmailConfig()
	fmt.Println(emailConfig)
	// 使用 go-mail 发送邮件
	m := gomail.NewMessage()
	m.SetHeader("From", emailConfig.User)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verification Code")
	m.SetBody("text/html", fmt.Sprintf("Your verification code is: %s", code))

	d := gomail.NewDialer(emailConfig.Host, emailConfig.Port, emailConfig.User, emailConfig.Password)

	// 配置TLS
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // 生产环境中，应该避免使用 InsecureSkipVerify 为 true

	if err := d.DialAndSend(m); err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	return nil
}

// 从 Viper 中获取邮件配置信息
func getEmailConfig() *confDetail.Email {
	emailConfig := &confDetail.Email{
		Host:     global.Config.Email.Host,
		Port:     global.Config.Email.Port,
		User:     global.Config.Email.User,
		Password: global.Config.Email.Password,
		UseSSL:   global.Config.Email.UseSSL,
		UserTls:  global.Config.Email.UserTls,
	}
	return emailConfig
}

func generateVerificationCode() string {
	// 生成6位随机数字验证码
	return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
}
