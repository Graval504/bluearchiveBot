package bluearchiveBot

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const Table = ".wiki-table > tbody:nth-child(1) > tr:nth-child(3) > td:nth-child(1) > div:nth-child(1) > div:nth-child(1) > dl:nth-child(1) > dd:nth-child(2) > div:nth-child(1) > div:nth-child(1) > table:nth-child(1) > tbody:nth-child(1) > "

func GetCharacterList() []Student {
	data := []Student{}
	runtime.GOMAXPROCS(runtime.NumCPU())
	htmlpy := exec.Command("python", "-c",
		"import requests; "+
			"from bs4 import BeautifulSoup as bs; "+
			"request_headers={'User-Agent':('PythonScraper'),}; "+
			"response = requests.get('https://namu.wiki/w/%ED%8B%80:%EB%B8%94%EB%A3%A8%20%EC%95%84%EC%B9%B4%EC%9D%B4%EB%B8%8C/%ED%95%99%EC%83%9D%EB%AA%85%EB%B6%80',headers=request_headers); "+
			"html = bs(response.text,'html.parser'); "+
			"print(str(html), end = '')")
	out, err := htmlpy.CombinedOutput()
	checkErr(err)
	reader := bytes.NewReader(out)
	html, err := goquery.NewDocumentFromReader(reader)
	checkErr(err)
	channel := make(chan []Student, 15)
	for schoolNum := 1; schoolNum <= 10; schoolNum++ {
		go getListFromHtml(schoolNum, html, channel)
	}
	for i := 1; i <= 10; i++ {
		data = append(data, <-channel...)
	}
	return data
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func getStudentSelector(schoolNum int, studentNum int) string {
	selector := Table +
		"tr:nth-child(" + fmt.Sprint(2*schoolNum) + ") > td > div > div > div:nth-child(" + fmt.Sprint(studentNum) + ") > span > div > a"
	return selector
}

func getStudentSubSelector(schoolNum int, studentNum int) string {
	selector := Table +
		"tr:nth-child(" + fmt.Sprint(2*schoolNum) + ") > td > div > div > div:nth-child(" + fmt.Sprint(studentNum) + ") > em > span > div > a"
	return selector
}

func getSchoolSelector(schoolNum int) string {
	selector := Table +
		"tr:nth-child(" + fmt.Sprint(2*schoolNum-1) + ") > td > div > strong > a > span"
	return selector
}

func getListFromHtml(schoolNum int, html *goquery.Document, c chan []Student) {
	var studentData []Student
	studentNum := 1
	try := 0
	schoolName := html.Find(getSchoolSelector(schoolNum)).Text()
	for {
		data := html.Find(getStudentSelector(schoolNum, studentNum))
		if data.Text() == "" {
			data = html.Find(getStudentSubSelector(schoolNum, studentNum))
		}
		if try >= 3 {
			break
		}
		if data.Text() == "" {
			try += 1
			continue
		}
		var name []string
		if strings.Contains(data.AttrOr("title", "THEREISNOLONGNAME"), "/") {
			name = []string{strings.ReplaceAll(strings.Split(data.AttrOr("title", "THEREISNOLONGNAME"), " ")[1], "/", "(") + ")", strings.ReplaceAll(data.AttrOr("title", "THEREISNOLONGNAME"), "/", "(") + ")"}
			nickname := getNickname(data.AttrOr("title", "THEREISNOLONGNAME"))
			name = append(name, nickname)
			if strings.Contains(nickname, "정") {
				name = append(name, getNewYearNickname(nickname))
			}
		} else if strings.Contains(data.AttrOr("title", "THEREISNOLONGNAME"), "(") {
			name = []string{data.Text(), strings.Split(data.AttrOr("title", "THEREISNOLONGNAME"), "(")[0]}
		} else {
			name = []string{data.Text(), data.AttrOr("title", "THEREISNOLONGNAME")}
		}
		link := data.AttrOr("href", "THEREISNOLINK")
		student := Student{}
		student.Name = name
		student.Link = "https://namu.wiki" + link
		student.School = schoolName
		studentData = append(studentData, student)
		studentNum += 1
	}
	c <- studentData
}

func getNickname(namedata string) string {
	tempname := strings.Split(namedata, "/")   // [Full name,Eventname]
	name := strings.Split(tempname[0], " ")[1] // Last name
	nick := tempname[1]                        // Event name
	if len(name) >= 9 {
		return nick[:3] + name[3:]
	} else if len(name) == 6 {
		return nick[:3] + name
	} else {
		return nick + name
	}
}

func getNewYearNickname(nickname string) string {
	return strings.ReplaceAll(nickname, "정", "뉴")
}
