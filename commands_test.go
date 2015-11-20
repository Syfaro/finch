package finch_test

import (
	"github.com/syfaro/finch"
	"testing"
)

func TestSimpleCommand(t *testing.T) {
	goodMessage := "/start"
	goodMessage2 := "/start param1 param2"
	goodMessage3 := "/start@FinchExampleBot"
	goodMessage4 := "/start@FinchExampleBot param1 param2"
	badMessage := "/stop"

	if !finch.SimpleCommand("start", goodMessage) {
		t.Error("goodMessage did not return true")
	}

	if !finch.SimpleCommand("start", goodMessage2) {
		t.Error("goodMessage2 did not return true")
	}

	if !finch.SimpleCommand("start", goodMessage3) {
		t.Error("goodMessage3 did not return true")
	}

	if !finch.SimpleCommand("start", goodMessage4) {
		t.Error("goodMessage4 did not return true")
	}

	if finch.SimpleCommand("start", badMessage) {
		t.Error("badMessage did not return false")
	}
}

func TestSimpleArgCommand(t *testing.T) {
	goodMessage := "/start arg1 arg2"
	goodMessage2 := "/start@FinchExampleBot arg1 arg2"
	badMessage := "/start arg1"
	badMessage2 := "/start arg1 arg2 arg3"

	if !finch.SimpleArgCommand("start", 2, goodMessage) {
		t.Error("goodMessage did not return true")
	}

	if !finch.SimpleArgCommand("start", 2, goodMessage2) {
		t.Error("goodMessage2 did not return true")
	}

	if finch.SimpleArgCommand("start", 2, badMessage) {
		t.Error("badMessage did not return false")
	}

	if finch.SimpleArgCommand("start", 2, badMessage2) {
		t.Error("badMessage2 did not return false")
	}
}
