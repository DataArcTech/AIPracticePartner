package Indexer

import (
	"context"
	"fmt"

	"strings"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/recursive"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
)

func newDocumentTransformer(ctx context.Context) (document.Transformer, error) {
	// 初始化结构化提取器
	/*structTrans, err := NewStructTransformer(ctx)
	if err != nil {
		//如果初始化失败，不中断流程做一个降级处理
		fmt.Printf("Warning: Init StructTransformer failed: %v\n", err)
	}*/

	config := &recursive.Config{
		ChunkSize:   1000,                                    //每个分块1000字符
		OverlapSize: 100,                                     //分块间重叠100字符
		Separators:  []string{"\n", "。", "！", "？", "!", "?"}, // 优先级从高到低
	}
	t, err := recursive.NewSplitter(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create recursive failed:%w", err)
	}

	return &JunkFilterTransformer{
		inner: t,
		//structTrans: structTrans,
	}, nil
}

type JunkFilterTransformer struct {
	inner document.Transformer
	//structTrans *StructTransformer
}

// Transform 转换内容
func (t *JunkFilterTransformer) Transform(ctx context.Context, docs []*schema.Document, opts ...document.TransformerOption) ([]*schema.Document, error) {
	/* 进行结构化处理 */
	/*if t.structTrans != nil {
		processedDocs, err := t.structTrans.Transform(ctx, docs, opts...)
		if err == nil {
			docs = processedDocs
		} else {
			fmt.Printf("Warning: StructTransform failed: %v\n", err)
		}
	}*/

	/* 调用内部的分块器进行分块 */
	chunks, err := t.inner.Transform(ctx, docs, opts...)
	if err != nil {
		return nil, err
	}

	/* 清理包含长分割线的垃圾内容 */
	var cleanChunks []*schema.Document
	for _, chunk := range chunks {
		cleaned := removeJunkLines(chunk.Content)
		if len(strings.TrimSpace(cleaned)) == 0 {
			continue
		}
		/* 创建新对象以避免副作用 */
		newDoc := *chunk
		newDoc.Content = cleaned
		/* 清空ID，确保Indexer为每个分块生成唯一的ID */
		newDoc.ID = ""
		cleanChunks = append(cleanChunks, &newDoc)
	}

	return cleanChunks, nil
}

/* 清楚长分割线 */
func removeJunkLines(content string) string {
	lines := strings.Split(content, "\n")
	var keepLines []string
	for _, line := range lines {
		if isJunkLine(line) {
			continue
		}
		keepLines = append(keepLines, line)
	}
	return strings.Join(keepLines, "\n")
}

/* 判断是否为长分割线 */
func isJunkLine(line string) bool {
	/* 检查行是否主要由破折号或类似分割符组成 */
	separators := 0
	for _, r := range line {
		/* 包含常见分割线字符 */
		if r == '-' || r == '—' || r == '–' || r == '_' || r == '=' || r == '*' {
			separators++
		}
	}
	/* 如果一行中有超过10个分割符，视为垃圾行 */
	if separators > 10 {
		return true
	}
	return false
}
