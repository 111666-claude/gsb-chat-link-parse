package chattext

import "testing"

func TestScanFindsItemRef(t *testing.T) {
	refs, _ := Scan("出 [item:sword_01] 一件")
	if len(refs) != 1 || refs[0].Kind != "item" || refs[0].ID != "sword_01" {
		t.Fatalf("应该抽到一条 item 引用，得到 %+v", refs)
	}
}

func TestScanWithoutRefs(t *testing.T) {
	refs, bad := Scan("今天天气不错")
	if len(refs) != 0 || bad != 0 {
		t.Fatalf("没有引用时不该有结果，得到 %+v bad=%d", refs, bad)
	}
}

func TestScanMultipleRefs(t *testing.T) {
	refs, _ := Scan("[item:a1] 和 [quest:b2]")
	if len(refs) != 2 {
		t.Fatalf("应该抽到两条引用，得到 %+v", refs)
	}
}
