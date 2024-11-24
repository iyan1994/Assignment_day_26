package service

import (
	"assignment_day_26/model"
	"assignment_day_26/repository"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreatePesanan(c *gin.Context) {

	var pesanan model.Pesanan

	err := c.ShouldBind(&pesanan)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.NewFailedResponse(fmt.Sprintf("failed to bind request %s :", err.Error())),
		)
		return
	}

	var existInventaris model.Inventaris
	err = repository.Db.First(&existInventaris, pesanan.IdProduk).Error

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			model.NewFailedResponse(fmt.Sprintf("id produk not found : %s", err.Error())),
		)
		return
	}

	if existInventaris.Jumlah < pesanan.Jumlah {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "count stok < jumlah pesanan"})
		return
	}
	update_stok := existInventaris.Jumlah - pesanan.Jumlah
	var inventaris model.Inventaris
	err = repository.Db.Model(&inventaris).Where("id_produk = ?", pesanan.IdProduk).Update("jumlah", update_stok).Error

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			model.NewFailedResponse(fmt.Sprintf("failed to update stok : %s", err.Error())),
		)
		return
	}

	tm_now := time.Now()
	pesanan.TanggalPesanan = tm_now

	err = repository.Db.Create(&pesanan).Error

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			model.NewFailedResponse(fmt.Sprintf("failed to save pesanan %s :", err.Error())),
		)
		return
	}
	result := gin.H{
		"id_produk":       pesanan.IdProduk,
		"id_pesanan":      pesanan.IdPesanan,
		"jumlah":          pesanan.Jumlah,
		"tanggal_pesanan": pesanan.TanggalPesanan,
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse("success", result))
}

func ViewByIdPesanan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.NewFailedResponse(fmt.Sprintf("invalid id : %s", err.Error())),
		)
		return
	}
	var pesanan model.Pesanan

	err = repository.Db.First(&pesanan, id).Error

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			model.NewFailedResponse(fmt.Sprintf("failed to get : %s", err.Error())),
		)
		return
	}

	var produk model.Produk

	err = repository.Db.First(&produk, pesanan.IdProduk).Error

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			model.NewFailedResponse(fmt.Sprintf("failed to get : %s", err.Error())),
		)
		return
	}

	result := gin.H{
		"id_pesanan":      pesanan.IdPesanan,
		"id_produk":       pesanan.IdProduk,
		"nama_produk":     produk.Nama,
		"kategori":        produk.Kategori,
		"jumlah":          pesanan.Jumlah,
		"tanggal_pesanan": pesanan.TanggalPesanan,
	}

	c.JSON(
		http.StatusOK,
		model.NewSuccessResponse("success", result),
	)

}
