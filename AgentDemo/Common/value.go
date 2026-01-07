package Common

const (
	FieldContent       = "content"
	FieldContentVector = "content_vector"
	FieldExtra         = "ext"
	FieldSummary       = "summary"
	FieldKeywords      = "keywords"
	KnowledgeName      = "_knowledge_name"
	RetrieverFieldKey  = "_retriever_field"
)

var (
	FrequencyPenalty float32 = 0.3 // 提高重复惩罚，避免无限复读
	PresencePenalty  float32 = 0.3 //存在惩罚
	Temperature      float32 = 0.8 //温度，控制随机生成
	TopP             float32 = 0.8 //核采样，控制候选词广度
	ExtKeys                  = []string{"_extension", "_file_name", "_source"}
)
