package bluearchiveBot

import (
	"log"
	"os/exec"
)

func GetCharacterList() string {
	htmlpy := exec.Command("python", "-c",
		"import requests; "+
			"from bs4 import BeautifulSoup as bs; "+
			"request_headers={'User-Agent':('PythonScraper'),}; "+
			"response = requests.get('https://namu.wiki/w/%ED%8B%80:%EB%B8%94%EB%A3%A8%20%EC%95%84%EC%B9%B4%EC%9D%B4%EB%B8%8C/%ED%95%99%EC%83%9D%EB%AA%85%EB%B6%80',headers=request_headers); "+
			"html = bs(response.text,'html.parser'); "+
			"print(str(html), end = '')")
	out, err := htmlpy.CombinedOutput()
	checkErr(err)
	return string(out)
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
