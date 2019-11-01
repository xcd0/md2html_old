package main

import ( // {{{
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	//"strings"

	//"./html2pdf"
	"./convertFromMarkdownToHtml"
) // }}}

func main() {
	flag.Parse()
	// 第一引数にマークダウンのファイルのパスを受け取る
	// 引数を元に構造体を作る
	fi := md2html.Argparse(flag.Arg(0))

	// そのまま印刷したら単純な文書になる印刷用htmlを出力
	makeHtmlByShurcooL(fi)
}

/*
	// 単純に生成した印刷用htmlをスライド用htmlにする
	makeHtmlForSlide(fi)

	// pdfをつくる
	fi.Pdfpath = fi.Dpath + fi.Basename + ".pdf"
	html2pdf.Html2pdf(fi)
}

*/
func makeHtmlByShurcooL(fi md2html.Fileinfo) { // {{{

	fi.Flavor = "gfm"
	// htmlを作成する
	html, err := md2html.Makehtml(fi)
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

/*
func makeHtmlForSlide(fi md2html.Fileinfo) { // {{{1

	header := Makeheader()
	body, err := Makebody(fi)
	footer := Makefooter()

	fi.Html = header + body + footer

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

	// 置換対象文字列
	targetArray := [][]string{
		{"</head>", "\n</head>"},
		{"<body>", "\n<body>"},
		{"</body>", "\n</body>"},
		{"<!-- pb -->", "</div>\n<hr id=\"pb\">\n<div id=\"child\">"},
	}

	// これが出力される
	output, err := replace(f, targetArray)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File %s couldn't create. : %v\n", fi.Htmlpath, err)
		return
	}

	htmlpath = fi.Dpath + fi.hasename + "_slide.html"
	err = ioutil.WriteFile(htmlpath, []byte(output), 0644)
	if err != nil {
		// Openエラー処理
		fmt.Fprintf(os.Stderr, "File %s could not create. : %v\n", htmlpath, err)
		fmt.Println(err)
		return
	}

} // }}}1

func replace(input, targetArray [][]string) (string, error) {

	lines := strings.Split(input, "\n")

	// 全ての置換対象文字列について回す
	for i := 0; i < len(targetArray); i++ {

		output := strings.NewRreplacer().Replace(lines)

	}
}
*/
