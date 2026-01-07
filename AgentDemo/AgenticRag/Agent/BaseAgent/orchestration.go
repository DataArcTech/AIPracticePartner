package BaseAgent

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func BuildAIDemo(ctx context.Context) (r compose.Runnable[map[string]any, *schema.Message], err error) {
	const (
		DemoAgent        = "DemoAgent"
		DemoChatTemplate = "DemoChatTemplate"
	)
	g := compose.NewGraph[map[string]any, *schema.Message]()
	demoAgentKeyOfLambda, err := newLambda(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddLambdaNode(DemoAgent, demoAgentKeyOfLambda)
	demoChatTemplateKeyOfChatTemplate, err := NewChatTemplate(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(DemoChatTemplate, demoChatTemplateKeyOfChatTemplate)
	_ = g.AddEdge(compose.START, DemoChatTemplate)
	_ = g.AddEdge(DemoAgent, compose.END)
	_ = g.AddEdge(DemoChatTemplate, DemoAgent)
	r, err = g.Compile(ctx, compose.WithGraphName("AIDemo"))
	if err != nil {
		return nil, err
	}
	return r, err
}
