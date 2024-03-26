package main

import (
	"encoding/csv"
	"fmt"
	"github.com/ahmetalpbalkan/go-linq"
	"github.com/xuri/excelize/v2"
	"golang.org/x/exp/slices"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const inputDir = "../../../docs/master/"
const outputDir = "../../../resource/master/"

func main() {
	fmt.Println("バージョンを指定してください。(ex: 1.0.0.0")

	var ver string
	_, err := fmt.Scan(&ver)
	if err != nil {
		log.Printf("バージョンの読み込みに失敗しました。%v", err)
		return
	}

	re, err := regexp.Compile(`^\d.\d.\d.\d$`)
	if err != nil {
		log.Printf("バージョン正規表現のコンパイルに失敗しました。%v", err)
		return
	}

	matches := re.FindAllString(ver, -1)
	if matches == nil {
		log.Printf("バージョンが正しいフォーマットではありません。%v", err)
		return
	}

	excelFiles, err := filepath.Glob(inputDir + "*.xlsx")
	if err != nil {
		log.Printf("エクセルファイルの読み込みに失敗しました。%v", err)
		return
	}

	for _, file := range excelFiles {
		log.Printf("%sをtsv出力します。", file)
		readAndWrite(ver, file)
		log.Printf("%sのtsv出力が完了しました。", file)
	}
}

func readAndWrite(ver string, file string) {
	f, err := excelize.OpenFile(file)
	defer func(f *excelize.File) {
		_ = f.Close()
	}(f)

	if err != nil {
		log.Fatalf("エクセルファイルの読み込みに失敗しました。%v", err)
	}

	sheet := getFileNameWithoutExt(file)
	rows, err := f.GetRows(sheet)

	p := fmt.Sprintf("%s%s.tsv", outputDir, sheet)
	data := parse(ver, rows)
	write(data, p)
}

func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

func parse(ver string, rows [][]string) [][]string {
	data := [][]string{}
	var ignoreColumnIdx []int
	for _, row := range rows {
		switch row[0] {
		case "#", "#N", ">=|memo":
			continue
		case "":
			d := []string{}
			for i, r := range row {
				if i == 0 {
					continue
				}
				d = append(d, r)
			}
			data = append(data, d)
			break
		case "#C":
			for i, r := range row {
				if i == 0 {
					continue
				}

				if r == "" {
					continue
				}

				if len(getIllegalVersions(ver, r)) > 0 {
					ignoreColumnIdx = append(ignoreColumnIdx, i)
				}
			}
			break
		case "##":
			d := []string{}
			for i, r := range row {
				if i == 0 {
					continue
				}

				if slices.Contains(ignoreColumnIdx, i) {
					continue
				}
				d = append(d, r)
			}
			data = append(data, d)
			break
		default:
			// バージョン指定は全ての条件に合致しているもののみを出力
			if len(getIllegalVersions(ver, row[0])) > 0 {
				continue
			}

			d := []string{}
			for i, r := range row {
				if i == 0 {
					continue
				}
				d = append(d, r)
			}
			data = append(data, d)
		}
	}

	return data
}

func write(data [][]string, p string) {
	var file *os.File
	_, err := os.Stat(p)
	if err == nil {
		file, err = os.Open(p)
	} else {
		file, err = os.Create(p)
	}
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	transform.NewWriter(file, japanese.ShiftJIS.NewEncoder())
	w := csv.NewWriter(file)
	w.Comma = '\t'

	// TODO: Access Deniedが出る。。。
	err = w.WriteAll(data)
	if err != nil {
		log.Printf("TSVファイルへの書き込みに失敗しました。%v", err)
		return
	}
}

func getIllegalVersions(ver string, tag string) []string {
	r := linq.
		From(strings.Split(tag, "&")).
		WhereT(func(s string) bool {
			t := strings.Split(s, "|")
			if len(t) != 2 {
				log.Printf("バージョンタグの指定ミス: %s", tag)
				return true
			}
			return !compareVersion(ver, t[1], t[0])
		})

	var v []string
	r.ToSlice(&v)

	return v
}

func compareVersion(ver1 string, ver2 string, ope string) bool {
	switch ope {
	case ">=":
		return ver1 >= ver2
	case ">":
		return ver1 > ver2
	case "<=":
		return ver1 <= ver2
	case "<":
		return ver1 < ver2
	case "=":
		return ver1 == ver2
	default:
		return false
	}
}
