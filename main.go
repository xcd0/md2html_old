package main

import ( // {{{
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	//"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	//"text/template"
) // }}}

func main() {
	flag.Parse()
	// 第一引数にマークダウンのファイルのパスを受け取る
	// 引数を元に構造体を作る
	fi := Argparse(flag.Arg(0))

	// wg4searchが終わるのを待ってwgをすすめる
	wg := sync.WaitGroup{}

	// 画像を探す
	searchTargetFile(&fi)

	// そのまま印刷したら単純な文書になる印刷用htmlを作成する 出力はしない
	wg.Add(1)
	go func() {
		defer wg.Done()
		//fmt.Println("call : makeHtmlNP")
		makeHtml(&fi)
	}()

	// スライド用htmlを生成する
	wg.Add(1)
	go func() {
		defer wg.Done()
		//fmt.Println("call : MakePdfForSlideNP")
		MakePdfForSlide(&fi)
	}()

	wg.Wait()

	// pdfをつくる
	/*
		fi.Pdfpath = fi.Dpath + fi.Basename + ".pdf"
		html2pdf.Html2pdf(fi)
	*/
}

func makeHtml(fi *Fileinfo) { // {{{

	fi.Html = Makeheader(*fi, "") + Makebody(fi.Apath, fi.RImgPath) + Makefooter()

	// 印刷用htmlを出力する
	err := ioutil.WriteFile(fi.Htmlpath, []byte(fi.Html), 0644)
	if err != nil {
		// Openエラー処理
		fmt.Fprintf(os.Stderr, "File %s could not open : %v\n", fi.Htmlpath, err)
		//fmt.Println(err)
		panic(err)
	}

} // }}}

func filter2body(in string) string { // {{{1

	lines := strings.Split(in, "\n")

	// 独自置換対象文字列
	var rep [][]string
	rep = append(rep, []string{`^===$`, "\n<div style='page-break-before:always'></div>\n"})

	output := ""
	for _, r := range rep { // 全ての置換対象文字列について回す
		reg := regexp.MustCompile(r[0])
		for _, line := range lines { // 一行ずつ
			output += reg.ReplaceAllString(line, r[1]) + "\n"
		}
	}
	// 上書きする
	return output

} // }}}1

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

func searchTargetFile(fi *Fileinfo) { // {{{1
	//fmt.Println("searchTargetFile")

	outputList := []string{}

	filepath.Walk(fi.Dpath, func(path string, _ os.FileInfo, _ error) error {
		// 相対パスを取得
		relativePath, _ := filepath.Rel(fi.Dpath, path)

		// jpg,jpeg,png,gif
		ext := filepath.Ext(relativePath)
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
			//fmt.Println(relativePath)
			outputList = append(outputList, relativePath)
		}
		return nil
	})

	// これを文字列の長い順に並び変える
	// でないと浅い階層に同じ名前のファイルがある場合誤った置換が発生する
	outputList = sortStirngsLen(outputList)

	fi.RImgPath = outputList
} // }}}1
