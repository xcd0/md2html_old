package cssMinify

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
	fname, err := filepath.Abs(inputFilePath)
	dname := filepath.Dir(fname)
	//log.Println(fname)
	//log.Println(dname)
	outputCssPath := filepath.Join(dname, "markdown.mini.css")
	//log.Println(outputCssPath)

	// css を開く
	inputCssFp, err := os.Open(fname)
	defer inputCssFp.Close()
	if err != nil {
		log.Println(err)
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
