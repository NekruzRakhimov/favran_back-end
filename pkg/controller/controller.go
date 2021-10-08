package controller

import (
	"fmt"
	"github.com/NekruzRakhimov/favran/models"
	"github.com/NekruzRakhimov/favran/pkg/repository"
	"github.com/NekruzRakhimov/favran/pkg/service"
	"github.com/NekruzRakhimov/favran/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetGoodInfo(c *gin.Context) {
	var (
		images    []string
		goodIdStr = c.Param("id")
		goodId    int
	)
	goodId, err := strconv.Atoi(goodIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	images, err = repository.GetImagesOfGood(goodId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"images": images})

}

func GetMainPageContent(c *gin.Context) {

	categories, err := repository.GetCategoryList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	if categories == nil {
		categories = []models.Category{}
	}

	banners, err := repository.GetAllBanners()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	if banners == nil {
		banners = []models.Banners{}
	}

	topCategoriesGoods, err := service.GetTopCategoriesGoods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	if topCategoriesGoods == nil {
		topCategoriesGoods = []models.TopCategory{}
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"categories":    categories,
			"banners":       banners,
			"topCategories": topCategoriesGoods,
		})

}

func GetCategoryGoods(c *gin.Context) {
	var (
		categoryIdStr = c.Param("category_id")
		pagingParams  models.PagingParams
		categoryId    int
		err           error
	)

	if c.Query("page") == "" {
		pagingParams.Page = 0
	} else {
		pagingParams.Page, err = strconv.Atoi(c.Query("page"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
			return
		}
	}

	if c.Query("limit") == "" {
		pagingParams.Limit = 0
	} else {
		pagingParams.Limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
			return
		}
	}

	categoryId, err = strconv.Atoi(categoryIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	goods, err := service.GetCategoryGoods(categoryId, pagingParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"goods": goods})
}

func GetImage(c *gin.Context) {
	imageName := c.Param("image_name")
	f1 := c.Param("f1") // first folder
	f2 := c.Param("f2") // second folder

	filePath := fmt.Sprintf("files/%s/%s/%s", f1, f2, imageName)
	fmt.Println(">>>>>>>>>>>>>>>>>>>", filePath)

	//c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename =% s", filename)) // fmt.Sprintf ("attachment; filename =% s", filename) Переименуйте загруженный файл
	//c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Header("Content-Type", "image/png")
	c.File(filePath)
}

func AddViewToTheGood(c *gin.Context) {
	var (
		goodIdStr = c.Param("good_id")
		goodId    int
	)

	goodId, err := strconv.Atoi(goodIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	if err := service.AddViewToTheGood(goodId, c.ClientIP()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reason": "кол-во просмотров данного товара увеличено на +1"})
}

func Search(c *gin.Context) {
	var (
		pagingParams models.PagingParams
		err          error
	)

	if c.Query("page") == "" {
		pagingParams.Page = 0
	} else {
		pagingParams.Page, err = strconv.Atoi(c.Query("page"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
			return
		}
	}

	if c.Query("limit") == "" {
		pagingParams.Limit = 0
	} else {
		pagingParams.Limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
			return
		}
	}

	params, err := utils.ConvertSearchParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	goods, err := service.Search(params, pagingParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"goods": goods})
}
