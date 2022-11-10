package services

import (
	"log"

	integ "github.com/AdiKhoironHasan/mister-aladin-articles/internal/integration"
	"github.com/AdiKhoironHasan/mister-aladin-articles/internal/repository"
	"github.com/AdiKhoironHasan/mister-aladin-articles/pkg/dto"
	"github.com/AdiKhoironHasan/mister-aladin-articles/pkg/dto/assembler"
)

type service struct {
	sqlRepo   repository.SqlRepository
	noSqlRepo repository.NoSqlRepository
	IntegServ integ.IntegServices
}

func NewService(sqlRepo repository.SqlRepository, noSqlRepo repository.NoSqlRepository, IntegServ integ.IntegServices) Services {
	return &service{sqlRepo, noSqlRepo, IntegServ}
}

func (s *service) CreateArticles(req *dto.ArticleReqDTO) error {
	dataArticleModel := assembler.ToCreateArticle(req)

	idArticle, err := s.sqlRepo.CreateArticles(dataArticleModel)

	if err != nil {
		return err
	}

	err = s.noSqlRepo.CreateArticles(idArticle, req)
	if err != nil {
		log.Println("Redis :", err)
	}

	var where string
	dataArticlesModels, err := s.sqlRepo.ShowArticles(where)
	if err != nil {
		return err
	}

	dataArticles := assembler.ToCreateArticleResponse(dataArticlesModels)

	err = s.noSqlRepo.CreateAllArticles(dataArticles)
	if err != nil {
		log.Println("Redis :", err)
	}

	return nil
}

func (s *service) ShowArticles(req *dto.ArticleParamReqDTO) ([]*dto.ArticleResDTO, error) {
	var dataArticles []*dto.ArticleResDTO
	var where string

	if req.Query == "" && req.Author == "" {
		dataArticles, err := s.noSqlRepo.ShowArticles()
		if err != nil {
			log.Println("Redis :", err)
		}

		if len(dataArticles) > 0 {
			log.Println("ShowArticles use Redis")
			return dataArticles, nil
		}
	}

	if req.Query != "" && req.Author != "" {
		where = "title LIKE '%" + req.Query + "%' AND body LIKE '%" + req.Query + "%' AND author LIKE '%" + req.Author + "%'"
	} else if req.Query != "" {
		where = "title LIKE '%" + req.Query + "%' AND body LIKE '%" + req.Query + "%'"
	} else if req.Author != "" {
		where = "author LIKE '%" + req.Author + "$'"
	}

	log.Println("ShowArticles use PostgreSQL")

	dataArticlesModels, err := s.sqlRepo.ShowArticles(where)
	if err != nil {
		return nil, err
	}

	dataArticles = assembler.ToShowArticlesResponse(dataArticlesModels)

	return dataArticles, nil
}

func (s *service) ShowArticlesByID(id int) (*dto.ArticleResDTO, error) {
	var dataArticle *dto.ArticleResDTO

	dataArticle, err := s.noSqlRepo.ShowArticlesByID(id)
	if err != nil {
		log.Println(err)
	}

	if dataArticle != nil {
		log.Println("ShowArticlesByID use Redis")
		return dataArticle, nil
	}

	log.Println("ShowArticlesByID use PostgreSQL")
	dataArticleModels, err := s.sqlRepo.ShowArticlesByID(id)

	if err != nil {
		return nil, err
	}

	if dataArticleModels == nil {
		return nil, nil
	}

	dataArticle = assembler.ToShowArticlesByIDResponse(dataArticleModels)

	return dataArticle, nil
}

func (s *service) UpdateArticle(req *dto.ArticleReqDTO) error {
	dataArticle := assembler.ToUpdateArticle(req)

	created, err := s.sqlRepo.UpdateArticle(dataArticle)
	if err != nil {
		return err
	}
	req.Created = created
	err = s.noSqlRepo.CreateArticles(req.ID, req)
	if err != nil {
		log.Println("Redis :", err)
	}

	var where string
	dataArticlesModels, err := s.sqlRepo.ShowArticles(where)
	if err != nil {
		return err
	}

	dataArticles := assembler.ToShowArticlesResponse(dataArticlesModels)

	err = s.noSqlRepo.CreateAllArticles(dataArticles)
	if err != nil {
		log.Println("Redis :", err)
	}

	return nil
}

func (s *service) DeleteArticle(id int) error {

	err := s.sqlRepo.DeleteArticle(id)
	if err != nil {
		return err
	}

	return nil
}
