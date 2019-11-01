package main

import (
	"./convertFromMarkdownToHtml"
)

func main() {
	flag.Parse()

	arg := flag.Arg(0)
	fi := md2html.Fileinfo{}
	//md2html.Argparse(flag.Arg(0))

	// 絶対パスを得る
	fi.Apath, _ = filepath.Abs(arg)
	// ファイルパスをディレクトリパスとファイル名に分割する
	fi.Dpath, fi.Filename = filepath.Split(fi.Apath)
	// 拡張子を得る
	fi.Ext = filepath.Ext(fi.Filename)
	// 拡張子なしの名前を得る
	fi.Basename = fi.Filename[:len(fi.Filename)-len(fi.Ext)]
	// 出力するhtmlのパスを得る
	fi.Htmlpath = fi.Dpath + fi.Basename + ".html"
	// 入力Cssのパスを得る
	fi.Csspath = fi.Dpath + "markdown.css"

	// htmlを作成する

	fi.Flavor = "gfm"
	testMakeHtml(fi)

	fi.Flavor = "github"
	testMakeHtml(fi)

	fi.Flavor = "other"
	testMakeHtml(fi)
}

func testMakeHtml(fi md2html.Fileinfo) {

	header := Makeheader(fi)
	body, err := Makebody(fi)
	footer := Makefooter()

	fi.Html = header + body + footer

	fi.Htmlpath = fi.Dpath + fi.Basename + "_" + fi.flavor + ".html"
	headerpath = fi.Dpath + fi.Basename + "_" + fi.flavor + "_header.txt"
	bodypath = fi.Dpath + fi.Basename + "_" + fi.flavor + "_body.txt"
	footerpath = fi.Dpath + fi.Basename + "_" + fi.flavor + "_footer.txt"
	html := fi.Html

	if err != nil {
		fmt.Println(err)
		return
	}

	// 出力する
	err = ioutil.WriteFile(fi.Htmlpath, []byte(html), 0644)
	if err != nil {
		errorFileOpen(err)
		return
	}
	err = ioutil.WriteFile(headerpath, []byte(header), 0644)
	if err != nil {
		errorFileOpen(err)
		return
	}
	err = ioutil.WriteFile(bodypath, []byte(body), 0644)
	if err != nil {
		errorFileOpen(err)
		return
	}
	err = ioutil.WriteFile(footerpath, []byte(footer), 0644)
	if err != nil {
		// Openエラー処理
		errorFileOpen(err)
		return
	}
}

func errorFileOpen(err error) {
	fmt.Fprintf(os.Stderr, "File %s could not open : %v\n", fi.Htmlpath, err)
	fmt.Println(err)
}
