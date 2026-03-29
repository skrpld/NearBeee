package routers

import (
	"net/http"

	"github.com/skrpld/NearBeee/internal/core/repository"
	"github.com/skrpld/NearBeee/internal/core/service"
	"github.com/skrpld/NearBeee/internal/transport/rest/handlers"
	"github.com/skrpld/NearBeee/internal/transport/rest/web"
)

func NewTopicsRouter(repo *repository.TopicsRepository) *http.ServeMux {
	srv := service.NewTopicsService(repo)
	controller := handlers.NewTopicsController(srv)
	router := http.NewServeMux()

	router.HandleFunc("POST /topics/", web.Handle(controller.CreateTopicHandler))
	router.HandleFunc("GET /topics/", web.Handle(controller.GetTopics))
	router.HandleFunc("GET /topics/{topic_id}", web.Handle(controller.GetTopics))
	router.HandleFunc("PUT /topics/{topic_id}", web.Handle(controller.UpdateTopicById))
	router.HandleFunc("DELETE /topics/{topic_id}", web.Handle(controller.DeleteTopicById))

	return router
}
