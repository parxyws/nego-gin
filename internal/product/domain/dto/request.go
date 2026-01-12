package dto

type ProductCreateRequest struct {
	CategoryID       int32   `json:"category_id" validate:"required"`
	ProductType      string  `json:"product_type" validate:"required"`
	Sku              string  `json:"sku" validate:"required"`
	Name             string  `json:"name" validate:"required"`
	Description      string  `json:"description"`
	ShortDescription string  `json:"short_description"`
	BasePrice        float64 `json:"base_price" validate:"required"`
	SalePrice        float64 `json:"sale_price"`
	StockQuantity    int32   `json:"stock_quantity"`
	IsUnlimitedStock bool    `json:"is_unlimited_stock"`
	WeightKg         float64 `json:"weight_kg"`
	DimensionsCm     string  `json:"dimensions_cm"`
	Brand            string  `json:"brand"`
	Condition        string  `json:"condition" validate:"required"`
	ShopCategoryID   *int32  `json:"shop_category_id"`
}

type ProductUpdateRequest struct {
	CategoryID     int32   `json:"category_id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	BasePrice      float64 `json:"base_price"`
	SalePrice      float64 `json:"sale_price"`
	StockQuantity  int32   `json:"stock_quantity"`
	Status         string  `json:"status"`
	ShopCategoryID *int32  `json:"shop_category_id"`
}

type CategoryCreateRequest struct {
	ParentCategoryID int32  `json:"parent_category_id"`
	CategoryName     string `json:"category_name" validate:"required"`
	Description      string `json:"description"`
}

type CategoryUpdateRequest struct {
	CategoryID       int32  `json:"category_id"`
	ParentCategoryID int32  `json:"parent_category_id"`
	CategoryName     string `json:"category_name" validate:"required"`
	Description      string `json:"description"`
}

type ShopCategoryCreateRequest struct {
	ParentShopCategoryID *int32 `json:"parent_shop_category_id"`
	CategoryName         string `json:"category_name" validate:"required"`
	Description          string `json:"description"`
	SortOrder            int32  `json:"sort_order"`
}

type ShopCategoryUpdateRequest struct {
	ParentShopCategoryID *int32 `json:"parent_shop_category_id"`
	CategoryName         string `json:"category_name"`
	Description          string `json:"description"`
	IsActive             *bool  `json:"is_active"`
	SortOrder            int32  `json:"sort_order"`
}

type TagCreateRequest struct {
	TagName string `json:"tag_name" validate:"required"`
}

type ProductVariantCreateRequest struct {
	Sku             string  `json:"sku" validate:"required"`
	VariantName     string  `json:"variant_name" validate:"required"`
	PriceAdjustment float64 `json:"price_adjustment"`
	Attributes      string  `json:"attributes" validate:"required"`
	StockQuantity   int32   `json:"stock_quantity"`
}

type ProductVariantUpdateRequest struct {
	Sku             string  `json:"sku"`
	VariantName     string  `json:"variant_name"`
	PriceAdjustment float64 `json:"price_adjustment"`
	Attributes      string  `json:"attributes"`
	StockQuantity   int32   `json:"stock_quantity"`
	IsActive        *bool   `json:"is_active"`
}

type ProductMediaUploadRequest struct {
	MediaType    string `json:"media_type" validate:"required"`
	MediaURL     string `json:"media_url" validate:"required"`
	ThumbnailURL string `json:"thumbnail_url"`
	AltText      string `json:"alt_text"`
	SortOrder    int32  `json:"sort_order"`
	IsPrimary    bool   `json:"is_primary"`
}

type ProductMediaUpdateRequest struct {
	SortOrder int32 `json:"sort_order"`
}

type ProductStatusUpdateRequest struct {
	Status string `json:"status" validate:"required"`
}

type ProductSearchRequest struct {
	Query      string  `json:"query"`
	CategoryID int32   `json:"category_id"`
	MinPrice   float64 `json:"min_price"`
	MaxPrice   float64 `json:"max_price"`
	SortBy     string  `json:"sort_by"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
}
