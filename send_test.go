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
	testStrong = &Strong{
		Dest:    strings.Split(testMobile, ","),
		Content: testContent,
	}
}

func Test_001_SendStrongUTF8(t *testing.T) {
	client := NewClient(testApiPrefix, testName, testPass)
	b, err := client.SendStrongUTF8(testStrong)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(b))
	}
}
