package web

type FormType string

const (
	UserForm     FormType = "user"
	LocationForm FormType = "location"
	PostForm     FormType = "post"
	MessageForm  FormType = "message"
	NullForm     FormType = ""
)

const (
	FormValue = "type"
	
	PostPathValue = "post_id"
	MsgPathValue  = "msg_id"
)
