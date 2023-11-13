package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Rosya-edwica/speller/db"
	"github.com/Rosya-edwica/speller/logger"
	api "github.com/Rosya-edwica/speller/speller-api"
)

const POOLS_LIMIT = 1000

var POLLS_LEN int

var data = []db.Element{}

func main() {
	start := time.Now().Unix()
	items := db.GetItems()
	grouped_skills := groupItems(items)
	fmt.Println(len(grouped_skills))
	for i, group := range grouped_skills {
		POLLS_LEN = len(group)
		correctAllItems(group)
		fmt.Println("group: ", i)
	}
	db.SaveResult(data)
	fmt.Println(time.Now().Unix()-start, "sec.")
}

func groupItems(items []db.Element) (grouped [][]db.Element) {
	for i := 0; i < len(items); i += POOLS_LIMIT {
		group := items[i:]
		if len(group) >= POOLS_LIMIT {
			grouped = append(grouped, group[:POOLS_LIMIT])
		} else {
			grouped = append(grouped, group)
		}
	}
	return
}

func correctAllItems(items []db.Element) {
	var wg sync.WaitGroup
	wg.Add(POLLS_LEN)

	for _, item := range items {
		go correct(item, &wg)
	}
	wg.Wait()

}

func correct(item db.Element, wg *sync.WaitGroup) {
	fixedName := strings.Clone(item.Name)
	wrongWords := api.CheckText(item.Name)
	for _, word := range wrongWords {
		fixedName = strings.ReplaceAll(fixedName, word.WrongVersion, word.CorrectVersion[0])
	}
	if fixedName != item.Name {
		data = append(data, db.Element{Id: item.Id, Name: item.Name, NewName: fixedName})
		logger.Log.Printf("Ошибка - %s -> %s", item.Name, fixedName)
	}
	wg.Done()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
