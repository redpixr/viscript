package process

import (
	"github.com/corpusc/viscript/msg"
)

func (st *State) UnpackMessage(msgType uint16, message []byte) []byte {
	switch msgType {

	case msg.TypeChar:
		var m msg.MessageChar
		msg.MustDeserialize(message, &m)
		st.onChar(m)

	case msg.TypeKey:
		var m msg.MessageKey
		msg.MustDeserialize(message, &m)
		st.onKey(m, message)

	case msg.TypeMouseScroll:
		var m msg.MessageMouseScroll
		msg.MustDeserialize(message, &m)
		st.onMouseScroll(m, message)

	default:
		println("UNKNOWN MESSAGE TYPE!")
	}

	if st.DebugPrintInputEvents {
		println()
	}

	return message
}
