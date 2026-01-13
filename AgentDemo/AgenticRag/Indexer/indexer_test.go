package Indexer

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/components/document"
)

func TestIndexer(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "_knowledge_name", "Insurance")

	/*知识库索引*/
	runnable, err := BuildIndexer(ctx)
	if err != nil {
		t.Fatalf("BuildIndexer failed: %v", err)
	}

	/*将文档索引到ES8*/
	testFiles := []string{
		`k:\go_projects\AIPracticePartner\test_docs\GlobalPowerMultiCurrencyPlan3-页面-12.pdf`,
	}

	for _, filePath := range testFiles {
		t.Logf("Start indexing file: %s", filePath)

		input := document.Source{
			URI: filePath,
		}

		/*梳理文档*/
		ids, err := runnable.Invoke(ctx, input)
		if err != nil {
			t.Errorf("Failed to index file %s: %v", filePath, err)
			continue
		}

		/*结果输出*/
		t.Logf("Successfully indexed file: %s", filePath)
		for _, msg := range ids {
			t.Logf("Generated IDs: %v", msg)
		}
	}
}
