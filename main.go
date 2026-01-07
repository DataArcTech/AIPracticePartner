package main

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/document/parser/pdf"
	"github.com/cloudwego/eino/components/document/parser"
)

func main() {
	//TestModel(context.Background())
	ctx := context.Background()

	// 初始化解析器
	p, err := pdf.NewPDFParser(ctx, &pdf.Config{
		ToPages: true, // 不按页面分割
	})
	if err != nil {
		panic(err)
	}

	// 打开 PDF 文件
	file, err := os.Open("k:\\go_projects\\AIPracticePartner\\test_docs\\GlobalPowerMultiCurrencyPlan3.pdf")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 解析文档
	docs, err := p.Parse(ctx, file,
		parser.WithURI("document.pdf"),
		parser.WithExtraMeta(map[string]any{
			"source": "./document.pdf",
		}),
	)
	if err != nil {
		panic(err)
	}

	// 使用解析结果
	for _, doc := range docs {
		println(doc.Content)
	}
}
