package dto

import (
	"github.com/AdiKhoironHasan/mister-aladin-articles/pkg/common/validator"
)

type ArticleReqDTO struct {
	ID      int
	Author  string `json:"author" valid:"required" validname:"author"`
	Title   string `json:"title" valid:"required" validname:"title"`
	Body    string `json:"body" valid:"required" validname:"body"`
	Created string
}

func (dto *ArticleReqDTO) Validate() error {
	v := validator.NewValidate(dto)

	return v.Validate()
}

type ArticleParamReqDTO struct {
	Query  string `json:"query" validname:"query" query:"query"`
	Author string `json:"author" validname:"author" query:"author"`
}

func (dto *ArticleParamReqDTO) Validate() error {
	v := validator.NewValidate(dto)

	return v.Validate()
}
