package htmlUtil

import (
  "testing"
  "container/list"
  "golang.org/x/net/html"
  "strings"
)



func TestParse(t *testing.T) {
    s := `<p>Links:</p><ul><li><a href="foo1">Foo</a><li><a href="/bar/baz">BarBaz</a><li><a href="www.domain.com/foo2">BarBaz</a></ul>`
    doc, _ := html.Parse(strings.NewReader(s))
    var foundUrls list.List
    Parse(doc, "domain.com", &foundUrls)
    if foundUrls.Len() != 2 {
        t.Error("Expected 2, got ", foundUrls.Len())
    }
    if !strings.HasPrefix(foundUrls.Front().Value.(string), "http://domain.com") {
        t.Error("We are return the front format of http request:  ", foundUrls.Front().Value.(string))
    }
}

func TestStrip(t *testing.T) {
    str := "https://www.domain.com/"
    str = stripUrl(str)
    if str != "domain.com" {
        t.Error("Expected domain.com, got ", str)
    }
}

func TestValidate(t *testing.T) {
    domain := "domain.com"
    s1 := "domain2.com/something"
    s2 := "domain.com/somethingelse"
    s3 := "/something/again"
    s1 = validateUrl(s1, domain)
    if s1 != "" {
        t.Error("Expected nothing, got ", s1)
    }
    s2 = validateUrl(s2, domain)
    if s2 != "http://domain.com/somethingelse" {
        t.Error("http://domain.com/somethingelse ", s2)
    }
    s3 = validateUrl(s3, domain)
    if s3 != "http://domain.com/something/again" {
        t.Error("http://domain.com/something/again ", s3)
    }
}
