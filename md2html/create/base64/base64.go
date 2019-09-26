package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	str := encode(flag.Arg(0))
	fmt.Print(str)
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
