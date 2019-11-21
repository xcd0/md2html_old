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

	var bodies []string

	page_num := info.absolute_page // 0から始まり実際のページ数で終了している
	for i := 0; i < page_num; i++ {
		// 1ページづつ
		htmlpath := filepath.Join(fi.Dpath, ".tmp", "html", fmt.Sprintf("%04d.html", i))
		mdpath := filepath.Join(fi.Dpath, ".tmp", "md", fmt.Sprintf("%04d.md", i))

		//fmt.Printf(">>> htmlpath = %v\n", htmlpath)

		// htmlのbodyを生成
		body := Makebody(mdpath, fi.RImgPath, "")

		bodies = append(bodies, body)

		// 1ページ分の完全なhtmlを作る
		html := Makeheader(*fi, fi.Dpath+"slide.css") + body + Makefooter()

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
			if err := os.RemoveAll(info.tmp_path); err != nil {
				fmt.Println(err)
			} else if info.debug {
				fmt.Println("debug:中間ディレクトリを削除していません")
			}
		}
	}

	// すべてのスライドを含むhtmlを生成
	genSlideHtml(bodies, filepath.Join(fi.Dpath, fi.Basename+"_slide.html"))

} // }}}

func returnPreBody(maxpage int, css, postHead string) string { // {{{
	return `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
	<script type="text/javascript" async
	src="https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.5/MathJax.js?config=TeX-MML-AM_CHTML" async>
	</script>
	<script type="text/x-mathjax-config">
	MathJax.Hub.Config({tex2jax: {inlineMath: [['$','$']]}});
	</script>
	<style type="text/css">
	<!--
	` + css + `
	-->
	</style>
	<script type="text/javascript">
<!--
function goto(prev, next) {
	document.getElementById(prev).style.visibility = "hidden";
	document.getElementById(next).style.visibility = "visible";
}

var current_num;
var current;
function zeroPadding(num,length){
	return ('0000000000' + num).slice(-length);
}

function setPage(p,n){
	goto(p,n)
	location.hash = current_num;
}

function init(){
	var urlHash = location.hash.slice(1);
	if( urlHash == false ){
		current_num = zeroPadding(0,4);
	} else {
		if( isNaN(urlHash) ){
			current_num = "0000"
			location.hash = current_num;
		} else {
			current_num = zeroPadding(urlHash,4);
			var n = Number(current_num);
			if( n < 0 || n > ` + fmt.Sprintf("%d", maxpage-1) + ` ){
				current_num = "0000"
				location.hash = current_num;
			}else{
			}
		}
	}
	current = "p" + String(current_num);
	main = document.getElementById("container");
	setPage( current, current )
}

function prev(){
	p = current
	if( current_num != "0000" ){
		current_num = zeroPadding( Number(current_num) - 1 ,4);
	}
	current = "p" + current_num;
	setPage( p, current )
}
function next(){
	p = current
	if( current_num != "` + fmt.Sprintf("%04d", maxpage-1) + `" ){
		current_num = zeroPadding( Number(current_num) + 1 ,4);
	}
	current = "p" + current_num;
	setPage( p, current )
}

function resizeFontSize(p,r) {
	//p.style.fontSize = x + '%';
	p.style.zoom = r;
	p.style.MozTransform = r;
	p.style.WebkitTransform = r;
}

function checkSize(c){
	var r = document.getElementById(c).getBoundingClientRect();
	console.log(
		"cH : " + c.clientHeight
		+ "cRH : " + r.Height
		+ ", clientH : " + document.documentElement.clientHeight
	);
}

function shurinkPage(c){
	//checkSize(c);
	if( c.clientHeight > document.documentElement.clientHeight ){
		i= 0.95;
		resizeFontSize(c, i)
		for(
			;
			i > 0.7
			&& c.clientHeight > document.documentElement.clientHeight;
			i -= 0.01
		) {
			resizeFontSize(c, i);
			//checkSize(c);
		}
	}
}

window.addEventListener("click", function(e) {
	var w = window.innerWidth;
	const ratio = 0.1;
	var left = w * ratio;
	var right = w * ( 1 - ratio );
	//console.log("offset:" + e.offsetX + "," + e.offsetY);
	if( e.offsetX < left ){
		prev()
		console.log("offsetX :" + e.offsetX);
	} else if( e.offsetX > right ){
		next()
		console.log("offsetX :" + e.offsetX);
	}
});

var mousewheelevent = 'onwheel' in document ? 'wheel' : 'onmousewheel' in document ? 'mousewheel' : 'DOMMouseScroll';
try{
	document.addEventListener (mousewheelevent, onWheel, false);
}catch(e){
	//for legacy IE
	document.attachEvent ("onmousewheel", onWheel);
}
function onWheel(e) {
	if(!e) e = window.event; //for legacy IE
	//e.preventDefault();
	var delta = e.deltaY ? -(e.deltaY) : e.wheelDelta ? e.wheelDelta : -(e.detail);
	if (delta < 0){
		//下にスクロールした場合の処理
		next();
	} else if (delta > 0){
		//上にスクロールした場合の処理
		prev();
	}
}

function keydownfunc( event ) {
	var key_code = event.keyCode;
	if( key_code === 33 ) { prev(); } // PageUp
	if( key_code === 34 ) { next(); } // PageDown
	if( key_code === 37 ) { prev(); } // ←
	if( key_code === 38 ) { prev(); } // ↑
	if( key_code === 39 ) { next(); } // →
	if( key_code === 40 ) { next(); } // ↓
	checkSize(current)
}

window.onload = init;

var countInterval = 0;

addEventListener("keydown", keydownfunc, false);
//window.addEventListener('load', resize, false);
//window.addEventListener('resize', resize, false);

` +
		//postHead +
		`
//-->
</script>
</head>
<body>
`
} // }}}

func genSlideHtml(bodies []string, slidePath string) { // {{{

	cssPath := filepath.Join(filepath.Dir(slidePath), "slide.css")

	slideCss := Minify(cssPath)
	if slideCss == "default" {
		slideCss = `/* builtin css*/
body{font-family:Helvetica,arial,sans-serif;font-size:3.5vmin;padding:0;margin:0;line-height:1.6;background-color:#fff;color:#444;word-wrap:break-word}body>:first-child{margin-top:0!important}body>:last-child{margin-bottom:0!important}a{color:#4183c4;text-decoration:none}a.absent{color:#c00}a.anchor{display:block;padding-left:30px;margin-left:-30px;cursor:pointer;position:absolute;top:0;left:0;bottom:0x}h1,h2,h3,h4,h5,h6{margin:20px 0 10px;padding:0;font-weight:700;-webkit-font-smoothing:antialiased;cursor:text;position:relative}h1:first-child,h1:first-child+h2,h2:first-child,h3:first-child,h4:first-child,h5:first-child,h6:first-child{margin-top:0;padding-top:0}h1:hover a.anchor,h2:hover a.anchor,h3:hover a.anchor,h4:hover a.anchor,h5:hover a.anchor,h6:hover a.anchor{text-decoration:none}h1 code,h1 tt,h2 code,h2 tt,h3 code,h3 tt,h4 code,h4 tt,h5 code,h5 tt,h6 code,h6 tt{font-size:inherit}h1{font-size:250%;margin-bottom:40px;padding-bottom:0}h1,h2{color:#000}h2{font-size:6vh;border-bottom:2px solid #ccc;margin-bottom:3%}h3{font-size:4vh;margin-top:-1.9%;border-bottom:1px solid #ddd;color:#555}h4{font-size:80%}h5{font-size:60%}h6{font-size:40%;color:#777}blockquote,dl,li,ol,p,pre,table{margin:2% 0}pre{white-space:pre-wrap}code{white-space:pre-wrap;word-wrap:break-word}li,ul{margin:.2em 0}hr{border:0 0 0;height:4px;padding:0}a:first-child h1,a:first-child h2,a:first-child h3,a:first-child h4,a:first-child h5,a:first-child h6,bo dy>h5:first-child,body>h1:first-child,body>h1:first-child+h2,body>h2:first-child,body>h3:first-child,body>h4:first-child,body>h6:first-child{margin-top:0;padding-top:0}h1 p,h2 p,h3 p,h4 p,h5 p,h6 p{margin-top:0}li p.first{display:inline-block}ol,ul{padding-left:30px}ol:first-child,ul:first-child{margin-top:0}dl,dl dt{padding:0}dl dt{font-size:14px;font-weight:700;font-weight:1400;font-style:italic;margin:15px 0 5px}dl dt:first-child{padding:0}dl dt>:first-child{margin-top:0}dl dt>:last-child{margin-bottom:0}dl dd{margin:0 0 15px;padding:0 15px}dl dd>:first-child{margin-top:0}dl dd>:last-child{margin-bottom:0}blockquote{border-left:4px solid #ddd;padding:0 15px;color:#777}blockquote>:first-child{margin-top:0}blockquote>:last-child{margin-bottom:0}table{padding:0;border-spacing:2px;border-collapse:collapse;max-width:90%;margin:auto}table,td,th{border:1px solid #ccc;font-size:90%}td,th{padding:0;margin:0}table tr{background-color:#fff;border-top:1px solid #c6cbd1;margin:0;padding:0}table tr:nth-child(2n){background-color:#f6f8fa}table tr th{font-weight:700;white-space:nowrap}table tr td,table tr th{border:1px solid #ccc;text-align:center;margin:0;padding:6px 13px}table tr td:first-child,table tr th:first-child{margin-top:0}img{max-height:100%;max-width:100%}code,tt{margin:1%;padding:.2% 1%;white-space:nowrap;border:1px solid #eaeaea;background-color:#f8f8f8;border-radius:3px}pre code{margin:0;padding:0;white-space:pre;border:0;background:0 0}.highlight pre,pre{border:1px solid #ccc;font-size:13px;line-height:19px;overflow:auto;padding:6px 10px;border-radius:3px}pre code,pre tt{background-color:transparent;border:0}.main-content{position:relative}.sub{position:absolute;padding:5vh 10vw;width:80vw;height:90vh;overflow-y:hidden}hr{border:0!important;color:#fff;height:4px}.page_num{border:0;position:absolute;right:10;bottom:10}.controller{width:50px;position:absolute;right:0;bottom:0}
`
	}

	jsDocArray := ""
	jsDocArray += "\n"
	//for i, p := range bodies {
	//	// hereに一旦入れる
	//	jsDocArray += fmt.Sprintf("\nvar here = function() {\n/*<!--start-->%v\n<!--fin-->\n*/};\n", p)
	//	//jsDocArray += fmt.Sprintf("\nvar here = function() {\n/*<!--start-->%v\n<!--fin-->\n*/};\n", p)
	//	// 不要部分を切り取ってグローバル変数に代入
	//	jsDocArray += fmt.Sprintf("const p%04v", i) + ` = here.toString().match(/\/\*([\s\S]*)\*\//).pop();` //.toString().match(/\/\*([\s\S])*\*\//).pop();`
	//	//jsDocArray += fmt.Sprintf("const p%04v", i) + ` = here.toString().match(/\/\*([\s\S]*)\*\//).pop();` //.toString().match(/\/\*([\s\S])*\*\//).pop();`
	//}
	jsDocArray += "\n"
	jsDocArray += "console.log(p0000);\n"

	output := returnPreBody(len(bodies), slideCss, jsDocArray)

	// bodyのみのsliceから1つにまとめたhtmlを生成する。
	output += "<div id=\"container\" class=\"main-content\">"

	for i, p := range bodies {
		output += fmt.Sprintf("<div id=\"p%04d\" id=\"container\" class=\"sub\" style=\"visibility: hidden;\">\n", i)
		output += fmt.Sprintf("%v\n", p)
		output += fmt.Sprintf("</div>\n")
	}
	output += "</div>" // main-content

	output += "</body>\n</html>\n"

	output = delEmptyLine(output)

	// スライド用htmlの1ページを出力する
	if err := ioutil.WriteFile(slidePath, []byte(output), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "File %s could not open : %v\n", slidePath, err)
		fmt.Println(err)
		panic(err)
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

	// ページ数
	if info.absolute_page > 0 && info.bool_print_page_num {
		// 最初のページではなく、ページ数を表示するように設定されていたら
		info.line += fmt.Sprintf("<div class=\"page_num>%v\"</div>", info.page)
	}

	// 表紙以外のページで引き継ぐかどうか処理する
	/*
		if info.absolute_page != 0 {
			switch info.state_title {
			case 2: // このページ内に ## があった h2 h3 両方とも引き継がない。
			// 既に###はリセットされているので無視してよい
			case 3: // このページ内に ### があった h2のみ引き継ぐ
				info.h3 = "none"
			default:
				// 何もしない
			}
		}
	*/

	// h2 と h3のみ前のページから引き継いで表記ことができる
	title := ""
	if info.absolute_page != 0 && info.bool_print_title {
		title += "\n<!-- 前のページから引き継いだタイトル ここから -->\n"
		if info.h2 != "none" { // h2が指定されてなかったら何もしない
			switch info.state_title {
			case -1: // このページ内に ## と ### がなかった
				// 前のページから引き継ぐ
				if info.h2 != "none" {
					title += "## " + info.h2 + "\n"
				}
				if info.h3 != "none" {
					title += "### " + info.h3 + "\n"
				}
			case 2: // このページ内に ## があった h2 h3 ともに引き継がない。
				// 何もしない
			case 3: // このページ内に ### があった h3は 引き継がない。h2を引き継ぐ
				if info.h2 != "none" {
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

	// ここまでを出力

	tmpfilename := filepath.Join(info.tmp_path, "md", fmt.Sprintf("%04d.md", info.absolute_page))

	file, err := os.Create(tmpfilename)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	if err = ioutil.WriteFile(tmpfilename, []byte(output), 0666); err != nil {
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
	info.debug = false
	info.state_code = false // これは```で囲まれている内側かどうかを保持する
	info.tmp_path = ""      // .tmpフォルダの絶対パス
	info.state_title = -1   // このページにh1~h3の表記があるか、あればその数値 なければ-1
	info.tmp_path = filepath.Join(fi.Dpath, ".tmp")
} // }}}

// 正規表現オブジェクトのコンパイルは1回でよいのでグローバル領域で行う
//regH1 := regexp.MustCompile(`^# `)
var regH2 = regexp.MustCompile(`^## `)
var regH3 = regexp.MustCompile(`^### `)
var regCode = regexp.MustCompile("^```")
var regPreamble = regexp.MustCompile(`^<!-- \$`)

//var regPageBreak1 = regexp.MustCompile(`^===$`)
//var regPageBreak2 = regexp.MustCompile(`^<!---->$`)
//var regPageBreak3 = regexp.MustCompile(`^<!-- === -->$`)
var regPageBreak = regexp.MustCompile(`^<!---->$`)

func parseMd(info *Info, fi *Fileinfo) { // {{{

	// 初期設定
	initInfo(info, fi)

	// 改行で分ける
	lines := strings.Split(ReadMd(fi.Filename), "\n")

	for _, info.line = range lines { // 一行ずつ
		//fmt.Printf(":%v", info.line)
		//for i, line := range lines { // 一行ずつ
		// <!-- $で始まっていたら独自記法
		if regPreamble.MatchString(info.line) {
			// マークダウンファイル内に記述されているプリアンブルを読み込む
			readPreamble(info)
		} else if regPageBreak.MatchString(info.line) {
			/*
				// 改ページを複数用意していたが削除した
					} else if regPageBreak1.MatchString(info.line) ||
						regPageBreak2.MatchString(info.line) ||
						regPageBreak3.MatchString(info.line) {
			*/
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
				if info.state_title != 2 {
					info.state_title = 3
				}
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

func delEmptyLine(in string) string { // {{{

	output := ""

	r := regexp.MustCompile(`^$`)
	lines := strings.Split(in, "\n")
	for _, line := range lines {
		if r.MatchString(line) == false {
			output += line + "\n"
		}
	}
	return output
} // }}}
