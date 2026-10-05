// Package chattext 从聊天文本里抽取物品与任务引用。
// 缺陷：不认转义、未知类型也算引用、ID 长度不校验、大小写不归一、一行引用数不封顶。
package chattext

import "strings"

// Ref 是一条引用。
type Ref struct {
	Kind string
	ID   string
}

var kinds = map[string]bool{"item": true, "quest": true, "achievement": true}

// Scan 抽取一行里的引用，返回引用列表与无法识别的片段数。
// 缺陷：`[[` 不转义、未知类型也收、ID 长度不查、大写不归一、一行不限量。
func Scan(line string) ([]Ref, int) {
	refs := []Ref{}
	bad := 0
	rest := line
	for {
		start := strings.Index(rest, "[")
		if start < 0 {
			break
		}
		end := strings.Index(rest[start:], "]")
		if end < 0 {
			break
		}
		body := rest[start+1 : start+end]
		kind, id, _ := strings.Cut(body, ":")
		if id == "" {
			bad++
		} else if kinds[kind] || kind != "" {
			refs = append(refs, Ref{Kind: kind, ID: id})
		} else {
			bad++
		}
		rest = rest[start+end+1:]
	}
	return refs, bad
}
