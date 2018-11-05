package autocorrect_test

import (
	"testing"

	"github.com/studygolang/autocorrect"
)

func TestAutoSpace(t *testing.T) {
	str := autocorrect.AutoSpace("Go语言中文网，Welcome you，gopher们")
	if str != "Go 语言中文网，Welcome you，gopher 们" {
		t.Error("error:", str)
	}
}

func TestAutoCorrect(t *testing.T) {
	str := autocorrect.AutoCorrect(" go语言中文网，Welcome you， gopher们，```go func")
	if str != " Go语言中文网，Welcome you， Gopher们，```go func" {
		t.Error("error:", str)
	}
}
