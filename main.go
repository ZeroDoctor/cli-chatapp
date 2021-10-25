package main

import (
	"github.com/awesome-gocui/gocui"
	"gitlab.com/smallwood/sw-chat/key"
	"gitlab.com/smallwood/sw-chat/view"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	g.Mouse = false
	g.Cursor = true
	g.Highlight = true
	g.SelFgColor = gocui.ColorCyan

	view.Handler().SetupViews([]string{"textbox", "screen", "header"})
	g.SetManagerFunc(view.Handler().Layout)
	key.SetBindings(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}

	view.Handler().Wait()
}
