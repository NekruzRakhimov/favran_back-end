package utils

import (
	"favran/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func ConvertSearchParams(c *gin.Context) (params models.SearchParams, err error){
	params.Query = c.Query("query")

	discountStr := c.Query("discount_amount")
	if discountStr == "" {
		params.DiscountAmount = 0
	} else {
		params.DiscountAmount, err = strconv.Atoi(discountStr)
		if err != nil {
			return models.SearchParams{}, err
		}
	}

	priceFrom := c.Query("price_from")
	if priceFrom == "" {
		params.PriceFrom = 0
	} else {
		params.PriceFrom, err = strconv.Atoi(priceFrom)
		if err != nil {
			return models.SearchParams{}, err
		}
	}

	priceTo := c.Query("price_to")
	if priceTo == "" {
		params.PriceTo = 0
	} else {
		params.PriceTo, err = strconv.Atoi(priceTo)
		if err != nil {
			return models.SearchParams{}, err
		}
	}

	return params, nil
}
