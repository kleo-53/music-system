package user

import (
	"context"
	"strings"

	"github.com/kleo-53/music-system/internal/controller/model"
	"github.com/kleo-53/music-system/internal/core"
	"github.com/kleo-53/music-system/pkg/postgres"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.SongStore {
	return &store{pg}
}

func convertToModelSong(song core.Song) model.Song {
	return model.Song{
		Song:        song.Song,
		Group:       song.Group,
		Text:        song.Text,
		ReleaseDate: song.ReleaseDate,
		Link:        song.Link,
	}
}

func (s *store) CreateSong(ctx context.Context, song model.SongCommon, details model.SongDetail) error {
	songToAdd := core.Song{
		Group: song.Group,
		Song:  song.Song,
	}
	if details.Text != "" {
		songToAdd.Text = details.Text
	}
	if details.ReleaseDate != "" {
		songToAdd.ReleaseDate = details.ReleaseDate
	}
	if details.Link != "" {
		songToAdd.Link = details.Link
	}
	return s.DB.WithContext(ctx).Create(&songToAdd).Error
}

func (s *store) UpdateSong(ctx context.Context, id int, newData model.SongFilters) error {
	var err error
	if newData.Text != "" {
		err = s.DB.WithContext(ctx).Model(&core.Song{}).Where("id = ?", id).Update("song_text", &newData.Text).Error
		if err != nil {
			return err
		}
	}
	if newData.Link != "" {
		err = s.DB.WithContext(ctx).Model(&core.Song{}).Where("id = ?", id).Update("link", &newData.Link).Error
		if err != nil {
			return err
		}
	}
	if newData.ReleaseDate != "" {
		err = s.DB.WithContext(ctx).Model(&core.Song{}).Where("id = ?", id).Update("release_date", &newData.ReleaseDate).Error
		if err != nil {
			return err
		}
	}
	if newData.Song != "" {
		err = s.DB.WithContext(ctx).Model(&core.Song{}).Where("id = ?", id).Update("song", &newData.Song).Error
		if err != nil {
			return err
		}
	}
	if newData.Group != "" {
		err = s.DB.WithContext(ctx).Model(&core.Song{}).Where("id = ?", id).Update("song_group", &newData.Group).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *store) DeleteSong(ctx context.Context, id int) error {
	return s.DB.WithContext(ctx).Delete(&core.Song{}, "id = ?", id).Error
}

func (s *store) GetSongText(ctx context.Context, id, page, pageSize int) ([]string, error) {
	var song core.Song
	if err := s.DB.WithContext(ctx).
		Model(core.Song{}).
		Where("id = ?", id).
		First(&song).Error; err != nil {
		return []string{}, err
	}
	couplets := strings.Split(song.Text, "\n\n")
	start := (page - 1) * pageSize
	end := min(len(couplets), start+pageSize)
	if start > len(couplets) {
		return []string{}, nil
	}
	return couplets[start:end], nil
}

func (s *store) GetSongsInfo(ctx context.Context, filters model.SongFilters, page, pageSize int) ([]model.Song, error) {
	var songs []core.Song
	query := s.DB.WithContext(ctx).Model(&core.Song{})
	if group := filters.Group; group != "" {
		query = query.Where("song_group like ?", "%"+group+"%")
	}
	if song := filters.Song; song != "" {
		query = query.Where("song LIKE ?", "%"+song+"%")
	}
	if text := filters.Text; text != "" {
		query = query.Where("song_text LIKE ?", "%"+text+"%")
	}
	if releaseDate := filters.ReleaseDate; releaseDate != "" {
		query = query.Where("release_date LIKE ?", "%"+releaseDate+"%")
	}
	if link := filters.Link; link != "" {
		query = query.Where("link LIKE ?", "%"+link+"%")
	}
	if err := query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&songs).Error; err != nil {
		return []model.Song{}, err
	}
	responce := []model.Song{}
	for _, song := range songs {
		responce = append(responce, convertToModelSong(song))
	}
	return responce, nil
}
