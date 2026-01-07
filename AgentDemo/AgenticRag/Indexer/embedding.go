package Indexer

import (
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino/components/embedding"
)

func newEmbedding(ctx context.Context) (eb embedding.Embedder, err error) {
	config := &openai.EmbeddingConfig{
		APIKey:  "sk-cmtnvcaupuoizcqogdbapkqyvdmyumolprmgwetjmxsxmwtk",
		BaseURL: "https://api.siliconflow.cn/v1",
		Model:   "Qwen/Qwen3-Embedding-8B",
	}
	eb, err = openai.NewEmbedder(ctx, config)
	if err != nil {
		return nil, err
	}
	return eb, nil
}
