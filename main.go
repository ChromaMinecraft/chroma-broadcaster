package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/robfig/cron/v3"
)

type jobs struct {
	Name  string     `json:"name"`
	Query string     `json:"query"`
	URL   string     `json:"url"`
	Data  dataObject `json:"data"`
}

type dataObject struct {
	Username string        `json:"username"`
	Content  string        `json:"content"`
	Embeds   []embedObject `json:"embeds"`
}

type embedObject struct {
	Color       uint              `json:"color"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Footer      embedFooterObject `json:"footer"`
}

type embedFooterObject struct {
	Text string `json:"text"`
}

func main() {
	var port string = ":8199"

	filePath := os.Getenv("FILE_PATH")
	fmt.Println("load cron file at:", filePath)

	if filePath == "" {
		log.Fatalln("FILE_PATH environment variable not set")
		return
	}

	file, err := ioutil.ReadFile("/static/" + filePath)
	if err != nil {
		log.Fatalln("error reading file: ", err)
		return
	}

	var jobList []jobs
	json.Unmarshal([]byte(file), &jobList)

	c := cron.New()
	c.AddFunc("@every 1h0m0s", func() { log.Println("log every 1 hour") })

	for _, each := range jobList {
		dataMarshal, _ := json.Marshal(each.Data)
		log.Printf("add cron job: %s, %s", each.Name, each.Query)
		c.AddFunc(each.Query, func() { webhookRequest(each.URL, dataMarshal) })
	}

	c.Start()

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic: %+v\n", r)
		}
	}()

	fmt.Printf("Server is running at port %s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}

func webhookRequest(url string, data []byte) {
	body := bytes.NewBuffer(data)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Printf("http NewRequest error: %v", err)
		return
	}
	req.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("http client Do error: %v", err)
		return
	}

	defer res.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(res.Body)

	log.Printf("status: %v", res.Status)
	log.Printf("body: %v", string(bodyBytes))
}
