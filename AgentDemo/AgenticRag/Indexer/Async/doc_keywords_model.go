package Async

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

// newChatModel 用于文档结构化处理
func newChatModel(ctx context.Context) (cm model.ToolCallingChatModel, err error) {
	config := &openai.ChatModelConfig{
		APIKey:  "sk-cmtnvcaupuoizcqogdbapkqyvdmyumolprmgwetjmxsxmwtk",
		BaseURL: "https://api.siliconflow.cn/v1",
		Model:   "Qwen/Qwen3-8B",
	}
	cm, err = openai.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}
