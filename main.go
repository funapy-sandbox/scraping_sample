package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// to get city name
	r := regexp.MustCompile("<(?:br/|BR/)>(.*)")

	doc, err := goquery.NewDocument("http://www.gsi.go.jp/KOKUJYOHO/CENTER/kendata/kumamoto_heso.htm")
	if err != nil {
		log.Fatal(err)
	}

	html, err := doc.Html()
	if err != nil {
		log.Fatal(err)
	}

	cityFinds := r.FindAllStringSubmatch(html, -1)

	doc.Find("table").Each(func(ti int, ts *goquery.Selection) {
		if ti > len(cityFinds) {
			log.Fatalf("ti is smaller than len(finds). ti is %v, len(finds) is %v\n", ti, len(cityFinds))
		}
		lat := ts.Find("TR").Eq(1).Find("TD").Eq(1).Text()
		lon := ts.Find("TR").Eq(2).Find("TD").Eq(1).Text()
		city := cityFinds[ti][1]

		if city == "" {
			city = "熊本県"
		}

		fmt.Println(city)
		fmt.Println(lat)
		fmt.Println(lon)
	})
}
