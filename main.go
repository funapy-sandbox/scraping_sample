package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

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

	pageTitle := doc.Find("title").Text()

	m := make(map[string]map[string]float64)
	doc.Find("table").Each(func(ti int, ts *goquery.Selection) {
		if ti > len(cityFinds) {
			log.Fatalf("ti is bigger than len(finds). ti is %v, len(finds) is %v\n", ti, len(cityFinds))
		}
		lat := dmsToDeg(ts.Find("TR").Eq(1).Find("TD").Eq(1).Text()[1:])
		lon := dmsToDeg(ts.Find("TR").Eq(2).Find("TD").Eq(1).Text()[1:])
		city := normalize(cityFinds[ti][1])

		// city: "" -> city: "熊本県"
		if city == "" {
			city = pageTitle
		}

		m[city] = map[string]float64{
			"lat": lat,
			"lon": lon,
		}
	})

	dumpJSON("./output/data.json", m)
}

// take out after "郡". e.g) 葦北郡芦北町 => 芦北町
// If there is "市" and "区", "区" is given priority. e.g) 熊本市中央区 => 中央区
func normalize(city string) string {
	idx := strings.Index(city, "郡")
	if idx != -1 {
		city = city[idx+3:]
	}

	idxShi := strings.Index(city, "市")
	idxKu := strings.Index(city, "区")

	if idxShi != -1 && idxKu != -1 {
		city = city[idxShi+3:]
	}
	return city
}

// 35°41′28.5576″ => 35.691266
func dmsToDeg(dms string) float64 {
	var degree float64
	var min float64
	var sec float64

	degSplit := strings.Split(dms, "°")
	if len(degSplit) != 2 {
		log.Fatalf("invalid dms(degeee) format: %v\n", dms)
	}
	degree, _ = strconv.ParseFloat(degSplit[0], 64)

	minSplit := strings.Split(degSplit[1], "′")
	if len(minSplit) != 2 {
		log.Fatalf("invalid dms(minute) format: %v\n", dms)
	}
	min, _ = strconv.ParseFloat(minSplit[0], 64)

	l := len(minSplit[1]) - 3
	if minSplit[1][l:] != "″" {
		log.Fatalf("invalid dms(second) format: %v\n", dms)
	}
	// 30″ -> 30
	sec, _ = strconv.ParseFloat(minSplit[1][:l], 64)

	return degree + (min / 60) + (sec / 60 / 60)
}

func dumpJSON(path string, m map[string]map[string]float64) {
	d, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(path, d, os.ModePerm)
}
