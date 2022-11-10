package repository

import (
	"github.com/AdiKhoironHasan/mister-aladin-articles/internal/models"
	"github.com/AdiKhoironHasan/mister-aladin-articles/pkg/dto"
)

type SqlRepository interface {
	CreateArticles(dataArticle *models.ArticleModels) (int, error)
	ShowArticles(where string) ([]*models.ArticleModels, error)
	ShowArticlesByID(id int) ([]*models.ArticleModels, error)
	UpdateArticle(dataArticle *models.ArticleModels) (string, error)
	DeleteArticle(id int) error
}

type NoSqlRepository interface {
	CreateArticles(id int, req *dto.ArticleReqDTO) error
	CreateAllArticles(articles []*dto.ArticleResDTO) error
	ShowArticles() ([]*dto.ArticleResDTO, error)
	ShowArticlesByID(id int) (*dto.ArticleResDTO, error)
}
