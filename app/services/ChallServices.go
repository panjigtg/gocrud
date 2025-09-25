package services

import (
    "crudprojectgo/app/models"
    "crudprojectgo/app/repository"
    // "database/sql"
    // "errors"
)

// type ChallServices struct {
//     repo *repository.ChallRepository
// }

// func NewChallServices(repo *repository.ChallRepository) *ChallServices {
//     return &ChallServices{repo: repo}
// }

// func (s *ChallServices) GetAllChall() ([]models.Chall, error) {
//     return s.repo.GetAll()
// }

type ChallServices struct {
    repo *repository.ChallRepository
}

func NewChallServices(repo *repository.ChallRepository) *ChallServices {
    return &ChallServices{repo: repo}
}

func (s *ChallServices) GetAllChall() ([]models.Chall, error) {
    return s.repo.GetAll()
}
