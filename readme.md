# md2htmlについて

## 概要

markdownのファイルからhtmlを生成します。  
序に画像ファイル(拡張子がjpg,png,gif)をbase64にエンコードして、  
テキストとしてhtmlに埋め込みます。

これにより、markdownと画像のファイル群から  
html単体の可搬性のある文書が生成できます。

## 使用方法

1. マークダウンのファイル(*.md)をmd2htmlに置く。
1. build.shを実行する。
	* windowsではmsysやwslから実行する必要がある。
	* msysがはいってるならショートカット置いてるのでbuild.shをD&Dするだけでよい。ただしmsysは`C:\msys64\mingw64.exe`に実行ファイルがあること。
1. htmlファイルが生成される。
	* このhtmlファイルには画像が埋め込まれているのでhtml単体で扱うことができる。

gif画像  
![](./build.gif)

## 既知の不具合

* ~~build.shと同じディレクトリに複数のマークダウンがある場合エラーになる。~~	← 修正した。

* `_`などのいくつかの文字がhugo によって誤変換される。←取り合えずざっくりとした置換処理を実装。それでも誤変換は残る。


## 使用条件

* shellscriptが動くこと。batではない。
	* windowsでは[msys2 ( https://www.msys2.org/ )](https://www.msys2.org/)などの端末環境をインストールする。
* [hugo](https://gohugo.io/)が動くこと。
	* windows用のバイナリは同梱している。
	* ほかのOSでは[Quick Start](https://gohugo.io/getting-started/quick-start/)を参考に導入する。
	* macでbrewが動くなら`brew install hugo`で終わり。
* 何かおかしい場合実行ファイルである`build.sh`や
`base64.go`を修正するやる気があること_(:3 」∠ )_。  
`base64.go`は画像をbase64にエンコードするプログラムのソースコードである。ささっとクロスコンパイルできるgolangで書いた。

## 変更履歴

|改定日		|版		|内容					|
|:---:|:---:|:---:|
|2019/08/28	|1		|新規作成				|