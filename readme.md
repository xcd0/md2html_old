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


## 使用方法

1. マークダウンのファイル(`*.md`)を`build.sh`のあるディレクトリに置く。  
画像なども同様に配置する。

1. `build.sh`を実行する。
	* ただのShellscriptなので、LinuxやMacならそのまま実行すればよい。  
	ただしHugoのインストールが必要。  
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
![](./build_on_win10.gif)

## 使用条件

* Shellscriptが動くこと。batではない。
	* windowsではWSLや[MSYS2 ( https://www.msys2.org/ )](https://www.msys2.org/)、
	[Git for Windows](https://gitforwindows.org/)などの端末環境をインストールする。
		* 拡張子`.sh`を開くプログラムとしてインストールした端末実行環境を指定するとよい。

* [hugo](https://gohugo.io/)が動くこと。
	* windows用のバイナリは同梱している。
	* ほかのOSでは[Quick Start](https://gohugo.io/getting-started/quick-start/)を参考に導入する。
	* macで[Homebrew](https://brew.sh/index_ja)が動くなら`brew install hugo`で終わり。
	* LinuxやWSLだと[Linuxbrew](https://docs.brew.sh/Homebrew-on-Linux)がおすすめ。`brew install hugo`で終わり。

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

## FAQ

* ~~build.shと同じディレクトリに複数のマークダウンがある場合エラーになる。~~
	* ちゃんと全部htmlを生成するように修正した。

* `_`がhugo によって誤変換される。
	* これはマークダウンの斜体記法(文字列の前後を`_`で囲む)が原因。
	* マークダウン上でアンダースコアの前にバックスラッシュをつければ回避できる。  
	`例）\_(┐「ε:)\_ → _(┐「ε:)_`
		* 取り合えずざっくりとした置換処理を実装した。つまり斜体にならない。  
		`例）_(┐「ε:)_ → _(┐「ε:)_`
		* この対応は正しくはないと思う。
		* しかし、斜体の使用率とアンダースコアの使用率とを比較すると斜体をほぼ使わないのでとりあえずこうする。
		* それでも誤変換は残る。

* htmlが文字化けした
	* mdファイルの文字コードが`Shift-jis`だと文字化けする。
	* `UTF-8`で保存する。

## 変更履歴

|改定日		|版		|内容					|
|:---:|:---:|:---:|
|2019/08/28	|1		|新規作成				|
|2019/09/26	|		|Makefileによるbase64.goのクロスコンパイル環境を作成|
|〃	|		|ログを出力するように変更|
|〃	|2		|いろいろ追記|
