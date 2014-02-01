package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Page struct {
	Title  string
	Query  string
	Result string
	Logs   []string
}

func readLogFile(filename string) chan string {
	ch := make(chan string)
	go func(ch chan string) {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Println("error opening file ", err)
			return
		}
		defer f.Close()
		r := bufio.NewReader(f)
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				ch <- ""
				return
			} else if err == io.EOF {
				ch <- ""
				return
			}
			ch <- fmt.Sprintf(line)
		}
	}(ch)
	return ch
}

func collectLogs(filepattern string, searchword string) []string {
	logfiles, _ := filepath.Glob(filepattern)
	var logs = []string{}
	for _, filename := range logfiles {
		for line := range readLogFile(filename) {
			if line == "" {
				break
			}
			if strings.Contains(line, searchword) {
				log := fmt.Sprintf("%s %s", path.Base(filename)[4:12], line)
				logs = append(logs, log)
			}
		}
	}
	return logs
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	f, _ := os.OpenFile("logfile.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	l := log.New(f, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	l.Printf("%s %s %s %s\n", r.RequestURI, r.RemoteAddr, r.UserAgent(), r.Referer())
	defer f.Close()

	searchword := r.FormValue("q")

	var logs = []string{}
	var result string = ""
	if len(searchword) == 0 {
		result = "You must input some words"
	} else {
		logs = collectLogs("*bb*.txt", searchword)
		if len(logs) == 0 {
			result = fmt.Sprintf("No result : %q", searchword)
		} else {
			result = fmt.Sprintf("%d results : %q\n", len(logs), searchword)
		}
	}

	var reverse_logs = []string{}
	for i := len(logs) - 1; i >= 0; i-- {
		reverse_logs = append(reverse_logs, logs[i])
	}

	var mainTemplate = template.Must(template.ParseFiles("templates/search.html"))

	err := mainTemplate.Execute(w, Page{
		Title:  "bb Log Search",
		Query:  searchword,
		Result: result,
		Logs:   reverse_logs,
	})

	if err != nil {
		fmt.Errorf("mainTemplate: %v", err)
	}
}

func main() {
	fmt.Println("Running on http://0.0.0.0:9000/")
	http.HandleFunc("/", viewHandler)
	http.ListenAndServe(":9000", nil)
}
