package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kleo-53/music-system/internal/controller/model"
	"github.com/kleo-53/music-system/pkg/logger"
)

type SongFilters struct {
	Group       string `json:"group,omitempty"`
	Song        string `json:"song,omitempty"`
	Text        string `json:"text,omitempty"`
	ReleaseDate string `json:"release_date,omitempty"`
	Link        string `json:"link,omitempty"`
}

// getSongsInfo godoc
//
// @Summary		Get songs info
// @Description	Get info about all songs with pagination and optional filters
// @Tags		songs
// @Accept		json
// @Produce		json
// @Param		group			query		string		false  "Filter by group name"
// @Param		song			query		string		false  "Filter by song name"
// @Param		text			query		string		false  "Filter by text content"
// @Param		release_date	query		string		false  "Filter by release date"
// @Param		link			query		string		false  "Filter by link"
// @Param		page			query		int			false	"Page number" 				default(1)
// @Param		page_size		query		int 		false 	"Number of songs per page" 	default(10)
// @Success		200				{object} 	[]model.SongCommon
// @Failure		400				{object} 	map[string]string 	"Invalid request payload"
// @Failure		500				{object} 	map[string]string 	"Failed to get any songs data"
// @Router		/api/v1/songs [get]
func (ro *Router) getSongsInfo(w http.ResponseWriter, r *http.Request) {
	filters := model.SongFilters{
		Group:       r.URL.Query().Get("group"),
		Song:        r.URL.Query().Get("song"),
		Text:        r.URL.Query().Get("text"),
		ReleaseDate: r.URL.Query().Get("release_date"),
		Link:        r.URL.Query().Get("link"),
	}
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	page_int, err := strconv.ParseInt(page, 10, strconv.IntSize)
	if err != nil {
		logger.Log().Error(r.Context(), "Failed to update song info: invalid page provided")
		JSONError(r.Context(), w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	pageSize := r.URL.Query().Get("page_size")
	if pageSize == "" {
		pageSize = "10"
	}
	page_size_int, err := strconv.ParseInt(pageSize, 10, strconv.IntSize)
	if err != nil {
		logger.Log().Error(r.Context(), "Failed to update song info: invalid page size provided")
		JSONError(r.Context(), w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	songs, err := ro.songService.GetSongsInfo(r.Context(), filters, int(page_int), int(page_size_int))
	if err != nil {
		logger.Log().Error(r.Context(), "Failed to get any songs data: "+err.Error())
		JSONError(r.Context(), w, http.StatusInternalServerError, "Failed to get any songs data")
		return
	}

	logger.Log().Info(r.Context(), "Get songs data")
	JSONResponse(r.Context(), w, http.StatusCreated, songs)
}

// @Summary 	Get song text
// @Description	Get text of song by ID with pagination
// @Tags 		songs
// @Accept 		json
// @Produce 	json
// @Param 		song_id 	path 		int 				true 	"Song ID"
// @Param 		page 		query 		int 				false 	"Page number" 				default(1)
// @Param 		page_size 	query 		int 				false 	"Number of verses per page"	default(10)
// @Success 	200 		{object} 	[]string
// @Failure 	400 		{object} 	map[string]string 	"Invalid request payload"
// @Failure 	500 		{object} 	map[string]string 	"Failed to get song text"
// @Router 		/api/v1/songs/{song_id} [get]
func (ro *Router) getSongText(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	songID := vars["song_id"]
	if songID == "" {
		logger.Log().Error(r.Context(), "Failed to get song text: no id provided")
		JSONError(r.Context(), w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	song_id, err := strconv.ParseInt(songID, 10, strconv.IntSize)
	if err != nil {
		logger.Log().Error(r.Context(), "Failed to update song info: invalid id provided")
		JSONError(r.Context(), w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	page_int, err := strconv.ParseInt(page, 10, strconv.IntSize)
	if err != nil {
		logger.Log().Error(r.Context(), "Failed to update song info: invalid page provided")
		JSONError(r.Context(), w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	pageSize := r.URL.Query().Get("page_size")
	if pageSize == "" {
		pageSize = "10"
	}
	page_size_int, err := strconv.ParseInt(pageSize, 10, strconv.IntSize)
	if err != nil {
		logger.Log().Error(r.Context(), "Failed to update song info: invalid page size provided")
		JSONError(r.Context(), w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	text, err := ro.songService.GetSongText(r.Context(), int(song_id), int(page_int), int(page_size_int))
	if err != nil {
		logger.Log().Error(r.Context(), "Failed to get song text: "+err.Error())
		JSONError(r.Context(), w, http.StatusInternalServerError, "Failed to get song text")
		return
	}
	logger.Log().Info(r.Context(), "Get song text")
	JSONResponse(r.Context(), w, http.StatusCreated, text)
}

// @Summary 	Delete song
// @Description	Delete a song by ID
// @Tags 		songs
// @Produce 	json
// @Param 		song_id path 		int 				true 				"Song ID"
// @Success		200 	{object} 	map[string]string 	"Song was deleted"
// @Failure 	400 	{object} 	map[string]string 	"Invalid request payload"
// @Failure 	500 	{object} 	map[string]string 	"Failed to delete song"
// @Router 		/api/v1/songs/{song_id} [delete]
func (ro *Router) deleteSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	songID := vars["song_id"]
	if songID == "" {
		logger.Log().Error(r.Context(), "Failed to get song text: no id provided")
		JSONError(r.Context(), w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	song_id, err := strconv.ParseInt(songID, 10, strconv.IntSize)
	if err != nil {
		logger.Log().Error(r.Context(), "Failed to update song info: invalid id provided")
		JSONError(r.Context(), w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	err = ro.songService.DeleteSong(r.Context(), int(song_id))
	if err != nil {
		logger.Log().Error(r.Context(), "Failed to delete song: "+err.Error())
		JSONError(r.Context(), w, http.StatusInternalServerError, "Failed to delete song")
		return
	}
	logger.Log().Info(r.Context(), "Song deleted successfully")
	JSONResponse(r.Context(), w, http.StatusOK, map[string]string{"message": "Song was deleted"})
}

// @Summary 	Update song
// @Description	Update song information by ID
// @Tags 		songs
// @Accept 		json
// @Produce 	json
// @Param 		song_id path 		int 				true 					"Song ID"
// @Param 		body 	body 		model.SongFilters 	true 					"Updated song data"
// @Success 	200 	{object} 	map[string]string 	"Song was updated"
// @Failure 	400 	{object} 	map[string]string 	"Invalid request payload"
// @Failure 	500 	{object} 	map[string]string 	"Failed to update song info"
// @Router 		/api/v1/songs/{song_id} [patch]
func (ro *Router) updateSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	songID := vars["song_id"]
	if songID == "" {
		logger.Log().Error(r.Context(), "Failed to update song info: no id provided")
		JSONError(r.Context(), w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	song_id, err := strconv.ParseInt(songID, 10, strconv.IntSize)
	if err != nil {
		logger.Log().Error(r.Context(), "Failed to update song info: invalid id provided")
		JSONError(r.Context(), w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	var newSongData model.SongFilters
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newSongData)
	if err != nil {
		logger.Log().Error(r.Context(), "Failed to update song info: "+err.Error())
		JSONError(r.Context(), w, http.StatusInternalServerError, "Failed to update song info")
	}
	err = ro.songService.UpdateSong(r.Context(), int(song_id), newSongData)
	if err != nil {
		logger.Log().Error(r.Context(), "Failed to update song info: "+err.Error())
		JSONError(r.Context(), w, http.StatusInternalServerError, "Failed to update song info")
		return
	}
	logger.Log().Info(r.Context(), "Song updated successfully")
	JSONResponse(r.Context(), w, http.StatusOK, map[string]string{"message": "Song was updated"})
}

// @Summary 	Add song
// @Description	Add a new song to the system
// @Tags 		songs
// @Accept 		json
// @Produce 	json
// @Param 		body 	body 		model.SongCommon 	true 				"New song data"
// @Success 	200 	{object} 	map[string]string 	"Song was added"
// @Failure 	400 	{object} 	map[string]string 	"Invalid request payload"
// @Failure 	500 	{object} 	map[string]string 	"Failed to add song"
// @Router 		/api/v1/songs [post]
func (ro *Router) addSong(w http.ResponseWriter, r *http.Request) {
	var req model.SongCommon
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log().Error(r.Context(), "Failed to add song: invalid request payload")
		JSONError(r.Context(), w, http.StatusBadRequest, "Failed to add song")
		return
	}
	if req.Group == "" || req.Song == "" {
		logger.Log().Error(r.Context(), "Failed to add song: invalid input")
		JSONError(r.Context(), w, http.StatusBadRequest, "Invalid input")
		return
	}
	details, err := getSongDetails(r, req.Group, req.Song)
	if err != nil {
		logger.Log().Error(r.Context(), "Failed to get song details:"+err.Error())
		details = model.SongDetail{}
		// JSONError(r.Context(), w, http.StatusNotImplemented, "Failed to get song details")
		// return
	}

	err = ro.songService.CreateSong(r.Context(), req, details)
	if err != nil {
		logger.Log().Error(r.Context(), "Failed to add song: "+err.Error())
		JSONError(r.Context(), w, http.StatusInternalServerError, "Failed to add song")
		return
	}
	logger.Log().Info(r.Context(), "Song was added")
	JSONResponse(r.Context(), w, http.StatusOK, map[string]string{"message": "Song was added"})
}

func getSongDetails(r *http.Request, group, song string) (model.SongDetail, error) {
	if err := godotenv.Load("config.env"); err != nil {
		logger.Log().Warn(r.Context(), "No .env file found, using environment variables")
	}
	fmt.Printf("%s?group=%s&song=%s", os.Getenv("EXTERNAL_API_URL"), group, song)
	url := fmt.Sprintf("%s?group=%s&song=%s", os.Getenv("EXTERNAL_API_URL"), group, song)
	resp, err := http.Get(url)
	if err != nil {
		return model.SongDetail{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.SongDetail{}, fmt.Errorf("bad request to exxternal api: %s", resp.Status)
	}

	var details model.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return model.SongDetail{}, err
	}

	return details, nil
}
