package model

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type Produk struct {
	ID        int             `json:"id"`
	Nama      string          `json:"nama"`
	Deksripsi sql.NullString  `json:"deksripsi"`
	Harga     decimal.Decimal `json:"harga"`
	Kategori  string          `json:"kategori"`
	Gambar    sql.NullString  `json:"gambar"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type ProdukDto struct {
	ID        int       `json:"id"`
	Nama      string    `json:"nama"`
	Deksripsi *string   `json:"deksripsi"`
	Harga     int       `json:"harga"`
	Kategori  string    `json:"kategori"`
	Gambar    *string   `json:"gambar"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoCreateTime"`
}

type TablerProduk interface {
	TableName() string
}

func (Produk) TableName() string {
	return "produk"
}

func (p *ProdukDto) FillFromModel(model Produk) {
	p.ID = model.ID
	p.Nama = model.Nama
	if model.Deksripsi.Valid {
		p.Deksripsi = &model.Deksripsi.String
	}
	p.Harga = int(model.Harga.IntPart())
	p.Kategori = model.Kategori
	if model.Gambar.Valid {
		p.Gambar = &model.Gambar.String
	}
	p.CreatedAt = model.CreatedAt
	p.UpdatedAt = model.UpdatedAt
}

func (p *ProdukDto) ToModel() Produk {
	model := Produk{
		ID:        p.ID,
		Nama:      p.Nama,
		Harga:     decimal.NewFromInt(int64(p.Harga)),
		Kategori:  p.Kategori,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}

	if p.Deksripsi != nil {
		model.Deksripsi.String = *p.Deksripsi
		model.Deksripsi.Valid = true
	}
	if p.Gambar != nil {
		model.Gambar.String = *p.Gambar
		model.Gambar.Valid = true
	}
	return model
}
