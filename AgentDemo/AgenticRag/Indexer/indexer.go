package Indexer

import (
	"AIPracticePartner/Common"
	"context"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/components/indexer/es8"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
)

// NewIndexer component initialization function of node 'FileIndexer' in graph 'Indexer'
func newIndexer(ctx context.Context) (idr indexer.Indexer, err error) {
	esClient, err := Common.EsClient(ctx)
	if err != nil {
		return nil, err
	}
	config := &es8.IndexerConfig{
		Client:    esClient,
		Index:     "agent_demo_knowledge",
		BatchSize: 50, // 增加批处理大小，减少Embedding网络请求次数，提升索引速度
		DocumentToFields: func(ctx context.Context, doc *schema.Document) (field2Value map[string]es8.FieldValue, err error) {
			var knowledgeName string
			if value, ok := ctx.Value(Common.KnowledgeName).(string); ok {
				knowledgeName = value
			} else {
				knowledgeName = "default"
			}
			/*生成文档ID*/
			if len(doc.ID) == 0 {
				doc.ID = uuid.New().String()
			}
			if doc.MetaData != nil {
				/*存储ext数据*/
				marshal, _ := sonic.Marshal(getExtData(doc))
				doc.MetaData[Common.FieldExtra] = string(marshal)
			}

			fields := map[string]es8.FieldValue{
				Common.FieldContent: {
					Value:    doc.Content,
					EmbedKey: Common.FieldContentVector,
				},
				Common.FieldExtra: {
					Value: doc.MetaData[Common.FieldExtra],
				},
				Common.KnowledgeName: {
					Value: knowledgeName,
				},
			}

			// TODO 添加结构化字段到索引
			if v, ok := doc.MetaData["summary"]; ok {
				fields[Common.FieldSummary] = es8.FieldValue{Value: v}
			}
			if v, ok := doc.MetaData["keywords"]; ok {
				fields[Common.FieldKeywords] = es8.FieldValue{Value: v}
			}
			return fields, nil
		},
	}
	embeddingIns11, err := newEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	config.Embedding = embeddingIns11
	idr, err = es8.NewIndexer(ctx, config)
	if err != nil {
		return nil, err
	}
	return idr, nil
}

func getExtData(doc *schema.Document) map[string]any {
	if doc.MetaData == nil {
		return nil
	}
	res := make(map[string]any)
	for _, key := range Common.ExtKeys {
		if v, e := doc.MetaData[key]; e {
			res[key] = v
		}
	}
	return res
}
