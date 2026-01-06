package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func main() {
	filename := flag.String("f", "ex1.html", "File name")
	flag.Parse()

	r, err := os.Open(*filename)
	if err != nil {
		exit("Failed to open the file!")
	}
	defer r.Close()

	output := htmlParse(r)

	for _, link := range output {
		fmt.Printf("href: %s, text: %s\n", link.Href, link.Text)
	}
}

func htmlParse(r *os.File) []Link {
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
				links = append(links, buildLink(currentHref, string(z.Text())))
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

func buildLink(href string, text string) Link {
	return Link{
		Href: href,
		Text: text,
	}
}

func exit(s string) {
	fmt.Println(s)
	os.Exit(1)
}
