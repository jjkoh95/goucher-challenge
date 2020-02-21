package goucherutils_test

import (
	"testing"

	"github.com/jjkoh95/goucher-challenge/goucherutils"
)

func TestIsEmail(t *testing.T) {
	var isEmailFlag bool
	isEmailFlag = goucherutils.IsEmail("jjkoh95@gmail.com")
	if !isEmailFlag {
		t.Error("Expected to return true given valid email")
	}

	isEmailFlag = goucherutils.IsEmail("")
	if isEmailFlag {
		t.Error("Expected to return false given empty string")
	}

	isEmailFlag = goucherutils.IsEmail("asiasd992183!@$@gmail.com")
	if isEmailFlag {
		t.Error("Expected to return false given invalid email")
	}
}
