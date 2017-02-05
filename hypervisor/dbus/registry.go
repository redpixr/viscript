package dbus

type DbusInstance struct {
	PubsubChannels map[ChannelId]PubsubChannel
	Resources      []ResourceMeta
}

func (self *DbusInstance) Init() {
	println("(dbus/instance.go).Init()")
	self.PubsubChannels = make(map[ChannelId]PubsubChannel)
	self.Resources = make([]ResourceMeta, 0)
}

func (self *DbusInstance) AddPubSubChannel(channelId ChannelId, pubSubChannel PubsubChannel) {
	println("(dbus/instance.go).AddPubSubChannel()")
	self.PubsubChannels[channelId] = pubSubChannel
}

//register that a resource exists
func (self *DbusInstance) ResourceRegister(ResourceId ResourceId, ResourceType ResourceType) {
	println("(dbus/instance.go).ResourceRegister()")
	x := ResourceMeta{}
	x.Id = ResourceId
	x.Type = ResourceType

	self.Resources = append(self.Resources, x)
}

//remove resource from list
func (self *DbusInstance) ResourceUnregister(ResourceID ResourceId, ResourceType ResourceType) {
	println("(dbus/instance.go).ResourceUnregister()")
	for i, resourceMeta := range self.Resources {
		if resourceMeta.Id == ResourceID {
			self.Resources = append(self.Resources[:i], self.Resources[i+1:]...)
		}
	}
}
