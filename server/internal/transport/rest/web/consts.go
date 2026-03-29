package web

type FormType string

const (
	UserForm     FormType = "user"
	LocationForm FormType = "location"
	TopicForm    FormType = "topic"
	MessageForm  FormType = "message"
	NullForm     FormType = ""
)

const (
	FormValue = "type"

	TopicPathValue = "topic_id"
	MsgPathValue   = "msg_id"
)
