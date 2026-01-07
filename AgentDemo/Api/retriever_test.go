package Api

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestRetriever(T *testing.T) {
	ctx := context.Background()
	retrieverMake, err := RetrieverMake(ctx)
	if err != nil {
		return
	}
	invoke, err := retrieverMake.Invoke(ctx, "保险")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("共检索 %d 文档\n", len(invoke))
}
