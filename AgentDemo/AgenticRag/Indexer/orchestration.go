package Indexer

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

func BuildIndexer(ctx context.Context) (r compose.Runnable[any, []string], err error) {
	const (
		FileIndexer         = "FileIndexer"
		FileLoader          = "FileLoader"
		DocumentTransformer = "DocumentTransformer"
	)
	g := compose.NewGraph[any, []string]()
	fileIndexerKeyOfIndexer, err := newIndexer(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddIndexerNode(FileIndexer, fileIndexerKeyOfIndexer)
	fileLoaderKeyOfLoader, err := newLoader(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddLoaderNode(FileLoader, fileLoaderKeyOfLoader)

	transformer, err := newDocumentTransformer(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddDocumentTransformerNode(DocumentTransformer, transformer)

	_ = g.AddEdge(compose.START, FileLoader)
	_ = g.AddEdge(FileLoader, DocumentTransformer)
	_ = g.AddEdge(DocumentTransformer, FileIndexer)
	_ = g.AddEdge(FileIndexer, compose.END)

	r, err = g.Compile(ctx, compose.WithGraphName("Indexer"))
	if err != nil {
		return nil, err
	}
	return r, err
}
