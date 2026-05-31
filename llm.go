package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// llm.go：把「解卦提示词」发给大模型 API，取回中文解读文本。
//
// 提示词本身已由各式的 GenerateAIPrompt / HuCanPrompt 生成（含卦象素材 +
// 解卦规则 + 输出结构要求），故这里把整段提示词作为 user 消息发出即可，
// 另配一句简短的 system 角色设定。支持 OpenAI 兼容与 Anthropic 原生两种协议。

// 解卦统一的 system 角色设定（两种协议共用）。
const interpretSystemPrompt = "你是一位精通周易、奇门遁甲、大六壬的中式术数顾问。" +
	"请严格依据用户给出的盘面素材与解卦规则作答，结论须标注依据，不臆造素材中没有的信息。" +
	"用中文回答，深入浅出，兼顾传统术数用语与现代通俗表达。"

// ErrLLMNotConfigured 表示未配置可用的解卦 API（apiKey 为空等）。
var ErrLLMNotConfigured = errors.New("未配置解卦 API：请在 config.json 的 llm 段填写 apiKey（或设置环境变量 ZHOUYI_LLM_API_KEY），并将 enabled 置为 true")

// Interpret 用配置好的大模型对一段解卦提示词作答，返回解读文本。
// cfg 为原始 LLMConfig（内部会自动补默认值）。
func Interpret(ctx context.Context, cfg LLMConfig, prompt string) (string, error) {
	rc, usable := cfg.resolved()
	if !usable {
		return "", ErrLLMNotConfigured
	}
	if strings.TrimSpace(prompt) == "" {
		return "", errors.New("提示词为空，无可解之卦")
	}

	reqCtx, cancel := context.WithTimeout(ctx, time.Duration(rc.TimeoutSec)*time.Second)
	defer cancel()

	switch rc.Provider {
	case "anthropic":
		return callAnthropic(reqCtx, rc, prompt)
	default:
		return callOpenAI(reqCtx, rc, prompt)
	}
}

// httpClient 复用一个带超时兜底的客户端（真正的超时由 context 控制）。
var httpClient = &http.Client{Timeout: 0}

// ===== OpenAI 兼容协议（/chat/completions） =====

type openAIChatRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatResponse struct {
	Choices []struct {
		Message openAIMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

func callOpenAI(ctx context.Context, cfg LLMConfig, prompt string) (string, error) {
	body, _ := json.Marshal(openAIChatRequest{
		Model:       cfg.Model,
		Temperature: cfg.Temperature,
		MaxTokens:   cfg.MaxTokens,
		Messages: []openAIMessage{
			{Role: "system", Content: interpretSystemPrompt},
			{Role: "user", Content: prompt},
		},
	})

	url := cfg.BaseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	raw, status, err := doRequest(req)
	if err != nil {
		return "", err
	}

	var out openAIChatResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return "", fmt.Errorf("解析响应失败（HTTP %d）：%v", status, err)
	}
	if out.Error != nil {
		return "", fmt.Errorf("API 报错（HTTP %d）：%s", status, out.Error.Message)
	}
	if status < 200 || status >= 300 {
		return "", fmt.Errorf("API 返回 HTTP %d：%s", status, truncate(string(raw), 300))
	}
	if len(out.Choices) == 0 {
		return "", errors.New("API 未返回任何结果")
	}
	text := strings.TrimSpace(out.Choices[0].Message.Content)
	if text == "" {
		return "", errors.New("API 返回内容为空")
	}
	return text, nil
}

// ===== Anthropic 原生协议（/v1/messages） =====

type anthropicRequest struct {
	Model       string             `json:"model"`
	MaxTokens   int                `json:"max_tokens"`
	Temperature float64            `json:"temperature"`
	System      string             `json:"system,omitempty"`
	Messages    []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

func callAnthropic(ctx context.Context, cfg LLMConfig, prompt string) (string, error) {
	body, _ := json.Marshal(anthropicRequest{
		Model:       cfg.Model,
		MaxTokens:   cfg.MaxTokens,
		Temperature: cfg.Temperature,
		System:      interpretSystemPrompt,
		Messages: []anthropicMessage{
			{Role: "user", Content: prompt},
		},
	})

	// baseURL 若已含 /v1 则不重复拼接，兼容用户直接填到 /v1 的情况。
	base := cfg.BaseURL
	path := "/v1/messages"
	if strings.HasSuffix(base, "/v1") {
		path = "/messages"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+path, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", cfg.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	raw, status, err := doRequest(req)
	if err != nil {
		return "", err
	}

	var out anthropicResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return "", fmt.Errorf("解析响应失败（HTTP %d）：%v", status, err)
	}
	if out.Error != nil {
		return "", fmt.Errorf("API 报错（HTTP %d）：%s", status, out.Error.Message)
	}
	if status < 200 || status >= 300 {
		return "", fmt.Errorf("API 返回 HTTP %d：%s", status, truncate(string(raw), 300))
	}
	var sb strings.Builder
	for _, c := range out.Content {
		if c.Type == "text" {
			sb.WriteString(c.Text)
		}
	}
	text := strings.TrimSpace(sb.String())
	if text == "" {
		return "", errors.New("API 返回内容为空")
	}
	return text, nil
}

// ===== 公共 =====

func doRequest(req *http.Request) ([]byte, int, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		// context 超时给一句更友好的提示
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, 0, errors.New("请求超时：模型未在限定时间内返回，可调大 config.json 的 llm.timeoutSec")
		}
		return nil, 0, fmt.Errorf("请求失败：%v", err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("读取响应失败：%v", err)
	}
	return raw, resp.StatusCode, nil
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}

// ===== 流式解卦 =====

// InterpretStream 以流式方式调用大模型解卦，每收到一段增量文本就回调 onDelta。
// 双协议各自解析自家的 SSE 流（OpenAI: data: {...delta.content}；
// Anthropic: content_block_delta {...delta.text}）。
// 返回 nil 表示正常结束；onDelta 收到的是增量片段（非全文）。
func InterpretStream(ctx context.Context, cfg LLMConfig, prompt string, onDelta func(string)) error {
	rc, usable := cfg.resolved()
	if !usable {
		return ErrLLMNotConfigured
	}
	if strings.TrimSpace(prompt) == "" {
		return errors.New("提示词为空，无可解之卦")
	}

	reqCtx, cancel := context.WithTimeout(ctx, time.Duration(rc.TimeoutSec)*time.Second)
	defer cancel()

	switch rc.Provider {
	case "anthropic":
		return streamAnthropic(reqCtx, rc, prompt, onDelta)
	default:
		return streamOpenAI(reqCtx, rc, prompt, onDelta)
	}
}

// openAIStreamRequest 在普通请求基础上加 stream=true。
type openAIStreamRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens"`
	Stream      bool            `json:"stream"`
}

func streamOpenAI(ctx context.Context, cfg LLMConfig, prompt string, onDelta func(string)) error {
	body, _ := json.Marshal(openAIStreamRequest{
		Model:       cfg.Model,
		Temperature: cfg.Temperature,
		MaxTokens:   cfg.MaxTokens,
		Stream:      true,
		Messages: []openAIMessage{
			{Role: "system", Content: interpretSystemPrompt},
			{Role: "user", Content: prompt},
		},
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.BaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Accept", "text/event-stream")

	resp, err := httpClient.Do(req)
	if err != nil {
		return streamReqErr(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return httpErrFromBody(resp)
	}

	got := false
	err = scanSSE(resp.Body, func(data string) bool {
		if data == "[DONE]" {
			return false
		}
		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if json.Unmarshal([]byte(data), &chunk) != nil {
			return true
		}
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			got = true
			onDelta(chunk.Choices[0].Delta.Content)
		}
		return true
	})
	if err != nil {
		return err
	}
	if !got {
		return errors.New("API 返回内容为空")
	}
	return nil
}

// anthropicStreamRequest 在普通请求基础上加 stream=true。
type anthropicStreamRequest struct {
	Model       string             `json:"model"`
	MaxTokens   int                `json:"max_tokens"`
	Temperature float64            `json:"temperature"`
	System      string             `json:"system,omitempty"`
	Messages    []anthropicMessage `json:"messages"`
	Stream      bool               `json:"stream"`
}

func streamAnthropic(ctx context.Context, cfg LLMConfig, prompt string, onDelta func(string)) error {
	body, _ := json.Marshal(anthropicStreamRequest{
		Model:       cfg.Model,
		MaxTokens:   cfg.MaxTokens,
		Temperature: cfg.Temperature,
		System:      interpretSystemPrompt,
		Stream:      true,
		Messages: []anthropicMessage{
			{Role: "user", Content: prompt},
		},
	})

	base := cfg.BaseURL
	path := "/v1/messages"
	if strings.HasSuffix(base, "/v1") {
		path = "/messages"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", cfg.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := httpClient.Do(req)
	if err != nil {
		return streamReqErr(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return httpErrFromBody(resp)
	}

	got := false
	err = scanSSE(resp.Body, func(data string) bool {
		var ev struct {
			Type  string `json:"type"`
			Delta struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"delta"`
		}
		if json.Unmarshal([]byte(data), &ev) != nil {
			return true
		}
		if ev.Type == "content_block_delta" && ev.Delta.Text != "" {
			got = true
			onDelta(ev.Delta.Text)
		}
		if ev.Type == "message_stop" {
			return false
		}
		return true
	})
	if err != nil {
		return err
	}
	if !got {
		return errors.New("API 返回内容为空")
	}
	return nil
}

// scanSSE 逐行读取 SSE 流，对每个「data: 」行的负载调用 fn。
// fn 返回 false 表示主动停止扫描（如收到 [DONE] / message_stop）。
func scanSSE(r io.Reader, fn func(data string) bool) error {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024) // 放大单行上限，防超长 data 行截断
	for sc.Scan() {
		line := sc.Text()
		if !strings.HasPrefix(line, "data:") {
			continue // 跳过 event:/空行/注释行
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "" {
			continue
		}
		if !fn(data) {
			return nil
		}
	}
	if err := sc.Err(); err != nil {
		return streamReqErr(err)
	}
	return nil
}

// streamReqErr 把底层请求错误翻译成友好提示。
func streamReqErr(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return errors.New("请求超时：模型未在限定时间内返回，可调大 config.json 的 llm.timeoutSec")
	}
	return fmt.Errorf("请求失败：%v", err)
}

// httpErrFromBody 读取非 2xx 响应体并构造错误（尝试解析常见错误结构）。
func httpErrFromBody(resp *http.Response) error {
	raw, _ := io.ReadAll(resp.Body)
	var e struct {
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if json.Unmarshal(raw, &e) == nil && e.Error != nil && e.Error.Message != "" {
		return fmt.Errorf("API 报错（HTTP %d）：%s", resp.StatusCode, e.Error.Message)
	}
	return fmt.Errorf("API 返回 HTTP %d：%s", resp.StatusCode, truncate(string(raw), 300))
}
