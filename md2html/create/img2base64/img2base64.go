package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	/*
		1: 入力テキストファイル
		2: デコードさあれるファイルパス
	*/
	flag.Parse()
	fmt.Fprintf(os.Stderr, "a0 %v\n", flag.Arg(0))
	fmt.Fprintf(os.Stderr, "a1 %v\n", flag.Arg(1))

	// ファイルをOpenする
	f, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "File %s could not read: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// <img src="./build_on_win10.gif" alt="" /></p>
	output := ""

	scanner := bufio.NewScanner(f)
	scanner.Buffer([]byte{}, math.MaxInt64)

	input1 := "<img src="
	input2 := flag.Arg(1)
	for scanner.Scan() {
		// ここで一行ずつ処理
		line := scanner.Text()
		b1 := strings.Index(line, input1)
		b2 := strings.Index(line, input2)
		if b1 == -1 {
			output += line + "\n"
		} else if b2 == -1 {
			output += line + "\n"
		} else {
			base64code := encode(flag.Arg(1))
			ext := filepath.Ext(flag.Arg(1))
			if ext == ".gif" {
				base64codetag := "<center><img src=\"data:image/gif;base64," + base64code + "\"></center>"
				output += base64codetag + "\n"
			} else if ext == ".png" {
				base64codetag := "<center><img src=\"data:image/png;base64," + base64code + "\"></center>"
				output += base64codetag + "\n"
			} else if ext == ".jpg" || ext == ".jpeg" {
				base64codetag := "<center><img src=\"data:image/jpeg;base64," + base64code + "\"></center>"
				output += base64codetag + "\n"
			} else {
				output += line + "\n"
			}
		}
	}
	if err = scanner.Err(); err != nil {
		// エラー処理
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(output)
}

//エンコード
func encode(str string) string {

	file, _ := os.Open(str)
	defer file.Close()

	fi, _ := file.Stat() //FileInfo interface
	size := fi.Size()    //ファイルサイズ

	data := make([]byte, size)
	file.Read(data)

	return base64.StdEncoding.EncodeToString(data)
}
