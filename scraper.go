package main

import (
	"fmt"
	"log"
	"github.com/PuerkitoBio/goquery"
	"os"
	"io/ioutil"
	"net/http"
	"bytes"
	"encoding/json"
)

func scrape() {

	myFile := "yourFile.txt" //Change
	res := crawlAndParse()

	if _, err := os.Stat(myFile); os.IsNotExist(err) {
		f, err := os.Create(myFile)
		checkErr(err)
		defer f.Close()
		fmt.Println("First Crawling")
		writeOnFile(res, f)
	} else {
		lines := readLines(myFile)
		if( lines == res){
			fmt.Println("Your Output if nothing changes") //Change
		} else {
			fmt.Println("Your Output if something Changes") //Change
			var err = os.Remove(myFile)
			checkErr(err)
			f, err := os.Create(myFile)
			checkErr(err)
			writeOnFile(res, f)
			notifyBot()
		}
	}

}

func crawlAndParse() string {

	doc, err := goquery.NewDocument("Your Url To Crawl")
	if err != nil {
		log.Fatal(err)
	}

	//Example to find all tr and then all td inside them and chaining them into a string
	var currS string
	doc.Find("tr").Each(func(i int, line *goquery.Selection) {
		if ((i > 1) && (i <= 30)) {
			line.Find("td").Each(func(j int, s *goquery.Selection) {
				currChar := s.Text()
				if (len(currChar) > 0) {
					if (j > 0) {
						currS += ", " + currChar
					} else {
						currS += currChar
					}
				}
			})
		}
	})
	return currS
}

func writeOnFile(whatToWrite string, f *os.File)  {
	f.WriteString(whatToWrite)
}

func readLines(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Print(err)
	}
	return string(b)
}

func checkErr(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)
	}
}

func notifyBot()  {
	request_url := "https://api.telegram.org/bot{TOKEN}/sendMessage?chat_id=@{channelName}" //Change

	client := &http.Client{}
	values := map[string]string{"text": "Your Text"} //Change
	jsonStr, _ := json.Marshal(values)
	req, _ := http.NewRequest("POST", request_url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	
	res, err := client.Do(req)
	if(err != nil){
		fmt.Println(err)
	} else {
		fmt.Println(res.Status)
	}
}

func main() {
	scrape()
}