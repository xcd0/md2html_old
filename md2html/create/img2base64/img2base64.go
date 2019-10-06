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

	// これが出力される
	output := ""

	scanner := bufio.NewScanner(f)
	// 1行がめっちゃ長い時ようにbufferを大きく取っている
	scanner.Buffer([]byte{}, math.MaxInt64)

	for scanner.Scan() {
		// ここで一行ずつ処理
		line := scanner.Text()
		position := strings.Index(line, flag.Arg(1))
		// flag.Arg(1)の文字列が含まれているかどうか調べる
		if position == -1 {
			output += line + "\n"
		} else {
			// flag.Arg(1)をbase64でエンコードしたデータに置き換える
			// base64でエンコードする
			code := flag.Arg(1)
			base64code := encode(code)
			// ファイルの拡張子ごとにヘッダをつける
			// gif,png,jpg,jpeg以外は元のファイル名になるようにしている
			ext := filepath.Ext(flag.Arg(1))
			if ext == ".gif" {
				code = "data:image/gif;base64," + base64code
			} else if ext == ".png" {
				code = "data:image/png;base64," + base64code
			} else if ext == ".jpg" || ext == ".jpeg" {
				code = "data:image/jpeg;base64," + base64code
			}
			// 前後を切り出す
			pre := line[:position]
			post := line[position+len(flag.Arg(1)):]

			// くっつけて上書きする
			line = pre + code + post
			// 出力する
			output += line + "\n"
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
