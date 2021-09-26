package models

type Banners struct {
	Id        int    `json:"id"`
	ImageName string `json:"image_name"`
}

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type TopCategory struct {
	Category Category `json:"category"`
	Goods    []Good   `json:"goods"`
}

type Good struct {
	Id             int      `json:"id"`
	CategoryId     int      `json:"category_id"`
	Images         []string `json:"images"`
	DiscountAmount float32  `json:"discount_amount"`
	IsTop          bool     `json:"is_top"`
	Name           string   `json:"name"`
	OldPrice       float32  `json:"old_price"`
	NewPrice       float32  `json:"new_price"`
	Shop
	Views int
}

type Shop struct {
	ShopName    string `json:"shop_name"`
	ShopPhone   string `json:"shop_phone"`
	ShopAddress string `json:"shop_address"`
}

type Image struct {
	Id        int    `json:"id" gorm:"id"`
	ImageName string `json:"image_name" gorm:"image_name"`
	GoodId    int    `json:"good_id" gorm:"good_id"`
}

type ViewedIps struct {
	Id     int    `json:"id" gorm:"id"`
	GoodId int    `json:"good_id" gorm:"good_id"`
	Ip     string `json:"ip" gorm:"ip"`
}

type SearchParams struct {
	Query string `json:"query"`
	DiscountAmount int `json:"discount_amount"`
	PriceFrom int `json:"price_from"`
	PriceTo int `json:"price_to"`
}

type PagingParams struct {
	Page int `json:"page"`
	Limit int `json:"limit"`
}
