package RatingModel

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
)

// newLambda component initialization function of node 'RatingLambda' in graph 'RatingModel'
func newLambda(ctx context.Context) (lba *compose.Lambda, err error) {
	config := &react.AgentConfig{
		MaxStep: 5,
	}
	chatModelIns11, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}
	config.ToolCallingModel = chatModelIns11
	// 集成全量知识库检索工具（评分专用）
	/*retrieverTool, err := BaseAgent.NewFullKnowledgeRetrieverTool(ctx)
	if err != nil {
		return nil, err
	}
	config.ToolsConfig.Tools = []tool.BaseTool{retrieverTool}*/
	ins, err := react.NewAgent(ctx, config)
	if err != nil {
		return nil, err
	}
	lba, err = compose.AnyLambda(ins.Generate, ins.Stream, nil, nil)
	if err != nil {
		return nil, err
	}
	return lba, nil
}
