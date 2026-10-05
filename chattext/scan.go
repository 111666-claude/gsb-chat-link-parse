// Package chattext 从聊天文本里抽取物品与任务引用。
package chattext

import "strings"

// maxRefsPerLine 是一行最多抽取的引用数，超出的计入 truncated。
const maxRefsPerLine = 8

// Ref 是一条引用。
type Ref struct {
	Kind string
	ID   string
}

var kinds = map[string]bool{"item": true, "quest": true, "achievement": true}

// Scan 抽取一行里的引用，返回引用列表与无法识别的片段数。
// 等价于 ScanLine 丢弃截断计数。
func Scan(line string) ([]Ref, int) {
	refs, bad, _ := ScanLine(line)
	return refs, bad
}

// ScanLine 抽取一行里的引用，返回引用列表、无法识别的片段数与超限额被截断的引用数。
// 规则：`[[` 是字面量左括号；类型只认 item、quest、achievement；
// ID 长度 2 到 16，只允许小写字母、数字与下划线，大写先归一成小写；
// 一行最多 8 条引用，超出的计入 truncated。单行扫描 O(行长)。
func ScanLine(line string) ([]Ref, int, int) {
	refs := make([]Ref, 0, maxRefsPerLine)
	bad := 0
	truncated := 0
	rest := line
	for {
		start := strings.Index(rest, "[")
		if start < 0 {
			break
		}
		if start+1 < len(rest) && rest[start+1] == '[' {
			rest = rest[start+2:]
			continue
		}
		end := strings.Index(rest[start:], "]")
		if end < 0 {
			break
		}
		body := rest[start+1 : start+end]
		kind, id, _ := strings.Cut(body, ":")
		if id == "" || !kinds[kind] {
			bad++
		} else {
			id = normalizeID(id)
			if validID(id) {
				if len(refs) < maxRefsPerLine {
					refs = append(refs, Ref{Kind: kind, ID: id})
				} else {
					truncated++
				}
			} else {
				bad++
			}
		}
		rest = rest[start+end+1:]
	}
	return refs, bad, truncated
}

// normalizeID 把 ID 里的大写字母归一成小写。
func normalizeID(id string) string {
	upper := false
	for i := 0; i < len(id); i++ {
		if id[i] >= 'A' && id[i] <= 'Z' {
			upper = true
			break
		}
	}
	if !upper {
		return id
	}
	lowered := []byte(id)
	for i, c := range lowered {
		if c >= 'A' && c <= 'Z' {
			lowered[i] = c + ('a' - 'A')
		}
	}
	return string(lowered)
}

// validID 校验归一后的 ID：长度 2 到 16，只允许小写字母、数字与下划线。
func validID(id string) bool {
	if len(id) < 2 || len(id) > 16 {
		return false
	}
	for i := 0; i < len(id); i++ {
		c := id[i]
		if (c < 'a' || c > 'z') && (c < '0' || c > '9') && c != '_' {
			return false
		}
	}
	return true
}
