package repository

import (
	// "github.com/jmoiron/sqlx"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/AdiKhoironHasan/mister-aladin-articles/internal/repository"
	"github.com/AdiKhoironHasan/mister-aladin-articles/pkg/dto"
	servErrors "github.com/AdiKhoironHasan/mister-aladin-articles/pkg/errors"
	"github.com/go-redis/redis/v8"
)

const ()

type RedisNoSqlRepo struct {
	Conn *redis.Client
}

func NewRepo(Conn *redis.Client) repository.NoSqlRepository {

	repo := &RedisNoSqlRepo{Conn}
	return repo
}

func (p *RedisNoSqlRepo) CreateArticles(id int, req *dto.ArticleReqDTO) error {
	key := fmt.Sprintf("article_%d", id)
	data := dto.ArticleResDTO{
		ID:      id,
		Author:  req.Author,
		Title:   req.Title,
		Body:    req.Body,
		Created: req.Created,
	}

	value, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed Marshal CreateArticles:", err.Error())
		return err
	}

	// ttl := time.Duration(600) * time.Second
	redis := p.Conn.Set(context.Background(), key, value, 0)

	result, err := redis.Result()
	if err != nil {
		log.Println("Failed Redis Get Result CreateArticles:", err.Error())
		return fmt.Errorf(servErrors.ErrorDB)
	}

	if result != "OK" {
		log.Println("Failed Redis Not OK CreateArticles: ", result)
		return fmt.Errorf(result)
	}

	return nil
}

func (p *RedisNoSqlRepo) CreateAllArticles(articles []*dto.ArticleResDTO) error {
	key := "article_all"
	value, err := json.Marshal(articles)
	if err != nil {
		log.Println("Failed Marshal CreateAllArticles:", err.Error())
		return err
	}

	// ttl := time.Duration(600) * time.Second
	redis := p.Conn.Set(context.Background(), key, value, 0)

	result, err := redis.Result()
	if err != nil {
		log.Println("Failed Redis Get Result CreateAllArticles:", err.Error())
		return fmt.Errorf(servErrors.ErrorDB)
	}

	if result != "OK" {
		log.Println("Failed Redis Not OK CreateAllArticles: ", result)
		return fmt.Errorf(result)
	}

	return nil
}

func (p *RedisNoSqlRepo) ShowArticles() ([]*dto.ArticleResDTO, error) {
	var dataArticles []*dto.ArticleResDTO
	key := "article_all"

	result := p.Conn.Get(context.Background(), key).Val()
	if result == "" {
		log.Println("Redis: No Result for key 'article_all' ShowArticles")
		return nil, nil
	}

	jsonData := []byte(result)

	err := json.Unmarshal(jsonData, &dataArticles)
	if err != nil {
		log.Println("Failed UnMarshal ShowArticles:", err.Error())
		return nil, err
	}

	return dataArticles, nil
}

func (p *RedisNoSqlRepo) ShowArticlesByID(id int) (*dto.ArticleResDTO, error) {
	var dataArticles *dto.ArticleResDTO
	key := fmt.Sprintf("article_%d", id)

	result := p.Conn.Get(context.Background(), key).Val()
	if result == "" {
		log.Println(fmt.Sprintf("Redis: No Result for key '%s' ShowArticlesByID", key))
		return nil, nil
	}

	jsonData := []byte(result)

	err := json.Unmarshal(jsonData, &dataArticles)
	if err != nil {
		log.Println("Failed UnMarshal ShowArticlesByID:", err.Error())
		return nil, err
	}

	return dataArticles, nil
}
