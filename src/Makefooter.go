package main

func Makefooter() string { // {{{
	footer := `<script>renderMathInElement(document.body,{delimiters: [{left: "$$", right: "$$", display: true},{left: "$", right: "$", display: false}]});</script>
</body>
</html>`
	return footer
} // }}}
