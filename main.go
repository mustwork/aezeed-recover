package main

import (
	"embed"
	"github.com/mustwork/aezeed-recover/cmd"
	"github.com/mustwork/aezeed-recover/internal/util"
)

//go:embed internal/ui/INFO.md
var readme embed.FS

func main() {
	content, err := readme.ReadFile("internal/ui/INFO.md")
	if err != nil {
		panic(err)
	}
	info, err := util.ParseMarkdownFirstSectionAsTviewText(content)
	if err != nil {
		panic(err)
	}
	cmd.Execute(info)
}
