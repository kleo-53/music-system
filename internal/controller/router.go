package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kleo-53/music-system/internal/core"
	"github.com/kleo-53/music-system/pkg/logger"
)

type Router struct {
	app         *mux.Router
	songService core.SongService
}

func NewRouter(
	app *mux.Router,
	songService core.SongService,
) *Router {
	router := &Router{
		app:         app,
		songService: songService,
	}
	router.initRoutes()
	return router
}

// func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
// 	r.app.ServeHTTP(w, req)
// }

func (r *Router) initRoutes() {

	s := r.app.PathPrefix("/api/v1").Subrouter()

	s.HandleFunc("/songs", r.getSongsInfo).Methods("GET")         // Получение данных библиотеки с фильтрацией по всем полям и пагинацией
	s.HandleFunc("/songs" , r.addSong).Methods("POST")           // Добавление новой песни в формате	JSON
	s.HandleFunc("/songs/{song_id}", r.getSongText).Methods("GET") // Получение текста песни с пагинацией по куплетам
	s.HandleFunc("/songs/{song_id}", r.updateSong).Methods("PATCH")         // Изменение данных песни
	s.HandleFunc("/songs/{song_id}", r.deleteSong).Methods("DELETE")        // Удаление песни

}

// func (r *Router) initRequestMiddlewares() {
// 	// r.mux.Use(logger.New())
// }

func JSONResponse(ctx context.Context, w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			logger.Log().Error(ctx, "Failed to encode response: "+err.Error())
		}
	}
}

func JSONError(ctx context.Context, w http.ResponseWriter, status int, message string) {
	JSONResponse(ctx, w, status, map[string]string{"error": message})
}
