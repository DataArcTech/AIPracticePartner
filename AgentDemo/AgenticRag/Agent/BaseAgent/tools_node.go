package BaseAgent

import (
	"AIPracticePartner/Api"
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/v2"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

type RetrieveInput struct {
	Query string `json:"query" jsonschema:"description=The search query to retrieve related documents"`
}

type RetrieveOutput struct {
	Documents []string `json:"documents"`
}

// duckduckgo 工具调用
func newTool(ctx context.Context) (bt tool.BaseTool, err error) {
	config := &duckduckgo.Config{}
	bt, err = duckduckgo.NewTextSearchTool(ctx, config)
	if err != nil {
		return nil, err
	}
	return bt, nil
}

// NewRetrieverTool 优化检索器
func NewRetrieverTool(ctx context.Context) (tool.BaseTool, error) {
	log.Println("优化检索器被调用")
	// 获取优化后的检索器函数
	retrieverFunc, err := Api.RetrieverMake(ctx)
	if err != nil {
		return nil, err
	}

	// 定义工具执行逻辑
	run := func(ctx context.Context, input *RetrieveInput) (*RetrieveOutput, error) {
		// 添加日志以验证工具调用
		fmt.Printf("\n[Tool: knowledge_retriever] 被调用，查询词: %s\n", input.Query)
		docs, err := retrieverFunc.Invoke(ctx, input.Query)
		if err != nil {
			fmt.Printf("[Tool: knowledge_retriever] 执行出错: %v\n", err)
			return nil, err
		}
		fmt.Printf("[Tool: knowledge_retriever] 检索成功，找到 %d 个文档\n", len(docs))
		var docContents []string
		for _, doc := range docs {
			// 简单的将文档内容提取出来
			docContents = append(docContents, doc.Content)
		}
		return &RetrieveOutput{Documents: docContents}, nil
	}

	// 创建工具实例
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "knowledge_retriever",
			Desc: "请使用此工具从知识库中检索特定的产品信息、条款和细则。作为培训合作伙伴（客户角色），您可以利用它来核实销售人员的说法、查找细节以提出具有挑战性的问题，或验证产品功能。输入内容应为与产品相关的关键词或具体条款。",
		},
		run,
	), nil
}

// NewFullKnowledgeRetrieverTool 全量知识库检索工具
func NewFullKnowledgeRetrieverTool(ctx context.Context) (tool.BaseTool, error) {
	log.Println("全量知识库检索器被调用")
	// 获取优化后的检索器函数
	retrieverFunc, err := Api.RetrieverMake(ctx)
	if err != nil {
		return nil, err
	}

	// 定义工具执行逻辑
	run := func(ctx context.Context, input *RetrieveInput) (*RetrieveOutput, error) {
		// 添加日志以验证工具调用
		fmt.Printf("\n[Tool: full_knowledge_retriever] (评分专用) 被调用，查询词: %s\n", "保险")
		docs, err := retrieverFunc.Invoke(ctx, input.Query)
		if err != nil {
			fmt.Printf("[Tool: full_knowledge_retriever] 执行出错: %v\n", err)
			return nil, err
		}
		fmt.Printf("[Tool: full_knowledge_retriever] 检索成功，找到 %d 个文档\n", len(docs))
		var docContents []string
		for i, doc := range docs {
			docContents = append(docContents, doc.Content)
			// 打印前3条文档内容作为示例
			if i < 3 {
				preview := doc.Content
				if len([]rune(preview)) > 100 {
					preview = string([]rune(preview)[:100]) + "..."
				}
				fmt.Printf("  [文档 %d]: %s\n", i+1, preview)
			}
		}
		return &RetrieveOutput{Documents: docContents}, nil
	}

	// 创建工具实例
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "full_knowledge_retriever",
			Desc: "【评分模型专用】上帝视角检索工具。可以搜索知识库中的所有内容（包括产品条款、销售话术、费率表、竞品分析等），用于获取“标准答案”以评估销售员的回答是否准确。使用此工具来查找最全面的信息。",
		},
		run,
	), nil
}
