#!/bin/bash

# gobump showで{"version":"5.3.1"} みたいなのが返ってくる
# コミットハッシュと.でつなげてそれっぽくする
version=`gobump show | gojq -r .version`.`git rev-parse --short HEAD`
# jqで取り出す
major=`echo $version | awk -F. '{ print $1 }'`
minor=`echo $version | awk -F. '{ print $2 }'`
patch=`echo $version | awk -F. '{ print $3 }'`
build=0

#echo verison : $version
#echo major   : $major
#echo minor   : $minor
#echo patch   : $patch
#echo build   : $build

# versioninfo.jsonを書き換える
cp versioninfo.json versioninfo_old.json
cat versioninfo_old.json \
| gojq ".FixedFileInfo.FileVersion.Major|=$major" \
| gojq ".FixedFileInfo.FileVersion.Minor|=$minor" \
| gojq ".FixedFileInfo.FileVersion.Patch|=$patch" \
| gojq ".FixedFileInfo.FileVersion.Build|=$build" \
| gojq ".FixedFileInfo.ProductVersion.Major|=$major" \
| gojq ".FixedFileInfo.ProductVersion.Minor|=$minor" \
| gojq ".FixedFileInfo.ProductVersion.Patch|=$patch" \
| gojq ".FixedFileInfo.ProductVersion.Build|=$build" \
| gojq ".StringFileInfo.FileVersion|=\"$version\"" \
| gojq ".StringFileInfo.ProductVersion|=\"$version\"" \
> versioninfo.json

