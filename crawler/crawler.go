package crawler

import (
  "fmt"
  "sync"
  "container/list"
  "Crawler_Golang/htmlUtil"
)


type crawler struct {
  url string
  domain string
  processedUrls *Map
}

//This structure is necessary to protect our map from
//concurrent reads and writes
type Map struct {
	sync.Mutex
	m map[string]struct{}
}

//Class constructor
func NewCrawler(url, domain string) *crawler {
  urls := make(map[string]struct{})
  var empty struct{}
  urls[url] = empty
  return &crawler{url, domain, &Map{m :urls}}
}

//Class constructor with existent urls map
func NewCrawlerUrls(url, domain string, processedUrls *Map) *crawler {
  return &crawler{url, domain, processedUrls}
}

// Access the processed urls map
func (c *crawler) GetProcessedUrls() map[string]struct{} {
  return c.processedUrls.m
}

// Add the current url to the processed urls
func (c *crawler) addProcessedUrl() {
  var empty struct{}
  c.processedUrls.Lock()
  c.processedUrls.m[c.url] = empty
  c.processedUrls.Unlock()
}

// check if the url is already processed
func (c *crawler) isUrlProcessed(url string) bool {
  defer c.processedUrls.Unlock()
  c.processedUrls.Lock()
  _, exists := c.processedUrls.m[url]
  if exists{
    return true
  }
  return false

}

// crawl starting from the current url
func (c *crawler) Crawl(g *sync.WaitGroup){
  defer g.Done()
  //http get url
  doc := htmlUtil.GetHtmlNode(c.url)
  if doc == nil {
    fmt.Println("this url didnt work", c.url)
    return
  }
  //add current url to processedUrls
  c.addProcessedUrl()
  //parse html response
  var foundUrls list.List
  htmlUtil.Parse(doc, c.domain, &foundUrls)
  //process the new found urls
  for e := foundUrls.Front(); e != nil; e = e.Next() {
    if c.isUrlProcessed(e.Value.(string)) {
      continue
    }
    g.Add(1)
    cr := NewCrawlerUrls(e.Value.(string), c.domain, c.processedUrls)
    fmt.Println("crawling", e.Value.(string))
    go cr.Crawl(g)
  }
}
