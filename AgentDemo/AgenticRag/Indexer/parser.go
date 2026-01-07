package Indexer

import (
	"context"

	"github.com/cloudwego/eino-ext/components/document/parser/docx"
	"github.com/cloudwego/eino-ext/components/document/parser/pdf"
	"github.com/cloudwego/eino/components/document/parser"
)

func newParser(ctx context.Context) (p parser.Parser, err error) {
	//pdf解析
	pdfParser, err := pdf.NewPDFParser(ctx, &pdf.Config{
		ToPages: true,
	})
	if err != nil {
		return nil, err
	}
	//docx解析
	docxParser, err := docx.NewDocxParser(ctx, &docx.Config{
		ToSections:     true,
		IncludeHeaders: true,
		IncludeFooters: true,
		IncludeTables:  true,
	})
	if err != nil {
		return nil, err
	}

	//初始化文本解析器
	textParser := &parser.TextParser{}

	return parser.NewExtParser(ctx, &parser.ExtParserConfig{Parsers: map[string]parser.Parser{
		".pdf":  pdfParser,
		".docx": docxParser,
		".doc":  docxParser,
	},
		FallbackParser: textParser,
	})

}
