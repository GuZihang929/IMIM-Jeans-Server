// @Author zhangjiaozhu 2024/2/27 18:47:00
package system

import (
	"IM-Server/global"
	"IM-Server/model/system"
	"IM-Server/service/user"
	"IM-Server/utils/jwts"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

type GptRequest struct {
	Msg string `json:"msg" binding:"required" msg:"请输入问题"` // 消息
}

// GptController 调gpt接口
func (s SysPublicController) GptController(c *gin.Context) {
	var gptRequest GptRequest
	if err := c.ShouldBindJSON(&gptRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	answer, err := user.GptService(gptRequest.Msg)
	if err != nil {
		c.JSON(http.StatusRequestTimeout, gin.H{"error": "密码哈希处理失败"})
		return
	}
	// 生成uuid
	sessionId := uuid.New().String()
	// 解析token获取用户id
	token := c.GetHeader("token")
	claims, err := jwts.ParseToken(token)
	if err != nil {
		log.Println("Token验证出错:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	sessionLog := system.SessionLog{
		SessionID: sessionId,
		Request:   gptRequest.Msg,
		Response:  answer,
		UserID:    int(claims.UserID),
	}
	// 异步存储会话记录
	go func() {
		global.DB.Create(&sessionLog).Debug()
	}()
	s.GptDemoController(c, answer)
	//c.JSON(http.StatusOK, answer)
}

// GptDemoController gpt消息流式输出
func (s SysPublicController) GptDemoController(c *gin.Context, answer string) {
	//	const Text = `
	//proxy_cache：通过这个模块，Nginx 可以缓存代理服务器从后端服务器请求到的响应数据。当下一个客户端请求相同的资源时，Nginx 可以直接从缓存中返回响应，而不必去请求后端服务器。这大大降低了代理服务器的负载，同时也能提高客户端访问速度。需要注意的是，使用 proxy_cache 模块时需要谨慎配置缓存策略，避免出现缓存不一致或者过期的情况。
	//
	//proxy_buffering：通过这个模块，Nginx 可以将后端服务器响应数据缓冲起来，并在完整的响应数据到达之后再将其发送给客户端。这种方式可以减少代理服务器和客户端之间的网络连接数，提高并发处理能力，同时也可以防止后端服务器过早关闭连接，导致客户端无法接收到完整的响应数据。
	//
	//综上所述， proxy_cache 和 proxy_buffering 都可以通过缓存技术提高代理服务器性能和安全性，但需要注意合理的配置和使用，以避免潜在的缓存不一致或者过期等问题。同时， proxy_buffering 还可以通过缓冲响应数据来提高代理服务器的并发处理能力，从而更好地服务于客户端。
	//`
	var Text = answer
	type ChatCompletionChunk struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int64  `json:"created"`
		Model   string `json:"model"`
		Choices []struct {
			Delta struct {
				Content string `json:"content"`
			} `json:"delta"`
			Index        int     `json:"index"`
			FinishReason *string `json:"finish_reason"`
		} `json:"choices"`
	}
	w := c.Writer
	// 设置Content-Type标头为text/event-stream
	w.Header().Set("Content-Type", "text/event-stream")
	// 设置缓存控制标头以禁用缓存
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Keep-Alive", "timeout=20")
	w.Header().Set("Transfer-Encoding", "chunked")
	// 生成一个uuid
	uid := uuid.NewString()
	created := time.Now().Unix()

	for i, v := range Text {
		eventData := fmt.Sprintf("%c", v)
		if eventData == "" {
			continue
		}
		var finishReason *string
		if i == len(Text)-1 {
			temp := "stop"
			finishReason = &temp
		}
		chunk := ChatCompletionChunk{
			ID:      uid,
			Object:  "chat.completion.chunk",
			Created: created,
			Model:   "gpt-3.5-turbo-0301",
			Choices: []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
				Index        int     `json:"index"`
				FinishReason *string `json:"finish_reason"`
			}{
				{Delta: struct {
					Content string `json:"content"`
				}{
					Content: eventData,
				},
					Index:        0,
					FinishReason: finishReason,
				},
			},
		}

		//fmt.Println("输出：" + eventData)
		marshal, err := json.Marshal(chunk)
		if err != nil {
			return
		}

		fmt.Fprintf(w, "data: %v\n\n", string(marshal))
		flusher, ok := w.(http.Flusher)
		if ok {
			flusher.Flush()
		} else {
			log.Println("Flushing not supported")
		}
		if i == len(Text)-1 {
			fmt.Fprintf(w, "data: [DONE]")
			flusher, ok := w.(http.Flusher)
			if ok {
				flusher.Flush()
			} else {
				log.Println("Flushing not supported")
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
