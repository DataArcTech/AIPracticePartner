package Rewrite

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

// NewRewriteModel 初始化重写模型
func NewRewriteModel(ctx context.Context) (cm model.ToolCallingChatModel, err error) {
	config := &openai.ChatModelConfig{
		APIKey:  "sk-cmtnvcaupuoizcqogdbapkqyvdmyumolprmgwetjmxsxmwtk",
		BaseURL: "https://api.siliconflow.cn/v1",
		Model:   "Qwen/Qwen2.5-7B-Instruct",
	}
	cm, err = openai.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}

var systemPrompt = "你是一个智能搜索助手，负责将用户的自然语言问题转化为精准的搜索引擎关键词。\n" +
	"目标：提取核心语义，去除无关词汇，生成适合 Elasticsearch 的查询字符串。\n" +
	"现在时间：{time_now} (仅当用户问题包含相对时间描述如'昨天'、'下周'时参考)\n\n" +
	"规则：\n" +
	"1. 输出仅包含 2-5 个核心关键词，用空格分隔。\n" +
	"2. 去除停用词（如'的'、'是'、'怎么'、'我想知道'等）。\n" +
	"3. 不要包含任何解释、标点符号或前缀。\n" +
	"4. 如果用户提到特定产品或条款，请保留准确名称。\n" +
	"5. 知识库领域为《{knowledgeBase}》，如果查询词本身隐含了该领域（如'保险'），可省略领域词以扩大匹配范围。\n\n" +
	"示例：\n" +
	"用户：缴费期内如果退保会有什么损失吗？\n" +
	"输出：缴费期 退保 损失 现金价值\n\n" +
	"用户：昨天的最新费率表\n" +
	"输出：费率表 2026-01-06\n\n" +
	"用户：重疾险有哪些除外责任\n" +
	"输出：重疾险 除外责任 条款\n"

func GetOptimizedQueryMessages(question, knowledgeBase string) ([]*schema.Message, error) {
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage(systemPrompt),
		schema.UserMessage("如下是用户的问题: {question}"),
	)

	data := map[string]any{
		"time_now":      time.Now().Format(time.RFC3339),
		"question":      question,
		"knowledgeBase": knowledgeBase,
	}

	messages, err := template.Format(context.Background(), data)
	if err != nil {
		return nil, fmt.Errorf("格式化模板失败: %w", err)
	}
	return messages, nil
}
