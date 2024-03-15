package service

import (
	"github.com/Warh40k/vk-intern-filmotecka/internal/api/repository"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"log/slog"
)

type FilmService struct {
	repos repository.Film
	log   *slog.Logger
}

func (s FilmService) PatchFilm(actor domain.Actor) error {
	//TODO implement me
	panic("implement me")
}

func NewFilmService(repos repository.Film, log *slog.Logger) *FilmService {
	return &FilmService{repos: repos, log: log}
}

func (s FilmService) CreateFilm(actor domain.Actor) error {
	//TODO implement me
	panic("implement me")
}

func (s FilmService) DeleteFilm(id int) error {
	//TODO implement me
	panic("implement me")
}

func (s FilmService) UpdateFilm(actor domain.Actor) error {
	//TODO implement me
	panic("implement me")
}

func (s FilmService) ListFilms() ([]domain.Actor, error) {
	//TODO implement me
	panic("implement me")
}

func (s FilmService) SearchFilm(film string, actor string) ([]domain.Film, error) {
	//TODO implement me
	panic("implement me")
}
