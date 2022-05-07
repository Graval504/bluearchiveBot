package bluearchiveBot

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func GetCharacterInfoFromList() {
	url := "https://google.com/search?q=site:https://codealone.tistory.com/65"
	res, _ := http.Get(url)
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	fmt.Println(doc.Find("h3").Text())
}

type Student struct{
	Name []string
	Link string
	School string
	Club string
	InitialFigures int
	InitialStarNum int
	AttackValue int
	StudentType []string
	AreaType []string
	Weapon string
	Equipment []string
	Ooparts []string
	ExUpMaterial [][2]int
	SkillUpMaterial [][2]int
	Skills struct {
		Ex skilltype
		Normal skilltype
		Passive skilltype
		sub skilltype
	}
	
}

type skilltype struct{
	SkillName string
	Skill []struct{
		SkillCost int
		SkillDescription string

	}
}