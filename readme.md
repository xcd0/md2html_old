# md2htmlについて

## 概要

markdownのファイルからhtmlを生成します。  

* markdown.cssファイルがあればminifyしてhtmlに埋め込みます。

* 序に画像ファイル(拡張子がjpg,png,gif)をbase64にエンコードして、  
テキストとしてhtmlに埋め込みます。

これらの処理により、markdownと画像のファイル群から  
html単体で画像やCSSを含んだ可搬性のある文書が生成できます。

例えば

* [readme.md](https://static.xcd0.com/2019/10/30/readme.md)
* [markdown.css](https://static.xcd0.com/2019/10/30/markdown.css)
* [build_on_win10.gif](https://static.xcd0.com/2019/10/30/build_on_win10.gif)

から

* [readme.html](https://static.xcd0.com/2019/10/30/readme.html)

を生成できます。

メインの処理はGolangで書いており、基本的にどのOSでもbashが動けば動作します。


## 使用方法

1. マークダウンの書式のテキスト文書を拡張子`.md`で保存する。  
マークダウンファイルと同じディレクトリにmarkdown.cssファイルを置く。  
`markdown.css`がない場合、デフォルトのcssが適応される。  
画像などはマークダウン内のパス記述に従って配置する。  
つまり同じディレクトリにある必要はない。

1. 実行ファイル `md2html` にマークダウンのファイルをドラッグアンドドロップして投げる。  
これでhtmlファイルが生成される。
	* このhtmlファイルには画像が埋め込まれているのでhtml単体で扱うことができる。

![](./build_on_win10.gif)


## 使用条件

* すべての使用条件を取り除いた。  
マウスがあってマークダウンファイルを実行ファイルにドラッグアンドドロップできれば良い。


## FAQ

* `_`がhugo によって誤変換される。
	* 誤変換ではない。
		* これはマークダウンの斜体記法(文字列の前後を`_`で囲む)が原因。
		* マークダウン上でアンダースコアの前にバックスラッシュをつければ回避できます。  
		例）マークダウン内記述`\_(┐「ε:)\_` → 出力`_(┐「ε:)_`

* htmlが文字化けした
	* mdファイルの文字コードが`Shift-jis`だと文字化けする。
	* `UTF-8`で保存する。

* 特殊なOS、CPUで使用するしたい  (32bit版やアーキテクチャがx86系じゃないやつ)
	* 環境に合わせてビルドしてください。
	* Golangでビルドできる環境なら動きます。
		* [https://golang.org/doc/install#requirements](https://golang.org/doc/install#requirements)を参照。
		* `If your OS or architecture is not on the list, you may be able to install from source or use gccgo instead.`だそうです
		* 2019/10/01時点の対応表

|Operating system	|	Architectures	|	Notes |
|---|---|---|
|FreeBSD 10.3 or later	amd64, 386	|	Debian GNU/kFreeBSD not supported |
|Linux 2.6.23 or later with glibc	|	amd64, 386, arm, arm64, s390x, ppc64le	|	CentOS/RHEL 5.x not supported.<br> Install from source for other libc. |
|macOS 10.10 or later	|	amd64	|	use the clang or gcc† that comes with Xcode‡ for cgo support |
|Windows 7, Server 2008R2 or later	|	amd64, 386	|	use MinGW (386) or MinGW-W64 (amd64) gcc†. No need for cygwin or msys. |


## 変更履歴

|改定日		|版		|内容					|
|:---:|:---:|:---:|
|2019/08/28	|1		|新規作成。				|
|2019/09/26	|2		|Makefileによるbase64.goのクロスコンパイル環境を作成。<br>ログを出力するように変更。<br>readme.mdにいろいろ追記した。|
|2019/10/01	|3		|hugoのバイナリを同梱するのをやめた<br>OSに合わせて自動でDLするようにした。|
|2019/10/30	|4		|shellscriptをやめた。<br>全部Golangで書きなおした。<br>これによりすべての使用条件を撤廃、実行時間が100倍はやくなった。|

