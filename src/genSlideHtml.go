package main

import ( // {{{
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	//"log"
	"os"
	//"sync"
	//"text/template"
	//"github.com/xcd0/go-nkf"
) // }}}

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

func returnBuiltinCss() string { // {{{
	css := `
	body {
		font-family: Helvetica, arial, sans-serif;
		font-size: 3.5vmin;
		padding: 0;
		margin: 0;
		line-height: 1.6;
		background-color: #fff;
		color: #444;
		word-wrap: break-word
	}

	body>:first-child {
		margin-top: 0!important
	}

	body>:last-child {
		margin-bottom: 0!important
	}

	a {
		color: #4183c4;
		text-decoration: none
	}

	a.absent {
		color: #c00
	}

	a.anchor {
		display: block;
		padding-left: 30px;
		margin-left: -30px;
		cursor: pointer;
		position: absolute;
		top: 0;
		left: 0;
		bottom: 0
	}

	h1, h2, h3, h4, h5, h6 {
		margin: 20px 0 10px;
		padding: 0;
		-webkit-font-smoothing: antialiased;
		cursor: text;
		position: relative
	}

	h1:first-child, h1:first-child+h2, h2:first-child, h3:first-child, h4:first-child, h5:first-child, h6:first-child {
		margin-top: 0;
		padding-top: 0
	}

	h1:hover a.anchor, h2:hover a.anchor, h3:hover a.anchor, h4:hover a.anchor, h5:hover a.anchor, h6:hover a.anchor {
		text-decoration: none
	}

	h1 code, h1 tt, h2 code, h2 tt, h3 code, h3 tt, h4 code, h4 tt, h5 code, h5 tt, h6 code, h6 tt {
		font-size: inherit
	}

	h1 {
		font-size: 6vh;
		font-weight: 700;
		margin-bottom: 40px;
		padding-bottom: 0
	}

	h1, h2 {
		color: #000
	}

	h2 {
		font-size: 6vh;
		border-bottom: 2px solid #ccc;
		margin-bottom: 3%
	}

	h3 {
		font-size: 4vh;
		margin-top: -1.9%;
		border-bottom: 1px solid #ddd;
		color: #555
	}

	h4 {
		font-size: 80%
	}

	h5 {
		font-size: 60%
	}

	h6 {
		font-size: 40%;
		color: #777
	}

	blockquote, dl, li, ol, p, pre, table {
		margin: 2% 0
	}

	code, pre {
		white-space: pre-wrap
	}

	code {
		word-wrap: break-word
	}

	li, ul {
		margin: .2em 0
	}

	hr {
		border: 0;
		height: 4px;
		padding: 0
	}

	a:first-child h1, a:first-child h2, a:first-child h3, a:first-child h4, a:first-child h5, a:first-child h6, bo dy>h5:first-child, body>h1:first-child, body>h1:first-child+h2, body>h2:first-child, body>h3:first-child, body>h4:first-child, body>h6:first-child {
		margin-top: 0;
		padding-top: 0
	}

	h1 p, h2 p, h3 p, h4 p, h5 p, h6 p {
		margin-top: 0
	}

	li p.first {
		display: inline-block
	}

	ol, ul {
		padding-left: 30px
	}

	ol:first-child, ul:first-child {
		margin-top: 0
	}

	dl, dl dt {
		padding: 0;
		font-size: 4vmin
	}

	dl dt {
		font-weight: 700;
		margin: 15px 0 5px
	}

	dl dt:first-child {
		padding: 0
	}

	dl dt>:first-child {
		margin-top: 0
	}

	dl dt>:last-child {
		margin-bottom: 0
	}

	dl dd {
		margin: 0 0 15px;
		padding: 0 15px
	}

	dl dd>:first-child {
		margin-top: 0
	}

	dl dd>:last-child {
		margin-bottom: 0
	}

	blockquote {
		border-left: 4px solid #ddd;
		padding: 0 15px;
		color: #777
	}

	blockquote>:first-child {
		margin-top: 0
	}

	blockquote>:last-child {
		margin-bottom: 0
	}

	table {
		padding: 0;
		border-spacing: 2px;
		border-collapse: collapse;
		max-width: 90%;
		margin: auto
	}

	table, td, th {
		border: 1px solid #ccc;
		font-size: 90%
	}

	table tr, td, th {
		padding: 0;
		margin: 0
	}

	table tr {
		background-color: #fff;
		border-top: 1px solid #c6cbd1
	}

	table tr:nth-child(2n) {
		background-color: #f6f8fa
	}

	table tr th {
		font-weight: 700;
		white-space: nowrap
	}

	table tr td, table tr th {
		border: 1px solid #ccc;
		text-align: center;
		margin: 0;
		padding: 6px 13px
	}

	table tr td:first-child, table tr th:first-child {
		margin-top: 0
	}

	img {
		max-height: 100%;
		max-width: 100%
	}

	code, tt {
		margin: 1%;
		padding: .2% 1%;
		white-space: nowrap;
		border: 1px solid #eaeaea;
		background-color: #f8f8f8;
		border-radius: 3px
	}

	pre code {
		margin: 0;
		padding: 0;
		white-space: pre;
		border: 0;
		background: 0
	}

	.highlight pre, pre {
		border: 1px solid #ccc;
		font-size: 13px;
		line-height: 19px;
		overflow: auto;
		padding: 6px 10px;
		border-radius: 3px
	}

	pre code, pre tt {
		background-color: transparent;
		border: 0
	}

	.main-content {
		position: relative
	}

	.sub {
		position: absolute;
		padding: 5vh 10vw;
		width: 80vw;
		height: 90vh;
		overflow-y: hidden
	}

	hr {
		border: 0!important;
		color: #fff;
		height: 4px
	}

	#page_num {
		border: 0;
		position: absolute;
		left: 95vw;
		top: 90vh;
		margin-right: 0;
		bottom: 10
	}

	.controller {
		width: 50px;
		position: absolute;
		right: 0;
		bottom: 0
	}
	`

	return css
} // }}}

//! @brief head内に埋め込むcssを出力する
func returnCssForSlide(slidePath string) string { // {{{
	cssPath := filepath.Join(filepath.Dir(slidePath), "slide.css")
	slideCss := Minify(cssPath)
	// slide.cssがないとき slideCssに"default"が入る
	log.Println(slideCss)
	if slideCss == "default" {
		slideCss = MinifyCssString(returnBuiltinCss())
		log.Println(slideCss)
	}
	return "<style type=\"text/css\">" + slideCss + "</style>\n"
} // }}}

//! @brief head内のjavascriptを出力する
func returnJavascriptForSlide(maxpage int) string { // {{{
	builtinJs := `
	function goto(prev, next) {
		document.getElementById(prev).style.visibility = "hidden";
		document.getElementById(next).style.visibility = "visible";
		if( current_num != "0000" ){
			document.getElementById("page_num").innerHTML = Number(current_num);
		}else{
			document.getElementById("page_num").innerHTML = "";
		}
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

	var bool_ctrl_key = false;

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
		if( bool_ctrl_key === false ){
			if (delta < 0){
				//下にスクロールした場合の処理
				next();
			} else if (delta > 0){
				//上にスクロールした場合の処理
				prev();
			}
		}
	}

	function keydownfunc( event ) {
		var key_code = event.keyCode;
		if( bool_ctrl_key === false ){
			if( key_code === 17 ) { bool_ctrl_key = true; } // ctrlキー
		}
		if( key_code === 33 ) { prev(); } // pageup
		if( key_code === 34 ) { next(); } // pagedown
		if( key_code === 37 ) { prev(); } // ←
		if( key_code === 38 ) { prev(); } // ↑
		if( key_code === 39 ) { next(); } // →
		if( key_code === 40 ) { next(); } // ↓
	}
	function keyupfunc( event ) {
		var key_code = event.keyCode;
		if( key_code === 17 ) { bool_ctrl_key = false; } // ctrlキー
	}

	window.onload = init;

	var countInterval = 0;

	addEventListener("keyup" , keyupfunc);
	addEventListener("keydown", keydownfunc, false);
	//window.addEventListener('load', resize, false);
	//window.addEventListener('resize', resize, false);
	`
	return "\t<script type=\"text/javascript\">" + MinifyJavascriptString(builtinJs) + "</script>\n"
} // }}}

//! @brief headを出力する
func returnHeadForSlide(maxpage int, slidePath string) string { // {{{
	// katex新しいの行内で動かなかった
	/*
		<!-- オンラインの時用 -->
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.12.0/dist/katex.min.css" integrity="sha384-AfEj0r4/OFrOo5t7NnNe46zW/tFgW6x/bCJG8FqQCEo3+Aro6EYUG4+cU+KJWu/X" crossorigin="anonymous">
		<script defer src="https://cdn.jsdelivr.net/npm/katex@0.12.0/dist/katex.min.js" integrity="sha384-g7c+Jr9ZivxKLnZTDUhnkOnsh30B4H0rpLUpJ4jAIKs4fnJI+sEnkvrMWph2EDg4" crossorigin="anonymous"></script>
		<script defer src="https://cdn.jsdelivr.net/npm/katex@0.12.0/dist/contrib/auto-render.min.js" integrity="sha384-mll67QQFJfxn0IYznZYonOWZ644AWYC+Pt2cHqMaRhXVrursRwvLnLaebdGIlYNa" crossorigin="anonymous" onload="renderMathInElement(document.body);"></script>

		<!-- オフラインの時用 katexフォルダがある前提 -->
		<link rel="stylesheet" href="katex/katex.min.css">
		<script src="katex/katex.min.js"></script>
		<script src="katex/auto-render.min.js" onload="renderMathInElement(document.body);"></script>
	*/
	return `
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/KaTeX/0.7.1/katex.min.css">
	<script src="https://cdnjs.cloudflare.com/ajax/libs/KaTeX/0.7.1/katex.min.js"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/KaTeX/0.7.1/contrib/auto-render.min.js"></script>
	` + returnCssForSlide(slidePath) + returnJavascriptForSlide(maxpage)
} // }}}

//! @brief スライド用htmlファイルを生成する
func genSlideHtml(bodies []string, slidePath string) { // {{{

	// ヘッダー部分を生成
	output := `<!DOCTYPE html>
	<html lang="en">
	<head>
	` + returnHeadForSlide(len(bodies), slidePath) + `</head>
	<body>
	<div id="page_num"></div>
	`
	// 1ページごとに区切られているbodiesからスライド用htmlを生成する
	output += "<div id=\"container\" class=\"main-content\">"

	// マークダウンから生成した1ページごとのhtmlをdivで囲んでhiddenの状態で出力する
	for i, p := range bodies {
		output += fmt.Sprintf("<div id=\"p%04d\" id=\"container\" class=\"sub\" style=\"visibility: hidden;\">\n", i)
		output += fmt.Sprintf("%v\n", p)
		output += fmt.Sprintf("</div>\n") // container
		log.Printf("--- %v ---\n", i)
		log.Printf("%v\n", p)
	}
	output += "</div>" // class main-content

	// katex用の記述を最後に出力する
	output += `<script>renderMathInElement(document.body,{delimiters: [{left: "$$", right: "$$", display: true},{left: "$", right: "$", display: false}]});</script>
	</body>
	</html>`

	output = delEmptyLine(output)

	// スライド用htmlの1ページを出力する
	if err := ioutil.WriteFile(slidePath, []byte(output), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "File %s could not open : %v\n", slidePath, err)
		fmt.Println(err)
		panic(err)
	}
} // }}}
