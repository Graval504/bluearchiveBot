package bluearchiveBot

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func GetCharacterList() {
	url := "https://namu.wiki/w/%ED%8B%80:%EB%B8%94%EB%A3%A8%20%EC%95%84%EC%B9%B4%EC%9D%B4%EB%B8%8C/%ED%95%99%EC%83%9D%EB%AA%85%EB%B6%80"
	res, _ := http.Get(url)
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	fmt.Println(doc.Find("h3").Text())
}
