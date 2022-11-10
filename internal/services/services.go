package services

import "github.com/AdiKhoironHasan/mister-aladin-articles/pkg/dto"

type Services interface {
	CreateArticles(req *dto.ArticleReqDTO) error
	ShowArticles(req *dto.ArticleParamReqDTO) ([]*dto.ArticleResDTO, error)
	ShowArticlesByID(id int) (*dto.ArticleResDTO, error)
	UpdateArticle(req *dto.ArticleReqDTO) error
	DeleteArticle(id int) error
}
