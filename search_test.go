package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func check200Status(t *testing.T, url string) {
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Error("Status code error")
	}
}

func checkSearchresponse(t *testing.T, url string, testFile string) {
	testLog := "Just a test!\r\n"
	fout, err := os.Create(testFile)
	if err != nil {
		t.Error("Test file create error", err)
	}
	defer fout.Close()
	for i := 0; i < 10; i++ {
		fout.WriteString(testLog)
	}

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Error("Status code error")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	if strings.Contains(string(body), testLog) == false {
		t.Error("Search log error")
	}

	if strings.Contains(string(body), "10 results") == false {
		t.Error("Search results error")
	}

	err = os.Remove(testFile)
	if err != nil {
		t.Error("Test file remove error", err)
	}
}

func TestServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(viewHandler))
	defer server.Close()

	check200Status(t, server.URL)
	check200Status(t, server.URL+"/bbsearch")
	checkSearchresponse(t, server.URL+"/bbsearch?q=test", "bb_00000000.txt")
	checkSearchresponse(t, server.URL+"/bbsearch?q=test", "general_00000000.txt")
}
