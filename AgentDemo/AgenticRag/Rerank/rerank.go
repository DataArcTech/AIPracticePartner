package rerank

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/schema"
)

type Conf struct {
	Model           string `json:"model"`
	ReturnDocuments bool   `json:"return_documents"`
	MaxChunksPerDoc int    `json:"max_chunks_per_doc"`
	OverlapTokens   int    `json:"overlap_tokens"`
	url             string
	apiKey          string
}
type Data struct {
	Query     string   `json:"query"`
	Documents []string `json:"documents"`
	TopN      int      `json:"top_n"`
}

type Req struct {
	*Data
	*Conf
}

type Result struct {
	Index          int     `json:"index"`
	RelevanceScore float64 `json:"relevance_score"`
}

type Resp struct {
	ID      string    `json:"id"`
	Results []*Result `json:"results"`
}

var rerankCfg *Conf

// NewRerank 创建Rerank
func NewRerank(ctx context.Context, query string, docs []*schema.Document, topK int) (output []*schema.Document, err error) {
	output, err = rerank(ctx, query, docs, topK)
	if err != nil {
		return
	}
	return
}

// GetConf 模型配置
func GetConf(ctx context.Context) *Conf {
	if rerankCfg != nil {
		return rerankCfg
	}
	apiKey := "sk-cmtnvcaupuoizcqogdbapkqyvdmyumolprmgwetjmxsxmwtk"
	baseUrl := "https://api.siliconflow.cn/v1"
	model := "Qwen/Qwen3-Reranker-0.6B"
	url := fmt.Sprintf("%s/rerank", baseUrl)
	rerankCfg = &Conf{
		apiKey:          apiKey,
		Model:           model,
		ReturnDocuments: false,
		MaxChunksPerDoc: 1024,
		OverlapTokens:   80,
		url:             url,
	}
	return rerankCfg
}

// rerank 处理
func rerank(ctx context.Context, query string, docs []*schema.Document, topK int) (output []*schema.Document, err error) {
	data := &Data{
		Query: query,
		TopN:  topK,
	}
	for _, doc := range docs {
		// 截断过长文本，减少网络传输和推理耗时
		content := doc.Content
		if len([]rune(content)) > 1500 {
			content = string([]rune(content)[:1500])
		}
		data.Documents = append(data.Documents, content)
	}
	// 重排
	results, err := rerankDoHttp(ctx, data)
	if err != nil {
		return
	}
	// 重新组装数据
	for _, result := range results {
		// 防止索引越界
		if result.Index >= len(docs) {
			continue
		}
		doc := docs[result.Index]
		doc.WithScore(result.RelevanceScore)
		output = append(output, docs[result.Index])
	}
	return
}

// rerankDoHttp 调用Rerank模型API
func rerankDoHttp(ctx context.Context, data *Data) ([]*Result, error) {
	cfg := GetConf(ctx)
	reqData := &Req{
		Data: data,
		Conf: cfg,
	}

	jsonData, err := sonic.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", cfg.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("rerank api failed, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result Resp
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = sonic.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Results, nil
}
