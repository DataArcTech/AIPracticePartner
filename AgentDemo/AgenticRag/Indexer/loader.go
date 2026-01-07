package Indexer

import (
	"context"

	"github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino/components/document"
)

// newLoader component initialization function of node 'FileLoader' in graph 'Indexer'
func newLoader(ctx context.Context) (ldr document.Loader, err error) {
	// 获取定义支持文档格式
	p, err := newParser(ctx)
	if err != nil {
		return nil, err
	}
	config := &file.FileLoaderConfig{
		Parser:      p,
		UseNameAsID: false, // 不使用文件名作为文档ID，由indexer生成
	}
	ldr, err = file.NewFileLoader(ctx, config)
	if err != nil {
		return nil, err
	}
	return ldr, nil
}
