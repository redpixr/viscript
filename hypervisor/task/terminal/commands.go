package process

import (
	"strconv"
	"strings"

	"bytes"

	"fmt"

	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/config"
	"github.com/skycoin/viscript/hypervisor"
	extTask "github.com/skycoin/viscript/hypervisor/task_ext"
	"github.com/skycoin/viscript/msg"
	"github.com/skycoin/viscript/signal"
)

const cp = "hypervisor/task/terminal/commands"

func (st *State) commandHelp() {
	st.PrintLn(app.GetBarOfChars("-", int(st.VisualInfo.NumColumns)))
	//st.PrintLn("Current commands:")
	st.PrintLn("------ Terminals ------")
	st.PrintLn("clear:                 Clears currently focused terminal.")
	st.PrintLn("close_term <id>:         Delete terminal with index to the terminal id.")
	st.PrintLn("list_terms:            List all terminal ids.")
	st.PrintLn("new_term:              Add new terminal (n for short).")
	st.PrintLn("------ Apps -----------")
	st.PrintLn("apps:                  Display all available apps with descriptions.")
	st.PrintLn("attach    <id>:        Attach external task with given id to terminal.")
	st.PrintLn("list_tasks (-f):       List running tasks (-f for full commands).")
	st.PrintLn("ping      <id>:        Ping app with given id.")
	st.PrintLn("res_usage <id>:        See resource usage for app with given id.")
	st.PrintLn("shutdown  <id>:        [TODO] Shutdown external task with given id.")
	st.PrintLn("start (-a) <command>:  Start external task. (-a to also attach).")
	// st.PrintLn("rpc:                   Issues command: \"go run rpc/cli/cli.go\"")
	// st.PrintLn("Current hotkeys:")
	st.PrintLn("CTRL+Z:                Detach currently attached process.")
	// st.PrintLn("    CTRL+C:           ___description goes here___")
	st.PrintLn(app.GetBarOfChars("-", int(st.VisualInfo.NumColumns)))
}

func (st *State) commandDisplayApps() {
	app.At(cp, "commandDisplayApps")
	apps := config.Global.Apps

	if len(apps) == 0 {
		st.PrintLn("No available apps found.")
		return
	}

	maxAppKeyLength := 0

	for appKey, _ := range apps {
		appKeyLength := len(appKey)

		if appKeyLength > maxAppKeyLength {
			maxAppKeyLength = appKeyLength
		}
	}

	var buffer bytes.Buffer

	maxAppKeyLength += 4 // Space after max string length app hash

	for appKey, app := range apps {
		buffer.WriteString(appKey)
		for i := 0; i < maxAppKeyLength-len(appKey); i++ {
			buffer.WriteString(" ")
		}
		buffer.WriteString(fmt.Sprintf("-%s\n", app.Desc))
	}

	st.PrintLn(buffer.String())
}

func (st *State) commandAppHelp(args []string) {
	app.At(cp, "commandAppHelp")

	appName := args[0]

	if !config.AppExistsWithName(appName) {
		st.PrintError("App with name: " + appName + "doesn't exist. " +
			"Try running 'apps'.")
		return
	}

	st.PrintLn(config.Global.Apps[appName].Help)
}

func (st *State) commandClearTerminal() {
	app.At(cp, "commandClearTerminal")

	st.VisualInfo.CurrRow = 0
	st.publishToOut(msg.Serialize(msg.TypeClear, msg.MessageClear{}))
	st.Cli.EchoWholeCommand(st.proc.OutChannelId)
}

func (st *State) commandStart(args []string) {
	app.At(cp, "commandStart")

	if len(args) < 1 {
		st.PrintError("Must pass a command into Start!")
		return
	}

	detached := args[0] != "-a"

	if !detached {
		args = args[1:]
	}

	appName := args[0]

	if !config.AppExistsWithName(appName) {
		st.PrintError("App with name: " + appName + "doesn't exist. " +
			"Try running 'apps'.")
		return
	}

	var tokens []string

	//if there are user passed args for the app override defaults set in config
	if len(args) > 1 {

		pathToApp := config.GetPathForApp(appName)

		tokens = append(tokens, pathToApp)

		for _, arg := range args {
			tokens = append(tokens, strings.ToLower(arg))
		}

	} else {
		tokens = config.GetPathWithDefaultArgsForApp(appName)
	}

	//if the app is daemon not allow to attach to it
	if config.Global.Apps[appName].Daemon {
		detached = true
	}

	newExtProc, err := extTask.MakeNewTaskExternal(tokens, detached)
	if err != nil {
		st.PrintError(err.Error())
		return
	}

	err = newExtProc.Start()
	if err != nil {
		st.PrintError(err.Error())
		return
	}

	extProcInterface := newExtProc.GetExtTaskInterface()

	procId := hypervisor.AddExtTask(extProcInterface)

	if !detached {
		err = st.proc.AttachExternalProcess(extProcInterface)
		if err != nil {
			st.PrintError(err.Error())
		}
	}

	st.PrintLn("Added External Process (ID: " +
		strconv.Itoa(int(procId)) + ", Command: " +
		newExtProc.CommandLine + ")")

}

func (st *State) commandAppPing(args []string) {
	app.At(cp, "commandAppPing")

	if len(args) < 1 {
		st.PrintError("No task id passed! e.g. ping 1")
		return
	}

	passedID, err := strconv.Atoi(args[0])
	if err != nil {
		st.PrintError("Task id must be an integer.")
		return
	}

	extProcID := msg.ExtProcessId(passedID)

	if !hypervisor.ExtProcessIsRunning(extProcID) {
		st.PrintError("Taks with given id is not running.")
		return
	}

	msgUserCommand := msg.MessageUserCommand{
		Sequence: signal.GetNextMessageID(),
		AppId:    uint32(extProcID),
		Payload:  msg.Serialize(msg.TypePing, msg.MessagePing{})}

	serializedCommand := msg.Serialize(msg.TypeUserCommand, msgUserCommand)

	signal.Monitor.Send(uint32(extProcID), serializedCommand)
}

func (st *State) commandShutDown(args []string) {
	app.At(cp, "commandShutDown")

	if len(args) < 1 {
		st.PrintError("No task id passed! e.g. shutdown 1")
		return
	}

	passedID, err := strconv.Atoi(args[0])
	if err != nil {
		st.PrintError("Task id must be an integer.")
		return
	}

	extProcID := msg.ExtProcessId(passedID)

	if !hypervisor.ExtProcessIsRunning(extProcID) {
		st.PrintError("Task with given id is not running.")
		return
	}

	msgUserCommand := msg.MessageUserCommand{
		Sequence: signal.GetNextMessageID(),
		AppId:    uint32(extProcID),
		Payload:  msg.Serialize(msg.TypeShutdown, msg.MessageShutdown{})}

	serializedCommand := msg.Serialize(msg.TypeUserCommand, msgUserCommand)

	signal.Monitor.Send(uint32(extProcID), serializedCommand)
}

func (st *State) commandResourceUsage(args []string) {
	app.At(cp, "commandResourceUsage")
	if len(args) < 1 {
		st.PrintError("No task id passed! e.g. res_usage 1")
		return
	}

	passedID, err := strconv.Atoi(args[0])
	if err != nil {
		st.PrintError("Task id must be an integer.")
		return
	}

	extProcID := msg.ExtProcessId(passedID)

	if !hypervisor.ExtProcessIsRunning(extProcID) {
		st.PrintError("Task with give id is not running.")
		return
	}

	msgUserCommand := msg.MessageUserCommand{
		Sequence: signal.GetNextMessageID(),
		AppId:    uint32(extProcID),
		Payload: msg.Serialize(msg.TypeResourceUsage,
			msg.MessageResourceUsage{})}

	serializedCommand := msg.Serialize(msg.TypeUserCommand, msgUserCommand)

	signal.Monitor.Send(uint32(extProcID), serializedCommand)
}

func (st *State) commandAttach(args []string) {
	app.At(cp, "commandAttach")

	if len(args) < 1 {
		st.PrintError("No task id passed! e.g. attach 1")
		return
	}

	passedID, err := strconv.Atoi(args[0])
	if err != nil {
		st.PrintError("Task id must be an integer.")
		return
	}

	extProcID := msg.ExtProcessId(passedID)

	extProc, err := hypervisor.GetExtProcess(extProcID)
	if err != nil {
		st.PrintError(err.Error())
		return
	}

	st.PrintLn(extProc.GetFullCommandLine())
	err = st.proc.AttachExternalProcess(extProc)
	if err != nil {
		st.PrintError(err.Error())
	}
}

func (st *State) commandListExternalTasks(args []string) {
	app.At(cp, "commandListExternalTasks")

	extTaskMap := hypervisor.ExtTaskListGlobal.ProcessMap
	if len(extTaskMap) == 0 {
		st.PrintLn("No external tasks running.\n" +
			"Try starting one with \"start\" command (\"help\" or \"h\" for help).")
		return
	}

	fullPrint := false

	if len(args) > 0 && args[0] == "-f" {
		fullPrint = true
	}

	for procId, extProc := range extTaskMap {
		procCommand := ""

		if fullPrint {
			procCommand = extProc.GetFullCommandLine()
		} else {
			procCommand = strings.Split(
				extProc.GetFullCommandLine(), " ")[0]
		}

		st.Printf("[ %d ] -> [ %s ]\n", int(procId), procCommand)
	}
}

func (st *State) commandCloseTerminalFirstStage(args []string) {
	if len(args) == 1 {
		//handle storedTerminalIds errors
		if /****/ len(st.storedTerminalIds) < 1 {
			st.PrintError("Use 'list_terms' command to see their IDs")
			return
		} else if len(st.storedTerminalIds) == 1 {
			st.PrintError("Shouldn't delete when only 1 remains (UNTIL GUI IS MADE)")
			return
		}

		//handle arg conversion errors
		storedId, err := strconv.Atoi(args[0])
		if err != nil {
			st.PrintError("Unable to convert passed index.")
			s := "err.Error(): \"" + err.Error() + "\""
			st.PrintError(s)
			return
		}

		//handle index errors
		if storedId < 0 ||
			storedId >= len(st.storedTerminalIds) {
			st.PrintError("Index not in range.")
			return
		}

		//everything should be valid here
		st.SendCommand("close_term", []string{strconv.Itoa(int(st.storedTerminalIds[storedId]))})
	} else { //args failure (too many/few passed)
		st.PrintError("Must supply ONE valid ID argument")
		//IDEALLY we'd use ANY VALID index at position 0,
		//but i think you'd want us to prioritize simplicity
		//above doing extra coding to handle nitpickiness like that.
		return
	}
}
