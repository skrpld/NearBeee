package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/skrpld/NearBeee/internal/core/models/dto"
	"github.com/skrpld/NearBeee/internal/transport/rest/web"
	"github.com/skrpld/NearBeee/pkg/errors"
)

type TopicsService interface {
	CreateTopic(rows *dto.CreateTopicRequest) (*dto.CreateTopicResponse, error)
	GetTopicsByUserId(rows *dto.GetTopicsByUserIdRequest) (*dto.GetTopicsByUserIdResponse, error)
	GetTopicsByLocation(rows *dto.GetTopicsByLocationRequest) (*dto.GetTopicsByLocationResponse, error)
	GetTopicById(rows *dto.GetTopicByTopicIdRequest) (*dto.GetTopicByTopicIdResponse, error)
	UpdateTopicById(rows *dto.UpdateTopicByIdRequest) (*dto.UpdateTopicByIdResponse, error)
	DeleteTopicById(rows *dto.DeleteTopicByIdRequest) (*dto.DeleteTopicResponse, error)
}
type TopicsController struct {
	topicsSrv TopicsService
}

func NewTopicsController(topicsSrv TopicsService) *TopicsController {
	return &TopicsController{topicsSrv: topicsSrv}
}

func (c *TopicsController) CreateTopicHandler(r *http.Request) (any, error) {
	var request dto.CreateTopicRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return nil, err
	}

	request.UserId = user.UserId

	return c.topicsSrv.CreateTopic(&request)
}

func (c *TopicsController) GetTopics(r *http.Request) (any, error) {
	switch web.FormType(r.FormValue(web.FormValue)) {
	case web.UserForm:
		return c.GetTopicsByUserId(r)
	case web.LocationForm:
		return c.GetTopicsByLocation(r)
	case web.TopicForm, web.NullForm:
		return c.GetTopicById(r)
	default:
		return nil, errors.ErrInvalidFormType
	}
}

func (c *TopicsController) GetTopicsByUserId(r *http.Request) (any, error) {
	var request dto.GetTopicsByUserIdRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return nil, err
	}

	request.UserId = user.UserId

	return c.topicsSrv.GetTopicsByUserId(&request)
}

func (c *TopicsController) GetTopicsByLocation(r *http.Request) (any, error) {
	var request dto.GetTopicsByLocationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	return c.topicsSrv.GetTopicsByLocation(&request)
}

func (c *TopicsController) GetTopicById(r *http.Request) (any, error) {
	var request dto.GetTopicByTopicIdRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	request.TopicId = r.PathValue(web.TopicPathValue)

	return c.topicsSrv.GetTopicById(&request)
}

func (c *TopicsController) UpdateTopicById(r *http.Request) (any, error) {
	var request dto.UpdateTopicByIdRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	request.TopicId = r.PathValue(web.TopicPathValue)

	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return nil, err
	}
	request.UserId = user.UserId

	return c.topicsSrv.UpdateTopicById(&request)
}

func (c *TopicsController) DeleteTopicById(r *http.Request) (any, error) {
	var request dto.DeleteTopicByIdRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	request.TopicId = r.PathValue(web.TopicPathValue)

	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return nil, err
	}
	request.UserId = user.UserId

	return c.topicsSrv.DeleteTopicById(&request)
}
