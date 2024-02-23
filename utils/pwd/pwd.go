package pwd

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPwd(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash), nil
}

func CheckPwd(hashPwd string, pwd string) bool {
	byteHash := []byte(hashPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(pwd))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

//func CheckPwd(hashPwd string, pwd string) bool {
//	// 直接比较已哈希的密码和用户输入的密码
//	err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(pwd))
//	if err != nil {
//		log.Println(err)
//		return false
//	}
//	return true
//}
