package main

import ( // {{{
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	makeHtmlByShurcooL(fi)

	// 単純に生成した印刷用htmlを
	// 置換などの処理を行う
	filter2Html(fi)
	// スライド用htmlにする
	//makeHtmlForSlide(fi)

	// pdfをつくる
	/*
		fi.Pdfpath = fi.Dpath + fi.Basename + ".pdf"
		html2pdf.Html2pdf(fi)
	*/
}

func makeHtmlByShurcooL(fi Fileinfo) { // {{{

	//fi.Flavor = "gfm"
	fi.Flavor = "github"
	// htmlを作成する
	html, err := Makehtml(fi)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 出力する
	err = ioutil.WriteFile(fi.Htmlpath, []byte(html), 0644)
	if err != nil {
		// Openエラー処理
		fmt.Fprintf(os.Stderr, "File %s could not open : %v\n", fi.Htmlpath, err)
		fmt.Println(err)
		return
	}

} // }}}

func filter2Html(fi Fileinfo) { // {{{1

	header := Makeheader(fi)
	body, err := Makebody(fi)
	if err != nil {
		panic(err)
	}
	footer := Makefooter()

	fi.Html = header + body + footer

	/*
			js := `<!--{{{-->
		<style>

		<script type="text/javascript">
			document.onkeydown = keydown;

		function keydown() {
			target.innerHTML = "キーが押されました KeyCode :" + event.keyCode;

			if (event.keyCode == 37) {
				// 左
			}
			else {
			}
			if (event.keyCode == 39) {
				// 右
			}
			else {
			}
			document.getElementById( "slider-main" ).onclick = function( event ) {
				var x = event.pageX ;	// 水平の位置座標
				var obj = document.getElementById("slide-main");
				var w = obj.getBoundingClientRect().width;
				console.log(w);
				if ( x < w / 2  ) {
					// 左
				}else{
					// 右
				}
			}
		}
		</script>
		<!--}}}-->
		`
	*/

	/*

		// 3階層分でjpg,jpeg,png,gifを検索してリストにする
		repFilePath := searchTargetFile(fi)
		if repFilePath == "err" {
			panic(repFilePath)
		}
	*/

	// 置換対象文字列
	rep := [][]string{}
	rep = append(rep, []string{"===", "<div style='page-break-before:always'></div>"})

	// fi.Htmlを置換していく
	// 改行で分割
	lines := strings.Split(fi.Html, "\n")

	output := ""
	// 一行ずつ
	for i, line := range lines {
		// 全ての置換対象文字列について回す
		for _, r := range rep {
			if strings.Contains(line, r[0]) {
				fmt.Printf("replace %4d : %v -> %v\n", i, r[0], r[1])
				output += strings.Replace(line, r[0], r[1], 1)
			} else {
				output += line
			}
		}
	}

	// これが出力される
	//output, err := replace(f, rep)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File %s couldn't create. : %v\n", fi.Htmlpath, err)
		return
	}

	htmlpath := fi.Dpath + fi.Basename + "_slide.html"
	err = ioutil.WriteFile(htmlpath, []byte(output), 0644)
	if err != nil {
		// Openエラー処理
		fmt.Fprintf(os.Stderr, "File %s could not create. : %v\n", htmlpath, err)
		fmt.Println(err)
		return
	}

} // }}}1

func searchTargetFile(fi Fileinfo) string { // {{{1

	var ch chan string

	go func() {
		filepath.Walk(fi.Dpath, func(path string, _ os.FileInfo, _ error) error {
			// 相対パスを取得
			relativePath, _ := filepath.Rel(fi.Dpath, path)

			// jpg,jpeg,png,gif
			ext := filepath.Ext(relativePath)
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
				fmt.Println(relativePath)
				ch <- relativePath
			}
			return nil
		})
		defer close(ch)
		defer fmt.Println("done")
	}()

	for str := range ch {
		fmt.Println(str)
	}

	return "test"
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
