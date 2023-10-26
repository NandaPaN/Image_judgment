package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// MySQLデータベースに接続
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/Image")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Ginのルーターを初期化
	r := gin.Default()

	// ルートエンドポイント
	// フォームを表示するハンドラ
	r.GET("/upload-form", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload-form.html", nil)
	})

	// 画像のアップロードを処理するハンドラ
	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("ファイルの取得に失敗しました: %s", err.Error()))
			return
		}

		// 画像をバイナリデータに変換
		imageFile, err := file.Open()
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("ファイルのオープンに失敗しました: %s", err.Error()))
			return
		}
		defer imageFile.Close()

		imageData, err := ioutil.ReadAll(imageFile)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("ファイルの読み込みに失敗しました: %s", err.Error()))
			return
		}

		// データベースに画像を挿入
		_, err = db.Exec("INSERT INTO images (name, data) VALUES (?, ?)", file.Filename, imageData)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("データベースへの挿入に失敗しました: %s", err.Error()))
			return
		}

		c.String(http.StatusOK, fmt.Sprintf("ファイル '%s' をアップロードしました。", file.Filename))
	})

	// サーバーを起動
	r.Run(":8080")
}
