package Retriever

import (
	"AIPracticePartner/AgentDemo/Common"
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/components/retriever/es8"
	"github.com/cloudwego/eino-ext/components/retriever/es8/search_mode"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// newRetriever component initialization function of node 'Retriever1' in graph 'Retriever'
func newRetriever(ctx context.Context) (rtr retriever.Retriever, err error) {
	vectorField := Common.FieldContentVector
	var score = 0.01
	esClient, err := Common.EsClient(ctx)
	if err != nil {
		return nil, err
	}
	if value, ok := ctx.Value(Common.RetrieverFieldKey).(string); ok {
		vectorField = value
	}
	config := &es8.RetrieverConfig{
		Client: esClient,
		Index:  "agent_demo_knowledge",
		//混合检索
		SearchMode: search_mode.SearchModeApproximate(&search_mode.ApproximateConfig{
			// QueryFieldName:  Common.FieldSummary, // 优先检索摘要，速度更快且语义更浓缩
			QueryFieldName:  Common.FieldContent,
			VectorFieldName: vectorField,
			Hybrid:          true,
		}),
		ResultParser:   EsHit2Document,
		TopK:           100,
		ScoreThreshold: &score,
	}
	embeddingIns11, err := newEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	config.Embedding = embeddingIns11
	rtr, err = es8.NewRetriever(ctx, config)
	if err != nil {
		return nil, err
	}
	return rtr, nil
}
func EsHit2Document(ctx context.Context, hit types.Hit) (doc *schema.Document, err error) {
	doc = &schema.Document{
		ID:       *hit.Id_,
		MetaData: map[string]any{},
	}

	var src map[string]any
	if err = sonic.Unmarshal(hit.Source_, &src); err != nil {
		return nil, err
	}

	for field, val := range src {
		switch field {
		case Common.FieldContent:
			doc.Content = val.(string)
		case Common.FieldContentVector:
			var v []float64
			for _, item := range val.([]interface{}) {
				v = append(v, item.(float64))
			}
			doc.WithDenseVector(v)
		case Common.FieldExtra:
			if val == nil {
				continue
			}
			doc.MetaData[Common.FieldExtra] = val.(string)
		case Common.KnowledgeName:
			doc.MetaData[Common.KnowledgeName] = val.(string)
		case Common.FieldSummary:
			doc.MetaData[Common.FieldSummary] = val.(string)
		case Common.FieldKeywords:
			doc.MetaData[Common.FieldKeywords] = val.([]interface{})
		default:
			return nil, fmt.Errorf("unexpected field=%s, val=%v", field, val)
		}
	}

	if hit.Score_ != nil {
		doc.WithScore(float64(*hit.Score_))
	}

	return doc, nil
}
