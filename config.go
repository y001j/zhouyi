package main

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"sync"
)

// 配置文件路径（与可执行文件同目录的 ./config.json）。
const configPath = "./config.json"

// LLMConfig 是解卦所用大模型 API 的配置。
//
// 支持两种协议（由 provider 决定）：
//   - openai     ：OpenAI 兼容的 /chat/completions 接口。
//     可对接 OpenAI、DeepSeek、月之暗面 Kimi、智谱、通义、Ollama、各类中转网关。
//   - anthropic  ：Anthropic 原生 /v1/messages 接口（Claude 官方或兼容网关）。
//
// 留空 provider 时默认按 openai 处理。enabled=false 或 apiKey 为空时视为「未配置」，
// 解卦功能优雅降级为「仅生成提示词，由用户自行复制到 AI」。
type LLMConfig struct {
	Enabled     bool    `json:"enabled"`               // 是否启用程序内解卦（false 则仅出提示词）
	Provider    string  `json:"provider,omitempty"`    // openai | anthropic，默认 openai
	BaseURL     string  `json:"baseURL,omitempty"`     // API 根地址，留空按 provider 取默认
	APIKey      string  `json:"apiKey,omitempty"`      // 密钥（也可用环境变量 ZHOUYI_LLM_API_KEY 覆盖）
	Model       string  `json:"model,omitempty"`       // 模型名，留空按 provider 取默认
	Temperature float64 `json:"temperature,omitempty"` // 采样温度，0 表示用默认 0.7
	MaxTokens   int     `json:"maxTokens,omitempty"`   // 最大生成 token，0 表示用默认 2048
	TimeoutSec  int     `json:"timeoutSec,omitempty"`  // 单次请求超时秒数，0 表示用默认 120
}

// Config 是整个程序的统一配置（对应 config.json）。
// 历史上只有 adminPassword 一个字段；现扩展出 llm 段用于程序内解卦。
type Config struct {
	AdminPassword string    `json:"adminPassword,omitempty"`
	LLM           LLMConfig `json:"llm,omitempty"`
}

var (
	cfgOnce   sync.Once
	cfgCached Config
)

// LoadConfig 读取并缓存 ./config.json（仅首次真正读盘）。
// 文件不存在或解析失败时返回零值 Config（不报错，保证程序可裸跑）。
// 随后叠加环境变量覆盖（便于容器化部署不落盘密钥）。
func LoadConfig() Config {
	cfgOnce.Do(func() {
		cfgCached = readConfigFile()
		applyEnvOverrides(&cfgCached)
	})
	return cfgCached
}

func readConfigFile() Config {
	var c Config
	b, err := os.ReadFile(configPath)
	if err != nil {
		return c
	}
	_ = json.Unmarshal(b, &c) // 解析失败也只是返回零值，不阻断程序
	return c
}

// applyEnvOverrides 让环境变量优先于配置文件，便于部署时不落盘敏感信息。
//
//	ADMIN_PASSWORD          管理员密码（沿用历史变量名）
//	ZHOUYI_LLM_ENABLED      1/true/yes 启用解卦
//	ZHOUYI_LLM_PROVIDER     openai | anthropic
//	ZHOUYI_LLM_BASE_URL     API 根地址
//	ZHOUYI_LLM_API_KEY      密钥
//	ZHOUYI_LLM_MODEL        模型名
func applyEnvOverrides(c *Config) {
	if v := os.Getenv("ADMIN_PASSWORD"); v != "" {
		c.AdminPassword = v
	}
	if v := os.Getenv("ZHOUYI_LLM_PROVIDER"); v != "" {
		c.LLM.Provider = v
	}
	if v := os.Getenv("ZHOUYI_LLM_BASE_URL"); v != "" {
		c.LLM.BaseURL = v
	}
	if v := os.Getenv("ZHOUYI_LLM_API_KEY"); v != "" {
		c.LLM.APIKey = v
	}
	if v := os.Getenv("ZHOUYI_LLM_MODEL"); v != "" {
		c.LLM.Model = v
	}
	if v := os.Getenv("ZHOUYI_LLM_ENABLED"); v != "" {
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "1", "true", "yes", "on":
			c.LLM.Enabled = true
		case "0", "false", "no", "off":
			c.LLM.Enabled = false
		}
	}
	if v := os.Getenv("ZHOUYI_LLM_MAX_TOKENS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			c.LLM.MaxTokens = n
		}
	}
}

// resolved 把 LLMConfig 补齐为可用配置（填默认值），并报告是否「可解卦」。
// 可解卦的最低条件：apiKey 非空（baseURL/model 可走 provider 默认）。
func (l LLMConfig) resolved() (LLMConfig, bool) {
	out := l
	out.Provider = strings.ToLower(strings.TrimSpace(out.Provider))
	if out.Provider == "" {
		out.Provider = "openai"
	}
	if out.BaseURL == "" {
		if out.Provider == "anthropic" {
			out.BaseURL = "https://api.anthropic.com"
		} else {
			out.BaseURL = "https://api.openai.com/v1"
		}
	}
	out.BaseURL = strings.TrimRight(out.BaseURL, "/")
	if out.Model == "" {
		if out.Provider == "anthropic" {
			out.Model = "claude-sonnet-4-6"
		} else {
			out.Model = "gpt-4o-mini"
		}
	}
	if out.Temperature <= 0 {
		out.Temperature = 0.7
	}
	if out.MaxTokens <= 0 {
		out.MaxTokens = 2048
	}
	if out.TimeoutSec <= 0 {
		out.TimeoutSec = 120
	}
	usable := strings.TrimSpace(out.APIKey) != ""
	return out, usable
}
