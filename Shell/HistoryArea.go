package shell

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewHistoryArea() *History {
	return &History{
		Box: tview.NewBox(),
	}
}

type History struct {
	// tview.Primitive
	*tview.Box
}

// Draw draws this primitive onto the screen. Implementers can call the
// screen's ShowCursor() function but should only do so when they have focus.
// (They will need to keep track of this themselves.)
func (hist *History) Draw(screen tcell.Screen) {
	hist.Box.DrawForSubclass(screen, hist)
}

// InputHandler returns a handler which receives key events when it has focus.
// It is called by the Application class.
//
// A value of nil may also be returned, in which case this primitive cannot
// receive focus and will not process any key events.
//
// The handler will receive the key event and a function that allows it to
// set the focus to a different primitive, so that future key events are sent
// to that primitive.
//
// The Application's Draw() function will be called automatically after the
// handler returns.
//
// The Box class provides functionality to intercept keyboard input. If you
// subclass from Box, it is recommended that you wrap your handler using
// Box.WrapInputHandler() so you inherit that functionality.
func (hist *History) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return hist.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		key := event.Key()
		_ = key
	})
}

// MouseHandler returns a handler which receives mouse events.
// It is called by the Application class.
//
// A value of nil may also be returned to stop the downward propagation of
// mouse events.
//
// The Box class provides functionality to intercept mouse events. If you
// subclass from Box, it is recommended that you wrap your handler using
// Box.WrapMouseHandler() so you inherit that functionality.
func (hist *History) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return hist.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		x, y := event.Position()
		if !hist.InRect(x, y) {
			return false, nil
		}

		switch action {
		case tview.MouseLeftDown:
			setFocus(hist)
			consumed = true
		case tview.MouseLeftClick:
			consumed = true
		case tview.MouseScrollUp:
			consumed = true
		case tview.MouseScrollDown:
			consumed = true
		}

		return
	})
}
