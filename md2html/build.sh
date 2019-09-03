#!/usr/bin/env bash


SCRIPT_DIR=$(cd $(dirname $0); pwd)
echo $SCRIPT_DIR
cd $SCRIPT_DIR
#pwd

# ./*.mdがあるか、1つでもあればいい
ls *.md > /dev/null 2>&1
#echo $? # 0なら1つ以上ある

if [ $? != 0 ]; then
	echo ".mdのファイルが存在しません。"
	echo "build.shと同じディレクトリに拡張子.mdのファイルを置いてください。"
	echo "終了します。"
	sleep 3
	exit 1
fi

# mdファイルの名前をリストにする

md=()
while read mdfile
do
	mdtmp1="${mdfile%.*}" # 拡張子削除
	mdtmp2="${mdtmp1#./}"  # 行頭の./を削除
	#mdtmp3=`echo -n $mdtmp2 | tr ' ' '_'`
	#md+=("$mdtmp3")   # 半角空白をアンダースコアに変換
	mdtmp3=`echo -n $mdtmp2 | tr ' ' '-'`
	md+=("$mdtmp3")   # 半角空白をハイフンに変換
	
	cp $mdfile $mdtmp3
done < <(find . -maxdepth 1 -name '*.md')
#echo "${md[@]}"



# contentの中身をクリア
mkdir create/content > /dev/null 2>&1
# contentにmdファイルを移動
cp -rf $SCRIPT_DIR/*.md create/content/ > /dev/null 2>&1
#cp -rf *.png create/content/ > /dev/null 2>&1
#cp -rf *.jpg create/content/ > /dev/null 2>&1


# htmlファイルを生成
cd $SCRIPT_DIR/create
	#which hugo
	#echo $?
	flag=`which hugo > /dev/null 2>&1; echo $?`
	if [ $flag -eq 1 ]; then
		bin=$SCRIPT_DIR/create/hugo.exe
	else
		bin=hugo
	fi
	echo $bin
	$bin > /dev/null
	# hugoはcontentの中にあるmdファイルから
	# publicディレクトリにあるmd名のディレクトリの中にindex.htmlを生成する
	# このindex.htmlを親ディレクトリ名、つまり元のmdファイルの名前に変更して
	# 元のmdのあるディレクトリに移動させる
	
	# 生成されたindex.htmlファイルの名前を書き換え、
	# mdファイルと同じディレクトリに持っていく

	for file in "${md[@]}"; do
		echo "cp $SCRIPT_DIR/create/public/$file/index.html $SCRIPT_DIR/${file}.html"
		cp $SCRIPT_DIR/create/public/$file/index.html $SCRIPT_DIR/${file}.html
	done

echo "htmlファイルを生成しました。"


# 中間ファイルをクリア
echo "中間ファイルをクリアします"
rm -rf $SCRIPT_DIR/create/content \
	$SCRIPT_DIR/create/resources \
	$SCRIPT_DIR/create/public \
	$SCRIPT_DIR/create/data \
	$SCRIPT_DIR/create/static \
	$SCRIPT_DIR/create/layouts \
	$SCRIPT_DIR/create/archetypes > /dev/null 2>&1

# _がhugoによって誤変換されてしまう場合がかなりあるので、
# とりあえず雑に置換する

echo "アンダースコアの誤変換を置換します"
for file in $md; do
	sed s@\&rsquo\;\<em\>@_@g ${file}.html > ${file}_1.html
	sed s@\<em\>@_@g ${file}_1.html > ${file}_2.html
	sed s@\</em\>@_@g ${file}_2.html > ${file}_3.html
	sed s@\&rsquo\;@_@g ${file}_3.html > ${file}_4.html
	sed s@__@_@g ${file}_4.html > ${file}_5.html
	cat ${file}_5.html > ${file}.html
	rm ${file}_*.html > /dev/null 2>&1
done

echo "画像をbase64に変換してhtmlに埋め込みます"

for html_ext in `\find . -maxdepth 1 -name '*.html'`; do
	html=${html_ext%.*}
	img=()
	base=()
	# png jpg gifファイルの名前をリストにする
	for file in `\find . -maxdepth 2 -name '*.jpg'`; do
		img+=($file)
		base+=("data:image/jpeg;base64,"`./create/base64.exe $file`)
	done
	for file in `\find . -maxdepth 2 -name '*.png'`; do
		img+=($file)
		base+=("data:image/png;base64,"`./create/base64.exe $file`)
	done
	for file in `\find . -maxdepth 2 -name '*.gif'`; do
		img+=($file)
		base+=("data:image/gif;base64,"`./create/base64.exe $file`)
	done
	
	imax=`expr ${#img[@]} - 1`
	
	cat ${html}.html > ${html}_tmp.html
	for i in `seq 0 1 $imax`
	do
		echo ${img[$i]}
		cat ${html}_tmp.html | sed "s@${img[$i]}@${base[$i]}@g" > ${html}_tmp_.html
		cat ${html}_tmp_.html > ${html}_tmp.html
	done
	mv ${html}_tmp.html ${html}.html
	rm -rf ${html}_tmp*.html

done

echo "終了します。"
sleep 2
exit
eep 2
exit

