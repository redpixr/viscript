package viewport

import (
	"fmt"
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/hypervisor/input/mouse"
	"github.com/corpusc/viscript/msg"
)

// triggered both by moving **AND*** by pressing buttons
func onMouseCursorPos(m msg.MessageMousePos) {
	if DebugPrintInputEvents {
		fmt.Print("TypeMousePos")
		showFloat64("X", m.X)
		showFloat64("Y", m.Y)
		println()
	}

	mouse.UpdatePosition(app.Vec2F{float32(m.X), float32(m.Y)}) // state update

	if mouse.HoldingLeftButton {
		println("TODO EVENTUALLY: implement something like 'ScrollFocusedTerm()'")
		//old logic is in ScrollTermThatHasMousePointer(mouse.PixelDelta.X, mouse.PixelDelta.Y),
		//which was a janky way to do it
	}
}

func onMouseScroll(m msg.MessageMouseScroll) {
	if DebugPrintInputEvents {
		print("TypeMouseScroll")
		showFloat64("X Offset", m.X)
		showFloat64("Y Offset", m.Y)
		showBool("HoldingControl", m.HoldingControl)
		println()
	}
}

// apparently every time this is fired, a mouse position event is ALSO fired
func onMouseButton(m msg.MessageMouseButton) {
	if DebugPrintInputEvents {
		fmt.Print("TypeMouseButton")
		showUInt8("Button", m.Button)
		showUInt8("Action", m.Action)
		showUInt8("Mod", m.Mod)
		println()
	}

	convertClickToTextCursorPosition(m.Button, m.Action)

	if msg.Action(m.Action) == msg.Press {
		switch msg.MouseButton(m.Button) {
		case msg.MouseButtonLeft:
			mouse.HoldingLeftButton = true

			// // detect clicks in rects
			// if mouse.CursorIsInside(ui.MainMenu.Rect) {
			// 	respondToAnyMenuButtonClicks()
			// } else { // respond to any panel clicks outside of menu
			focusOnTopmostRectThatContainsPointer()
			// }
		}
	} else if msg.Action(m.Action) == msg.Release {
		switch msg.MouseButton(m.Button) {
		case msg.MouseButtonLeft:
			mouse.HoldingLeftButton = false
		}
	}
}

func focusOnTopmostRectThatContainsPointer() {
	var topmostZ float32
	var topmostId msg.TerminalId

	for id, t := range Terms.Terms {
		if mouse.CursorIsInside(t.Bounds) {
			if topmostZ < t.Depth {
				topmostZ = t.Depth
				topmostId = id
			}
		}
	}

	if topmostZ > 0 {
		Terms.FocusedId = topmostId
		Terms.Focused = Terms.Terms[topmostId]
	}
}

func convertClickToTextCursorPosition(button, action uint8) {
	// if glfw.MouseButton(button) == glfw.MouseButtonLeft &&
	// 	glfw.Action(action) == glfw.Press {

	// 	foc := Focused

	// 	if foc.IsEditable && foc.Content.Contains(mouse.GlX, mouse.GlY) {
	// 		if foc.MouseY < len(foc.TextBodies[0]) {
	// 			foc.CursY = foc.MouseY

	// 			if foc.MouseX <= len(foc.TextBodies[0][foc.CursY]) {
	// 				foc.CursX = foc.MouseX
	// 			} else {
	// 				foc.CursX = len(foc.TextBodies[0][foc.CursY])
	// 			}
	// 		} else {
	// 			foc.CursY = len(foc.TextBodies[0]) - 1
	// 		}
	// 	}
	// }
}

func respondToAnyMenuButtonClicks() {
	// for _, bu := range ui.MainMenu.Buttons {
	// 	if mouse.CursorIsInside(bu.Rect.Rectangle) {
	// 		bu.Activated = !bu.Activated

	// 		switch bu.Name {
	// 		case "Run":
	// 			if bu.Activated {
	// 				//script.Process(true)
	// 			}
	// 			break
	// 		case "Testing Tree":
	// 			if bu.Activated {
	// 				//script.Process(true)
	// 				//script.MakeTree()
	// 			} else { // deactivated
	// 				// remove all terminals with trees
	// 				b := Terms[:0]
	// 				for _, t := range Terms {
	// 					if len(t.Trees) < 1 {
	// 						b = append(b, t)
	// 					}
	// 				}
	// 				Terms = b
	// 				//fmt.Printf("len of b (from Terms) after removing ones with trees: %d\n", len(b))
	// 				//fmt.Printf("len of Terms: %d\n", len(Terms))
	// 			}
	// 			break
	// 		}

	// 		app.Con.Add(fmt.Sprintf("%s toggled\n", bu.Name))
	// 	}
	// }
}

// the rest of these funcs are almost identical, just top 2 vars customized (and string format)
func showBool(s string, x bool) {
	fmt.Printf("   [%s: %t]", s, x)
}

func showUInt8(s string, x uint8) {
	fmt.Printf("   [%s: %d]", s, x)
}

func showSInt32(s string, x int32) {
	fmt.Printf("   [%s: %d]", s, x)
}

func showUInt32(s string, x uint32) {
	fmt.Printf("   [%s: %d]", s, x)
}

func showFloat64(s string, f float64) {
	fmt.Printf("   [%s: %.1f]", s, f)
}
