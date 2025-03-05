package config

import (
	"fmt"
	"strings"
	"time"
)

// TTL 表示缓存的生存时间
type TTL struct {
	Value string `yaml:"value"` // 可以是 "5m", "2h", "1d", "permanent" 等
}

// Duration 将 TTL 转换为 time.Duration
// 支持的格式：
// - 数字+单位: "30s", "5m", "2h", "1d", "1w"
// - "permanent" 或 "0": 表示永久存储
func (t TTL) Duration() (time.Duration, bool) {
	if t.Value == "" {
		return 0, false
	}

	// 处理特殊值
	if t.Value == "permanent" || t.Value == "0" || t.Value == "-1" {
		return -1, true // -1 表示永久
	}

	// 从字符串中提取数字和单位
	var value string = t.Value
	var unit string

	// 找到最后一位数字的位置
	lastDigitIndex := -1
	for i, char := range value {
		if char >= '0' && char <= '9' {
			lastDigitIndex = i
		}
	}

	if lastDigitIndex >= 0 && lastDigitIndex < len(value)-1 {
		unit = strings.TrimSpace(value[lastDigitIndex+1:])
		value = strings.TrimSpace(value[:lastDigitIndex+1])
	}

	// 根据单位解析持续时间
	var duration time.Duration
	var err error

	if unit == "" {
		// 没有单位，默认为秒
		duration, err = time.ParseDuration(value + "s")
	} else {
		switch strings.ToLower(unit) {
		case "s", "sec", "second", "seconds":
			duration, err = time.ParseDuration(value + "s")
		case "m", "min", "minute", "minutes":
			duration, err = time.ParseDuration(value + "m")
		case "h", "hour", "hours":
			duration, err = time.ParseDuration(value + "h")
		case "d", "day", "days":
			duration, err = time.ParseDuration(value + "h")
			if err == nil {
				duration *= 24 // 转换为小时
			}
		case "w", "week", "weeks":
			duration, err = time.ParseDuration(value + "h")
			if err == nil {
				duration *= 24 * 7 // 转换为小时
			}
		default:
			err = fmt.Errorf("unknown time unit: %s", unit)
		}
	}

	if err != nil {
		return 0, false
	}

	return duration, true
}

// IsPermanent 检查是否是永久存储
func (t TTL) IsPermanent() bool {
	duration, ok := t.Duration()
	return ok && duration < 0
}
