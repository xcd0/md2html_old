package main

import (

	//"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
)

func MinifyJavascriptString(inputJs string) string {
	mediatype := "application/javascript"
	m := minify.New()
	m.AddFunc(mediatype, js.Minify)
	minified, _ := m.String(mediatype, inputJs)
	return minified
}

func MinifyCssString(inputCss string) string {
	mediatype := "text/css"
	m := minify.New()
	m.AddFunc(mediatype, css.Minify)
	minifiedCss, _ := m.String(mediatype, inputCss)
	return minifiedCss
}

func Minify(inputFilePath string) string {
	fname, _ := filepath.Abs(inputFilePath)
	//bname := inputFilePath[:len(inputFilePath)-len(filepath.Ext(inputFilePath))]
	//dname := filepath.Dir(fname)
	//log.Println(fname)
	//log.Println(dname)
	//outputCssName := bname + ".mini.css"
	//outputCssPath = filepath.Join(dname, outputCssName)
	//log.Println(outputCssName)
	//log.Println(outputCssPath)

	// css をがあるか調べる
	_, err := os.Stat(fname)
	if err != nil {
		// cssファイルがない
		// デフォルトのCSSを使う
		// minifyしない
		return "default"
	}
	// ファイル読み込み
	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}

	inputCss := string(bytes)

	return MinifyCssString(inputCss)
}
