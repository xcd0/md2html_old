package main

import ( // {{{
	"fmt"
	"io/ioutil"
	//"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	//"sync"
	//"text/template"
) // }}}

type Info struct {
	width            float64
	height           float64
	flagSize         bool
	page             int
	flagPrintPageNum bool
	footer           string
	h1               string
	h2               string
	h3               string
	output           string
	line             string
	path             string
}

func MakeHtmlForSlide(fi *Fileinfo) {

	fi.TmpPath = fi.Dpath + "/.tmp"
	if f, err := os.Stat(fi.TmpPath); os.IsNotExist(err) || !f.IsDir() {
		// ない
		// 作る
		if _, err := os.Stat(fi.TmpPath); os.IsNotExist(err) {
			os.Mkdir(fi.TmpPath, 0777)
		}
	} else {
		// ある
		// 消す
		if err := os.Remove(fi.TmpPath); err != nil {
			fmt.Println(err)
		}
		// 作る
		if _, err := os.Stat(fi.TmpPath); os.IsNotExist(err) {
			os.Mkdir(fi.TmpPath, 0777)
		}
	}

	// ファイル開く
	f, err := os.Open(fi.Dpath)
	if err != nil {
		fmt.Println("error")
	}
	defer f.Close()

	// 一気に全部読み取り
	b, err := ioutil.ReadAll(f)
	md := string(b)

	// 改行で分ける
	lines := strings.Split(md, "\n")

	info := Info{
		// 初期値A4
		width:            841.89,
		height:           595.28,
		flagSize:         false,
		page:             0,
		flagPrintPageNum: true,
		footer:           "",
		h1:               "none",
		h2:               "none",
		h3:               "none",
		output:           "",
		line:             "",
		path:             fi.TmpPath,
	}

	regH1 := regexp.MustCompile(`^# `)
	regH2 := regexp.MustCompile(`^## `)
	regH3 := regexp.MustCompile(`^### `)
	for _, info.line = range lines { // 一行ずつ
		//for i, line := range lines { // 一行ずつ
		// <!-- $で始まっていたら独自記法
		if strings.Contains(info.line, "<!-- $") {
			preamble(&info)
		} else if strings.Contains(info.line, "===") {
			// mdを===で分割する
			dividePage(&info)
		} else {
			if regH1.MatchString(info.line) {
				info.h1 = info.line[2:]
			}
			if regH2.MatchString(info.line) {
				info.h2 = info.line[3:]
			}
			if regH3.MatchString(info.line) {
				info.h3 = info.line[4:]
			}
			info.output += info.line + "\n"
		}
	}

}

func preamble(info *Info) {
	if info.flagSize == false {
		tmp := strings.Index(info.line, "<!-- $width:\"")
		if tmp >= 0 {
			info.width, _ = strconv.ParseFloat(info.line[tmp+len("<!-- $width:\""):strings.Index(info.line, "\" -->")], 64)
			return
		}
		tmp = strings.Index(info.line, "<!-- $height:\"")
		if tmp >= 0 {
			info.height, _ = strconv.ParseFloat(info.line[tmp+len("<!-- $height:\""):strings.Index(info.line, "\" -->")], 64)
			return
		}
	}
	tmp := strings.Index(info.line, "<!-- $size:\"16:9\" -->")
	if tmp >= 0 {
		info.flagSize = true
		info.width = 841.89
		info.height = 841.89 * 9.0 / 16.0
		return
	}
	tmp = strings.Index(info.line, "<!-- $size:\"4:3\" -->")
	if tmp >= 0 {
		info.flagSize = true
		info.width = 841.89
		info.height = 841.89 * 3.0 / 4.0
		return
	}
	tmp = strings.Index(info.line, "<!-- $page_number:\"")
	if tmp >= 0 {
		info.page, _ = strconv.Atoi(info.line[tmp+len("<!-- $page_number:\"") : strings.Index(info.line, "\" -->")])
		return
	}
	tmp = strings.Index(info.line, "<!-- $set_page_number:\"true\" -->")
	if tmp >= 0 {
		info.flagPrintPageNum = true
		return
	}
	tmp = strings.Index(info.line, "<!-- $set_page_number:\"false\" -->")
	if tmp >= 0 {
		info.flagPrintPageNum = false
		return
	}
	tmp = strings.Index(info.line, "<!-- $footer:\"")
	if tmp >= 0 {
		return
	}
}

func dividePage(info *Info) {

	if info.page > 0 || info.flagPrintPageNum {
		// ページ数を出力しやすいようにする
		info.line += fmt.Sprintf("", info.page)
	}
	// ここまでを出力
	//fmt.Printf( "replace %4d : %v -> %v\n%v\n %v\n", i, r[0], r[1], line, strings.Replace(line, r[0], r[1], 1))
	filename := fmt.Sprintf("%v/%04d.md", info.path, info.page)
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filename, []byte(info.output), 0666)
	if err != nil {
		panic(err)
	}

	// リセット
	info.output = ""
	if info.page > 0 {
		info.output += "<!-- h1:\"" + info.h1 + "\" -->"
		info.output += "<!-- h2:\"" + info.h2 + "\" -->"
		info.output += "<!-- h3:\"" + info.h3 + "\" -->"
		info.output += "<!-- width:\"" + strconv.FormatFloat(info.width, 'f', 2, 64) + "\" -->"
		info.output += "<!-- height:\"" + strconv.FormatFloat(info.height, 'f', 2, 64) + "\" -->"
	}
}
