package BaseAgent

import (
	"context"
	"log"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
)

// newLambda component initialization function of node 'DemoAgent' in graph 'AIDemo'
func newLambda(ctx context.Context) (lba *compose.Lambda, err error) {
	config := &react.AgentConfig{
		MaxStep: 5,
	}
	chatModelIns11, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}
	config.ToolCallingModel = chatModelIns11
	// toolIns21, err := newTool(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// 集成知识库检索工具
	retrieverTool, err := NewRetrieverTool(ctx)
	if err != nil {
		return nil, err
	}
	config.ToolsConfig.Tools = []tool.BaseTool{retrieverTool}
	log.Println("Agent已启动")
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
