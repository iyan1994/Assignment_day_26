package service

import (
	"assignment_day_26/model"
	"assignment_day_26/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateInventaris(c *gin.Context) {
	var inventaris model.Inventaris

	err := c.ShouldBind(&inventaris)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.NewFailedResponse(fmt.Sprintf("failed to bind request %s :", err.Error())),
		)
		return
	}

	err = repository.Db.Create(&inventaris).Error

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			model.NewFailedResponse(fmt.Sprintf("failed to save inventaris %s :", err.Error())),
		)
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse("success", inventaris))

}

func ViewStok(c *gin.Context) {
	var inventaris model.Inventaris

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(
			http.StatusBadGateway,
			model.NewFailedResponse(fmt.Sprintf("invalid id : %s", err.Error())),
		)
		return
	}

	err = repository.Db.First(&inventaris, id).Error
	if err != nil {
		c.JSON(
			http.StatusNotFound,
			model.NewFailedResponse(fmt.Sprintf("inventaris not found : %s", err.Error())),
		)
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse("success", inventaris))

}

func UpdateStok(c *gin.Context) {
	var inventaris model.Inventaris

	err := c.ShouldBind(&inventaris)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.NewFailedResponse(fmt.Sprintf("failed to bind request : %s", err.Error())),
		)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.NewFailedResponse(fmt.Sprintf("invalid id : %s", err.Error())),
		)
		return

	}
	var existInventaris model.Inventaris
	err = repository.Db.First(&existInventaris, id).Error
	if err != nil {
		c.JSON(
			http.StatusNotFound,
			model.NewFailedResponse(fmt.Sprintf("inventaris not found : %s", err.Error())),
		)
		return
	}

	tambahStok := existInventaris.Jumlah + inventaris.Jumlah

	err = repository.Db.Model(&inventaris).Where("id_produk = ?", id).Update("jumlah", tambahStok).Error

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			model.NewFailedResponse(fmt.Sprintf("failed to update stok inventaris : %s", err.Error())),
		)
		return
	}
	result := gin.H{
		"id_produk":  existInventaris.IdProduk,
		"jumlah":     inventaris.Jumlah,
		"lokasi":     existInventaris.Lokasi,
		"created_at": existInventaris.CreatedAt,
		"updated_t":  existInventaris.UpdatedAt,
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse("success update stok", result))

}
