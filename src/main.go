package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
)

type Score struct {
	Sum   int
	Count int
}

type Player map[string]*Score

type MeanScore map[int][]string

func main() {
	if len(os.Args) < 2 {
		log.Fatal("処理対象のゲームプレイログCSVファイルを指定してください。")
	}
	csvFile := os.Args[1]

	f, err := os.Open(csvFile)
	if err != nil {
		log.Fatal("指定されたファイルにアクセスできません。")
	}
	defer f.Close()

	r := csv.NewReader(f)

	// 1行目はヘッダー
	header, err := r.Read()
	if err != nil {
		log.Fatal("CSVファイルの読み込みに失敗しました。")
	}
	if !checkHeader(header) {
		log.Fatal("不正なCSVファイルです。")
	}

	// プレイヤーのスコアを集計
	p, err := LoadScore(r)
	if err != nil {
		log.Fatal(err)
	}

	// 平均スコアを計算
	m := CalcMeanScore(p)

	// ランキングを出力
	PrintRank(m)

}

// Headerのチェック
func checkHeader(header []string) bool {
	return header[0] == "create_timestamp" && header[1] == "player_id" && header[2] == "score"
}

func (p Player) NewPlayer(id string) {
	p[id] = &Score{}
}

func (p Player) AddScore(id string, score int) {
	if _, ok := p[id]; !ok {
		p.NewPlayer(id)
	}
	p[id].Sum += score
	p[id].Count++
}

func (s Score) GetAverageScore(id string) int {
	return int(math.Round(float64(s.Sum) / float64(s.Count)))
}
func LoadScore(r *csv.Reader) (Player, error) {
	p := Player{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return p, errors.New("CSVファイルの読み込みに失敗しました。")
		}
		score, err := strconv.Atoi(record[2])
		if err != nil {
			return p, errors.New("スコアの読み込みに失敗しました。")
		}
		p.AddScore(record[1], score)
	}
	return p, nil
}

func (m MeanScore) NewMeanScore(score int) {
	m[score] = []string{}
}

func (m MeanScore) AddPlayer(score int, id string) {
	if _, ok := m[score]; !ok {
		m.NewMeanScore(score)
	}
	m[score] = append(m[score], id)
}

func CalcMeanScore(p Player) MeanScore {
	m := MeanScore{}
	for id, s := range p {
		score := s.GetAverageScore(id)
		m.AddPlayer(score, id)
	}
	return m
}

func SortMeanScore(m MeanScore) []int {
	meanScore := make([]int, 0, len(m))
	for s := range m {
		meanScore = append(meanScore, s)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(meanScore)))
	return meanScore
}

func PrintRank(m MeanScore) {
	keys := SortMeanScore(m)
	rank := 1
	rankLimit := 10

	fmt.Println("rank,player_id,mean_score")

	for _, k := range keys {
		sort.Strings(m[k])
		for _, id := range m[k] {
			fmt.Printf("%d,%s,%d\n", rank, id, k)
		}
		rank += len(m[k])
		if rank > rankLimit {
			break
		}
	}
}
