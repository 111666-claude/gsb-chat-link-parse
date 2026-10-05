# chat-link-parse

聊天引用抽取：从 `testdata` 下的聊天文本里识别 `[类型:ID]` 引用，用对齐表格输出，只用 Go 标准库。

```
go run ./cmd/chat-link-parse --sample=escape
go run ./cmd/chat-link-parse --sample=kind
go run ./cmd/chat-link-parse --sample=case
go run ./cmd/chat-link-parse --sample=quota
go run ./cmd/chat-link-parse --input testdata/chat.txt
go test ./...
go vet ./...
```

## 口径（README 为准）

- **引用形式**：`[类型:ID]`，类型只允许 `item`、`quest`、`achievement`，其它类型算无法识别。
- **转义**：`[[` 是字面量左括号，不开始一条引用。
- **ID 规范**：长度 2 到 16，只允许小写字母、数字与下划线；输入里的大写要先归一成小写再判。
- **限额**：一行最多抽取 8 条引用，超出的计入 `truncated`，不进入结果。
- **不变量**：同一行重复扫描结果相同；同一行的引用顺序按出现位置升序。
- **代价**：单行扫描是 O(行长)，不许对每条引用重扫整行；10 万行每秒 10 万行。

## 输出契约（不改格式）

```
refs=..
refs=.. bad=..
id=..
refs=.. truncated=..
line  refs  bad      （--input 时的对齐表格）
```

## 目录

```
chattext/scan.go                引用抽取
cmd/chat-link-parse/main.go     命令行入口与表格输出
testdata/chat.txt               聊天样例
chattext/scan_test.go           抽取用例
```
