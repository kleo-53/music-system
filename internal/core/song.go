package core

import (
	"context"

	"github.com/kleo-53/music-system/internal/controller/model"
)

type (
	Song struct {
		ID          int    `gorm:"column:id;primaryKey"`
		Group       string `gorm:"column:song_group"`
		Song        string `gorm:"column:song"`
		Text        string `gorm:"column:song_text"`
		ReleaseDate string `gorm:"column:release_date"`
		Link        string `gorm:"column:link"`
	}

	SongStore interface {
		CreateSong(ctx context.Context, song model.SongCommon, details model.SongDetail) error
		UpdateSong(ctx context.Context, id int, newData model.SongFilters) error
		DeleteSong(ctx context.Context, id int) error
		GetSongText(ctx context.Context, id, page, pageSize int) ([]string, error)
		GetSongsInfo(ctx context.Context, filters model.SongFilters, page, pageSize int) ([]model.Song, error)
	}

	SongService interface {
		CreateSong(ctx context.Context, song model.SongCommon, details model.SongDetail) error
		UpdateSong(ctx context.Context, id int, newData model.SongFilters) error
		DeleteSong(ctx context.Context, id int) error
		GetSongText(ctx context.Context, id, page, pageSize int) ([]string, error)
		GetSongsInfo(ctx context.Context, filters model.SongFilters, page, pageSize int) ([]model.Song, error)
	}
)

func (Song) TableName() string {
	return "songs"
}
