package main

import (
  "fmt"
  "Crawler_Golang/crawler"
  "sync"
  "time"
)

func main() {
  url := "http://monzo.com"
  domain := "monzo.com"
  c := crawler.NewCrawler(url, domain)

  fmt.Println("Web crawler starting with url:", url)

  g := &sync.WaitGroup{}
  start := time.Now()
  g.Add(1)
  c.Crawl(g)
  g.Wait()
  t := time.Now()
  elapsed := t.Sub(start)
  fmt.Println("in", elapsed, "seconds, we found", len(c.GetProcessedUrls()))

  for k, _ := range c.GetProcessedUrls() {
    fmt.Println(k)
  }
}
