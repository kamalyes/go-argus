/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-25 17:35:19
 * @FilePath: \go-argus\varslot_test.go
 * @Description: Var/VarString 分槽旁路缓存测试
 *
 * 守护点：规则 tag 高度同构、普遍共享前后缀，槽位哈希必须采样内部字节。
 * 若退化为仅用长度+首尾字节，形如 "min=2,max=50" / "min=3,max=40" 的 tag
 * 会全部撞进同一槽位，分槽退化为单槽——每次轮转都未命中、每次都分配。
 */

package validator

import (
	"testing"
)

// TestVarSlotIndex_PrefixSuffixSharedTagsSpread 共享前后缀的 tag 必须散开
func TestVarSlotIndex_PrefixSuffixSharedTagsSpread(t *testing.T) {
	// 长度、首字节、末字节两两相同，差异仅在中间——最容易撞槽的真实形态
	tags := []string{
		"min=2,max=50",
		"min=3,max=40",
		"min=4,max=30",
		"min=5,max=20",
	}

	seen := make(map[uint][]string, len(tags))
	for _, tag := range tags {
		idx := varSlotIndex(tag)
		seen[idx] = append(seen[idx], tag)
	}

	if len(seen) < 3 {
		t.Fatalf("4 个共享前后缀的 tag 仅散到 %d 个槽位，哈希未采样内部字节: %v",
			len(seen), seen)
	}
}

// TestVarSlotIndex_InRange 下标必须落在槽位范围内
func TestVarSlotIndex_InRange(t *testing.T) {
	tags := []string{
		"", "a", "ab", "required", "required,email",
		"min=1,max=100", "omitempty,url", "gte=0,lte=130",
		"oneof=red green blue", "uuid4", "len=8", "ipv4",
	}
	for _, tag := range tags {
		if idx := varSlotIndex(tag); idx >= varSlotCount {
			t.Fatalf("tag %q 槽位下标 %d 越界（上限 %d）", tag, idx, varSlotCount)
		}
	}
}

// TestVarSlotIndex_Deterministic 同一 tag 必须稳定映射到同一槽位
func TestVarSlotIndex_Deterministic(t *testing.T) {
	const tag = "required,email"
	want := varSlotIndex(tag)
	for i := 0; i < 100; i++ {
		if got := varSlotIndex(tag); got != want {
			t.Fatalf("tag %q 第 %d 次映射到 %d，首次为 %d", tag, i, got, want)
		}
	}
}

// TestVarSlotCount_PowerOfTwo 槽位数必须为 2 的幂（下标用掩码取）
func TestVarSlotCount_PowerOfTwo(t *testing.T) {
	if varSlotCount == 0 || varSlotCount&(varSlotCount-1) != 0 {
		t.Fatalf("varSlotCount=%d 不是 2 的幂，掩码取下标会越界", varSlotCount)
	}
}

// TestVarSlot_RotatingTagsCorrectness 轮转多 tag 的判定结果必须正确
//
// 分槽改造不得让某个 tag 读到另一个 tag 的规则（槽内 tag 比对是正确性闸门）。
func TestVarSlot_RotatingTagsCorrectness(t *testing.T) {
	v := New()

	cases := []struct {
		tag     string
		value   string
		wantErr bool
	}{
		{"required,email", "k@example.com", false},
		{"required,email", "not-an-email", true},
		{"required,numeric", "12345", false},
		{"required,numeric", "abc", true},
		{"required,alpha", "abc", false},
		{"required,alpha", "123", true},
		{"min=2,max=50", "hello", false},
		{"min=2,max=50", "x", true},
	}

	// 轮转多轮，确保槽位被反复覆写后判定仍正确
	for round := 0; round < 5; round++ {
		for _, c := range cases {
			err := v.VarString(c.value, c.tag)
			if (err != nil) != c.wantErr {
				t.Fatalf("第 %d 轮 tag=%q value=%q: wantErr=%v 实得 %v",
					round, c.tag, c.value, c.wantErr, err)
			}
		}
	}
}
