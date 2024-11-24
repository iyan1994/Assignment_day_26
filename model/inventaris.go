package model

import "time"

type Inventaris struct {
	IdProduk  int       `json:"id_produk"`
	Jumlah    int       `json:"jumlah"`
	Lokasi    string    `json:"lokasi"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TablerInventaris interface {
	TableName() string
}

func (Inventaris) TableName() string {
	return "inventaris"
}
