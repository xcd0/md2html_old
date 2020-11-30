package main

import (
	"regexp"
	"strings"

	"gopkg.in/russross/blackfriday.v2"
)

func filter2body(in string) string { // {{{

	lines := strings.Split(in, "\n")

	// 独自置換対象文字列
	var rep [][]string
	rep = append(rep, []string{`^===$`, "\n<div style='page-break-before:always'></div>\n"})

	output := ""
	for _, r := range rep { // 全ての置換対象文字列について回す
		reg := regexp.MustCompile(r[0])
		for _, line := range lines { // 一行ずつ
			output += reg.ReplaceAllString(line, r[1]) + "\n"
		}
	}
	// 上書きする
	return output

} // }}}

func Makebody(mdpath string, rImgPath []string, t string) string { // {{{1

	stringmd := ReadMd(mdpath)

	stringmd = ReplaceImg4mdPre(rImgPath, stringmd)

	// 生成に使うライブラリに合わせて生成する
	//body := string(bluemonday.UGCPolicy().SanitizeBytes(blackfriday.MarkdownBasic([]byte(stringmd))))
	//body := string(blackfriday.MarkdownBasic([]byte(stringmd)))
	//body := string(bluemonday.UGCPolicy().SanitizeBytes(blackfriday.Run([]byte(stringmd))))

	//bytebody, _ := shurcooL_GFM([]byte(stringmd))
	//bytebody := string(blackfriday.MarkdownBasic([]byte(stringmd)))

	commonHtmlFlags := 0 |
		blackfriday.UseXHTML |
		blackfriday.Smartypants |
		blackfriday.SmartypantsFractions |
		blackfriday.SmartypantsDashes |
		blackfriday.SmartypantsLatexDashes

	extensions := 0 |
		blackfriday.NoIntraEmphasis |
		blackfriday.Tables |
		blackfriday.FencedCode |
		blackfriday.Autolink |
		blackfriday.Strikethrough |
		blackfriday.AutoHeadingIDs |
		blackfriday.HeadingIDs |
		blackfriday.BackslashLineBreak |
		blackfriday.DefinitionLists

	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: commonHtmlFlags,
	})
	opt := []blackfriday.Option{
		blackfriday.WithRenderer(renderer),
		blackfriday.WithExtensions(extensions),
	}
	bytebody := string(blackfriday.Run([]byte(stringmd), opt...))
	body := string(bytebody)

	// 独自記法の置換
	body = filter2body(body)

	return body

	/*
		readMd(fi)
		bytemd := []byte(fi.Md)

		// 生成に使うライブラリに合わせて生成する
		var bytebody []byte
		if fi.Flavor == "github" {
			bytebody, _ = gitHubAPI(bytemd)
		} else if fi.Flavor == "gfm" {
			bytebody, _ = shurcooL_GFM(bytemd)
		} else {
			bytebody = blackfriday.MarkdownBasic(bytemd)
		}

		body := string(bytebody)

		return body
	*/

} //}}}1
