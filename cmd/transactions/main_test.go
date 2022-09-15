package main

import (
	"os"
	"testing"
)

func TestinitFlags(t *testing.T) {
	if err := initFlags(); err != nil {
		t.Fail()
	}
	if os.Getenv("NODE") == "" {
		t.Fail()
	}
	if from == "" {
		t.Log("from flag empty")
		t.Fail()
	}
}

func BenchmarkLenTimeInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lenTimeInt()
	}
}

func BenchmarkLenTimeString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lenTimeString()
	}
}

func lenTimeInt() bool {
	if len(os.Getenv("NODE")) == 0 {
		return true
	}
	return false
}

func lenTimeString() bool {
	if os.Getenv("NODE") == "" {
		return true
	}
	return false
}
