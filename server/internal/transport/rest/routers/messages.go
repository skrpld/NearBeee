package routers

import (
	"net/http"

	"github.com/skrpld/NearBeee/internal/core/repository"
	"github.com/skrpld/NearBeee/internal/core/service"
	"github.com/skrpld/NearBeee/internal/transport/rest/handlers"
	"github.com/skrpld/NearBeee/internal/transport/rest/web"
)

func NewMessagesRouter(repo *repository.MongodbRepository) *http.ServeMux {
	srv := service.NewMessagesService(repo)
	controller := handlers.NewMessagesController(srv)
	router := http.NewServeMux()

	router.HandleFunc("POST /messages/", web.Handle(controller.CreateMessage))
	router.HandleFunc("GET /messages/", web.Handle(controller.GetMessage))
	router.HandleFunc("GET /messages/{msg_id}", web.Handle(controller.GetMessage))
	router.HandleFunc("PUT /messages/{msg_id}", web.Handle(controller.UpdateMessageById))
	router.HandleFunc("DELETE /messages/{msg_id}", web.Handle(controller.DeleteMessageById))

	return router
}
