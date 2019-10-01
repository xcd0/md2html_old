#!/usr/bin/env bash


SCRIPT_DIR=$(cd $(dirname $0); pwd)
cd $SCRIPT_DIR
rm -rf log.txt error_log.txt

echo "--------------------------------------------------------------------------------" | tee -a log.txt
echo "Current directory : $SCRIPT_DIR" | tee -a log.txt
echo "--------------------------------------------------------------------------------" | tee -a log.txt

# contentの中身をクリア
mkdir create/content > /dev/null 2>&1
# contentにmdファイルを移動
cp -rf $SCRIPT_DIR/*.md create/content/ > /dev/null 2>&1

function search_markdown () { # {{{
	# ./*.mdがあるか、1つでもあればいい
	ls *.md > /dev/null 2>&1
	#echo $? # 0なら1つ以上ある

	if [ $? != 0 ]; then
		echo ".mdのファイルが存在しません。" 1>&2
		echo "build.shと同じディレクトリに拡張子.mdのファイルを置いてください。" 1>&2
		echo "終了するにはEnterキーを押してください。" 1>&2
		read
		# エラーログを出力して終了
		cat << EOS > $SCRIPT_DIR/error_log.txt
.mdのファイルが存在しません。
build.shと同じディレクトリに拡張子.mdのファイルを置いてください。
EOS
		exit 1
	fi

	# mdファイルの名前をリストにする

	md=()
	md_ori=()
	while read mdfile
	do
		md_ori+=("$mdfile")
		mdtmp1="${mdfile%.*}" # 拡張子削除
		mdtmp2="${mdtmp1#./}"  # 行頭の./を削除
		
		# 名前に半角空白を含むmdファイルの名前を書き換える
		## "a b.md" -> "a-b"
		mdtmp3=`echo -n $mdtmp2 | tr ' ' '-'`
		md+=("$mdtmp3")   # 半角空白をハイフンに変換
		
	done < <(find . -maxdepth 1 -name '*.md')

	echo "処理するマークダウンファイルは以下のファイルです。" | tee -a log.txt
	echo "${md_ori[@]}" | tee -a log.txt

} # }}}

search_markdown

echo "--------------------------------------------------------------------------------" | tee -a log.txt

# dl_hugo search_hugoから呼ばれる
function dl_hugo () { # {{{
	# hugoをダウンロードする
	cd $SCRIPT_DIR/create
	# OS毎にダウンロードするバイナリが違うので分岐
	if [ "$COMSPEC" != "" ]; then
		echo "Windows用の" | tee -a log.txt
		echo "Hugo 0.5.8.3をダウンロードします。(約12MB)" | tee -a log.txt
		wget https://github.com/gohugoio/hugo/releases/download/v0.58.3/hugo_0.58.3_Windows-64bit.zip
		unzip -j hugo_0.58.3_Windows-64bit.zip hugo.exe
		rm hugo_0.58.3_Windows-64bit.zip
		## winだけzipなのめんどい...
		## Windowsだけ拡張子がついているが 面倒なので $SCRIPT_DIR/create/hugo にシンボリックリンクを張る
		## これならユーザーがhugo.exeをどこからか取ってきて置いた場合も動作するはず
		#ln -sf hugo.exe hugo
	elif [ "$(uname)" == 'Darwin' ]; then
		echo -n "macOS用の" | tee -a log.txt
		echo "Hugo 0.5.8.3をダウンロードします。(約12MB)" | tee -a log.txt
		wget -O - https://github.com/gohugoio/hugo/releases/download/v0.58.3/hugo_0.58.3_macOS-64bit.tar.gz | tar xzvf - hugo
	elif [ "$(expr substr $(uname -s) 1 5)" == 'Linux' ]; then
		echo "Linux用の" | tee -a log.txt
		echo "Hugo 0.5.8.3をダウンロードします。(約12MB)" | tee -a log.txt
		wget -O - https://github.com/gohugoio/hugo/releases/download/v0.58.3/hugo_0.58.3_Linux-64bit.tar.gz | tar xzvf - hugo
	else
		# Windows/mac/linuxじゃないのでエラー ユーザーにhugoをインストールするよう促す
		echo "Hugoが見つかりません。エラーです。" 1>&2
		echo "Hugoをインストールしてください。" 1>&2
		echo "最新版のバイナリは以下のURLからダウンロードできます。"
		echo "https://github.com/gohugoio/hugo/releases"
		echo "終了するにはEnterキーを押してください。" 1>&2
		read
		exit 1
	fi
	echo "最新版のバイナリは以下のURLからダウンロードできます。" | tee -a log.txt
	echo "https://github.com/gohugoio/hugo/releases" | tee -a log.txt
	echo "最新版を取得した場合PATHを通すか $hugo においてください。" | tee -a log.txt
} # }}}

# hugo があるかどうか調べる
function search_hugo () { #{{{
	which hugo
	flag=$?
	#if [ $flag -eq 0 ]; then
	if [ $flag -ne 0 ]; then
		# hugo がインストールされている
		hugo=`which hugo`
		$hugo version | tee -a log.txt
	else
		# hugo がインストールされていない
		# DLされているか調べる
		#if [ -e $SCRIPT_DIR/create/hugo.exe ]; then
		#	# Windowsだけめんどい
		#	ln -sf hugo.exe hugo
		#	read
		#fi

		hugo="$SCRIPT_DIR/create/hugo"

		if [ -e $hugo ]; then
			# Hugoのバイナリが置いてある
			# 正しく動くバイナリかどうかたたいてみる
			$hugo version > /dev/null 2>&1
			check=`$?`
			if [ $check -eq 0 ]; then
				# DLされているHugoのバイナリを使う
				$hugo version | tee -a log.txt
				# hugoダウンロード済み
				echo "Hugoがダウンロードされています。" | tee -a log.txt
				echo "最新版のバイナリは以下のURLからダウンロードできます。" | tee -a log.txt
				echo "https://github.com/gohugoio/hugo/releases" | tee -a log.txt
				echo "最新版を取得した場合PATHを通すか $hugo においてください。" | tee -a log.txt
			else
				# DLされているHugoのバイナリが...
				# 削除する
				echo "$SCRIPT_DIR/create/にあるHugoのバイナリが正しくありません" | tee -a log.txt
				if [ -e $SCRIPT_DIR/create/hugo.exe ]; then
					rm $SCRIPT_DIR/create/hugo.exe | tee -a log.txt
				fi
				rm $SCRIPT_DIR/create/hugo | tee -a log.txt
				echo "システムにHugoがインストールされていません。" | tee -a log.txt
				dl_hugo
			fi
		else
			echo "システムにHugoがインストールされていません。" | tee -a log.txt
			dl_hugo
		fi
	fi
} # }}}

search_hugo

echo "hugo path : $hugo" | tee -a log.txt

echo "--------------------------------------------------------------------------------" | tee -a log.txt

# markdownファイルからhtmlファイルを生成
function create_html () { # {{{

	cd $SCRIPT_DIR/create

	# hugoの実行
	echo $hugo
	out=`$hugo`
	echo "$out" >> $SCRIPT_DIR/log.txt
	
	if [ `echo $?` -ne 0 ]; then
		# hugo 変換失敗
		echo "Hugoでのmarkdown -> htmlの変換に失敗しました。" 1>&2
		echo $out
		echo "終了するにはEnterキーを押してください。" 1>&2
		# エラーログを出力して終了
		cat << EOS > $SCRIPT_DIR/error_log.txt
Hugoでのmarkdown -> htmlの変換に失敗しました。
$e
EOS
		exit 1
	fi
	echo "`$hugo`"
	
	# 実行後の処理
	# 生成されたindex.htmlファイルの名前を書き換え、
	# mdファイルと同じディレクトリに持っていく
	## hugoはcontentの中にあるmdファイルから
	## publicディレクトリにあるmd名のディレクトリの中にindex.htmlを生成する
	## このindex.htmlを親ディレクトリ名、つまり元のmdファイルの名前に変更して
	## 元のmdのあるディレクトリに移動させる
	
	for file in "${md[@]}"; do
		#echo "cp $SCRIPT_DIR/create/public/$file/index.html $SCRIPT_DIR/${file}.html"
		cp $SCRIPT_DIR/create/public/$file/index.html $SCRIPT_DIR/${file}.html
	done

	# 中間ファイルをクリア
	#echo "中間ファイルをクリアします" | tee -a log.txt
	rm -rf $SCRIPT_DIR/create/content \
		$SCRIPT_DIR/create/resources \
		$SCRIPT_DIR/create/public \
		$SCRIPT_DIR/create/data \
		$SCRIPT_DIR/create/static \
		$SCRIPT_DIR/create/layouts \
		$SCRIPT_DIR/create/archetypes > /dev/null 2>&1


} # }}}

echo "htmlファイルを生成します" | tee -a log.txt

create_html


echo "--------------------------------------------------------------------------------" | tee -a log.txt

# _がhugoによって誤変換されてしまう場合がかなりあるので、
# とりあえず雑に置換する

# そもそも _斜体_ という記法がある
# ただまあ斜体使わないので...


function convert_underscore () { # {{{

	cd $SCRIPT_DIR
	for file in $md; do
		sed s@\&rsquo\;\<em\>@_@g ${file}.html > ${file}_1.html
		sed s@\<em\>@_@g ${file}_1.html > ${file}_2.html
		sed s@\</em\>@_@g ${file}_2.html > ${file}_3.html
		sed s@\&rsquo\;@_@g ${file}_3.html > ${file}_4.html
		sed s@__@_@g ${file}_4.html > ${file}_5.html
		cat ${file}_5.html > ${file}.html
		rm ${file}_*.html > /dev/null 2>&1
	done

} # }}}

echo "アンダースコアの誤変換を置換します" | tee -a log.txt

# そのままがいいならコメントアウトする
convert_underscore

echo "--------------------------------------------------------------------------------" | tee -a log.txt

function embed_image () { # {{{

	for html_ext in `\find . -maxdepth 1 -name '*.html'`; do
		html=${html_ext%.*}
		img=()
		base=()
		if [ "$COMSPEC" != "" ]; then
			base64=./create/base64/base64_win
		elif [ "$(uname)" == "Darwin" ]; then
			base64=./create/base64/base64_mac
		elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
			base64=./create/base64/base64_linux
		fi
		# png jpg gifファイルの名前をリストにする
		for file in `\find . -maxdepth 2 -name '*.jpg'`; do
			img+=($file)
			base+=("data:image/jpeg;base64,"`$base64 $file`)
		done
		for file in `\find . -maxdepth 2 -name '*.png'`; do
			img+=($file)
			base+=("data:image/png;base64,"`$base64 $file`)
		done
		for file in `\find . -maxdepth 2 -name '*.gif'`; do
			img+=($file)
			base+=("data:image/gif;base64,"`$base64 $file`)
		done
		
		imax=`expr ${#img[@]} - 1`
		
		cat ${html}.html > ${html}_tmp.html
		for i in `seq 0 1 $imax`
		do
			echo ${img[$i]} | tee -a log.txt
			cat ${html}_tmp.html | sed "s@${img[$i]}@${base[$i]}@g" > ${html}_tmp_.html
			cat ${html}_tmp_.html > ${html}_tmp.html
		done
		mv ${html}_tmp.html ${html}.html
		rm -rf ${html}_tmp*.html
	done

} # }}}

echo "画像をbase64に変換してhtmlに埋め込みます" | tee -a log.txt

embed_image

echo "--------------------------------------------------------------------------------" | tee -a log.txt
echo "マークダウンからhtmlファイルの生成処理が完了しました。" | tee -a log.txt
echo "終了します。" | tee -a log.txt
exit


