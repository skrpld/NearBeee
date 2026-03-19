package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/skrpld/NearBeee/internal/core/models/dto"
	"github.com/skrpld/NearBeee/internal/transport/rest/web"
	"github.com/skrpld/NearBeee/pkg/errors"
)

type MessagesService interface {
	CreateMessage(ctx context.Context, rows *dto.CreateMessageRequest) (*dto.CreateMessageResponse, error)
	GetMessageByMessageId(ctx context.Context, rows *dto.GetMessageByMessageIdRequest) (*dto.GetMessageByMessageIdResponse, error)
	GetMessageByUserId(ctx context.Context, rows *dto.GetMessageByUserIdRequest) (*dto.GetMessageByUserIdResponse, error)
	GetMessagesByPostId(ctx context.Context, rows *dto.GetMessagesByPostIdRequest) (*dto.GetMessagesByPostIdResponse, error)
	UpdateMessageById(ctx context.Context, rows *dto.UpdateMessageByIdRequest) (*dto.UpdateMessageByIdResponse, error)
	DeleteMessageById(ctx context.Context, rows *dto.DeleteMessageByIdRequest) (*dto.DeleteMessageByIdResponse, error)
}

type MessagesController struct {
	messagesSrv MessagesService
}

func NewMessagesController(messagesSrv MessagesService) *MessagesController {
	return &MessagesController{messagesSrv: messagesSrv}
}

func (c *MessagesController) CreateMessage(r *http.Request) (any, error) {
	var request dto.CreateMessageRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return nil, err
	}

	request.UserId = user.UserId

	return c.messagesSrv.CreateMessage(r.Context(), &request)
}

func (c *MessagesController) GetMessage(r *http.Request) (any, error) {
	switch web.FormType(r.FormValue(web.FormValue)) {
	case web.UserForm:
		return c.GetMessageByUserId(r)
	case web.PostForm:
		return c.GetMessagesByPostId(r)
	case web.MessageForm, web.NullForm:
		return c.GetMessageByMessageId(r)
	default:
		return nil, errors.ErrInvalidFormType
	}
}

func (c *MessagesController) GetMessageByMessageId(r *http.Request) (any, error) {
	var request dto.GetMessageByMessageIdRequest

	request.MessageId = r.PathValue(web.MsgPathValue)

	return c.messagesSrv.GetMessageByMessageId(r.Context(), &request)
}

func (c *MessagesController) GetMessageByUserId(r *http.Request) (any, error) {
	var request dto.GetMessageByUserIdRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return nil, err
	}

	request.UserId = user.UserId

	return c.messagesSrv.GetMessageByUserId(r.Context(), &request)
}

func (c *MessagesController) GetMessagesByPostId(r *http.Request) (any, error) {
	var request dto.GetMessagesByPostIdRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	return c.messagesSrv.GetMessagesByPostId(r.Context(), &request)
}

func (c *MessagesController) UpdateMessageById(r *http.Request) (any, error) {
	var request dto.UpdateMessageByIdRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return nil, err
	}

	request.UserId = user.UserId
	request.MessageId = r.PathValue(web.MsgPathValue)

	return c.messagesSrv.UpdateMessageById(r.Context(), &request)
}

func (c *MessagesController) DeleteMessageById(r *http.Request) (any, error) {
	var request dto.DeleteMessageByIdRequest

	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return nil, err
	}

	request.UserId = user.UserId
	request.MessageId = r.PathValue(web.MsgPathValue)

	return c.messagesSrv.DeleteMessageById(r.Context(), &request)
}
