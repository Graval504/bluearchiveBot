package bluearchiveBot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//const div string = "body > div:nth-child(1) > div > div:nth-child(2) > article > div:nth-last-child(2)"

func GetCharacterInfoFromData(data []Student) []Student {
	info := []Student{}
	runtime.GOMAXPROCS(runtime.NumCPU())
	channel := make(chan Student, len(data))
	for _, stud := range data {
		go getOneCharacterInfo(stud, channel)
	}
	for i := 0; i < len(data); i++ {
		info = append(info, <-channel)
	}
	return info
}

func CreateJsonFileFromData(data []Student) {
	jsonfile, err := json.Marshal(data)
	checkErr(err)
	ioutil.WriteFile("./jsonfile.json", jsonfile, os.FileMode(0644))
}

type Student struct {
	Name            []string
	Link            string
	School          string
	Club            string
	InitialFigures  int
	InitialStarNum  int
	AttackValue     int
	StudentType     [5]string
	AreaType        [3]string
	Weapon          string
	Equipment       [3]string
	Ooparts         [2]string
	ExUpMaterial    [4][2]int
	SkillUpMaterial [9][2]int
	Skills          struct {
		Ex      exType
		Normal  skillType
		Passive skillType
		Sub     skillType
	}
}

type exType struct {
	SkillName string
	Skill     [5]struct {
		Level       int
		Description string
		Cost        int
	}
}
type skillType struct {
	SkillName string
	Skill     [10]struct {
		Level       int
		Description string
	}
}

func requestWithPython(url string) *goquery.Selection {
	python := exec.Command("python", "-c",
		"import sys; "+
			"import requests; "+
			"from bs4 import BeautifulSoup as bs; "+
			"sys.stdout.reconfigure(encoding='utf-8'); "+
			"sys.stdin.reconfigure(encoding='utf-8'); "+
			"request_headers={'User-Agent':('PythonScraper'),}; "+
			"response = requests.get('"+url+"',headers=request_headers); "+
			"html = bs(response.text,'html.parser'); "+
			"print(str(html), end = '')")
	out, err := python.CombinedOutput()
	checkErr(err)
	reader := bytes.NewReader(out)
	html, err := goquery.NewDocumentFromReader(reader)
	checkErr(err)
	selection := html.Find(".wiki-table-wrap.table-center")
	div := selection.Parent()
	return div
}

func getOneCharacterInfo(stud Student, c chan Student) {
	html := requestWithPython(stud.Link)
	stud.Club = getClubFromHtml(html)
	stud.InitialStarNum = getStarNum(html)
	stud.InitialFigures = getInitialFigures(html)
	stud.StudentType = getStudentType(html)
	stud.Weapon = getWeapon(html)
	stud.AreaType = getAreaType(html)
	stud.Equipment = getEquipment(html)
	stud.Ooparts = getOoparts(html)
	stud.ExUpMaterial = getExUpMaterial(html)
	stud.SkillUpMaterial = getSkillUpMaterial(html)
	stud.Skills = getSkills(html)
	stud.AttackValue = getAttackValue(stud.InitialFigures, stud.InitialStarNum)
	c <- stud
}

func getClubFromHtml(html *goquery.Selection) string {
	selected := html.Find("* > div > table > tbody > tr:nth-child(6) > td:nth-child(2) > div > a").First().Text()
	if selected == "?????????" {
		selected = html.Find("* > div > table > tbody > tr:nth-child(5) > td:nth-child(2) > div > a").First().Text()
	}
	if selected == "?????????" {
		selected = "???????????? ??????"
	}
	return strings.TrimSpace(selected)
}

func getStarNum(html *goquery.Selection) int {
	selected := html.Find("div:nth-child(11) > div > table > tbody > tr:nth-child(1) > td > div > div:nth-child(5) > a")
	if selected.Children().Length() == 0 {
		selected = html.Find("div:nth-child(12) > div > table > tbody > tr:nth-child(1) > td > div > div:nth-child(5) > a")
	}
	if selected.Children().Length() == 0 {
		selected = html.Find("div:nth-child(13) > div > table > tbody > tr:nth-child(1) > td > div > div:nth-child(5) > a")
	}
	if selected.Children().Length() == 0 {
		selected = html.Find("div:nth-child(14) > div > table > tbody > tr:nth-child(1) > td > div > div:nth-child(5) > a")
	}
	return selected.Length()
}

func getInitialFigures(html *goquery.Selection) int {
	selected := html.Find("div:nth-child(11) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(3) > div > div > table > tbody > tr:nth-child(2) > td:nth-child(2) > div")
	if selected.Text() == "" {
		selected = html.Find("div:nth-child(12) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(3) > div > div > table > tbody > tr:nth-child(2) > td:nth-child(2) > div")
	}
	if selected.Text() == "" {
		selected = html.Find("div:nth-child(13) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(3) > div > div > table > tbody > tr:nth-child(2) > td:nth-child(2) > div")
	}
	if selected.Text() == "" {
		selected = html.Find("div:nth-child(14) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(3) > div > div > table > tbody > tr:nth-child(2) > td:nth-child(2) > div")
	}
	return stringToInt(strings.TrimSpace(strings.Trim(selected.Text(), "?????????")))
}

func getStudentType(html *goquery.Selection) [5]string {
	j := 1
	selected := html.Find("div:nth-child(11) > div > table > tbody > tr:nth-child(1) > td > div > div:nth-child(1) > div > div > table > tbody > tr:nth-child(2)")
	studtype, exist := selected.Attr("title")
	if !exist {
		selected = html.Find("div:nth-child(12) > div > table > tbody > tr:nth-child(1) > td > div > div:nth-child(1) > div > div > table > tbody > tr:nth-child(2)")
		studtype, exist = selected.Attr("title")
	}
	if !exist {
		selected = html.Find("div:nth-child(13) > div > table > tbody > tr:nth-child(1) > td > div > div:nth-child(1) > div > div > table > tbody > tr:nth-child(2)")
		studtype, _ = selected.Attr("title")
	}
	if !exist {
		selected = html.Find("div:nth-child(14) > div > table > tbody > tr:nth-child(1) > td > div > div:nth-child(1) > div > div > table > tbody > tr:nth-child(2)")
		studtype, _ = selected.Attr("title")
	}
	if studtype == "???????????? ?????????" {
		studtype = "SPECIAL"
	} else {
		studtype = "STRIKER"
	}
	var role, position, attackType, defenceType string
	for i := 0; i < 3; i++ {
		if j == 1 {
			selected := html.Find("div:nth-child(11) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(5) > div > div > table > tbody > tr:nth-child(2) > td > div > div:nth-child(" + fmt.Sprint(i+1) + ") > div > div:nth-child(" + fmt.Sprint(j) + ")")
			if selected.Text() == "" {
				selected = html.Find("div:nth-child(12) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(5) > div > div > table > tbody > tr:nth-child(2) > td > div > div:nth-child(" + fmt.Sprint(i+1) + ") > div > div:nth-child(" + fmt.Sprint(j) + ")")
			}
			if selected.Text() == "" {
				selected = html.Find("div:nth-child(13) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(5) > div > div > table > tbody > tr:nth-child(2) > td > div > div:nth-child(" + fmt.Sprint(i+1) + ") > div > div:nth-child(" + fmt.Sprint(j) + ")")
			}
			if selected.Text() == "" {
				selected = html.Find("div:nth-child(14) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(5) > div > div > table > tbody > tr:nth-child(2) > td > div > div:nth-child(" + fmt.Sprint(i+1) + ") > div > div:nth-child(" + fmt.Sprint(j) + ")")
			}
			role = strings.TrimSpace(selected.Text())
			j += 1
		}
		selected := html.Find("div:nth-child(11) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(5) > div > div > table > tbody > tr:nth-child(2) > td > div > div:nth-child(" + fmt.Sprint(i+1) + ") > div > div:nth-child(" + fmt.Sprint(j) + ")")
		if selected.Text() == "" {
			selected = html.Find("div:nth-child(12) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(5) > div > div > table > tbody > tr:nth-child(2) > td > div > div:nth-child(" + fmt.Sprint(i+1) + ") > div > div:nth-child(" + fmt.Sprint(j) + ")")
		}
		if selected.Text() == "" {
			selected = html.Find("div:nth-child(13) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(5) > div > div > table > tbody > tr:nth-child(2) > td > div > div:nth-child(" + fmt.Sprint(i+1) + ") > div > div:nth-child(" + fmt.Sprint(j) + ")")
		}
		if selected.Text() == "" {
			selected = html.Find("div:nth-child(14) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(5) > div > div > table > tbody > tr:nth-child(2) > td > div > div:nth-child(" + fmt.Sprint(i+1) + ") > div > div:nth-child(" + fmt.Sprint(j) + ")")
		}
		switch i {
		case 0:
			position = selected.Text()
		case 1:
			attackType = selected.Text()
		case 2:
			defenceType = selected.Text()
		}
	}
	return [5]string{studtype, role, position, attackType, defenceType}
}

func getAreaType(html *goquery.Selection) [3]string {
	areaType := [3]string{}
	for i := 0; i < 3; i++ {
		areaType[i] = html.Find("div:nth-child(11) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(5) > div > div > table > tbody > tr:nth-child(2) > td > div > div:nth-child("+fmt.Sprint(i+4)+") > div:nth-child(2) > a").AttrOr("title", "??????X")[6:]
		if areaType[i] == "X" {
			areaType[i] = html.Find("div:nth-child(12) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(5) > div > div > table > tbody > tr:nth-child(2) > td > div > div:nth-child("+fmt.Sprint(i+4)+") > div:nth-child(2) > a").AttrOr("title", "??????X")[6:]
		}
		if areaType[i] == "X" {
			areaType[i] = html.Find("div:nth-child(13) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(5) > div > div > table > tbody > tr:nth-child(2) > td > div > div:nth-child("+fmt.Sprint(i+4)+") > div:nth-child(2) > a").AttrOr("title", "??????X")[6:]
		}
		if areaType[i] == "X" {
			areaType[i] = html.Find("div:nth-child(14) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(5) > div > div > table > tbody > tr:nth-child(2) > td > div > div:nth-child("+fmt.Sprint(i+4)+") > div:nth-child(2) > a").AttrOr("title", "??????X")[6:]
		}
	}
	return areaType
}

func getWeapon(html *goquery.Selection) string {
	selected := html.Find("div:nth-child(11) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(11) > div > div > table > tbody > tr:nth-child(2) > td:nth-child(1)")
	if selected.Text() == "" {
		selected = html.Find("div:nth-child(12) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(11) > div > div > table > tbody > tr:nth-child(2) > td:nth-child(1)")
	}
	if selected.Text() == "" {
		selected = html.Find("div:nth-child(13) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(11) > div > div > table > tbody > tr:nth-child(2) > td:nth-child(1)")
	}
	if selected.Text() == "" {
		selected = html.Find("div:nth-child(14) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(11) > div > div > table > tbody > tr:nth-child(2) > td:nth-child(1)")
	}
	return strings.TrimSpace(selected.Text())
}

func getEquipment(html *goquery.Selection) [3]string {
	equipment := [3]string{}
	for i := 0; i < 3; i++ {
		equipment[i] = html.Find("div:nth-child(11) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(11) > div > div > table > tbody > tr:nth-child(2) > td:nth-child(" + fmt.Sprint(2+i) + ")").Text()
		if equipment[i] == "" {
			equipment[i] = html.Find("div:nth-child(12) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(11) > div > div > table > tbody > tr:nth-child(2) > td:nth-child(" + fmt.Sprint(2+i) + ")").Text()
		}
		if equipment[i] == "" {
			equipment[i] = html.Find("div:nth-child(13) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(11) > div > div > table > tbody > tr:nth-child(2) > td:nth-child(" + fmt.Sprint(2+i) + ")").Text()
		}
		if equipment[i] == "" {
			equipment[i] = html.Find("div:nth-child(14) > div > table > tbody > tr:nth-child(2) > td > div:nth-child(3) > div:nth-child(11) > div > div > table > tbody > tr:nth-child(2) > td:nth-child(" + fmt.Sprint(2+i) + ")").Text()
		}
	}
	return equipment
}

func getOoparts(html *goquery.Selection) [2]string {
	ooparts := [2]string{}
	for i := 0; i < 2; i++ {
		ooparts[i] = html.Find("div:nth-child(13) > div:nth-child(1) > div:nth-child(2) > div > table > tbody > tr:nth-child(3) > td:nth-child(" + fmt.Sprint(3+i) + ")").Text()
		if ooparts[i] == "" {
			ooparts[i] = html.Find("div:nth-child(14) > div:nth-child(1) > div:nth-child(2) > div > table > tbody > tr:nth-child(3) > td:nth-child(" + fmt.Sprint(3+i) + ")").Text()
		}
		if ooparts[i] == "" {
			ooparts[i] = html.Find("div:nth-child(15) > div:nth-child(1) > div:nth-child(2) > div > table > tbody > tr:nth-child(3) > td:nth-child(" + fmt.Sprint(3+i) + ")").Text()
		}
		if i == 1 {
			ooparts[i] = strings.Split(ooparts[i], "(")[0]
		} else {
			ooparts[i] = strings.Split(ooparts[i], "(")[0]
		}
	}
	return ooparts
}

func getExUpMaterial(html *goquery.Selection) [4][2]int {
	exUpMaterial := [4][2]int{}
	j := 0
	for i := 0; i < 4; i++ {
		material0 := strings.Trim(html.Find("div:nth-child(13) > div:nth-child(1) > div:nth-child(3) > div > table > tbody > tr > td > div > div > dl > dd > div > div > table > tbody > tr:nth-child("+fmt.Sprint(2*i+2)+") > td:nth-child("+fmt.Sprint(3+j)+")").Text(), " ??")
		if material0 == "" {
			material0 = strings.Trim(html.Find("div:nth-child(14) > div:nth-child(1) > div:nth-child(3) > div > table > tbody > tr > td > div > div > dl > dd > div > div > table > tbody > tr:nth-child("+fmt.Sprint(2*i+2)+") > td:nth-child("+fmt.Sprint(3+j)+")").Text(), " ??")
		}
		if material0 == "" {
			material0 = strings.Trim(html.Find("div:nth-child(15) > div:nth-child(1) > div:nth-child(3) > div > table > tbody > tr > td > div > div > dl > dd > div > div > table > tbody > tr:nth-child("+fmt.Sprint(2*i+2)+") > td:nth-child("+fmt.Sprint(3+j)+")").Text(), " ??")
		}
		exUpMaterial[i][0] = stringToInt(material0)
		if i == 0 {
			j = 1
			exUpMaterial[i][1] = 0
		} else {
			material1 := strings.Trim(html.Find("div:nth-child(13) > div:nth-child(1) > div:nth-child(3) > div > table > tbody > tr > td > div > div > dl > dd > div > div > table > tbody > tr:nth-child("+fmt.Sprint(2*i+2)+") > td:nth-child("+fmt.Sprint(4+j)+")").Text(), " ??")
			if material1 == "" {
				material1 = strings.Trim(html.Find("div:nth-child(14) > div:nth-child(1) > div:nth-child(3) > div > table > tbody > tr > td > div > div > dl > dd > div > div > table > tbody > tr:nth-child("+fmt.Sprint(2*i+2)+") > td:nth-child("+fmt.Sprint(4+j)+")").Text(), " ??")
			}
			if material1 == "" {
				material1 = strings.Trim(html.Find("div:nth-child(15) > div:nth-child(1) > div:nth-child(3) > div > table > tbody > tr > td > div > div > dl > dd > div > div > table > tbody > tr:nth-child("+fmt.Sprint(2*i+2)+") > td:nth-child("+fmt.Sprint(4+j)+")").Text(), " ??")
			}
			exUpMaterial[i][1] = stringToInt(material1)
		}
	}
	return exUpMaterial
}

func getSkillUpMaterial(html *goquery.Selection) [9][2]int {
	skillUpMaterial := [9][2]int{}
	skillUpMaterial[0] = [2]int{0, 0}
	skillUpMaterial[1] = [2]int{0, 0}
	skillUpMaterial[8] = [2]int{0, 0}
	skillUpMaterial[2][1] = 0
	material0 := strings.Trim(html.Find("div:nth-child(13) > div:nth-child(1) > div:nth-child(4) > div > table > tbody > tr > td > div > div > dl > dd > div > div > table > tbody > tr:nth-child(6) > td:nth-child(4)").Text(), " ??")
	if material0 == "" {
		material0 = strings.Trim(html.Find("div:nth-child(14) > div:nth-child(1) > div:nth-child(4) > div > table > tbody > tr > td > div > div > dl > dd > div > div > table > tbody > tr:nth-child(6) > td:nth-child(4)").Text(), " ??")
	}
	if material0 == "" {
		material0 = strings.Trim(html.Find("div:nth-child(15) > div:nth-child(1) > div:nth-child(4) > div > table > tbody > tr > td > div > div > dl > dd > div > div > table > tbody > tr:nth-child(6) > td:nth-child(4)").Text(), " ??")
	}
	skillUpMaterial[2][0] = stringToInt(material0)
	for i := 0; i < 5; i++ {
		material0 := strings.Trim(html.Find("div:nth-child(13) > div:nth-child(1) > div:nth-child(4) > div > table > tbody > tr > td > div > div > dl > dd > div > div > table > tbody > tr:nth-child("+fmt.Sprint(8+2*i)+") > td:nth-child("+fmt.Sprint(3+i%2)+")").Text(), " ??")
		material1 := strings.Trim(html.Find("div:nth-child(13) > div:nth-child(1) > div:nth-child(4) > div > table > tbody > tr > td > div > div > dl > dd > div > div > table > tbody > tr:nth-child("+fmt.Sprint(8+2*i)+") > td:nth-child("+fmt.Sprint(4+i%2)+")").Text(), " ??")
		if material0 == "" {
			material0 = strings.Trim(html.Find("div:nth-child(14) > div:nth-child(1) > div:nth-child(4) > div > table > tbody > tr > td > div > div > dl > dd > div > div > table > tbody > tr:nth-child("+fmt.Sprint(8+2*i)+") > td:nth-child("+fmt.Sprint(3+i%2)+")").Text(), " ??")
			material1 = strings.Trim(html.Find("div:nth-child(14) > div:nth-child(1) > div:nth-child(4) > div > table > tbody > tr > td > div > div > dl > dd > div > div > table > tbody > tr:nth-child("+fmt.Sprint(8+2*i)+") > td:nth-child("+fmt.Sprint(4+i%2)+")").Text(), " ??")
		}
		if material0 == "" {
			material0 = strings.Trim(html.Find("div:nth-child(15) > div:nth-child(1) > div:nth-child(4) > div > table > tbody > tr > td > div > div > dl > dd > div > div > table > tbody > tr:nth-child("+fmt.Sprint(8+2*i)+") > td:nth-child("+fmt.Sprint(3+i%2)+")").Text(), " ??")
			material1 = strings.Trim(html.Find("div:nth-child(15) > div:nth-child(1) > div:nth-child(4) > div > table > tbody > tr > td > div > div > dl > dd > div > div > table > tbody > tr:nth-child("+fmt.Sprint(8+2*i)+") > td:nth-child("+fmt.Sprint(4+i%2)+")").Text(), " ??")
		}
		skillUpMaterial[3+i][0] = stringToInt(material0)
		skillUpMaterial[3+i][1] = stringToInt(material1)
	}

	return skillUpMaterial
}

func getSkills(html *goquery.Selection) struct {
	Ex      exType
	Normal  skillType
	Passive skillType
	Sub     skillType
} {
	var skills struct {
		Ex      exType
		Normal  skillType
		Passive skillType
		Sub     skillType
	}
	skills.Ex.SkillName = html.Find("div:nth-child(13) > div:nth-last-child(8) > table > tbody > tr:nth-child(1) > td:nth-child(2) > div > strong").Text()
	if skills.Ex.SkillName == "" {
		skills.Ex.SkillName = html.Find("div:nth-child(14) > div:nth-last-child(8) > table > tbody > tr:nth-child(1) > td:nth-child(2) > div > strong").Text()
	}
	if skills.Ex.SkillName == "" {
		skills.Ex.SkillName = html.Find("div:nth-child(15) > div:nth-last-child(8) > table > tbody > tr:nth-child(1) > td:nth-child(2) > div > strong").Text()
	}
	for i := 0; i < 5; i++ {
		skills.Ex.Skill[i].Level = i + 1
		skills.Ex.Skill[i].Description = html.Find("div:nth-child(13) > div:nth-last-child(8) > table > tbody > tr:nth-child(" + fmt.Sprint(2+i) + ") > td:nth-child(2)").Text()
		cost := strings.Trim(html.Find("div:nth-child(13) > div:nth-last-child(8) > table > tbody > tr:nth-child("+fmt.Sprint(2+i)+") > td:nth-child(3)").Text(), "COST:")
		if skills.Ex.Skill[i].Description == "" {
			skills.Ex.Skill[i].Description = html.Find("div:nth-child(14) > div:nth-last-child(8) > table > tbody > tr:nth-child(" + fmt.Sprint(2+i) + ") > td:nth-child(2)").Text()
			cost = strings.Trim(html.Find("div:nth-child(14) > div:nth-last-child(8) > table > tbody > tr:nth-child("+fmt.Sprint(2+i)+") > td:nth-child(3)").Text(), "COST:")
		}
		if skills.Ex.Skill[i].Description == "" {
			skills.Ex.Skill[i].Description = html.Find("div:nth-child(15) > div:nth-last-child(8) > table > tbody > tr:nth-child(" + fmt.Sprint(2+i) + ") > td:nth-child(2)").Text()
			cost = strings.Trim(html.Find("div:nth-child(15) > div:nth-last-child(8) > table > tbody > tr:nth-child("+fmt.Sprint(2+i)+") > td:nth-child(3)").Text(), "COST:")
		}
		skills.Ex.Skill[i].Cost = stringToInt(cost)
	}
	for i := 0; i < 10; i++ {
		skills.Normal.SkillName = html.Find("div:nth-child(13) > div:nth-last-child(6) > table > tbody > tr:nth-child(1) > td:nth-child(2) > div > strong").Text()
		skills.Passive.SkillName = html.Find("div:nth-child(13) > div:nth-last-child(4) > table > tbody > tr:nth-child(1) > td:nth-child(2) > div > strong").Text()
		skills.Sub.SkillName = html.Find("div:nth-child(13) > div:nth-last-child(2) > table > tbody > tr:nth-child(1) > td:nth-child(2) > div > strong").Text()
		if skills.Normal.SkillName == "" {
			skills.Normal.SkillName = html.Find("div:nth-child(14) > div:nth-last-child(6) > table > tbody > tr:nth-child(1) > td:nth-child(2) > div > strong").Text()
			skills.Passive.SkillName = html.Find("div:nth-child(14) > div:nth-last-child(4) > table > tbody > tr:nth-child(1) > td:nth-child(2) > div > strong").Text()
			skills.Sub.SkillName = html.Find("div:nth-child(14) > div:nth-last-child(2) > table > tbody > tr:nth-child(1) > td:nth-child(2) > div > strong").Text()
		}
		if skills.Normal.SkillName == "" {
			skills.Normal.SkillName = html.Find("div:nth-child(15) > div:nth-last-child(6) > table > tbody > tr:nth-child(1) > td:nth-child(2) > div > strong").Text()
			skills.Passive.SkillName = html.Find("div:nth-child(15) > div:nth-last-child(4) > table > tbody > tr:nth-child(1) > td:nth-child(2) > div > strong").Text()
			skills.Sub.SkillName = html.Find("div:nth-child(15) > div:nth-last-child(2) > table > tbody > tr:nth-child(1) > td:nth-child(2) > div > strong").Text()
		}
		skills.Normal.Skill[i].Level = i + 1
		skills.Passive.Skill[i].Level = i + 1
		skills.Sub.Skill[i].Level = i + 1
		skills.Normal.Skill[i].Description = html.Find("div:nth-child(13) > div:nth-last-child(6) > table > tbody > tr:nth-child(" + fmt.Sprint(2+i) + ") > td:nth-child(2)").Text()
		skills.Passive.Skill[i].Description = html.Find("div:nth-child(13) > div:nth-last-child(4) > table > tbody > tr:nth-child(" + fmt.Sprint(2+i) + ") > td:nth-child(2)").Text()
		skills.Sub.Skill[i].Description = html.Find("div:nth-child(13) > div:nth-last-child(2) > table > tbody > tr:nth-child(" + fmt.Sprint(2+i) + ") > td:nth-child(2)").Text()
		if skills.Normal.Skill[i].Description == "" {
			skills.Normal.Skill[i].Description = html.Find("div:nth-child(14) > div:nth-last-child(6) > table > tbody > tr:nth-child(" + fmt.Sprint(2+i) + ") > td:nth-child(2)").Text()
			skills.Passive.Skill[i].Description = html.Find("div:nth-child(14) > div:nth-last-child(4) > table > tbody > tr:nth-child(" + fmt.Sprint(2+i) + ") > td:nth-child(2)").Text()
			skills.Sub.Skill[i].Description = html.Find("div:nth-child(14) > div:nth-last-child(2) > table > tbody > tr:nth-child(" + fmt.Sprint(2+i) + ") > td:nth-child(2)").Text()
		}
		if skills.Normal.Skill[i].Description == "" {
			skills.Normal.Skill[i].Description = html.Find("div:nth-child(15) > div:nth-last-child(6) > table > tbody > tr:nth-child(" + fmt.Sprint(2+i) + ") > td:nth-child(2)").Text()
			skills.Passive.Skill[i].Description = html.Find("div:nth-child(15) > div:nth-last-child(4) > table > tbody > tr:nth-child(" + fmt.Sprint(2+i) + ") > td:nth-child(2)").Text()
			skills.Sub.Skill[i].Description = html.Find("div:nth-child(15) > div:nth-last-child(2) > table > tbody > tr:nth-child(" + fmt.Sprint(2+i) + ") > td:nth-child(2)").Text()
		}
	}
	return skills
}

func stringToInt(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return -1
	}
	return val
}

func getAttackValue(initialValue int, starNum int) int {
	table := [3]float32{1.22, 1.1, 1.}
	return int(float32(initialValue) * table[starNum-1])
}
