package cmd

import (
	"os"
	"strings"
)

var operatorLocale = normalizeLocale(os.Getenv("PASTURESTACK_LOCALE"))

var operatorMessages = map[string]map[string]string{
	"en-US": {
		"app.usage":      "PastureStack CLI for compatible control-platform resources",
		"flag.config":    "Client configuration file (default ${HOME}/.pasturestack/cli.json)",
		"flag.url":       "Specify the compatible API endpoint URL",
		"flag.accessKey": "Specify the API access key",
		"flag.secretKey": "Specify the API secret key",
		"error.noEnv":    "Failed to find the current environment",
		"error.noURL":    "PLATFORM_URL or --url is not set; run 'config'",
		"config.saved":   "Saved configuration to %s",
		"select":         "Select",
		"environments":   "Environments",
	},
	"zh-TW": {
		"app.usage":      "管理相容控制平台資源的 PastureStack 命令列工具",
		"flag.config":    "用戶端設定檔（預設 ${HOME}/.pasturestack/cli.json）",
		"flag.url":       "指定相容 API 端點網址",
		"flag.accessKey": "指定 API 存取金鑰",
		"flag.secretKey": "指定 API 私密金鑰",
		"error.noEnv":    "找不到目前環境",
		"error.noURL":    "尚未設定 PLATFORM_URL 或 --url；請先執行 'config'",
		"config.saved":   "設定已儲存至 %s",
		"select":         "請選擇",
		"environments":   "環境",
	},
}

func normalizeLocale(value string) string {
	if strings.EqualFold(strings.TrimSpace(value), "zh-TW") {
		return "zh-TW"
	}
	return "en-US"
}

func SetLocale(value string) {
	operatorLocale = normalizeLocale(value)
}

func T(key string) string {
	if message := operatorMessages[operatorLocale][key]; message != "" {
		return message
	}
	if message := operatorMessages["en-US"][key]; message != "" {
		return message
	}
	return key
}
