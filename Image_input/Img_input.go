package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Ginのルーターを初期化
	r := gin.Default()

	// ルートエンドポイント
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, Gin!")
	})

	// サーバーを起動
	r.Run(":8080")
}
