package cxtonsms

import (
	"flag"
	"strings"
	"testing"
)

var (
	testApiPrefix string
	testName      string
	testPass      string
	testMobile    string
	testContent   string
	testStrong    *Strong
)

func init() {
	flag.StringVar(&testApiPrefix, "host", "", "Api Host, etc http://xxxxxxxxx")
	flag.StringVar(&testName, "name", "", "name")
	flag.StringVar(&testPass, "password", "", "password")
	flag.StringVar(&testMobile, "mobile", "", "mobile number")
	flag.StringVar(&testContent, "content", "", "sms content")
	flag.Parse()
	if testMobile != "" && testContent != "" {
		testStrong = &Strong{
			Dest:    strings.Split(testMobile, ","),
			Content: testContent,
		}
	} else {
		testStrong = nil
	}
}

func Test_001_SendStrongUTF8(t *testing.T) {
	if testStrong == nil {
		return
	}
	client := NewClient(testApiPrefix, testName, testPass)
	b, err := client.SendStrongUTF8(testStrong)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(b))
	}
}

func Test_002_GetBalanceUTF8(t *testing.T) {
	client := NewClient(testApiPrefix, testName, testPass)
	b, err := client.GetBalanceUTF8(nil)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(b))
	}
}
