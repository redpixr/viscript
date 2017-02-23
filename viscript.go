/*

------- NEXT THINGS TODO: -------

* Unfocusing
* Move terminals
	* "option to snap to grid"
* Resize terminals
	* "snap to grid or char size"

* put console feedback into focused terminal



------- NEWER TODO: -------

* put console feedback into focused terminal



------- OLDER TODO: ------- (everything below was for the text editor)

* KEY-BASED NAVIGATION
	* CTRL-HOME/END - PGUP/DN
* BACKSPACE/DELETE at the ends of lines
	pulls us up to prev line, or pulls up next line
* when there is no scrollbar, should be able to see/interact with text in that area
* when auto appending to the end of a terminal, scroll all the way down
		(manual activity in the middle could increase size, so do this only when appending to body)


------- LOWER PRIORITY POLISH: -------

* if typing goes past right of screen, auto-horizontal-scroll as you type
* same for when newlines/enters/returns push cursor past the bottom of visible space
* scrollbars should have a bottom right corner, and a thickness sized background
		for void space, reserved for only that, so the bar never covers up the rightmost
		character/cursor
* when pressing delete at/after the end of a line, should pull up the line below
* vertical scrollbars could have a smaller rendering of the first ~40 chars?
		however not if we map the whole vertical space (when scrollspace is taller than screen),
		because this requires scaling the text.  and keeping the aspect ratio means ~40 (max)
		would alter the width of the scrollbar

*/

package main

import (
	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/rpc/terminalmanager"
	"github.com/corpusc/viscript/viewport"
)

func main() {
	println("Starting...")

	hypervisor.Init()

	viewport.DebugPrintInputEvents = true //print input events
	viewport.ViewportInit()               //runtime.LockOSThread(), InitCanvas()
	viewport.ViewportScreenInit()
	viewport.InitEvents()
	viewport.ViewportTerminalsInit() //start the terminal

	// rpc
	go func() {
		rpcInstance := terminalmanager.NewRPC()
		rpcInstance.Serve()
	}()

	println("Start Loop;")
	for viewport.CloseWindow == false {
		viewport.DispatchEvents() //event channel
		hypervisor.ProcessTick()  //processes, handle incoming events
		viewport.PollUiInputEvents()
		viewport.Tick()
		viewport.UpdateDrawBuffer()
		viewport.SwapDrawBuffer() //with new frame
	}

	println("Closing down viewport")
	viewport.ViewportScreenTeardown()
	hypervisor.HypervisorTeardown()
}
