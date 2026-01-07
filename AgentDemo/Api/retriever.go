package Api

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"

	rerank "AIPracticePartner/AgentDemo/AgenticRag/Rerank"
	"AIPracticePartner/AgentDemo/AgenticRag/Retriever"
	"AIPracticePartner/AgentDemo/AgenticRag/Rewrite"
)

// RetrieverFunc 定义检索函数类型
type RetrieverFunc func(ctx context.Context, input string, opts ...compose.Option) ([]*schema.Document, error)

// Invoke 实现 Invoke 方法
func (f RetrieverFunc) Invoke(ctx context.Context, input string, opts ...compose.Option) ([]*schema.Document, error) {
	return f(ctx, input, opts...)
}

// RetrieverMake 构建优化后的检索器
/*流程：Rewrite（重写） -> Retrieve（检索） -> Rerank（重排）*/
func RetrieverMake(ctx context.Context) (RetrieverFunc, error) {
	// 重写搜索词模型
	rewriteModel, err := Rewrite.NewRewriteModel(ctx)
	if err != nil {
		return nil, fmt.Errorf("创建重写模型失败: %w", err)
	}
	// 获取基础检索器
	baseRetriever, err := Retriever.BuildRetriever(ctx)
	if err != nil {
		return nil, fmt.Errorf("构建基础检索器失败: %w", err)
	}

	// 返回核心处理逻辑
	return func(ctx context.Context, input string, opts ...compose.Option) ([]*schema.Document, error) {
		// 通过知识库名称查询重写
		start := time.Now()
		fmt.Printf("原始查询词: %s\n", input)
		msgs, err := Rewrite.GetOptimizedQueryMessages(input, "Insurance")
		if err != nil {
			return nil, fmt.Errorf("生成重写提示词失败: %w", err)
		}
		resp, err := rewriteModel.Generate(ctx, msgs)
		if err != nil {
			return nil, fmt.Errorf("执行查询重写失败: %w", err)
		}
		fmt.Printf("Rewrite总时间: %v\n", time.Since(start))

		rewrittenQuery := resp.Content
		// 失败则回退到原始查询
		if rewrittenQuery == "" {
			rewrittenQuery = input
		}

		// 知识库检索
		start = time.Now()
		docs, err := baseRetriever.Invoke(ctx, rewrittenQuery)
		if err != nil {
			return nil, fmt.Errorf("检索失败: %w", err)
		}
		fmt.Printf("Retrieve总时间：%v, 文档数量: %d\n", time.Since(start), len(docs))

		// 如果没有检索到文档，直接返回空
		if len(docs) == 0 {
			return nil, nil
		}

		// --- 结果重排 ---
		// 使用重写后的查询词对检索结果进行重排
		start = time.Now()
		topK := 20

		// 限制参与重排的文档数量
		rerankCandidates := docs
		if len(docs) > 50 {
			rerankCandidates = docs[:50]
		}

		rerankedDocs, err := rerank.NewRerank(ctx, rewrittenQuery, rerankCandidates, topK)
		if err != nil {
			// 如果重排失败，降级处理：直接返回原始检索结果
			fmt.Printf("[RetrieverMake] Rerank 失败，回退到原始检索结果: %v\n", err)
			return docs, nil
		}
		fmt.Printf("Rerank总时间: %v\n", time.Since(start))

		return rerankedDocs, nil
	}, nil
}
