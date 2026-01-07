package Retriever

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func BuildRetriever(ctx context.Context) (r compose.Runnable[string, []*schema.Document], err error) {
	const Retriever1 = "Retriever1"
	g := compose.NewGraph[string, []*schema.Document]()
	retriever1KeyOfRetriever, err := newRetriever(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddRetrieverNode(Retriever1, retriever1KeyOfRetriever)
	_ = g.AddEdge(compose.START, Retriever1)
	_ = g.AddEdge(Retriever1, compose.END)
	r, err = g.Compile(ctx, compose.WithGraphName("Retriever"))
	if err != nil {
		return nil, err
	}
	return r, err
}
