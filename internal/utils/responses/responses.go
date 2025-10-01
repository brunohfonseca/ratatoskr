package responses

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{
		"status": "ok",
		"data":   data,
	})
}

func Error(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{
		"status": "error",
		"error":  err.Error(),
	})
}

func ErrorMsg(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{
		"status": "error",
		"error":  msg,
	})
}
