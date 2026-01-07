package Async

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type StructTransformer struct {
	chatModel model.ToolCallingChatModel
}

func NewStructTransformer(ctx context.Context) (*StructTransformer, error) {
	cm, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}
	return &StructTransformer{chatModel: cm}, nil
}

func (t *StructTransformer) Transform(ctx context.Context, docs []*schema.Document, opts ...document.TransformerOption) ([]*schema.Document, error) {
	for _, doc := range docs {
		// 提取元数据 (产品名称, 关键词, 摘要)
		// 节省token和避免超长，获取前5000个字符进行提取
		contentPreview := doc.Content
		if len([]rune(contentPreview)) > 5000 {
			contentPreview = string([]rune(contentPreview)[:5000])
		}

		metaData, err := t.extractMetadata(ctx, contentPreview)
		if err != nil {
			log.Printf("提取元数据失败 (docID: %s): %v", doc.ID, err)
			continue
		}

		// 将提取的元数据注入到文档的 MetaData 中
		if doc.MetaData == nil {
			doc.MetaData = make(map[string]any)
		}

		// 确保summary存在
		if summary, ok := metaData["summary"].(string); ok && summary != "" {
			doc.MetaData["summary"] = summary
		}

		// 确保keywords存在
		if keywords, ok := metaData["keywords"].([]interface{}); ok {
			// 转换为 []string
			var keywordStrs []string
			for _, k := range keywords {
				if s, ok := k.(string); ok {
					keywordStrs = append(keywordStrs, s)
				}
			}
			doc.MetaData["keywords"] = keywordStrs
		}

		log.Printf("关键词: %v", doc.MetaData["keywords"])
	}
	return docs, nil
}

func (t *StructTransformer) extractMetadata(ctx context.Context, content string) (map[string]any, error) {
	systemPrompt := `你是一个专业的文档分析助手。请分析给定的产品文档片段，提取以下关键信息，并以严格的 JSON 格式返回：
	1. "keywords": 核心关键词列表（字符串数组，提取3-5个最关键的词，如产品类型、核心功能、适用人群）
	2. "summary": 文档摘要（字符串，200字以内，一定要抓取最核心的内容和重点内容）
	请直接返回 JSON 字符串，不要包含 Markdown 代码块标记（如 '''json）。`

	userPrompt := fmt.Sprintf("文档内容片段：\n%s", content)

	msgs := []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(userPrompt),
	}

	resp, err := t.chatModel.Generate(ctx, msgs)
	if err != nil {
		return nil, err
	}

	result := resp.Content
	// 清理可能的 markdown 标记
	result = strings.TrimPrefix(result, "```json")
	result = strings.TrimPrefix(result, "```")
	result = strings.TrimSuffix(result, "```")
	result = strings.TrimSpace(result)

	var meta map[string]any
	err = json.Unmarshal([]byte(result), &meta)
	if err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w, 原文: %s", err, result)
	}

	return meta, nil
}
