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

func (s FilmService) PatchFilm(input domain.NullableFilm, actorIds []int) (domain.Film, error) {
	return s.repos.PatchFilm(input, actorIds)
}

func NewFilmService(repos repository.Film, log *slog.Logger) *FilmService {
	return &FilmService{repos: repos, log: log}
}

func (s FilmService) CreateFilm(film domain.Film, actorIds []int) (int, error) {
	return s.repos.CreateFilm(film, actorIds)
}

func (s FilmService) DeleteFilm(id int) error {
	//TODO implement me
	panic("implement me")
}

func (s FilmService) UpdateFilm(film domain.Film, actorIds []int) error {
	return s.repos.UpdateFilm(film, actorIds)
}

func (s FilmService) ListFilms(sortBy, sortDir string, actorId int) ([]domain.Film, error) {
	if actorId != -1 {
		return s.repos.ListFilmsByActor(sortBy, sortDir, actorId)
	} else {
		return s.repos.ListFilms(sortBy, sortDir)
	}
}

func (s FilmService) SearchFilm(query string) ([]domain.Film, error) {
	return s.repos.SearchFilm(query)
}
