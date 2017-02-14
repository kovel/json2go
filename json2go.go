package main

import (
	"bufio"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"regexp/syntax"
)

func main() {
	jsonUrl := flag.String("url", "", "URL to JSON file/response")
	filePath := flag.String("file", "", "Path to JSON file")
	help := flag.Bool("help", false, "Print help")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	client := http.DefaultClient

	var jsonString string

	if len(*filePath) > 0 {
		data, err := ioutil.ReadFile(*filePath)
		if err != nil {
			log.Fatal(err)
		}

		jsonString = string(data)
	} else if len(*jsonUrl) > 0 {
		resp, err := client.Get(*jsonUrl)
		if err != nil {
			log.Fatal(err)
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		jsonString = string(data)
	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter you JSON:")
		data, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		jsonString = data
	}

	formData := &url.Values{}
	formData.Set("json", jsonString)
	formData.Set("submit", "generate")
	resp, err := client.PostForm("http://json2struct.mervine.net/", *formData)
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	re, err := syntax.Parse("<textarea class=\"form-control\" name=\"struct\" readonly=\"true\">(.*)</textarea>", syntax.DotNL|syntax.OneLine)
	matches := regexp.MustCompile(re.String()).FindAllStringSubmatch(string(data), -1)
	log.Println(html.UnescapeString(matches[0][1]))

}
