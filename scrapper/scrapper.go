package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type jobData struct {
	id       string
	location string
	title    string
	salary   string
	summary  string
}

// Scrapping function
func Scrap(word string) {
	var LIMIT int = 50
	var MAIN_URL string = "https://kr.indeed.com/jobs?q=" + word + "&limit=" + strconv.Itoa(LIMIT)
	var jobs []jobData
	// 2-1
	c2 := make(chan []jobData)
	pageNum := getPageNum(MAIN_URL)
	for i := 0; i < pageNum; i++ {
		go getPage(i, MAIN_URL, LIMIT, c2)
	}
	// 2-4
	for i := 0; i < pageNum; i++ {
		jobSlice := <-c2
		jobs = append(jobs, jobSlice...)
	}
	saveJobData(jobs)
	fmt.Println("Save done !")
}

// 2-2
func getPage(i int, url string, limit int, c2 chan<- []jobData) {
	var jobs []jobData
	// 1-1
	c := make(chan jobData)
	PAGE_URL := url + "&start=" + strconv.Itoa(i*limit)
	fmt.Println("Requesting: ", PAGE_URL)
	res, err := http.Get(PAGE_URL)
	if err != nil {
		log.Fatalln("err")
	}
	if res.StatusCode != 200 {
		log.Fatalln("err: ", res.StatusCode)
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln("err")
	}
	jobBox := doc.Find(".tapItem")
	jobBox.Each(func(i int, s *goquery.Selection) {
		// 1-2
		go getJobData(s, c)
	})
	// 1-5
	for i := 0; i < jobBox.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}
	// 2-3
	c2 <- jobs
}

func getPageNum(url string) int {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalln("err")
	}
	if res.StatusCode != 200 {
		log.Fatalln("err: ", res.StatusCode)
	}
	defer res.Body.Close()
	pageNum := 0
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln("err")
	}
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pageNum = s.Find("a").Length()
	})

	return pageNum
}

// 1-3
func getJobData(s *goquery.Selection, c chan<- jobData) {
	id, _ := s.Attr("data-jk")
	title := s.Find("h2>span").Text()
	location := s.Find(".companyLocation").Text()
	salary := s.Find(".salary-snippet").Text()
	summary := s.Find(".job-snippet").Text()
	// 1-4
	c <- jobData{
		id:       id,
		location: location,
		title:    title,
		salary:   salary,
		summary:  summary}
}

func saveJobData(jobs []jobData) {
	// file 만들기
	f, err := os.Create("jobs.csv")
	if err != nil {
		log.Fatalln("err")
	}
	// writer 만들기
	w := csv.NewWriter(f)

	header := []string{"ID", "Location", "Title", "Salary", "Summary"}
	err2 := w.Write(header)
	if err2 != nil {
		log.Fatalln("err")
	}
	for _, job := range jobs {
		jobDetail := []string{job.id, job.location, job.title, job.salary, job.summary}
		err3 := w.Write(jobDetail)
		if err3 != nil {
			log.Fatalln("err")
		}

	}

	// file 저장
	defer w.Flush()
}
