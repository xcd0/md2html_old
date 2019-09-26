#!/usr/bin/env bash

SCRIPT_DIR=$(cd $(dirname $0); pwd)
cd $SCRIPT_DIR
rm -rf log.txt error_log.txt

echo "--------------------------------------------------------------------------------" | tee -a log.txt
echo "Current directory : $SCRIPT_DIR" | tee -a log.txt

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

echo "--------------------------------------------------------------------------------" | tee -a log.txt

# contentの中身をクリア
mkdir create/content > /dev/null 2>&1
# contentにmdファイルを移動
cp -rf $SCRIPT_DIR/*.md create/content/ > /dev/null 2>&1
#cp -rf *.png create/content/ > /dev/null 2>&1
#cp -rf *.jpg create/content/ > /dev/null 2>&1

# hugo があるかどうか調べる
flag=`which hugo > /dev/null 2>&1; echo $?`
if [ $flag -eq 0 ]; then
	# hugo がインストールされている
	hugo=`which hugo`
else
	# hugo がインストールされていない
	# 基本的にエラー終了
	# ただしOSがWindowsなら同梱されているexeを使う
	if [ "$COMSPEC" != "" ]; then
		# Windowsなので同梱されているexeを使う
		echo "OS check : Windows" | tee -a log.txt
		hugo="$SCRIPT_DIR/create/hugo.exe"
		
		# exeを同梱しているつもりだが、一応あるかどうかチェック
		if [ -e $bin ]; then
			echo "Hugoがインストールされていない為、同梱しているhugo.exeを使用します。" | tee -a log.txt
		else
			echo "Hugoが見つかりません。エラーです。" 1>&2
			echo "Hugoをインストールしてください。" 1>&2
			echo "終了するにはEnterキーを押してください。" 1>&2
			read
			# エラーログを出力して終了
			cat << EOS > $SCRIPT_DIR/error_log.txt
Hugoが見つかりません。
Hugoをインストールしてください。
EOS
			exit 1
		fi
	else
		# Windowsじゃないのでエラー終了
		echo "Hugoが見つかりません。エラーです。" 1>&2
		echo "Hugoをインストールしてください。" 1>&2
		echo "終了するにはEnterキーを押してください。" 1>&2
		read
		exit 1
	fi
fi
echo "hugo path : $hugo" | tee -a log.txt

echo "--------------------------------------------------------------------------------" | tee -a log.txt

# htmlファイルを生成
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
			exit
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
	cd ..

echo "htmlファイルを生成しました。" | tee -a log.txt

# 中間ファイルをクリア
echo "中間ファイルをクリアします" | tee -a log.txt
rm -rf $SCRIPT_DIR/create/content \
	$SCRIPT_DIR/create/resources \
	$SCRIPT_DIR/create/public \
	$SCRIPT_DIR/create/data \
	$SCRIPT_DIR/create/static \
	$SCRIPT_DIR/create/layouts \
	$SCRIPT_DIR/create/archetypes > /dev/null 2>&1

# _がhugoによって誤変換されてしまう場合がかなりあるので、
# とりあえず雑に置換する

echo "アンダースコアの誤変換を置換します" | tee -a log.txt
for file in $md; do
	sed s@\&rsquo\;\<em\>@_@g ${file}.html > ${file}_1.html
	sed s@\<em\>@_@g ${file}_1.html > ${file}_2.html
	sed s@\</em\>@_@g ${file}_2.html > ${file}_3.html
	sed s@\&rsquo\;@_@g ${file}_3.html > ${file}_4.html
	sed s@__@_@g ${file}_4.html > ${file}_5.html
	cat ${file}_5.html > ${file}.html
	rm ${file}_*.html > /dev/null 2>&1
done

echo "--------------------------------------------------------------------------------" | tee -a log.txt

echo "画像をbase64に変換してhtmlに埋め込みます" | tee -a log.txt

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

echo "--------------------------------------------------------------------------------" | tee -a log.txt
echo "マークダウンからhtmlファイルの生成処理が完了しました。" | tee -a log.txt
echo "終了します。" | tee -a log.txt
exit


