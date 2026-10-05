// Command chat-link-parse 扫描聊天文本里的引用，用对齐表格输出。
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"example.com/chattext/chattext"
)

const maxRefs = 8

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("chat-link-parse", flag.ContinueOnError)
	flags.SetOutput(stderr)
	input := flags.String("input", "testdata/chat.txt", "聊天文本路径")
	sample := flags.String("sample", "", "escape / kind / case / quota / file")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	switch *sample {
	case "escape":
		refs, _ := chattext.Scan("[[item:sword_01]] 字面量")
		fmt.Fprintf(stdout, "refs=%d\n", len(refs))
	case "kind":
		refs, bad := chattext.Scan("[mount:horse_01] 未知类型")
		fmt.Fprintf(stdout, "refs=%d bad=%d\n", len(refs), bad)
	case "case":
		refs, _ := chattext.Scan("[item:SWORD_01]")
		id := ""
		if len(refs) > 0 {
			id = refs[0].ID
		}
		fmt.Fprintf(stdout, "id=%s\n", id)
	case "quota":
		line := strings.TrimSpace(strings.Repeat("[item:a1] ", 12))
		refs, _, truncated := chattext.ScanLine(line)
		fmt.Fprintf(stdout, "refs=%d truncated=%d\n", len(refs), truncated)
	case "file", "":
		handle, err := os.Open(*input)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		defer handle.Close()
		writer := tabwriter.NewWriter(stdout, 0, 4, 2, ' ', 0)
		fmt.Fprintln(writer, "line\trefs\tbad")
		scanner := bufio.NewScanner(handle)
		index := 0
		for scanner.Scan() {
			index++
			refs, bad := chattext.Scan(scanner.Text())
			fmt.Fprintf(writer, "%d\t%d\t%d\n", index, len(refs), bad)
		}
		writer.Flush()
	default:
		fmt.Fprintln(stderr, "需要 --sample escape|kind|case|quota|file")
		return 2
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
