package bluearchiveBot

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func GetCharacterInfo() {
	url := "https://google.com/search?q=site:https://codealone.tistory.com/65"
	res, _ := http.Get(url)
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	fmt.Println(doc.Find("h3").Text())
}
