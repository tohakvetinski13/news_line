package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"test_task1/store"

	"github.com/sirupsen/logrus"
)

// Handler is http handlers for exchange
type Handler struct {
	storeService store.Service
	logger       *logrus.Logger
}

func New(storeService store.Service, logger *logrus.Logger) *Handler {
	return &Handler{
		storeService: storeService,
		logger:       logger,
	}
}

var OpenCursors map[string]string = make(map[string]string)

func (h *Handler) GetLikes(w http.ResponseWriter, r *http.Request) {

	newsID := r.URL.Query().Get("id")

	likes, err := h.storeService.GetLikesByNewsID(newsID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp := &GetLikeResp{
		Likes: likes,
	}

	res, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(res)
}

func (h *Handler) GetStartNewsPage(w http.ResponseWriter, r *http.Request) {

	userID := r.URL.Query().Get("id")
	count := r.URL.Query().Get("count")

	_, ok := OpenCursors["cursor"+userID]
	if ok {

		err := h.storeService.CloseCursor(userID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = h.storeService.GetCursor(userID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	} else {
		err := h.storeService.GetCursor(userID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		OpenCursors["cursor"+userID] = "cursor"
	}

	news, err := h.storeService.GetFetch(userID, count)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := &GetNewsResp{
		News: news,
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(respJSON)
}

func (h *Handler) GetFetchPage(w http.ResponseWriter, r *http.Request) {

	userID := r.URL.Query().Get("id")
	count := r.URL.Query().Get("count")

	_, ok := OpenCursors["cursor"+userID]
	if !ok {
		err := h.storeService.GetCursor(userID)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		OpenCursors[userID] = "cursor"
	}

	news, err := h.storeService.GetFetch(userID, count)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := &GetNewsResp{
		News: news,
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(respJSON)
}
