package hypervisor

import (
	"fmt"
	//_ "image/png"
	//"log"
	/*
		"go/build"
		"runtime"
	*/
	"bytes"
	"math"
	//"strconv"

	"encoding/binary"

	"github.com/corpusc/viscript/gfx"
	"github.com/corpusc/viscript/msg"
	"github.com/corpusc/viscript/script"
	"github.com/go-gl/glfw/v3.2/glfw"
)

/*
var prevMousePixelX float64
var prevMousePixelY float64
var mousePixelDeltaX float64
var mousePixelDeltaY float64
*/

// this can also be triggered by onMouseButton
func onMouseCursorPos(m msg.MessageMousePos) {

	// gfx.Curs.UpdatePosition(float32(x), float32(y)) //state update

	// mousePixelDeltaX = x - prevMousePixelX
	// mousePixelDeltaY = y - prevMousePixelY
	// prevMousePixelX = x
	// prevMousePixelY = y

	// //rendering update
	// if /* LMB held */ w.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
	// 	gfx.Rend.ScrollPanelThatIsHoveredOver(mousePixelDeltaX, mousePixelDeltaY)
	// }

}

func onMouseScroll(m msg.MessageMouseScroll) {
	/*
		var delta float64 = 30

		// if horizontal
		//state update
		if w.GetKey(glfw.KeyLeftShift) == glfw.Press || w.GetKey(glfw.KeyRightShift) == glfw.Press {
			gfx.Rend.ScrollPanelThatIsHoveredOver(yOff*-delta, 0)
		} else {
			gfx.Rend.ScrollPanelThatIsHoveredOver(xOff*delta, yOff*-delta)
		}
	*/
}

//Working
//func onChar(w *glfw.Window, char rune) {
func onChar(m msg.MessageOnCharacter) {
	InsertRuneIntoDocument("Rune", m.Rune)
	script.Process(false)
}

// WEIRD BEHAVIOUR OF KEY EVENTS.... for a PRESS, you can detect a
// shift/alt/ctrl/super key through the "mod" variable,
// (see the top of "action == glfw.Press" section for an example)
// regardless of left/right key used.
// BUT for RELEASE, the "mod" variable will NOT tell you what key it is!
// so you will have to handle both left & right mod keys via the "action" variable!

// func onKey(
// 	w *glfw.Window,
// 	key glfw.Key,
// 	scancode int,
// 	action glfw.Action,
// 	mod glfw.ModifierKey) {
func onKey(m msg.MessageKey) {
	/*
		foc := gfx.Rend.Focused

		if action == glfw.Release {
			switch key {

			case glfw.KeyEscape:
				w.SetShouldClose(true)

			case glfw.KeyLeftShift:
				fallthrough
			case glfw.KeyRightShift:
				fmt.Println("done selecting")
				foc.Selection.CurrentlySelecting = false // TODO?  possibly flip around if selectionStart comes after selectionEnd in the page flow?

			case glfw.KeyLeftControl:
				fallthrough
			case glfw.KeyRightControl:
				fmt.Println("Control RELEASED")
			case glfw.KeyLeftAlt:
				fallthrough
			case glfw.KeyRightAlt:
				fmt.Println("Alt RELEASED")
			case glfw.KeyLeftSuper:
				fallthrough
			case glfw.KeyRightSuper:
				fmt.Println("'Super' modifier key RELEASED")
			}
		} else {
			// glfw.Repeat   or   glfw.Press
			b := foc.TextBodies[0]

			CharWid := int32(gfx.Rend.CharWidInPixels)
			CharHei := int32(gfx.Rend.CharHeiInPixels)
			numOfCharsV := gfx.CurrAppHeight / CharHei
			numOfCharsH := gfx.CurrAppWidth / CharWid

			s := strconv.Itoa(int(numOfCharsV))

			fmt.Printf("Rectangle Right %s\n\n\n", s)

			switch mod {
			case glfw.ModShift:
				fmt.Println("started selecting")
				foc.Selection.CurrentlySelecting = true
				foc.Selection.StartX = foc.CursX
				foc.Selection.StartY = foc.CursY
			}

			switch key {
			case glfw.KeyEnter:
				startOfLine := b[foc.CursY][:foc.CursX]
				restOfLine := b[foc.CursY][foc.CursX:len(b[foc.CursY])]
				b[foc.CursY] = startOfLine
				b = insert(b, foc.CursY+1, restOfLine)

				foc.CursX = 0
				foc.CursY++
				foc.TextBodies[0] = b

				if foc.CursY >= len(b) {
					foc.CursY = len(b) - 1
				}
			case glfw.KeyHome:
				if eitherControlKeyHeld(w) {
					foc.CursY = 0
				}

				foc.CursX = 0
				movedCursorSoUpdateDependents()
			case glfw.KeyEnd:
				if eitherControlKeyHeld(w) {
					foc.CursY = len(b) - 1
				}

				foc.CursX = len(b[foc.CursY])
				movedCursorSoUpdateDependents()
			case glfw.KeyUp:
				if foc.CursY > 0 {
					foc.CursY--

					if foc.CursX > len(b[foc.CursY]) {
						foc.CursX = len(b[foc.CursY])
					}
				}

				movedCursorSoUpdateDependents()
			case glfw.KeyDown:
				if foc.CursY < len(b)-1 {
					if numOfCharsV < (int32(foc.CursY) + 1) {
						gfx.Rend.ScrollPanelThatIsHoveredOver(0, float64(CharHei))
					}
					foc.CursY++

					if foc.CursX > len(b[foc.CursY]) {
						foc.CursX = len(b[foc.CursY])
					}
				}

				movedCursorSoUpdateDependents()
			case glfw.KeyLeft:
				if foc.CursX == 0 {
					if foc.CursY > 0 {
						foc.CursY--
						foc.CursX = len(b[foc.CursY])
					}
				} else {
					if mod == glfw.ModControl {
						foc.CursX = getWordSkipPos(foc.CursX, -1)
					} else {
						if (numOfCharsH - int32(foc.CursX)) > (int32(foc.CursX) + 4) {
							gfx.Rend.ScrollPanelThatIsHoveredOver(float64(-CharWid), 0)
						}
						foc.CursX--
					}
				}

				movedCursorSoUpdateDependents()
			case glfw.KeyRight:
				if foc.CursX < len(b[foc.CursY]) {
					if mod == glfw.ModControl {
						foc.CursX = getWordSkipPos(foc.CursX, 1)
					} else {
						fmt.Println(numOfCharsH)
						if numOfCharsH < (int32(foc.CursX) + 4) {
							gfx.Rend.ScrollPanelThatIsHoveredOver(float64(CharWid), 0)
						}
						foc.CursX++
					}
				}

				movedCursorSoUpdateDependents()
			case glfw.KeyBackspace:
				if foc.CursX == 0 {
					b = remove(b, foc.CursY, b[foc.CursY])
					foc.TextBodies[0] = b
					foc.CursY--
					foc.CursX = len(b[foc.CursY])

				} else {
					foc.RemoveCharacter(false)
				}

			case glfw.KeyDelete:
				foc.RemoveCharacter(true)
				fmt.Println("Key Deleted")

			}

			script.Process(false)
		}
	*/
}

func eitherControlKeyHeld(w *glfw.Window) bool {
	if w.GetKey(glfw.KeyLeftControl) == glfw.Press || w.GetKey(glfw.KeyRightControl) == glfw.Press {
		return true
	} else {
		return false
	}
}

// must be in range
func insert(slice []string, index int, value string) []string {
	slice = slice[0 : len(slice)+1]      // grow the slice by one element
	copy(slice[index+1:], slice[index:]) // move the upper part of the slice out of the way and open a hole
	slice[index] = value
	return slice
}

// similar to insert method, instead moves current slice element and appends to one above
func remove(slice []string, index int, value string) []string {
	slice = append(slice[:index], slice[index+1:]...)
	slice[index-1] = slice[index-1] + value
	return slice
}

func getWordSkipPos(xIn int, change int) int {
	peekPos := xIn
	foc := gfx.Rend.Focused
	b := foc.TextBodies[0]

	for {
		peekPos += change

		if peekPos < 0 {
			return 0
		}

		if peekPos >= len(b[foc.CursY]) {
			return len(b[foc.CursY])
		}

		if string(b[foc.CursY][peekPos]) == " " {
			return peekPos
		}
	}
}

func movedCursorSoUpdateDependents() {
	// --- Always-Visible-Cursor Autoscrolling ---
	//
	// TODO ^

	//
	// --- Selection Marking ---
	//
	// when SM is made functional,
	// we should probably detect whether cursor
	// position should update Start_ or End_ at this point.
	// rather than always making that the "end".
	// i doubt marking forwards or backwards will ever alter what is
	// done with the selection

	foc := gfx.Rend.Focused

	if foc.Selection.CurrentlySelecting {
		foc.Selection.EndX = foc.CursX
		foc.Selection.EndY = foc.CursY
	} else { // moving cursor without shift gets rid of selection
		foc.Selection.StartX = math.MaxUint32
		foc.Selection.StartY = math.MaxUint32
		foc.Selection.EndX = math.MaxUint32
		foc.Selection.EndY = math.MaxUint32
	}
}

func getSlice(wBuf *bytes.Buffer, err error) (data []byte) {
	data = make([]byte, 0)

	if err != nil {
		fmt.Println("binary.Write failed:", err)
	} else {
		b := wBuf.Bytes()

		for i := 0; i < wBuf.Len(); i++ {
			data = append(data, b[i])
		}
	}

	return
}
func getBytesOfRune(value rune) (data []byte) {
	wBuf := new(bytes.Buffer)
	err := binary.Write(wBuf, binary.LittleEndian, value)
	data = getSlice(wBuf, err)
	return
}
