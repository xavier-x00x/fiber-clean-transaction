package entity

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ProductCode      string  `gorm:"size:30;not null;uniqueIndex:idx_product_deleted"`
	Name             string  `gorm:"size:250;not null"`
	Description      string  `gorm:"type:text"`
	PurchasePrice    float64 `gorm:"type:decimal(15,2);default:0;not null"`
	Margin           float64 `gorm:"type:decimal(5,2);default:0;not null;comment: 'Margin percentage to calculate selling price'"`
	SellPrice        float64 `gorm:"type:decimal(15,2);default:0;not null;comment: 'Price before tax'"`
	SellPriceWithTax float64 `gorm:"type:decimal(15,2);default:0;not null;comment: 'Price after tax'"`
	Stock            float64 `gorm:"type:decimal(15,2);default:0;not null"`
	Status           bool    `gorm:"type:tinyint(1);default:1;comment: 'Product status (active/inactive)'"`
	StoreID          uint    `gorm:"not null"`
	Store            Store

	TaxID uint `gorm:"not null"`
	Tax   Tax

	UnitID uint `gorm:"not null"`
	Unit   Unit

	CategoryID uint `gorm:"not null"`
	Category   Category

	StockAlerts float64        `gorm:"type:decimal(15,2);default:0;not null;comment: 'Minimum stock level to trigger alerts'"`
	StockMax    float64        `gorm:"type:decimal(15,2);default:0;not null;comment: 'Maximum stock level'"`
	DeletedAt   gorm.DeletedAt `gorm:"uniqueIndex:idx_product_deleted"`
}
