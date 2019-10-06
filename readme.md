# md2htmlについて

## 概要

markdownのファイルからhtmlを生成します。  
序に画像ファイル(拡張子がjpg,png,gif)をbase64にエンコードして、  
テキストとしてhtmlに埋め込みます。

これにより、markdownと画像のファイル群から  
html単体で画像を含む、つまり可搬性のある文書が生成できます。

例えば

* [readme.md](https://xcd0.com/static/20190926/readme.md)
* [build_on_win10.gif](https://xcd0.com/static/20190926/build.gif)

から

* [readme.html](https://xcd0.com/static/20190926/readme.html)

を生成できます。

メインの処理はGolangで書いており、単体で動作するバイナリを同梱しているので  
基本的にどのOSでもbashが動けば動作します。


## 使用方法

1. マークダウンのファイル(`*.md`)を`build.sh`のあるディレクトリに置く。  
画像なども同様に配置する。

1. `build.sh`を実行する。
	* ただのShellscriptなので、LinuxやMacならそのまま実行すればよい。  
	ただしHugoのインストールが必要。  
	[Hugoのバイナリ](https://github.com/gohugoio/hugo/releases)があればよい。
	詳細は下記の[使用条件](#使用条件)を参照。
	
	* WindowsではMSYSやWSLから実行する必要がある。  
	Git Bashとかでも大丈夫なはず。  
	WSLの場合、WSLの中のOSにHugoをインストールする必要がある。  
	詳細は下記の[使用条件](#使用条件)を参照。
	
	* MSYS2の64bit版がインストールされているなら  
	ショートカットを置いてるので`build.sh`をD&Dするだけでよい。  
	ただしmsysは`C:\msys64\mingw64.exe`に実行ファイルがあること。

1. htmlファイルが生成される。
	* このhtmlファイルには画像が埋め込まれているのでhtml単体で扱うことができる。

Windows10 MSYS2 での実行例 gif 

<center> ![](./build_on_win10.gif) </center>

## 使用条件

* Shellscriptが動くこと。batではない。
	* windowsではWSLや[MSYS2 ( https://www.msys2.org/ )](https://www.msys2.org/)、
	[Git for Windows](https://gitforwindows.org/)などの端末環境をインストールする。
		* 拡張子`.sh`を開くプログラムとしてインストールした端末実行環境を指定するとよい。

* [hugo](https://gohugo.io/)が動くこと。
	* ~~windows用のバイナリは同梱している。~~
	* ~~[Quick Start](https://gohugo.io/getting-started/quick-start/)を参考に導入してもよい。~~
	* ~~macで[Homebrew](https://brew.sh/index_ja)が動くなら`brew install hugo`で終わり。~~
	* ~~LinuxやWSLだと[Linuxbrew](https://docs.brew.sh/Homebrew-on-Linux)がおすすめ。`brew install hugo`で終わり。~~
	* システムにHugoをインストールしなくても勝手にバイナリをダウンロードするように変更した。(第3版)
	* [https://github.com/gohugoio/hugo/releases](https://github.com/gohugoio/hugo/releases)から最新版が手に入る。
		* 別に新しい必要はない。

* 生成するhtmlの外観をカスタマイズしたい場合、Hugoの知識が必要。  
すごく雑に言えば`md2html/create/themes/github-markdown/layouts`あたりを編集すればよい。  
htmlは`md2html/create/themes/github-markdown/layouts`以下の
	* `default/single.html`
	* `partials/markdown.mini.css`  
から生成される。
例えばCSSの編集だけでいい場合、  
`md2html/create/themes/github-markdown/layouts/partials/markdown.mini.css`  
を編集すればよい。  
ただしminifyしているので[CSSのコード整形ツール](https://lab.syncer.jp/Tool/CSS-PrettyPrint/)などを使用して見やすくするとよい。  
保存するときは再度minifyすることをお勧めする。

* 何かおかしい場合実行ファイルである`build.sh`や
`base64.go`を修正するやる気があること_(:3 」∠ )_。  
`base64.go`は画像をbase64にエンコードするプログラムのソースコードである。  
ささっとクロスコンパイルできるGolangで書いた。  
`Makefile`を置いているので、GoのコンパイラにPATHが通っていて、`make`ができればコンパイルできる。

## メモ

* ~~build.shと同じディレクトリに複数のマークダウンがある場合エラーになる。~~
	* ちゃんと全部htmlを生成するように修正した。

* `_`がhugo によって誤変換される。
	* 誤変換ではない可能性が高いです。
		* これはマークダウンの斜体記法(文字列の前後を`_`で囲む)が原因っぽいです。
		* マークダウン上でアンダースコアの前にバックスラッシュをつければ回避できます。  
		`例）\_(┐「ε:)\_ → _(┐「ε:)_`
	* 取り合えずざっくりとした置換処理を実装しました。  
	`例）_(┐「ε:)_ → _(┐「ε:)_`
		* この対応は正しくはないと思います。
		* しかし、斜体の使用率とアンダースコアの使用率とを比較すると  
		斜体をほぼ使わないのでとりあえずこれで行くことにしました。
		* それでも誤変換は残っています。
	* `./build.sh`の238行目付近`convert_underscore`をコメントアウトすれば  
	`_`の置換処理をスキップするので正式なmarkdown記法になります。(第3版)

* htmlが文字化けした
	* mdファイルの文字コードが`Shift-jis`だと文字化けする。
	* `UTF-8`で保存する。

* 特殊なOS、CPUで使用する場合(32bit版やアーキテクチャがx86系じゃないやつ)
	* 32bit版とかなら動きます。そのままでは動かないですです。それ以外の環境はまあ環境次第です。
	* markdown→htmlの変換に使用しているHugoのバイナリを
	[https://github.com/gohugoio/hugo/releases](https://github.com/gohugoio/hugo/releases)
	から取得してパスを通すか、md2html/create/以下においてください。
	* Golangで書いている md2html/create/base64/ 以下のバイナリを環境に合わせてビルドする必要があります。
		* この場合build.shの中でOS判定している部分において正しく動作しない可能性が高いので修正が必要になります。
			* 250行目付近(第3版)の処理の後に`base64=ビルドしたバイナリへのパス`を書けば問題ないはずです。
	* Golangが動けば動くはずです。動かないなら動作しません。(Hugoもbase64も同様)
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
|2019/08/28	|1		|新規作成				|
|2019/09/26	|		|Makefileによるbase64.goのクロスコンパイル環境を作成|
|〃	|		|ログを出力するように変更|
|〃	|2		|いろいろ追記|
|2019/10/01	|3		|hugoのバイナリを同梱するのをやめた<br>OSに合わせて自動でDLするようにした|

