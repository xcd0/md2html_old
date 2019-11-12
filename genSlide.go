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
		body := Makebody(mdpath, fi.RImgPath)

		bodies = append(bodies, body)

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

	genSlideHtml(bodies, filepath.Join(fi.Dpath, fi.Basename+"_slide.html"))

} // }}}

func returnPreBody(maxpage int) string { // {{{
	return `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
	<link rel="stylesheet" type="text/css" href="./github-markdown.css" />
	<script type="text/javascript" async
	src="https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.5/MathJax.js?config=TeX-MML-AM_CHTML" async>
	</script>
	<script type="text/x-mathjax-config">
	MathJax.Hub.Config({tex2jax: {inlineMath: [['$','$']]}});
	</script>
	<style type="text/css">
	<!--
body{font-family:Helvetica,arial,sans-serif;font-size:14px;line-height:1.6;padding:0;background-color:#fff;color:#333x}a{color:#4183c4;text-decoration:nonex}a.absent{color:#c00x}a.anchor{display:block;padding-left:30px;margin-left:-30px;cursor:pointer;position:absolute;top:0;left:0;bottom:0x}h1,h2,h3,h4,h5,h6{margin:20px 0 10px;padding:0;font-weight:700;-webkit-font-smoothing:antialiased;cursor:text;position:relativex}h1:first-child,h1:first-child+h2,h2:first-child,h3:first-child,h4:first-child,h5:first-child,h6:first-child{margin-top:0;padding-top:0x}h1:hover a.anchor,h2:hover a.anchor,h3:hover a.anchor,h4:hover a.anchor,h5:hover a.anchor,h6:hover a.anchor{text-decoration:nonex}h1 code,h1 tt,h2 code,h2 tt,h3 code,h3 tt,h4 code,h4 tt,h5 code,h5 tt,h6 code,h6 tt{font-size:inheritx}h1{font-size:30px;border-bottom:2px solid #bbb;margin-bottom:20pxx}h1,h2{color:#000x}h2{font-size:24px;border-bottom:1px solid #cccx}h3{font-size:18pxx}h4{font-size:1pcx}h5,h6{font-size:14pxx}h6{color:#777x}blockquote,dl,li,ol,p,pre,table,ul{margin:15px 0x}hr{border:0 0 0;color:#ccc;height:4px;padding:0x}a:first-child h1,a:first-child h2,a:first-child h3,a:first-child h4,a:first-child h5,a:first-child h6,body>h1:first-child,body>h1:first-child+h2,body>h2:first-child,body>h3:first-child,body>h4:first-child,body>h5:first-child,body>h6:first-child{margin-top:0;padding-top:0x}h1 p,h2 p,h3 p,h4 p,h5 p,h6 p{margin-top:0x}li p.first{display:inline-blockx}ol,ul{padding-left:30pxx}ol:first-child,ul:first-child{margin-top:0x}ol:last-child,ul:last-child{margin-bottom:0x}dl,dl dt{padding:0x}dl dt{font-size:14px;font-weight:700;font-style:italic;margin:15px 0 5pxx}dl dt:first-child{padding:0x}dl dt>:first-child{margin-top:0x}dl dt>:last-child{margin-bottom:0x}dl dd{margin:0 0 15px;padding:0 15pxx}dl dd>:first-child{margin-top:0x}dl dd>:last-child{margin-bottom:0x}blockquote{border-left:4px solid #ddd;padding:0 15px;color:#777x}blockquote>:first-child{margin-top:0x}blockquote>:last-child{margin-bottom:0x}table{padding:0;border-spacing:2px;border-collapse:collapsex}table,td,th{border:1px solid #cccx}td,th{padding:0;margin:0}table tr{background-color:#fff;border-top:1px solid #c6cbd1;margin:0;padding:0x}table tr:nth-child(2n){background-color:#f6f8fax}table tr th{font-weight:700x}table tr td,table tr th{border:1px solid #ccc;text-align:left;margin:0;padding:6px 13pxx}table tr td:first-child,table tr th:first-child{margin-top:0x}table tr td:last-child,table tr th:last-child{margin-bottom:0x}img{max-width:100%;max-height:80%}span.frame,span.frame>span{display:block;overflow:hiddenx}span.frame>span{border:1px solid #ddd;float:left;margin:13px 0 0;padding:7px;width:autox}span.frame span img{display:block;float:leftx}span.frame span span{clear:both;color:#333;display:block;padding:5px 0 0x}span.align-center{display:block;overflow:hidden;clear:bothx}span.align-center>span{display:block;overflow:hidden;margin:13px auto 0;text-align:centerx}span.align-center span img{margin:0 auto;text-align:centerx}span.align-right{display:block;overflow:hidden;clear:bothx}span.align-right>span{display:block;overflow:hidden;margin:13px 0 0;text-align:rightx}span.align-right span img{margin:0;text-align:rightx}span.float-left{display:block;margin-right:13px;overflow:hidden;float:leftx}span.float-left span{margin:13px 0 0x}span.float-right{display:block;margin-left:13px;overflow:hidden;float:rightx}span.float-right>span{display:block;overflow:hidden;margin:13px auto 0;text-align:rightx}code,tt{margin:0 2px;padding:0 5px;white-space:nowrap;border:1px solid #eaeaea;background-color:#f8f8f8;border-radius:3pxx}pre code{margin:0;padding:0;white-space:pre;border:0;background:transparentx}.highlight pre,pre{background-color:#f8f8f8;border:1px solid #ccc;font-size:13px;line-height:19px;overflow:auto;padding:6px 10px;border-radius:3px}pre code,pre tt{background-color:transparent;border:0x}.sub{position:absolute}.main-content{max-width:100%;max-height:100%;margin:auto;padding:50px}hr{page-break-before:always;page-break-after:always;border:0!important;color:#fff;height:4px;padding:0}
	-->
	</style>
<script type="text/javascript">
<!--

function goto(to_page) {
	var page = document.getElementById(current);
	page.style.visibility = "hidden";
	current = to_page;
	page = document.getElementById(current);
	page.style.visibility = "visible";
}

function resize() {
	//document.documentElement.clientWidth
	//document.documentElement.clientHeight

	var c1 = document.getElementById('container');
	var h = window.innerHeight;
	var w = window.innerWidth;
	h = (h - 100) * 0.9;
	w = (w - 100) * 0.9;
	c1.style.Width= w + 'px';
	c1.style.Height = h + 'px';

	/*
	var c document.getElementById(current);
	var tags c.getElementByTagName("img");
	for(var i = 0; i < tags.length; i++){
		tag[i].style.Width = w + 'px';
		tag[i].style.Height = h + 'px';
	}
	*/
}

var current_num;
var current;
function zeroPadding(num,length){
	return ('0000000000' + num).slice(-length);
}
function main() {
	resize();
	current_num = zeroPadding(0,4);
	current = "p" + String(current_num);
	goto(current);
}
window.onload = main;


window.addEventListener('load', resize, false);
window.addEventListener('resize', resize, false);

function prev(){

	var page = document.getElementById(current);
	page.style.visibility = "hidden";

	if( current_num == "0000" ){
		current_num = zeroPadding( ` + fmt.Sprintf("%04d", maxpage-1) + ` ,4);
	} else {
		current_num = zeroPadding( Number(current_num) - 1 ,4);
	}
	current = "p" + current_num;

	page = document.getElementById(current);
	page.style.visibility = "visible";
}
function next(){
	var page = document.getElementById(current);
	page.style.visibility = "hidden";

	if( current_num == "` + fmt.Sprintf("%04d", maxpage-1) + `" ){
		current_num = "0000"
		current_num = zeroPadding( Number(current_num) + 1 ,4);
	} else {
		current_num = zeroPadding( Number(current_num) + 1 ,4);
	}
	current = "p" + current_num;

	page = document.getElementById(current);
	page.style.visibility = "visible";
}
window.addEventListener("click", function(e) {
	console.log("offset:" + e.offsetX + "," + e.offsetY); 
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
	var delta = e.deltaY ? -(e.deltaY) : e.wheelDelta ? e.wheelDelta : -(e.detail);
	if (delta < 0){
		e.preventDefault();
		//下にスクロールした場合の処理
		next();
	} else if (delta > 0){
		e.preventDefault();
		//上にスクロールした場合の処理
		prev();
	}
}

//-->
</script>
</head>
<body>
<div id="container" class="main-content">
`
} // }}}

func genSlideHtml(bodies []string, slidePath string) { // {{{

	// bodyのみのsliceから1つにまとめたhtmlを生成する。
	output := returnPreBody(len(bodies))
	for i, b := range bodies {
		output += fmt.Sprintf("\n<div id=\"p%04d\" id=\"container\" class=\"sub\" style=\"visibility: hidden;\">", i)
		output += b
		if i == 0 {
			output += fmt.Sprintf("<p><a href=\"javascript:goto('p%04d')\">次へ</a></p>", i+1)
		} else if i < len(bodies)-1 {
			output += fmt.Sprintf("<p><a href=\"javascript:goto('p%04d')\">戻る | <a href=\"javascript:goto('p%04d')\">次へ</a></p>", i-1, i+1)
		} else {
			output += fmt.Sprintf("<p><a href=\"javascript:goto('p%04d')\">表紙へ</a></p>", 0)
		}
		output += "</div>"
	}
	output += "</div>\n</body>\n</html>\n"

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
		tmpdir = filepath.Join(info.tmp_path, "body")
		if _, err := os.Stat(tmpdir); os.IsNotExist(err) {
			os.Mkdir(tmpdir, 0777)
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
