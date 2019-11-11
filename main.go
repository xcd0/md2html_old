package main

import ( // {{{
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	//"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	//"sync"
	//"text/template"
) // }}}

func main() {
	flag.Parse()
	// 第一引数にマークダウンのファイルのパスを受け取る
	// 引数を元に構造体を作る
	fi := Argparse(flag.Arg(0))

	// そのまま印刷したら単純な文書になる印刷用htmlを出力
	makeHtml(&fi)

	// 独自置換対象文字列
	makeReplaceStrings(&fi)

	// 単純に生成した印刷用htmlに対して 置換などの処理を行った後出力する
	// 独自拡張の置換と画像のbase64置換
	filter2Html(&fi)

	// スライド用htmlを生成する
	MakeHtmlForSlide(&fi)

	// pdfをつくる
	/*
		fi.Pdfpath = fi.Dpath + fi.Basename + ".pdf"
		html2pdf.Html2pdf(fi)
	*/
}

func makeHtml(fi *Fileinfo) { // {{{

	fi.Flavor = ""
	fi.Html = Makehtml(fi)

	//log.Println(fi.Html)

	// 出力する
	err := ioutil.WriteFile(fi.Htmlpath, []byte(fi.Html), 0644)
	if err != nil {
		// Openエラー処理
		fmt.Fprintf(os.Stderr, "File %s could not open : %v\n", fi.Htmlpath, err)
		fmt.Println(err)
		return
	}

} // }}}

func filter2Html(fi *Fileinfo) { // {{{1

	//csspath := fi.Dpath + "slide.css"
	csspath := fi.Dpath + "markdown.css"
	header := Makeheader(*fi, csspath)
	body := Makebody(*fi)
	footer := Makefooter()

	// fi.Htmlを置換していく
	// 改行で分割
	fi.Html = header + body + footer
	lines := strings.Split(fi.Html, "\n")

	output := ""
	for _, line := range lines { // 一行ずつ
		//for i, line := range lines { // 一行ずつ
		for _, r := range fi.Rep { // 全ての置換対象文字列について回す
			if strings.Contains(line, r[0]) {
				//fmt.Printf( "replace %4d : %v -> %v\n%v\n %v\n", i, r[0], r[1], line, strings.Replace(line, r[0], r[1], 1))
				output += strings.Replace(line, r[0], r[1], 1) + "\n"
			} else {
				output += line + "\n"
			}
		}
	}
	// 上書きする
	fi.Html = output

	// 画像を上書きする
	replaceImg(fi)

	err := ioutil.WriteFile(fi.Htmlpath, []byte(fi.Html), 0644)
	if err != nil {
		// Openエラー処理
		fmt.Fprintf(os.Stderr, "File %s could not open : %v\n", fi.Htmlpath, err)
		fmt.Println(err)
		panic(err)
	}

} // }}}1

func makeReplaceStrings(fi *Fileinfo) { // {{ {
	// 独自拡張
	fi.Rep = append(fi.Rep, []string{"===", "<div style='page-break-before:always'></div>"})
} // }} }

func replaceImg(fi *Fileinfo) { // {{{
	// 画像の置き換え
	// 画像を探す
	outputList := searchTargetFile(*fi)
	// これを文字列の長い順に並び変える
	// でないと浅い階層に同じ名前のファイルがある場合誤った置換が発生する
	outputList = sortStirngsLen(outputList)

	output := "" // これが出力される

	// ここで一行ずつ処理
	lines := strings.Split(fi.Html, "\n")

	for _, line := range lines { // 一行ずつ
		//for j, line := range lines { // 一行ずつ
		// <img src=が含まれるはずなので前もって判別しておく
		if strings.Contains(line, "<img src=") == false {
			// <img src=が含まれていない
			output += line + "\n"
		} else {
			// <img src=が含まれる

			// リストにある画像のパスが含まれるか調べる
			for _, path := range outputList {
				// もしoutputListにおいて先にマッチしていたら
				// パスが\で区切られている場合を考えてすべて/にする
				line = strings.Replace(line, "\\", "/", -1)
				path = strings.Replace(path, "\\", "/", -1)
				// チェック
				if strings.Contains(line, path) == false {
					// 違うので次のpathに
					continue
				}
				// マッチ

				//fmt.Printf("---- line : %v , match : %v ----\n%v\n", j, path, line)

				// lineに含まれるパスの前後を切り出す
				// 例
				// <li><p><img src="./img/build_on_win10.gif" alt=""></p></li>
				//                 ↑                       ↑
				// このダブルクォーテーションの位置を調べる
				// ./が混じらないよう 一文字づつ前方に調べて"を探す

				// <img>にリンクが張ってあればそれを消す
				tmpLinkNum := strings.Index(line, "<a")
				if tmpLinkNum >= 0 {
					tmpLinkNumPost := -1
					for i := tmpLinkNum; i < len(line); i++ {
						if line[i] == '>' {
							// 見つかったのでそれをpostDQとする
							tmpLinkNumPost = i + 1
							break
						}
					}
					if tmpLinkNumPost < 0 {
						panic(1)
					} else {

						//fmt.Println("--->>\n" + line)
						line = line[:tmpLinkNum] + line[tmpLinkNumPost:]
						//fmt.Println("---<<\n" + line)
					}
					tmpLinkNum = strings.Index(line, "</a>")

					if tmpLinkNum >= 0 {
						//fmt.Println("--->>\n" + line)
						line = line[:tmpLinkNum] + line[tmpLinkNum+len("</a>"):]
						//fmt.Println("---<<\n" + line)
					}
				}

				// <imgより前にマッチしないようにする
				preDQ := -1
				tmpNum := strings.Index(line, "<img")
				tmpLine := line[tmpNum:]
				//fmt.Println("---tmpline\n" + tmpLine)

				for i := strings.Index(tmpLine, path); i >= 0; i-- {
					if tmpLine[i] == '"' {
						// 見つかったのでそれをpreDQとする
						preDQ = tmpNum + i + 1
						//fmt.Println("---match!\n" + line[:preDQ])
						break
					}
				}
				if preDQ < 0 {
					// ダブルクォーテーションが見つからなかった
					fmt.Println("マッチしましたが前にダブルクォーテーションがありませんでした")
					output += line + "\n"
					break
				}

				postDQ := -1
				for i := preDQ; i < len(line); i++ {
					if line[i] == '"' {
						// 見つかったのでそれをpostDQとする
						postDQ = i
						break
					}
				}
				if postDQ < 0 {
					// ダブルクォーテーションが見つからなかった
					fmt.Println("マッチしましたが後ろにダブルクォーテーションがありませんでした")
					/*
						fmt.Println("ori " + line)
						fmt.Println("cut pre  :" + line[preDQ:])
						fmt.Println("cut post :" + line[:])
					*/
					output += line + "\n"
					break
				}

				//
				/*
					fmt.Println("<<\n" + line)
					fmt.Printf(">> pre : %v '%v', post : %v '%v'\n", preDQ, line[preDQ], postDQ, line[postDQ])
					fmt.Printf(">> %v\n", line[preDQ:postDQ])

					code := "@@@@@@"
					/ */
				///*
				// 文字列が含まれているかどうか調べる
				// base64でエンコードしたデータに置き換える
				base64code := EncodeBase64(path)
				// ファイルの拡張子ごとにヘッダをつける
				// gif,png,jpg,jpegのみ
				ext := filepath.Ext(path) // 一致した画像ファイルのパスから拡張子を調べる
				code := ""
				if ext == ".gif" {
					code = "data:image/gif;base64," + base64code
				} else if ext == ".png" {
					code = "data:image/png;base64," + base64code
				} else if ext == ".jpg" || ext == ".jpeg" {
					code = "data:image/jpeg;base64," + base64code
				}

				// */
				// 前後を切り出す
				pre := line[:preDQ]
				post := line[postDQ:]

				// くっつけて上書きする
				line = pre + code + post

				//fmt.Println(">>\n" + line)

				// 置換したらその行は終わる
				break

			}
			// すべてのpathをチェックした
			//fmt.Println(line)
			// 出力する
			output += line + "\n"
		}
	}

	// 上書き
	fi.Html = output
} // }}}

func searchTargetFile(fi Fileinfo) []string { // {{{1

	outputList := []string{}

	filepath.Walk(fi.Dpath, func(path string, _ os.FileInfo, _ error) error {
		// 相対パスを取得
		relativePath, _ := filepath.Rel(fi.Dpath, path)

		// jpg,jpeg,png,gif
		ext := filepath.Ext(relativePath)
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
			fmt.Println(relativePath)
			outputList = append(outputList, relativePath)
		}
		return nil
	})

	return outputList
} // }}}1

/*
func replace(input, targetArray [][]string) (string, error) { // {{{

	lines := strings.Split(input, "\n")

	// 全ての置換対象文字列について回す
	for i := 0; i < len(targetArray); i++ {

		output := strings.NewRreplacer().Replace(lines)

	}
} // }}}
*/
//エンコード
func EncodeBase64(str string) string { // {{{

	file, err := os.Open(str)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	fi, _ := file.Stat() //FileInfo interface
	size := fi.Size()    //ファイルサイズ

	data := make([]byte, size)
	file.Read(data)

	return base64.StdEncoding.EncodeToString(data)
} // }}}

func sortStirngsLen(in []string) []string { // {{{

	type imgPath struct {
		path   string
		length int
	}

	lengthCount := make([]imgPath, len(in))
	for i, str := range in {
		lengthCount[i] = imgPath{path: str, length: len(str)}
	}

	// 大きい順に並べる
	sort.Slice(lengthCount, func(i, j int) bool { return lengthCount[i].length > lengthCount[j].length })

	out := make([]string, len(in))
	for i := 0; i < len(in); i++ {
		out[i] = lengthCount[i].path
	}
	return out
} // }}}
