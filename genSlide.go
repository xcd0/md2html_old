package main

import ( // {{{
	"fmt"
	"io/ioutil"
	//"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	//"sync"
	//"text/template"
	//"github.com/xcd0/go-nkf"
) // }}}

type Info struct { // {{{
	width               float64
	height              float64
	flagSize            bool
	page                int
	absolute_page       int
	bool_print_page_num bool
	bool_print_title    bool
	//h1                  string // 廃止
	h2          string
	h3          string
	footer      string
	preamble    string
	output      string
	line        string
	debug       bool
	state_code  bool
	tmp_path    string
	state_title int // このページにh1~h3の表記があるか、あればその数値 なければ-1
} // }}}

func MakePdfForSlide(fi *Fileinfo) { // {{{
	//fmt.Println("MakePdfForSlide")

	var info Info

	// マークダウンファイルを調べ、ページごとに分割する
	parseMd(&info, fi)

	// ページごとに分割したmdをhtmlに変換する

	page_num := info.absolute_page // 0から始まり実際のページ数で終了している
	for i := 0; i < page_num; i++ {
		// 1ページづつ
		htmlpath := filepath.Join(fi.Dpath, ".tmp", "html", fmt.Sprintf("%04d.html", i))
		mdpath := filepath.Join(fi.Dpath, ".tmp", "md", fmt.Sprintf("%04d.md", i))

		//fmt.Printf(">>> htmlpath = %v\n", htmlpath)

		// htmlのbodyを生成
		body := Makebody(mdpath, fi.RImgPath)

		// 画像を置換してhtmlを作る
		html := Makeheader(*fi, fi.Dpath+"markdown.css") + body + Makefooter()

		// スライド用htmlの1ページを出力する
		if err := ioutil.WriteFile(htmlpath, []byte(html), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "File %s could not open : %v\n", htmlpath, err)
			fmt.Println(err)
			panic(err)
		}
	}

	// 中間ディレクトリを削除
	if f, err := os.Stat(info.tmp_path); os.IsNotExist(err) || f.IsDir() {
		// 存在するので消す
		if info.debug == false {
			// デバッグ中でないとき
			if err := os.Remove(info.tmp_path); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("debug:中間ディレクトリを削除していません")
			}
		}
	}

} // }}}

func readPreamble(info *Info) { // {{{1
	if info.absolute_page == 0 { // 表紙ページのみで指定できる項目 {{{
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
	} // }}}

	// どのページでも指定できる項目
	// ページ内に表記するページ数
	tmp := strings.Index(info.line, "<!-- $page_number:\"")
	if tmp >= 0 {
		info.page, _ = strconv.Atoi(info.line[tmp+len("<!-- $page_number:\"") : strings.Index(info.line, "\" -->")])
		return
	}
	// ページ内にページ数を表記する
	tmp = strings.Index(info.line, "<!-- $set_page_number:\"true\" -->")
	if tmp >= 0 {
		info.bool_print_page_num = true
		return
	}
	// ページ内にページ数を表記しない
	tmp = strings.Index(info.line, "<!-- $set_page_number:\"false\" -->")
	if tmp >= 0 {
		info.bool_print_page_num = false
		return
	}
	// フッターを設定
	tmp = strings.Index(info.line, "<!-- $footer:\"")
	if tmp >= 0 {
		return
	}
	// 各ページにタイトルを表示する
	tmp = strings.Index(info.line, "<!-- $title:\"true\" -->")
	if tmp >= 0 {
		info.bool_print_title = true
		return
	}
	// 各ページにタイトルを表示しない
	tmp = strings.Index(info.line, "<!-- $title:\"false\" -->")
	if tmp >= 0 {
		info.bool_print_title = false
		return
	}
} // }}}1

func dividePage(info *Info) { // {{{

	//fmt.Printf("ap : %v", info.absolute_page)

	if info.absolute_page == 0 {
		// 表紙ページのみ
		if f, err := os.Stat(info.tmp_path); os.IsNotExist(err) || f.IsDir() {
			// ある
			// 消す
			if err := os.RemoveAll(info.tmp_path); err != nil {
				fmt.Println(err)
				if info.debug {
					fmt.Println("debug:既存の中間ディレクトリを削除しました")
				}
			}
		}
		// 作る
		if _, err := os.Stat(info.tmp_path); os.IsNotExist(err) {
			os.Mkdir(info.tmp_path, 0777)
			//fmt.Printf("%vを作成しました", info.tmp_path)
		}
		tmpdir := filepath.Join(info.tmp_path, "md")
		if _, err := os.Stat(tmpdir); os.IsNotExist(err) {
			os.Mkdir(tmpdir, 0777)
			//fmt.Printf("%vを作成しました", tmpdir)
		}
		tmpdir = filepath.Join(info.tmp_path, "html")
		if _, err := os.Stat(tmpdir); os.IsNotExist(err) {
			os.Mkdir(tmpdir, 0777)
		}
		tmpdir = filepath.Join(info.tmp_path, "pdf")
		if _, err := os.Stat(tmpdir); os.IsNotExist(err) {
			os.Mkdir(tmpdir, 0777)
		}
		if info.debug {
			//fmt.Println("debug:中間ディレクトリを作成しました")
			//fmt.Println(info.tmp_path)
		}
	}

	outputOnePage(info)

	// 出力終わったので次のページ用に設定する
	info.page++
	info.absolute_page++
	info.state_code = false
	info.state_title = -1

	// リセット
	info.output = ""
	info.preamble = ""

} // }}}

func outputOnePage(info *Info) { // {{{

	// フッタ
	info.line += fmt.Sprintf("<footer>%v</footer>", info.footer)

	if info.absolute_page >= 0 || info.bool_print_page_num {
		// ページ数
	}

	if info.absolute_page > 0 && info.bool_print_page_num {
		// 最初のページではなく、ページ数を表示するように設定されていたら
		info.line += fmt.Sprintf("<div class=\"page_num>%v\"</div>", info.page)
	}

	// ここまでを出力

	tmpfilename := filepath.Join(info.tmp_path, "md", fmt.Sprintf("%04d.md", info.absolute_page))

	/*
		if info.debug {
			fmt.Println(">> " + tmpfilename)
		}
	*/

	file, err := os.Create(tmpfilename)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	switch info.state_title {
	case 2: // このページ内に ## があった 引き継がない。
		info.h2 = "none"
		fallthrough
	case 3: // このページ内に ### があった 引き継がない。
		info.h3 = "none"
	default:
		// 何もしない
	}

	// h2 と h3のみ前のページから引き継いで表記ことができる
	title := ""
	if info.bool_print_title {
		title += "\n<!-- 前のページから引き継いだタイトル ここから -->\n"
		if info.h2 != "none" { // h2が指定されてなかったら何もしない

			switch info.state_title {
			case -1: // このページ内に ## と ### がなかった
				// 前のページから引き継ぐ
				title += "## " + info.h2 + "\n"
				if info.h3 != "none" {
					title += "### " + info.h3 + "\n"
				}
			case 2: // このページ内に ## があった 引き継がない。
				// 何もしない
			case 3: // このページ内に ### があった 引き継がない。
				if info.h3 != "none" {
					title += "## " + info.h2 + "\n"
				}
			default:
				// 何もしない
			}
		}
		title += "<!-- 前のページから引き継いだタイトル ここまで -->\n"
	}

	// 現在の設定を確認用に出力する 特に使わない
	info.preamble += "<!-- 自動生成されたプリアンブル ここから -->\n"
	info.preamble += fmt.Sprintf("<!-- // $width:\"%v\" -->\n", info.width)
	info.preamble += fmt.Sprintf("<!-- // $height:\"%v\" -->\n", info.height)
	info.preamble += fmt.Sprintf("<!-- // $page_number:\"%v\" -->\n", info.bool_print_page_num)
	info.preamble += fmt.Sprintf("<!-- // $page:\"%v\" -->\n", info.page)
	info.preamble += fmt.Sprintf("<!-- // $absolute_page:\"%v\" -->\n", info.absolute_page)
	//info.preamble += fmt.Sprintf("<!-- // $h1:\"%v\" -->\n", info.h1)
	info.preamble += fmt.Sprintf("<!-- // $h2:\"%v\" -->\n", info.h2)
	info.preamble += fmt.Sprintf("<!-- // $h3:\"%v\" -->\n", info.h3)
	info.preamble += fmt.Sprintf("<!-- // $title:\"%v\" -->\n", info.bool_print_title)
	info.preamble += fmt.Sprintf("<!-- // $state_title:\"%v\" -->\n", info.state_title)
	info.preamble += "<!-- 自動生成されたプリアンブル ここまで -->\n"

	// プリアンブルとタイトルをくっつける
	output := info.preamble + title + info.output

	err = ioutil.WriteFile(tmpfilename, []byte(output), 0666)
	if err != nil {
		panic(err)
	}

} // }}}

func initInfo(info *Info, fi *Fileinfo) { // {{{
	info.width = 841.89
	info.height = 595.28
	info.flagSize = false
	info.page = 0          // これはページ内での表記に使用されるページ数
	info.absolute_page = 0 // これは表記とは関係ないページ数
	info.bool_print_page_num = true
	info.bool_print_title = true // 各ページ毎にタイトルを表示するかどうか
	//h1=                  "none" // 廃止
	info.h2 = "none"
	info.h3 = "none"
	info.footer = ""
	info.preamble = ""
	info.output = ""
	info.line = ""
	info.debug = true
	info.state_code = false // これは```で囲まれている内側かどうかを保持する
	info.tmp_path = ""      // .tmpフォルダの絶対パス
	info.state_title = -1   // このページにh1~h3の表記があるか、あればその数値 なければ-1
	info.tmp_path = filepath.Join(fi.Dpath, ".tmp")
} // }}}

func parseMd(info *Info, fi *Fileinfo) { // {{{

	// 初期設定
	initInfo(info, fi)

	// 改行で分ける
	lines := strings.Split(ReadMd(fi.Filename), "\n")

	//regH1 := regexp.MustCompile(`^# `)
	regH2 := regexp.MustCompile(`^## `)
	regH3 := regexp.MustCompile(`^### `)
	regCode := regexp.MustCompile("^```")
	regPreamble := regexp.MustCompile(`^<!-- \$`)
	regPageBreak := regexp.MustCompile(`^===$`)
	for _, info.line = range lines { // 一行ずつ
		//fmt.Printf(":%v", info.line)
		//for i, line := range lines { // 一行ずつ
		// <!-- $で始まっていたら独自記法
		if regPreamble.MatchString(info.line) {
			// マークダウンファイル内に記述されているプリアンブルを読み込む
			readPreamble(info)
		} else if regPageBreak.MatchString(info.line) {
			// mdを===で分割する
			// === は出力されない
			dividePage(info)
		} else {
			// h2,h3を設定
			if regH2.MatchString(info.line) {
				// 上書き
				info.h2 = info.line[3:]
				info.state_title = 2
				// h3はリセット
				info.h3 = "none"
			}
			if regH3.MatchString(info.line) {
				// 上書き
				info.h3 = info.line[4:]
				info.state_title = 3
			}
			if regCode.MatchString(info.line) {
				// 論理反転
				info.state_code = !info.state_code
			}

			if info.state_code == true {
				// そのまま出力する ただし行末に半角空白を付与する
				info.output += info.line + "  \n"
			} else {
				// ```で囲まれている場合は空白を入れない
				info.output += info.line + "\n"
			}
		}
	}
	// 最後の1ページを出力する
	dividePage(info)

} // }}}
