package BaseAgent

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

func newChatModel(ctx context.Context) (cm model.ToolCallingChatModel, err error) {

	config := &openai.ChatModelConfig{
		APIKey:  "1f5aeabe-11dc-4a52-8b44-c5aa143b0f7f",
		BaseURL: "https://ark.cn-beijing.volces.com/api/v3",
		Model:   "doubao-seed-1-6-251015",
		/*Model:   "doubao-seed-1-6-flash-250828",*/
	}
	cm, err = openai.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}
