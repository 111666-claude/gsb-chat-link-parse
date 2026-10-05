// Package chattext 从聊天文本里抽取物品与任务引用。
package chattext

import "strings"

// Ref 是一条引用。
type Ref struct {
	Kind string
	ID   string
}

var kinds = map[string]bool{"item": true, "quest": true, "achievement": true}

// maxRefsPerLine 是一行最多抽取的引用数，超出的计入 truncated。
const maxRefsPerLine = 8

// Scan 抽取一行里的引用，返回引用列表与无法识别的片段数。
func Scan(line string) ([]Ref, int) {
	refs, bad, _ := ScanLine(line)
	return refs, bad
}

// ScanLine 与 Scan 相同，但额外返回因超出单行限额而被截断的引用数。
// 单行一趟扫描，O(行长)。
func ScanLine(line string) (refs []Ref, bad, truncated int) {
	refs = []Ref{}
	for i := 0; i < len(line); {
		open := strings.IndexByte(line[i:], '[')
		if open < 0 {
			break
		}
		start := i + open
		if start+1 < len(line) && line[start+1] == '[' {
			// `[[` 是字面量左括号，不开始一条引用。
			i = start + 2
			continue
		}
		end := strings.IndexByte(line[start+1:], ']')
		if end < 0 {
			break
		}
		body := line[start+1 : start+1+end]
		i = start + 1 + end + 1
		kind, id, ok := strings.Cut(body, ":")
		if !ok || !kinds[kind] {
			bad++
			continue
		}
		id = strings.ToLower(id)
		if !validID(id) {
			bad++
			continue
		}
		if len(refs) >= maxRefsPerLine {
			truncated++
			continue
		}
		refs = append(refs, Ref{Kind: kind, ID: id})
	}
	return refs, bad, truncated
}

// validID 校验归一后的 ID：长度 2 到 16，只允许小写字母、数字与下划线。
func validID(id string) bool {
	if len(id) < 2 || len(id) > 16 {
		return false
	}
	for i := 0; i < len(id); i++ {
		c := id[i]
		if c != '_' && (c < 'a' || c > 'z') && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}
