package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
	vision "google.golang.org/api/vision/v1"
)

type ReceiptOcrClient struct {
	client *vision.Service
}

func NewReceiptOcrClient(credentialPath string) (*ReceiptOcrClient, error) {
	ctx := context.Background()

	// Google Cloud Vision AI クライアントの初期化
	client, err := vision.NewService(ctx, option.WithCredentialsFile(credentialPath))
	if err != nil {
		return nil, err
	}

	return &ReceiptOcrClient{client}, nil
}

func (r *ReceiptOcrClient) GetPaymentAmount(keyWord string, imagePath string) (int, error) {
	// Vision AIにリクエストを送信
	image := &vision.Image{Content: imageToBase64(imagePath)}
	feature := &vision.Feature{Type: "TEXT_DETECTION"}
	req := &vision.AnnotateImageRequest{
		Image:    image,
		Features: []*vision.Feature{feature},
	}
	batch := &vision.BatchAnnotateImagesRequest{Requests: []*vision.AnnotateImageRequest{req}}
	resp, err := r.client.Images.Annotate(batch).Do()
	if err != nil {
		return -999, err
	}

	// OCR結果を取得
	text := resp.Responses[0].FullTextAnnotation.Text
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		if strings.Contains(line, keyWord) {
			// 合計金額が含まれる行を見つけたら、数字の部分を抽出
			parts := strings.Fields(line)
			for _, part := range parts {
				if number := extractNumber(part); number != "" {
					return parseNumber(number)
				}
			}
		}
	}

	return -999, nil

}
func (r *ReceiptOcrClient) GetPaymentAmount(keyWord string, imagePath string) (int, error) {
	// Vision AIにリクエストを送信
	image := &vision.Image{Content: imageToBase64(imagePath)}
	feature := &vision.Feature{Type: "TEXT_DETECTION"}
	req := &vision.AnnotateImageRequest{
		Image:    image,
		Features: []*vision.Feature{feature},
	}
	batch := &vision.BatchAnnotateImagesRequest{Requests: []*vision.AnnotateImageRequest{req}}
	resp, err := r.client.Images.Annotate(batch).Do()
	if err != nil {
		return -999, err
	}

	// OCR結果を取得
	text := resp.Responses[0].FullTextAnnotation.Text
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		if strings.Contains(line, keyWord) {
			// 合計金額が含まれる行を見つけたら、数字の部分を抽出
			parts := strings.Fields(line)
			for _, part := range parts {
				if number := extractNumber(part); number != "" {
					return parseNumber(number)
				}
			}
		}
	}

	return -999, nil
}

func extractNumber(input string) string {
	// 数字と通貨記号以外を削除
	var result strings.Builder
	for _, r := range input {
		if (r >= '0' && r <= '9') || r == '¥' || r == ',' || r == '.' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func parseNumber(input string) (int, error) {
	// 数字の文字列を整数に変換
	// 通貨記号とカンマを削除し、小数点をドットに変換
	input = strings.ReplaceAll(input, "¥", "")
	input = strings.ReplaceAll(input, ",", "")
	input = strings.ReplaceAll(input, ".", "")
	return fmt.Sprintf("%s", input)
}

func imageToBase64(imagePath string) string {
	// 画像ファイルをBase64エンコードして返す
	// ここに画像ファイルをBase64エンコードする処理を追加
	// この部分は具体的な実装に合わせて修正が必要です
	return ""
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <imagePath>")
	}

	credentialPath := "config/service_account_key.json" // サービスアカウントキーのパス
	keyWord := "合計" // 合計金額を意味するキーワード
	imagePath := os.Args[1] // 画像ファイルのパス

	client, err := NewReceiptOcrClient(credentialPath)
	if err != nil {
		log.Fatalf("Error initializing OCR client: %v", err)
	}

	amount, err := client.GetPaymentAmount(keyWord, imagePath)
	if err != nil {
		log.Fatalf("Error extracting payment amount: %v", err)
	}

	if amount == -999 {
		log.Println("Payment amount not found in the receipt.")
	} else {
		log.Printf("Payment amount is %d yen.", amount)
	}
}
