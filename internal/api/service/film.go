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

func (s FilmService) PatchFilm(film domain.FilmInput) error {
	//TODO implement me
	panic("implement me")
}

func NewFilmService(repos repository.Film, log *slog.Logger) *FilmService {
	return &FilmService{repos: repos, log: log}
}

func (s FilmService) CreateFilm(film domain.Film) (int, error) {
	return s.repos.CreateFilm(film)
}

func (s FilmService) DeleteFilm(id int) error {
	//TODO implement me
	panic("implement me")
}

func (s FilmService) UpdateFilm(film domain.Film) error {
	//TODO implement me
	panic("implement me")
}

func (s FilmService) ListFilms(sortBy, sortDir string) ([]domain.Film, error) {
	return s.repos.ListFilms(sortBy, sortDir)
}

func (s FilmService) SearchFilm(film string, actor string) ([]domain.Film, error) {
	//TODO implement me
	panic("implement me")
}
