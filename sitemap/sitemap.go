package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []loc    `xml:"url"`
}

type Link struct {
	Href string
	Text string
}

func main() {
	url := flag.String("l", "https://books.toscrape.com/", "URL of the main page")
	itemLimits := flag.Int("n", 16, "max urls")
	flag.Parse()

	urls := crawl(*url, *itemLimits)
	outputXML(urls)
}

func crawl(initialURL string, itemLimits int) []string {
	visited := make(map[string]struct{})
	var queue []string
	queue = append(queue, initialURL)
	for i := 0; i < itemLimits; i++ {
		if len(queue) == 0 {
			break
		}
		url := queue[0]
		queue = queue[1:]

		r, err := http.Get(url)
		if err != nil {
			fmt.Printf("请求失败：%v\n", err)
			continue
		}

		if r.StatusCode != http.StatusOK {
			fmt.Printf("错误状态码：%v\n", r.Status)
			continue
		}

		links := htmlParse(r.Body)

		for _, link := range links {
			linkResolved := resolve(link.Href, initialURL)

			if !isSameDomain(linkResolved, initialURL) {
				continue
			}

			if _, exists := visited[linkResolved]; exists {
				continue
			}

			visited[url] = struct{}{} // Add url to the visited
			queue = append(queue, linkResolved)
		}
		r.Body.Close()
	}

	// map to string
	var ret []string
	for key := range visited {
		ret = append(ret, key)
	}
	return ret
}

func resolve(href, baseStr string) string {
	base, _ := url.Parse(baseStr)
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	resolved := base.ResolveReference(uri)
	resolved.Fragment = ""
	return resolved.String()
}

func isSameDomain(link, base string) bool {
	l, _ := url.Parse(link)
	b, _ := url.Parse(base)
	return l.Host == b.Host
}

// 给定 html 内容，输出解析好的链接
func htmlParse(r io.Reader) []Link {
	z := html.NewTokenizer(r)
	var links []Link

	var currentHref string
	depth := 0
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return links
		case html.TextToken:
			if depth > 0 {
				link := Link{
					Href: currentHref,
					Text: string(z.Text()),
				}
				links = append(links, link)
			}
		case html.StartTagToken, html.EndTagToken:
			tn, hasAttr := z.TagName()
			if len(tn) == 1 && tn[0] == 'a' {
				if tt == html.StartTagToken {
					depth++

					if hasAttr {
						for {
							key, val, more := z.TagAttr()
							if string(key) == "href" {
								currentHref = string(val)
								break
							}
							if !more {
								break
							}
						}
					}
				} else {
					depth--
				}
			}
		}
	}
}

func outputXML(pages []string) {
	toXml := urlset{
		Xmlns: xmlns,
	}

	for _, page := range pages {
		toXml.Urls = append(toXml.Urls, loc{Value: page})
	}

	// 打印 XML 头部声明
	fmt.Print(xml.Header)

	// 序列化结构体
	// MarshalIndent 会生成带有缩进的易读格式
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXml); err != nil {
		fmt.Printf("XML 编码失败: %v\n", err)
	}
}
