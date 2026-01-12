package dto

import "time"

type ProductResponse struct {
	ProductID         string    `json:"product_id"`
	SellerID          string    `json:"seller_id"`
	CategoryID        int32     `json:"category_id"`
	ShopCategoryID    int32     `json:"shop_category_id"`
	ProductType       string    `json:"product_type"`
	Sku               string    `json:"sku"`
	Name              string    `json:"name"`
	Slug              string    `json:"slug"`
	Description       string    `json:"description"`
	ShortDescription  string    `json:"short_description"`
	BasePrice         float64   `json:"base_price"`
	SalePrice         float64   `json:"sale_price"`
	CostPrice         float64   `json:"cost_price"`
	StockQuantity     int32     `json:"stock_quantity"`
	LowStockThreshold int32     `json:"low_stock_threshold"`
	IsUnlimitedStock  bool      `json:"is_unlimited_stock"`
	Status            string    `json:"status"`
	IsFeatured        bool      `json:"is_featured"`
	WeightKg          float64   `json:"weight_kg"`
	DimensionsCm      string    `json:"dimensions_cm"`
	Brand             string    `json:"brand"`
	Condition         string    `json:"condition"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CategoryResponse struct {
	CategoryID       int32             `json:"category_id"`
	ParentCategoryID int32             `json:"parent_category_id"`
	CategoryName     string            `json:"category_name"`
	Slug             string            `json:"slug"`
	Description      string            `json:"description"`
	ImageURL         string            `json:"image_url"`
	IsActive         bool              `json:"is_active"`
	CreatedAt        time.Time         `json:"created_at"`
	Products         []ProductResponse `json:"products"`
}

type TagResponse struct {
	TagID     int32     `json:"tag_id"`
	TagName   string    `json:"tag_name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

type ShopCategoryResponse struct {
	ShopCategoryID       int32     `json:"shop_category_id"`
	SellerID             string    `json:"seller_id"`
	ParentShopCategoryID int32     `json:"parent_shop_category_id"`
	CategoryName         string    `json:"category_name"`
	Slug                 string    `json:"slug"`
	Description          string    `json:"description"`
	IsActive             bool      `json:"is_active"`
	SortOrder            int32     `json:"sort_order"`
	CreatedAt            time.Time `json:"created_at"`
}

type ProductVariantResponse struct {
	VariantID       int32     `json:"variant_id"`
	ProductID       string    `json:"product_id"`
	Sku             string    `json:"sku"`
	VariantName     string    `json:"variant_name"`
	PriceAdjustment float64   `json:"price_adjustment"`
	Attributes      string    `json:"attributes"`
	StockQuantity   int32     `json:"stock_quantity"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ProductMediaResponse struct {
	MediaID      int32     `json:"media_id"`
	ProductID    string    `json:"product_id"`
	MediaType    string    `json:"media_type"`
	MediaURL     string    `json:"media_url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	AltText      string    `json:"alt_text"`
	SortOrder    int32     `json:"sort_order"`
	IsPrimary    bool      `json:"is_primary"`
	CreatedAt    time.Time `json:"created_at"`
}
