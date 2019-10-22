package main

import (
    "fmt"
    "net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"sync"
)

const PORT = "5000"

type Name struct {
	Name string `json:"name"`
	Surname string `json:"surname"`
	Gender string `json:"gender"`
    Region string  `json:"region"`
}

type Joke struct {
	Type string `json:"type"`
	Value JokeValue `json:"value"`
}

type JokeValue struct {
	Id int `json:"id"`
	Joke string `json:"joke"`
	Categories []string `json:"categories"`
}

type CombinedResp struct {
	Web1 interface{}
	Web2 interface{}
	mux sync.Mutex
}

func main() {
	log.Println("Starting HTTP Server.")

	http.HandleFunc("/", handleFunc)
	var err = http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.Fatalf("Server failed starting. Error: %s", err)
	}
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	combinedResp := &CombinedResp{}
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go getName(combinedResp, wg)
	go getJoke(combinedResp, wg)
	wg.Wait()

	data, err := json.MarshalIndent(combinedResp, "", "    ")
	if err != nil {
		log.Printf("JSON marshaling failed: %s", err)
		fmt.Fprintf(w, "Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func getName(combinedResp *CombinedResp, wg *sync.WaitGroup) {
	defer wg.Done()
	nameURL := "https://uinames.com/api/"
	raw, err := do(nameURL)
	if err != nil {
		log.Printf("Error getting name from %s: %+v", nameURL, err)
		setCombinedResp("name", combinedResp, err.Error())
		return
	}
	name := &Name{}
	err = json.Unmarshal(raw, name)
	if err != nil {
		log.Printf("Error unmarshalling name json %+v\n", err)
		setCombinedResp("name", combinedResp, err.Error())
		return
	}
	setCombinedResp("name", combinedResp, name)
}

func getJoke(combinedResp *CombinedResp, wg *sync.WaitGroup) {
	defer wg.Done()
	jokeURL := "http://api.icndb.com/jokes/random?firstName=Costel&lastName=Sassu&limitTo=[nerdy]"
	raw, err := do(jokeURL)
	if err != nil {
		log.Printf("Error getting joke from %s: %+v\n", jokeURL, err)
		setCombinedResp("joke", combinedResp, err.Error())
		return
	}
	joke := &Joke{}
	err = json.Unmarshal(raw, joke)
	if err != nil {
		log.Printf("Error unmarshalling joke json %+v\n", err)
		setCombinedResp("joke", combinedResp, err.Error())
		return 
	}
    setCombinedResp("joke", combinedResp, joke)
}

func setCombinedResp(field string, c *CombinedResp, msg interface{}) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if field == "name" {
		c.Web1 = msg
	} else {
		c.Web2 = msg
	}
}

func do(url string) ([]byte, error) {
	var raw []byte
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating http request: %+v\n", err) 
		return raw, err
	}
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil && resp == nil {
		log.Printf("Error getting from %s: %+v\n", url, err)
		return raw, err
	}
	log.Printf("response status: %d\n", resp.StatusCode)
	defer resp.Body.Close()
	raw, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Could not parse response body. %+v\n", err)
		return raw, err
	}
	return raw, nil
}
