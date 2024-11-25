package user

import (
	"context"

	"github.com/kleo-53/music-system/internal/controller/model"
	"github.com/kleo-53/music-system/internal/core"
)

type service struct {
	songStore core.SongStore
}

func New(store core.SongStore) core.SongService {
	return &service{
		songStore: store,
	}
}

func (s *service) GetSongsInfo(ctx context.Context, filters model.SongFilters, page, pageSize int) ([]model.Song, error) {
	return s.songStore.GetSongsInfo(ctx, filters, page, pageSize)
}

func (s *service) GetSongText(ctx context.Context, id, page, pageSize int) ([]string, error) {
	return s.songStore.GetSongText(ctx, id, page, pageSize)
}

func (s *service) DeleteSong(ctx context.Context, id int) error {
	return s.songStore.DeleteSong(ctx, id)
}

func (s *service) UpdateSong(ctx context.Context, id int, newData model.SongFilters) error {
	return s.songStore.UpdateSong(ctx, id, newData)
}

func (s *service) CreateSong(ctx context.Context, song model.SongCommon, details model.SongDetail) error {
	return s.songStore.CreateSong(ctx, song, details)
}
