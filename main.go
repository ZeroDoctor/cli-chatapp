package main

import (
	"errors"
	"time"

	"github.com/awesome-gocui/gocui"
	"gitlab.com/smallwood/sw-chat/channel"
	"gitlab.com/smallwood/sw-chat/key"
	"gitlab.com/smallwood/sw-chat/view"
)

func updateClock() {
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()

	for range tick.C {
		nowStr := time.Now().Format("02/01/2006 15:04:05")
		channel.HeaderChan <- channel.Data{Type: "clock", Object: nowStr}
	}
}

func update() {
	go updateClock()
	// other code here
}

func start() {
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

	go update()

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		panic(err)
	}

	view.Handler().Wait()
}

func main() {
	start()
}
