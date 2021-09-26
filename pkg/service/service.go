package service

import (
	"favran/models"
	"favran/pkg/repository"
)

func GetTopCategoriesGoods() (topCategoriesGoods []models.TopCategory, err error) {
	categories, err := repository.GetTopCategoryList()
	if err != nil {
		return nil, err
	}

	topCategoriesGoods, err = repository.MatchTopCategoriesWithGoods(categories)
	if err != nil {
		return nil, err
	}

	return topCategoriesGoods, nil
}

func GetCategoryGoods(categoryId int, pagingParams models.PagingParams) (goods []models.Good, err error) {
	return  repository.GetCategoryGoods(categoryId, pagingParams)
}

func AddViewToTheGood(goodId int, ip string) error {
	exists, err := repository.IsInViewedIps(goodId, ip)
	if err != nil {
		return err
	}

	if exists {
		return nil
	} else {
		return repository.AddViewToTheGood(goodId, ip)
	}
}

func Search(params models.SearchParams, pagingParams  models.PagingParams) (goods []models.Good, err error) {
	return repository.Search(params, pagingParams)
}