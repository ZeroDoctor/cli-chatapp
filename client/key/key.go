package key

import (
	"github.com/awesome-gocui/gocui"
	"gitlab.com/smallwood/sw-chat/view"
)

func SetBindings(g *gocui.Gui) {
	// global binding
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, view.Handler().Quit); err != nil {
		panic(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, view.Handler().Quit); err != nil {
		panic(err)
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, view.Handler().NextView); err != nil {
		panic(err)
	}

	// screen bindings
	if err := g.SetKeybinding("screen", rune('j'), gocui.ModNone, downScreen); err != nil {
		panic(err)
	}

	if err := g.SetKeybinding("screen", rune('k'), gocui.ModNone, upScreen); err != nil {
		panic(err)
	}
}

func upScreen(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy-5); err != nil && oy > 0 {
		if err := v.SetOrigin(ox, oy-5); err != nil {
			return err
		}
	}

	return nil
}

func downScreen(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy+5); err != nil {
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy-5); err != nil {
			return err
		}
	}

	return nil
}

