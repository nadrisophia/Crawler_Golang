package htmlUtil

import (
  "fmt"
  "strings"
  "net/http"
  "container/list"
  "golang.org/x/net/html"
  "io/ioutil"
)

// Parse an html node and finds all the links
// Could be concurrent i guess but CBA for now
func Parse( node *html.Node, domain string, foundUrls *list.List){
  if node.Type == html.ElementNode && node.Data == "a" {
    for _, att := range node.Attr {
      if att.Key == "href" {
        var newUrl = validateUrl(att.Val, domain)
        if newUrl != "" {
          foundUrls.PushBack(newUrl)
        }
        break
      }
    }
  }
  for n := node.FirstChild; n != nil; n = n.NextSibling {
    Parse(n, domain, foundUrls)
  }
}

//It basically filters external links and return a usable url
func validateUrl(s, domain string) string {
  s = stripUrl(s)
  if strings.HasPrefix(s, domain) {
    s = strings.Join([]string{"http://", s},  "")
  }else if strings.HasPrefix(s, "/") {
    s = strings.Join([]string{"http://", domain, s}, "")
  } else {
    s = ""
  }
  return s
}

// To avoid some duplications
func stripUrl(s string) string {
  if strings.Contains(s, "cdn-cgi/l/email-protection"){
    return ""
  }
  res := strings.TrimSuffix(s, "/")
  res = strings.TrimPrefix(res, "http://")
  res = strings.TrimPrefix(res, "https://")
  res = strings.TrimPrefix(res, "www.")
  return res
}

//http get to return the html reponse to parse
func GetHtmlNode(url string) *html.Node {
  response, err := http.Get(url)
  if err != nil {
    fmt.Println(err)
    return nil
  }
  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    fmt.Println(err)
    return nil
  }
  doc, err := html.Parse(strings.NewReader(string(body)))
  if err != nil {
    fmt.Println(err)
    return nil
  }
  return doc
}
