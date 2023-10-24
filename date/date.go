package main

import (
	"fmt"
	"sort"
)

type BoundingBox struct {
	X, Y int
}

type Symbol struct {
	Text        string
	BoundingBox BoundingBox
}

func getSortedLines(response VisionResponse, threshold int) [][]Symbol {
	// 1. テキスト抽出とソート
	document := response.FullTextAnnotation
	var bounds []Symbol

	for _, page := range document.Pages {
		for _, block := range page.Blocks {
			for _, paragraph := range block.Paragraphs {
				for _, word := range paragraph.Words {
					for _, symbol := range word.Symbols {
						bbox := BoundingBox{
							X: symbol.BoundingBox.Vertices[0].X,
							Y: symbol.BoundingBox.Vertices[0].Y,
						}
						s := Symbol{
							Text:        symbol.Text,
							BoundingBox: bbox,
						}
						bounds = append(bounds, s)
					}
				}
			}
		}
	}

	sort.Slice(bounds, func(i, j int) bool {
		return bounds[i].BoundingBox.Y < bounds[j].BoundingBox.Y
	})

	// 2. 同じ高さのものをまとめる
	oldY := -1
	var line []Symbol
	var lines [][]Symbol

	for _, bound := range bounds {
		x := bound.BoundingBox.X
		y := bound.BoundingBox.Y

		if oldY == -1 {
			oldY = y
		} else if oldY-threshold <= y && y <= oldY+threshold {
			oldY = y
		} else {
			oldY = -1
			sort.Slice(line, func(i, j int) bool {
				return line[i].BoundingBox.X < line[j].BoundingBox.X
			})
			lines = append(lines, line)
			line = nil
		}

		line = append(line, bound)
	}

	sort.Slice(line, func(i, j int) bool {
		return line[i].BoundingBox.X < line[j].BoundingBox.X
	})
	lines = append(lines, line)

	return lines
}

type VisionResponse struct {
	FullTextAnnotation struct {
		Pages []struct {
			Blocks []struct {
				Paragraphs []struct {
					Words []struct {
						Symbols []struct {
							Text        string
							BoundingBox struct {
								Vertices []struct {
									X int
									Y int
								}
							}
						}
					}
				}
			}
		}
	}
}

func main() {
	// この部分にVisionResponseを実際のデータで初期化するコードを追加してください

	threshold := 5
	lines := getSortedLines(VisionResponse{}, threshold)
	fmt.Println(lines)
}
