package model

import "time"

type Pesanan struct {
	IdPesanan      int       `json:"id_pesanan"`
	IdProduk       int       `json:"id_produk"`
	Jumlah         int       `json:"jumlah"`
	TanggalPesanan time.Time `json:"tanggal_pesanan"`
}

type TablerPesanan interface {
	TableName() string
}

func (Pesanan) TableName() string {
	return "pesanan"
}
