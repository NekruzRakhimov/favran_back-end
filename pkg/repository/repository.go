package repository

import (
	"favran/db"
	"favran/models"
	"fmt"
	"github.com/jinzhu/gorm"
)

func GetAllBanners() (banners []models.Banners, err error) {
	sqlQuery := "SELECT id, image_name FROM banners ORDER BY id"

	return banners, db.GetDBConn().Table("banners").Raw(sqlQuery).Scan(&banners).Error
}

func GetCategoryList() (categories []models.Category, err error) {
	sqlQuery := "SELECT id, name FROM categories ORDER BY id"

	return categories, db.GetDBConn().Table("categories").Raw(sqlQuery).Scan(&categories).Error
}

func GetCategoryGoods(categoryId int, pagingParams models.PagingParams) (goods []models.Good, err error) {
	var (
		sqlQuery = "SELECT goods.*, shops.shop_name, shops.shop_phone, shops.shop_address FROM goods, shops WHERE category_id = ? AND shop_id = shops.id AND is_active = true ORDER BY goods.id"
		pagingQuery string
	)

	if pagingParams.Page > 0 && pagingParams.Limit > 0 {
		pagingQuery = fmt.Sprintf(" OFFSET %d LIMIT %d", (pagingParams.Page - 1) * pagingParams.Limit, pagingParams.Limit)
	}

	if err := db.GetDBConn().Raw(sqlQuery + pagingQuery, categoryId).Scan(&goods).Error; err != nil {
		return nil, err
	}

	for i := range goods {
		images, err := GetImagesOfGood(goods[i].Id)
		if err != nil {
			return nil, err
		}
		goods[i].Images = images
		fmt.Println(">>>>>>>>>>>>>", goods[i].Images)
	}

	return goods, nil
}

func AddViewToTheGood(goodId int, ip string) error {
	var good models.Good
	good.Id = goodId

	if err := GetGoodCountOfViews(&good); err != nil {
		return err
	}

	fmt.Println(">>ID", good.Id)
	fmt.Println(">>views", good.Views)
	trans := db.GetDBConn().Begin()
	if err := IncrementGoodViewsCount(trans, good.Views, good.Id); err != nil {
		trans.Rollback()
		return err
	}

	if err := AddToViewedIpsTable(trans, good.Id, ip); err != nil {
		trans.Rollback()
		return err
	}


	return trans.Commit().Error
}

func AddToViewedIpsTable(db *gorm.DB, goodId int, ip string) error {
	sqlQuery := "INSERT INTO viewed_ips (good_id, ip) VALUES(?,?)"
	return db.Table("viewed_ips").Exec(sqlQuery, goodId, ip).Error
}

func IncrementGoodViewsCount (db *gorm.DB, views, goodId int) error {
	views++
	fmt.Println(views)
	sqlQuery := "UPDATE goods SET views = ? WHERE id = ?"

	if err := db.Table("goods").Exec(sqlQuery, views, goodId).Error; err != nil {
		return err
	}

	return nil
}

func GetGoodCountOfViews(good *models.Good) error {
	sqlQuery := "SELECT views FROM goods WHERE id = ?"

	if err := db.GetDBConn().Table("goods").Raw(sqlQuery, good.Id).Scan(&good).Error; err != nil {
		return err
	}

	return nil
}

func MatchTopCategoriesWithGoods(categories []models.Category) (topCategoriesGoods []models.TopCategory, err error) {
	var topCategory models.TopCategory

	for _, category := range categories {
		topCategory.Category = category
		topGoods, err := GetCategoryTopGoodsFullInfo(category.Id)
		if err != nil {
			return nil, err
		}
		topCategory.Goods = topGoods
		topCategoriesGoods = append(topCategoriesGoods, topCategory)
	}

	return topCategoriesGoods, nil
}

func GetCategoryTopGoodsFullInfo(categoryId int) (topGoods []models.Good, err error) {
	topGoods, err = GetCategoryTopGoodsInfoWithoutImages(categoryId)
	if err != nil {
		return nil, err
	}

	for i := range topGoods {
		images, err := GetImagesOfGood(topGoods[i].Id)
		if err != nil {
			return nil, err
		}
		topGoods[i].Images = images
		fmt.Println(">>>>>>>>>>>>>GoodID", topGoods[i].Id)
		fmt.Println(">>>>>>>>>>>>>GoodImages", topGoods[i].Images)
		fmt.Println(">>>>>>>>>>>>>Good", topGoods)

	}

	return topGoods, nil
}

func GetCategoryTopGoodsInfoWithoutImages(categoryId int) (topGoods []models.Good, err error) {
	sqlQuery := "SELECT goods.*, shops.shop_name, shops.shop_phone, shops.shop_address FROM goods, shops WHERE is_top = true AND category_id = ? AND shop_id = shops.id AND goods.is_active = true ORDER BY goods.id"

	if err := db.GetDBConn().Raw(sqlQuery, categoryId).Scan(&topGoods).Error; err != nil {
		return nil, err
	}

	return topGoods, nil
}

func GetTopCategoryList() (categories []models.Category, err error) {
	sqlQuery := "SELECT id, name FROM categories WHERE is_top=true ORDER BY id"

	return categories, db.GetDBConn().Table("categories").Raw(sqlQuery).Scan(&categories).Error
}

func GetImagesOfGood(goodId int) (images []string, err error) {
	var imagesStruct []models.Image
	sqlQuery := "SELECT image_name FROM images WHERE good_id = ? ORDER BY id"

	if err := db.GetDBConn().Table("images").Raw(sqlQuery, goodId).Scan(&imagesStruct).Error; err != nil {
		return nil, err
	}

	for _, imageStruct := range imagesStruct {
		images = append(images, imageStruct.ImageName)
	}


	return images, nil
}

func IsInViewedIps(goodId int, ip string) (exists bool, err error) {
	var viewedIp models.ViewedIps
	sqlQuery := "SELECT * FROM viewed_ips WHERE good_id = ? AND ip = ? ORDER BY id"
	if err = db.GetDBConn().Table("viewed_ips").Raw(sqlQuery, goodId, ip).Scan(&viewedIp).Error; err != nil && err.Error() != "record not found"{
		return exists, err
	}
	fmt.Println(">>>exists", viewedIp.Id)

	if viewedIp.Id > 0 {
		exists = true
	}

	return exists, nil
}

func Search(params models.SearchParams, pagingParams  models.PagingParams) (goods []models.Good, err error) {
	var (
		sqlQuery = "SELECT goods.*, shops.shop_name, shops.shop_phone, shops.shop_address FROM goods, shops WHERE shop_id = shops.id AND is_active = true "
		pagingQuery = " ORDER BY goods.id "
		filterQuery = "AND (name like '%" + params.Query + "%' OR shops.shop_name like '%" + params.Query + "%' OR shops.shop_address like '%" + params.Query + "%')"
	)

	if params.DiscountAmount > 0 {
		filterQuery += fmt.Sprintf(" AND discount_amount = %d ", params.DiscountAmount)
	}

	if params.PriceFrom > 0 {
		filterQuery += fmt.Sprintf(" AND new_price >= %d ", params.PriceFrom)
	}

	if params.PriceTo > 0 {
		filterQuery += fmt.Sprintf(" AND new_price <= %d ", params.PriceTo)
	}

	if pagingParams.Page > 0 && pagingParams.Limit > 0 {
		pagingQuery += fmt.Sprintf(" OFFSET %d LIMIT %d ", (pagingParams.Page - 1) * pagingParams.Limit, pagingParams.Limit)
	}

	if err := db.GetDBConn().Raw(sqlQuery + filterQuery + pagingQuery).Scan(&goods).Error; err != nil {
		return nil, err
	}

	for i := range goods {
		images, err := GetImagesOfGood(goods[i].Id)
		if err != nil {
			return nil, err
		}
		goods[i].Images = images
		fmt.Println(">>>>>>>>>>>>>", goods[i].Images)
	}

	return goods, nil
}