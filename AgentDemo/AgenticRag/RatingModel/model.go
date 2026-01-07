package RatingModel

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

func newChatModel(ctx context.Context) (cm model.ToolCallingChatModel, err error) {
	config := &openai.ChatModelConfig{
		APIKey:  "sk-cmtnvcaupuoizcqogdbapkqyvdmyumolprmgwetjmxsxmwtk",
		BaseURL: "https://api.siliconflow.cn/v1",
		Model:   "Pro/deepseek-ai/DeepSeek-V3.2",
	}
	cm, err = openai.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}
