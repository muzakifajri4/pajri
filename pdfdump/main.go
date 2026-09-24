package main

import (
	"fmt"
	"os"

	"github.com/ledongthuc/pdf"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run pdfdump <pdfpath>")
		return
	}
	path := os.Args[1]
	f, r, err := pdf.Open(path)
	if err != nil {
		fmt.Printf("failed to open PDF: %v\n", err)
		return
	}
	defer f.Close()

	totalPage := r.NumPage()
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		content, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}
		fmt.Printf("===== PAGE %d =====\n", pageIndex)
		fmt.Printf("%s\n", content)
	}
}