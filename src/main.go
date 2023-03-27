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

// rankLimit は出力するランキングの上限です。
const rankLimit = 10

type Score struct {
	Sum   int
	Count int
}

type PlayerData map[string]*Score

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

// checkHeader はヘッダーの内容をチェックします。
func checkHeader(header []string) bool {
	return header[0] == "create_timestamp" && header[1] == "player_id" && header[2] == "score"
}

// AddScore はプレイヤーのスコアを集計します。
func (d PlayerData) AddScore(id string, score int) {
	if _, ok := d[id]; !ok {
		// プレイヤーが存在しない場合は初期化
		d[id] = &Score{}
	}
	d[id].Sum += score
	d[id].Count++
}

// CalcMeanScore はプレイヤーの平均スコアを計算します。
// math.Round で四捨五入しています。
func (s Score) CalcMeanScore() int {
	return int(math.Round(float64(s.Sum) / float64(s.Count)))
}

// LoadScore はCSVファイルからスコアを読み込みます。
func LoadScore(r *csv.Reader) (PlayerData, error) {
	p := PlayerData{}
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

// AddPlayer はプレイヤーのスコアを追加します。
func (m MeanScore) AddPlayer(score int, id string) {
	if _, ok := m[score]; !ok {
		m[score] = []string{}
	}
	m[score] = append(m[score], id)
}

// CalcMeanScore はプレイヤーの平均スコアを計算します。
func CalcMeanScore(d PlayerData) MeanScore {
	m := MeanScore{}
	for id, s := range d {
		score := s.CalcMeanScore()
		m.AddPlayer(score, id)
	}
	return m
}

// SortMeanScore は平均スコアの降順でソートします。
func SortMeanScore(m MeanScore) []int {
	meanScore := make([]int, 0, len(m))
	for s := range m {
		meanScore = append(meanScore, s)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(meanScore)))
	return meanScore
}

// PrintRank はランキングを出力します。
func PrintRank(m MeanScore) {
	keys := SortMeanScore(m)
	rank := 1

	fmt.Println("rank,player_id,mean_score")

	for _, k := range keys {
		// 同率の場合プレイヤーID順にソート
		sort.Strings(m[k])
		for _, id := range m[k] {
			fmt.Printf("%d,%s,%d\n", rank, id, k)
		}
		rank += len(m[k])
		// rankLimit に達したら終了
		if rank > rankLimit {
			break
		}
	}
}
