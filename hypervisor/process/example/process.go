package process

import (
	"github.com/corpusc/viscript/msg"
)

//non-instanced
func NewProcess() *Process {
	println("(process/example/process.go).NewProcess()")
	var p Process

	p.Id = msg.NextProcessId()

	p.MessageIn = make(chan []byte, msg.ChannelCapacity)
	//p.OutChannel = make(chan []byte, msg.ChannelCapacity)

	p.State.Init(&p)

	return &p
}

//Example process
type Process struct {
	Id msg.ProcessId

	MessageIn chan []byte
	//OutChannel chan []byte

	State State
}

func (self *Process) GetProcessInterface() msg.ProcessInterface {
	println("(process/example/process.go).GetProcessInterface()")
	return msg.ProcessInterface(self)
}

func (self *Process) DeleteProcess() {
	println("(process/example/process.go).DeleteProcess()")
	close(self.MessageIn)
	//close(self.OutChannel)
	self.State.proc = nil
	self = nil
}

// Implement ProcessInterface

func (self *Process) GetId() msg.ProcessId {
	println("(process/example/process.go).GetId()")
	return self.Id
}

func (self *Process) GetIncomingChannel() chan []byte {
	println("(process/example/process.go).GetIncomingChannel()")
	return self.MessageIn
}

// func (self *Process) GetOutChannel() chan []byte {
// 	//println("(process/example/process.go).GetOutChannel()")
// 	return self.OutChannel
// }

//Business logic
func (self *Process) Tick() {
	//println("(process/example/process.go).Tick()")
	self.State.HandleMessages()
}
