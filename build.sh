#!/bin/bash

if [ ! -e *.md ]; then
	echo ".mdのファイルが存在しません。"
	echo "build.shと同じディレクトリに拡張子.mdのファイルを置いてください。"
	echo "終了します。"
	sleep 3
	exit 1
fi

# mdファイルの名前をリストにする
for file in *.md; do
	md+=${file%.*}
done

# contentの中身をクリア
mkdir create/content
# contentにmdファイルを移動
cp -rf *.md create/content/

# htmlファイルを生成
cd create
	./hugo > /dev/null
	# publicディレクトリにある$mdディレクトリ内にindex.htmlが生成される
	
	# 生成されたindex.htmlファイルの名前を書き換え、
	# mdファイルと同じディレクトリに持っていく
	
	cd public
		for file in $md; do
			cd $file
				cp index.html ${file}.html
				cp	${file}.html ../../../
				cd ..
		done
		cd ..
	cd ..

# 中間ファイルをクリア
rm -rf create/content \
	create/resources \
	create/public \
	create/data \
	create/static \
	create/layouts \
	create/archetypes
