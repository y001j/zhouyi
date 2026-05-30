package qimen

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateAIPrompt_Smoke(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	pan, err := BuildPan(time.Date(2024, 6, 21, 12, 0, 0, 0, loc))
	if err != nil {
		t.Fatalf("BuildPan err: %v", err)
	}
	prompt := GenerateAIPrompt(pan, "近期工作运势如何", "", "career")
	if prompt == "" {
		t.Fatal("prompt 为空")
	}
	// 一些关键小节必须存在
	mustHave := []string{
		"## 起局时间与四柱",
		"## 值符 · 值使",
		"## 九宫盘",
		"## 九宫逐格详情",
		"## 请按以下结构解局",
		"夏至",
		"值符",
		"值使",
	}
	for _, m := range mustHave {
		if !strings.Contains(prompt, m) {
			t.Errorf("prompt 缺失 %q", m)
		}
	}
	t.Logf("prompt chars=%d lines=%d", len(prompt), strings.Count(prompt, "\n")+1)
}
