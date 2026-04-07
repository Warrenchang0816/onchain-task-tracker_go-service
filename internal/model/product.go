package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID                  int64
	ProductID           string      // UUID
	MerchantAddress     string      // 商家钱包地址
	CustomerAddress     *string     // 购买者钱包地址(可空)

	// 商品信息
	Name                string
	Description         string
	Category            string      // 食物类别(主食/饮料/甜品)
	ImageURL            string

	// 价格信息
	OriginalPrice       decimal.Decimal  // 原价
	DiscountedPrice     decimal.Decimal  // 折扣价
	DiscountRate        int              // 折扣率(基点) 500 = 5%

	// 库存与状态
	Quantity            int              // 总库存数量
	AvailableQty        int              // 可用数量
	Status              string           // AVAILABLE|RESERVED|SOLD|EXPIRED
	ExpiryTime          *time.Time       // 食物过期时间

	// 销售信息
	SalesTarget         int              // 目标销售数量
	SoldQuantity        int              // 已售数量

	// 区块链字段
	ChainID             int64
	VaultAddress        string
	ContractProductID   string
	PaymentStatus       string           // NOT_FUNDED|FUNDED|PAID|REFUNDED

	// 时间戳
	CreatedAt           time.Time
	UpdatedAt           time.Time
	ListedAt            time.Time
	SoldOutAt           *time.Time
}