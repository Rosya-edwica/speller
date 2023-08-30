package main

import (
	"fmt"
	"strings"
	"time"
	"sync"

	api "github.com/Rosya-edwica/speller/speller-api"
	"github.com/Rosya-edwica/speller/logger"
	"github.com/Rosya-edwica/speller/db"
)

const POOLS_LIMIT = 1000
var POLLS_LEN int

var data = []db.Skill{}

func main(){
	start := time.Now().Unix()
	skills := db.GetSkills()
	grouped_skills := groupSkills(skills)
	fmt.Println(len(grouped_skills))
	for i, group := range grouped_skills{
		POLLS_LEN = len(group)
		correctAllSkills(group)
		fmt.Println("group: ", i)
	}
	db.SaveResult(data)
	fmt.Println(time.Now().Unix() - start, "sec.")
}

func groupSkills(skills []string) (grouped_skills [][]string){
	for i:=0; i<len(skills); i+=POOLS_LIMIT {
		group := skills[i:]
		if len(group) >= POOLS_LIMIT{
			grouped_skills = append(grouped_skills, group[:POOLS_LIMIT])
		} else {
			grouped_skills = append(grouped_skills, group)
		}
	}
	return
}

func correctAllSkills(skills []string) {
	var wg sync.WaitGroup
	wg.Add(POLLS_LEN)

	for _, skill := range skills {
		go correct(skill, &wg)
	}
	wg.Wait()

}

func correct(skill string, wg *sync.WaitGroup) {
	correctedSkill := strings.Clone(skill)
	wrongWords := api.CheckText(skill)
	for _, word := range wrongWords {
		correctedSkill = strings.ReplaceAll(correctedSkill, word.WrongVersion, word.CorrectVersion[0])
	}
	if correctedSkill != skill {
		data = append(data, db.Skill{CurrentName: skill, FixedName: correctedSkill})
		logger.Log.Printf("Ошибка - %s -> %s", skill, correctedSkill)
	} else {
		data = append(data, db.Skill{CurrentName: skill, FixedName: ""})
		logger.Log.Printf("Нет ошибок - %s", skill)
	}
	wg.Done()
}

func checkErr(err error) {
	if err != nil{
		panic(err)
	}
}