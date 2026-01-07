package main

import (
	"AIPracticePartner/AgentDemo/AgenticRag/Agent/BaseAgent"
	"AIPracticePartner/AgentDemo/AgenticRag/Agent/GentleAgent"
	"AIPracticePartner/AgentDemo/AgenticRag/Agent/GrumpyAgent"
	"AIPracticePartner/AgentDemo/AgenticRag/RatingModel"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

var (
	BaseAgentTemplate = map[string]interface{}{
		"name":                "陈晓明",
		"age":                 30,
		"gender":              "男性",
		"job":                 "程序员",
		"family_status":       "未婚未育，居住在一线城市，有房贷压力。父母退休且养老金微薄，家庭抗风险能力较弱。", //家庭状态
		"psychological_state": "因职业特性存在裁员焦虑，具备极强的理财与忧患意识，正考虑构建安全底座。",      //心理状态
		"personality_traits":  "直爽理性，脾气好但极度注重细节。对逻辑闭环、数据准确性要求极高。",         //性格特征
		"chat_history":        "",                                         //历史记录
		"question":            "",                                         //提问
		"mission": "有较大额的闲置资金，想了解一个稳健收益的保险理财渠道，" +
			"为长期养老做打算、做多元财富配置、同时避免失业风险等。" +
			"希望了解产品的缴费模式、年化收益率（分红）、支持缴费的货币币种类型、是否允许中途取现、犹豫期和缴费期内退保的条款说明。", //任务
	}
	GentleAgentPersona = map[string]interface{}{
		"name":                "李温和",
		"age":                 28,
		"gender":              "女性",
		"job":                 "行政助理",
		"family_status":       "已婚未育，家庭收入一般，主要依靠丈夫收入。性格内向，不喜欢做决定。",                  //家庭状态
		"psychological_state": "对未来充满不确定性，想买保险但又怕钱打水漂，耳根子软，容易被别人的意见左右，需要极强的信任建立过程。", //心理状态
		"personality_traits":  "性格温和、优柔寡断、胆小慎重。不敢直接拒绝人，但会用各种理由推脱，需要反复确认安全感。",        //性格特征
		"chat_history":        "",                                                   //历史记录
		"question":            "",                                                   //提问
		"mission": "场景任务：向一位优柔寡断、耳根子软但又极其谨慎的客户推销“守护一生终身寿险”。" +
			"关键异议 (Key Objection)：“我得回去跟我老公商量一下，万一以后急用钱取不出来怎么办？”" +
			"通关标准 (Success Criteria)：建立足够的信任感，通过“保单贷款”消除其对资金灵活性的担忧，并引导其迈出决策的一步。", //任务
	}
	GrumpyAgentPersona = map[string]interface{}{
		"name":                "张暴躁",
		"age":                 45,
		"gender":              "男性",
		"job":                 "小企业主",
		"family_status":       "离异独居，经济条件较好但资金周转压力大。曾经被理财产品坑过，对金融产品有天然的抵触。",          //家庭状态
		"psychological_state": "像个火药桶，认为保险都是骗人的。只相信白纸黑字的条款，对任何营销话术都极其反感，随时准备挂电话或骂人。", //心理状态
		"personality_traits":  "暴躁、多疑、挑剔、极度不耐烦。喜欢打断别人，说话带刺，攻击性强。",                    //性格特征
		"chat_history":        "",                                                    //历史记录
		"question":            "",                                                    //提问
		"mission": "场景任务：向一位对保险充满敌意、极其不耐烦的客户推销“守护一生终身寿险”。" +
			"关键异议 (Key Objection)：“少跟我扯那些没用的，直接说重点！你们这些卖保险的嘴里没一句实话！”" +
			"通关标准 (Success Criteria)：在客户的高压攻击下保持专业，用精准的数据和条款条款（而非话术）折服客户，证明产品的真实收益。", //任务
	}
)

func TestModel(ctx context.Context) {
	// 选择 Agent
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("请选择要模拟的客户类型：")
	fmt.Println("1. 标准客户 (陈晓明 - 理性慎重)")
	fmt.Println("2. 温和客户 (李温和 - 优柔寡断)")
	fmt.Println("3. 暴躁客户 (张暴躁 - 多疑挑剔)")
	fmt.Print("请输入选项 (1-3): ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var runnable compose.Runnable[map[string]any, *schema.Message]
	var selectedAgent map[string]interface{}
	var err error

	switch choice {
	case "2":
		fmt.Println(">>> 已选择：温和客户 (GentleAgent)")
		runnable, err = GentleAgent.BuildGentleAgent(ctx)
		selectedAgent = GentleAgentPersona
	case "3":
		fmt.Println(">>> 已选择：暴躁客户 (GrumpyAgent)")
		runnable, err = GrumpyAgent.BuildGrumpyAgent(ctx)
		selectedAgent = GrumpyAgentPersona
	default:
		fmt.Println(">>> 已选择：标准客户 (BaseAgent)")
		runnable, err = BaseAgent.BuildAIDemo(ctx)
		selectedAgent = BaseAgentTemplate
	}

	if err != nil {
		log.Fatalf("Build graph failed: %v", err)
	}

	// 初始化历史记录
	var history []*schema.Message

	fmt.Println(">>> 练习开始。请输入您的问题（输入 'exit' 或 '结束' 停止练习）：")

	// ==========================================
	// 场景一：循环对话
	// ==========================================
	for {
		fmt.Print("\nUser: ")
		userInput, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Println("读取输入错误:", err)
			continue
		}
		userInput = strings.TrimSpace(userInput)

		if userInput == "" {
			continue
		}

		if userInput == "exit" || userInput == "quit" || userInput == "结束" {
			break
		}

		selectedAgent["question"] = userInput
		selectedAgent["chat_history"] = history
		fmt.Print("AI: \n")
		log.Println("\n正在生成中。。。。")
		start := time.Now()
		var stream *schema.StreamReader[*schema.Message]
		stream, err = runnable.Stream(ctx, selectedAgent)
		if err != nil {
			log.Println("Stream error:", err)
			continue
		}
		fullResponse := ""
		// 读取流
		for {
			chunk, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				log.Println("Stream recv error:", err)
				break
			}
			// 流式输出
			fmt.Print(chunk.Content)
			fullResponse += chunk.Content
		}
		fmt.Println("\n总时间：", time.Since(start))
		// 更新历史记录
		history = append(history, &schema.Message{
			Role:    schema.User,
			Content: userInput,
		})
		history = append(history, &schema.Message{
			Role:    schema.Assistant,
			Content: fullResponse,
		})
	}

	// ==========================================
	// 场景二：用户点击“结束练习”
	// ==========================================
	fmt.Println("\n>>> 用户点击结束...")

	//记录对话次数
	historyLen := 0
	for _, msg := range history {
		if msg.Role == schema.Assistant {
			historyLen++
		}
	}

	fmt.Println("对话次数：", historyLen)

	//进入对话评测
	if historyLen >= 3 {
		selectedAgent["question"] = ""
		selectedAgent["chat_history"] = ""
		log.Println("正在评分中。。。。。")
		model, err := RatingModel.BuildRatingModel(ctx)
		if err != nil {
			return
		}
		input := map[string]interface{}{
			"persona_summary": selectedAgent,
			"history":         history,
		}
		invoke, err := model.Invoke(ctx, input)
		if err != nil {
			fmt.Println(err)
			return
		}
		var prettyJSON bytes.Buffer
		if err = json.Indent(&prettyJSON, []byte(invoke.Content), "", "  "); err != nil {
			// 直接打印原始内容
			fmt.Println(invoke.Content)
		} else {
			fmt.Println(prettyJSON.String())
		}
	}
	fmt.Println("\n>>> 练习已完全结束")
}
