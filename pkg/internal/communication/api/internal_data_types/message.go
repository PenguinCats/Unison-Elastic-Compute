package internal_data_types

type MessageDataType int

const (
	MessageDataTypeContainerCreate MessageDataType = iota + 1
	MessageDataTypeContainerStart
	MessageDataTypeContainerStop
	MessageDataTypeContainerRemove
	MessageDataTypeContainerProfile
	MessageDataTypeContainerStatus
)

type Message struct {
	MessageType MessageDataType `json:"message_type"`
	Value       []byte          `json:"value"`
}
