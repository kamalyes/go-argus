/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
 * @FilePath: \go-argus\rule\time_test.go
 * @Description: time.go 测试，覆盖时间值识别、表达式解析和跨字段时间比较
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package rule

import (
	"reflect"
	"testing"
	"time"
)

func TestTimeValueFromTime(t *testing.T) {
	now := time.Now()
	result, ok := TimeValue(reflect.ValueOf(now), "")
	if !ok || !result.Equal(now) {
		t.Fatalf("expected time.Time to be recognized, ok=%v", ok)
	}
}

func TestTimeValueFromString(t *testing.T) {
	result, ok := TimeValue(reflect.ValueOf("2024-01-01T00:00:00Z"), "")
	if !ok {
		t.Fatal("expected RFC3339 string to be parsed")
	}
	if result.Year() != 2024 {
		t.Fatalf("unexpected year: %d", result.Year())
	}
}

func TestTimeValueFromStringWithLayout(t *testing.T) {
	result, ok := TimeValue(reflect.ValueOf("2024-06-15"), "2006-01-02")
	if !ok {
		t.Fatal("expected custom layout to be parsed")
	}
	if result.Year() != 2024 || result.Month() != 6 {
		t.Fatalf("unexpected date: %v", result)
	}
}

func TestTimeValueFromStringEmpty(t *testing.T) {
	_, ok := TimeValue(reflect.ValueOf(""), "")
	if ok {
		t.Fatal("expected empty string to fail")
	}
}

func TestTimeValueInvalid(t *testing.T) {
	_, ok := TimeValue(reflect.ValueOf(12345), "")
	if ok {
		t.Fatal("expected int to fail time parsing")
	}
}

func TestTimeValueInvalidReflect(t *testing.T) {
	_, ok := TimeValue(reflect.Value{}, "")
	if ok {
		t.Fatal("expected invalid reflect.Value to fail")
	}
}

func TestTimeValueCannotInterfacePath(t *testing.T) {
	v := reflect.ValueOf(42)
	_, ok := TimeValue(v, "")
	if ok {
		t.Fatal("expected int to fail TimeValue")
	}
}

type mockTimestamp struct {
	seconds int64
	nanos   int64
}

func (m mockTimestamp) GetSeconds() int64 { return m.seconds }
func (m mockTimestamp) GetNanos() int64   { return m.nanos }

func TestTimeValueCallSecondsNanosValid(t *testing.T) {
	ts := mockTimestamp{seconds: 1700000000, nanos: 500}
	result, ok := TimeValue(reflect.ValueOf(ts), "")
	if !ok {
		t.Fatal("expected callSecondsNanos to work")
	}
	if result.Unix() != 1700000000 {
		t.Fatalf("unexpected unix: %d", result.Unix())
	}
}

type mockTimestampNoNanos struct {
	seconds int64
}

func (m mockTimestampNoNanos) GetSeconds() int64 { return m.seconds }

func TestTimeValueCallSecondsNanosNoNanos(t *testing.T) {
	ts := mockTimestampNoNanos{seconds: 1700000000}
	result, ok := TimeValue(reflect.ValueOf(ts), "")
	if !ok {
		t.Fatal("expected callSecondsNanos to work without GetNanos")
	}
	if result.Unix() != 1700000000 {
		t.Fatalf("unexpected unix: %d", result.Unix())
	}
}

type mockAsTime struct {
	t time.Time
}

func (m mockAsTime) AsTime() time.Time { return m.t }

func TestTimeValueCallAsTimeValidMethod(t *testing.T) {
	now := time.Now()
	ts := mockAsTime{t: now}
	result, ok := TimeValue(reflect.ValueOf(ts), "")
	if !ok {
		t.Fatal("expected AsTime method to work")
	}
	if !result.Equal(now) {
		t.Fatalf("expected %v, got %v", now, result)
	}
}

type mockAsTimeBadReturn struct{}

func (m mockAsTimeBadReturn) AsTime() string { return "not-time" }

func TestTimeValueCallAsTimeBadReturn(t *testing.T) {
	ts := mockAsTimeBadReturn{}
	_, ok := TimeValue(reflect.ValueOf(ts), "")
	if ok {
		t.Fatal("expected AsTime returning non-time.Time to fail")
	}
}

type mockGetSecondsBadReturn struct{}

func (m mockGetSecondsBadReturn) GetSeconds() string { return "bad" }

func TestTimeValueCallSecondsNanosBadReturn(t *testing.T) {
	ts := mockGetSecondsBadReturn{}
	_, ok := TimeValue(reflect.ValueOf(ts), "")
	if ok {
		t.Fatal("expected GetSeconds returning non-int to fail")
	}
}

type timeField struct{}

func (f timeField) AsTime() time.Time { return time.Time{} }

func TestCallAsTimeCannotInterface(t *testing.T) {
	type unsafeStruct struct {
		unexported timeField
	}
	s := unsafeStruct{}
	v := reflect.ValueOf(s).Field(0)
	_, ok := TimeValue(v, "")
	if ok {
		t.Fatal("expected unexported field to fail TimeValue")
	}
}

func TestCallSecondsNanosRecover(t *testing.T) {
	type withGetSeconds struct {
		unexported mockTimestamp
	}
	s := withGetSeconds{unexported: mockTimestamp{seconds: 100}}
	v := reflect.ValueOf(s).Field(0)
	_, ok := TimeValue(v, "")
	if ok {
		t.Fatal("expected unexported GetSeconds to fail via recover")
	}
}

func TestTimeValueReadSecondsNanos(t *testing.T) {
	type fakeTimestamp struct {
		Seconds int64
		Nanos   int64
	}
	ts := fakeTimestamp{Seconds: 1700000000, Nanos: 0}
	result, ok := TimeValue(reflect.ValueOf(ts), "")
	if !ok {
		t.Fatal("expected readSecondsNanos to work")
	}
	if result.Unix() != 1700000000 {
		t.Fatalf("unexpected unix timestamp: %d", result.Unix())
	}
}

func TestTimeValueReadSecondsNanosNoNanos(t *testing.T) {
	type partialTimestamp struct {
		Seconds int64
	}
	ts := partialTimestamp{Seconds: 1700000000}
	result, ok := TimeValue(reflect.ValueOf(ts), "")
	if !ok {
		t.Fatal("expected readSecondsNanos to work without Nanos")
	}
	if result.Unix() != 1700000000 {
		t.Fatalf("unexpected unix timestamp: %d", result.Unix())
	}
}

func TestTimeValueReadSecondsNanosInvalid(t *testing.T) {
	type notTimestamp struct {
		Name string
	}
	_, ok := TimeValue(reflect.ValueOf(notTimestamp{Name: "bad"}), "")
	if ok {
		t.Fatal("expected non-timestamp struct to fail")
	}
}

func TestTimeValueCallAsTimeNoMethod(t *testing.T) {
	type protoTimestamp struct{}
	ts := protoTimestamp{}
	result, ok := TimeValue(reflect.ValueOf(ts), "")
	if ok {
		t.Fatalf("expected no AsTime method, but got ok=true result=%v", result)
	}
}

func TestResolveTimeExprNow(t *testing.T) {
	now := time.Now()
	result, ok := ResolveTimeExpr("now", now)
	if !ok || !result.Equal(now) {
		t.Fatalf("expected 'now' to return current time, ok=%v", ok)
	}
}

func TestResolveTimeExprEmpty(t *testing.T) {
	_, ok := ResolveTimeExpr("", time.Now())
	if ok {
		t.Fatal("expected empty expr to fail")
	}
}

func TestResolveTimeExprPlusDuration(t *testing.T) {
	now := time.Now()
	result, ok := ResolveTimeExpr("now+5m", now)
	if !ok {
		t.Fatal("expected now+5m to parse")
	}
	expected := now.Add(5 * time.Minute)
	if !result.Equal(expected) {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestResolveTimeExprMinusDays(t *testing.T) {
	now := time.Now()
	result, ok := ResolveTimeExpr("now-3d", now)
	if !ok {
		t.Fatal("expected now-3d to parse")
	}
	expected := now.Add(-3 * 24 * time.Hour)
	if !result.Equal(expected) {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestResolveTimeExprInvalidDuration(t *testing.T) {
	_, ok := ResolveTimeExpr("now+abc", time.Now())
	if ok {
		t.Fatal("expected invalid duration to fail")
	}
}

func TestResolveTimeExprInvalidPrefix(t *testing.T) {
	_, ok := ResolveTimeExpr("tomorrow", time.Now())
	if ok {
		t.Fatal("expected invalid prefix to fail")
	}
}

func TestCompareTimeExpr(t *testing.T) {
	now := time.Now()
	future := now.Add(time.Hour)
	result := CompareTimeExpr(reflect.ValueOf(future), "now", "gt", now)
	if !result {
		t.Fatal("expected future > now")
	}
}

func TestCompareTimeExprInvalidField(t *testing.T) {
	now := time.Now()
	result := CompareTimeExpr(reflect.ValueOf("not-a-time"), "now", "gt", now)
	if result {
		t.Fatal("expected invalid field to fail")
	}
}

func TestParseDurationDays(t *testing.T) {
	d, ok := parseDuration("5d")
	if !ok || d != 5*24*time.Hour {
		t.Fatalf("expected 5 days, got %v ok=%v", d, ok)
	}
}

func TestParseDurationInvalidDays(t *testing.T) {
	_, ok := parseDuration("abcd")
	if ok {
		t.Fatal("expected invalid days to fail")
	}
}

func TestParseDurationEmpty(t *testing.T) {
	_, ok := parseDuration("")
	if ok {
		t.Fatal("expected empty duration to fail")
	}
}

func TestParseDurationStandard(t *testing.T) {
	d, ok := parseDuration("1h30m")
	if !ok || d != 90*time.Minute {
		t.Fatalf("expected 1h30m, got %v ok=%v", d, ok)
	}
}

func TestParseTimeStringEmpty(t *testing.T) {
	_, ok := parseTimeString("", "")
	if ok {
		t.Fatal("expected empty string to fail")
	}
}

func TestParseTimeStringWithLayout(t *testing.T) {
	result, ok := parseTimeString("2024-01-15", "2006-01-02")
	if !ok || result.Year() != 2024 {
		t.Fatalf("expected 2024, ok=%v", ok)
	}
}

func TestParseTimeStringDefaultFormats(t *testing.T) {
	formats := []string{
		"2024-01-01T00:00:00Z",
		"2024-01-01T00:00:00+08:00",
		"2024-01-01 00:00:00",
		"2024-01-01",
	}
	for _, f := range formats {
		_, ok := parseTimeString(f, "")
		if !ok {
			t.Fatalf("expected %q to parse", f)
		}
	}
}

func TestParseTimeStringInvalid(t *testing.T) {
	_, ok := parseTimeString("not-a-date", "")
	if ok {
		t.Fatal("expected invalid date to fail")
	}
}
