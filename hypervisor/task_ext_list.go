package hypervisor

import (
	"errors"
	"strconv"

	"github.com/skycoin/viscript/msg"
)

var ExtTaskListGlobal ExtTaskList

type ExtTaskList struct {
	ProcessMap map[msg.ExtProcessId]msg.ExtTaskInterface
}

func initExtTaskList() {
	ExtTaskListGlobal.ProcessMap = make(map[msg.ExtProcessId]msg.ExtTaskInterface)
}

func teardownExtTaskList() {
	ExtTaskListGlobal.ProcessMap = nil
	// TODO: Further cleanup
}

func ExtProcessIsRunning(procId msg.ExtProcessId) bool {
	_, exists := ExtTaskListGlobal.ProcessMap[procId]
	return exists
}

func AddExtTask(ep msg.ExtTaskInterface) msg.ExtProcessId {
	id := ep.GetId()

	if !ExtProcessIsRunning(id) {
		ExtTaskListGlobal.ProcessMap[id] = ep
	}

	return id
}

func GetExtProcess(id msg.ExtProcessId) (msg.ExtTaskInterface, error) {
	extProc, exists := ExtTaskListGlobal.ProcessMap[id]
	if exists {
		return extProc, nil
	}

	err := errors.New("External process with id " +
		strconv.Itoa(int(id)) + " doesn't exist!")

	return nil, err
}

func RemoveExtProcess(id msg.ExtProcessId) {
	delete(ExtTaskListGlobal.ProcessMap, id)
}

func TickExtTasks() {
	// TODO: Read from response channels if they contain any new messages
	// for _, p := range ExtTaskListGlobal.ProcessMap {
	// data, err := monitor.Monitor.ReadFrom(p.GetId())
	// if err != nil {
	// 	// println(err.Error())
	// 	// monitor.Monitor.PrintAll()
	// 	continue
	// }

	// ackType := msg.GetType(data)

	// switch ackType {
	// case msg.TypeUserCommandAck:

	// }

	// select {
	// case <-p.GetTaskExitChannel():
	// 	println("Got the exit in task ext list")
	// default:
	// }
	// p.Tick()
	// }

}
