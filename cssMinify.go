package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
)

func Minify(inputFilePath string) string {
	flag.Parse()
	// 第一引数にマークダウンのファイルのパスを受け取る
	// 引数を元に構造体を作る
	fname, _ := filepath.Abs(inputFilePath)
	dname := filepath.Dir(fname)
	//log.Println(fname)
	//log.Println(dname)
	outputCssPath := filepath.Join(dname, "markdown.mini.css")
	//log.Println(outputCssPath)

	// css を開く
	_, err := os.Stat(fname)
	if err != nil {
		// cssファイルがない
		// デフォルトのCSSを使う
		// minifyしない
		log.Println("error : do not exist css file")
		log.Println(fname)
		log.Println(err)
		return "error : do not exist css file"
	}

	inputCssFp, err := os.Open(fname)
	defer inputCssFp.Close()
	if err != nil {
		panic(err)
	}

	outputCssFp, err := os.Create(outputCssPath)
	if err != nil {
		panic(err)
	}
	defer outputCssFp.Close()

	outputWriter := bufio.NewWriter(outputCssFp)

	mediatype := "text/css"
	m := minify.New()
	m.AddFunc(mediatype, css.Minify)

	if err := m.Minify(mediatype, outputWriter, inputCssFp); err != nil {
		panic(err)
	}

	return outputCssPath
}
