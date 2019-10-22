package main

import (
	"testing"
	"net/http"
	"time"
	"fmt"
	"io/ioutil"
)

func TestWebservice(t *testing.T) {
	start := time.Now()
	ch := make(chan string)
	i := 0
	for i < 50 {
		go MakeRequest(ch, t)
		i++
	}

	i = 0
	for i < 50 {
		fmt.Println(<-ch)
		i++
	}
	fmt.Printf("%.2f elapsed\n", time.Since(start).Seconds())
}

func MakeRequest(ch chan<-string, t *testing.T) {
	start := time.Now()
	resp, err := http.Get("http://localhost:5000")
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}
	defer resp.Body.Close()
	secs := time.Since(start).Seconds()
	body, _ := ioutil.ReadAll(resp.Body)
	ch <- fmt.Sprintf("%.2f elapsed with response length: %d", secs, len(body))
	t.Log(string(body))
}
