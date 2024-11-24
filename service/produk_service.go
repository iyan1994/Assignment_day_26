package service

import (
	"assignment_day_26/model"
	"assignment_day_26/repository"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateProduk(c *gin.Context) {
	var produkDto model.ProdukDto
	err := c.ShouldBind(&produkDto)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.NewFailedResponse(fmt.Sprintf("failed to bind request: %s", err.Error())),
		)
		return
	}

	produk := produkDto.ToModel()

	err = repository.Db.Create(&produk).Error

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			model.NewFailedResponse(fmt.Sprintf("failed to save produk : %s", err.Error())),
		)
		return
	}
	produkDto.ID = produk.ID
	produkDto.CreatedAt = produk.CreatedAt
	produkDto.UpdatedAt = produk.UpdatedAt

	c.JSON(http.StatusOK, model.NewSuccessResponse("success to save produk", produk))

}

func UpdateProduk(c *gin.Context) {
	var produkDto model.ProdukDto

	err := c.ShouldBind(&produkDto)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.NewFailedResponse(fmt.Sprintf("failed to bind request %s :", err.Error())),
		)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.NewFailedResponse(fmt.Sprintf("invalid id : %s :", err.Error())),
		)
		return
	}

	var existProduk model.Produk

	err = repository.Db.First(&existProduk, id).Error

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			model.NewFailedResponse(fmt.Sprintf("not exist id = %s", err.Error())),
		)
		return
	}

	produk := produkDto.ToModel()
	produk.ID = existProduk.ID
	produk.CreatedAt = existProduk.CreatedAt
	produk.UpdatedAt = existProduk.UpdatedAt

	err = repository.Db.Save(&produk).Error

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			model.NewFailedResponse(fmt.Sprintf("failed to update produk %s", err.Error())),
		)
		return
	}
	produkDto.ID = existProduk.ID
	produkDto.CreatedAt = existProduk.CreatedAt
	produkDto.UpdatedAt = existProduk.UpdatedAt
	c.JSON(http.StatusOK, model.NewSuccessResponse("success", produkDto))
}

func DeleteProduk(c *gin.Context) {
	var existProduk model.Produk
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.NewFailedResponse(fmt.Sprintf("invalid id : %s", err.Error())),
		)
		return
	}

	err = repository.Db.First(&existProduk, id).Error
	if err != nil {
		c.JSON(
			http.StatusNotFound,
			model.NewFailedResponse(fmt.Sprintf("not exist id : %s", err.Error())),
		)
		return
	}
	produk := model.Produk{ID: id}
	result := repository.Db.Delete(produk)

	if result.Error != nil {
		c.JSON(
			http.StatusInternalServerError,
			model.NewFailedResponse(fmt.Sprintf("failed to delete produk : %s", err.Error())),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		model.NewSuccessResponse(fmt.Sprintf("%d produk delete", result.RowsAffected), nil),
	)

}

func ViewProduk(c *gin.Context) {
	var produk []model.Produk
	query := `select * from produk`
	filter := c.Query("filter")
	var args []any

	if filter != "" {
		query = fmt.Sprintf(
			"%s %s",
			query,
			"where kategori = ?",
		)
		args = append(args, filter)
	}

	err := repository.Db.Raw(query, args...).Scan(&produk).Error

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			model.NewFailedResponse(fmt.Sprintf("failed to search produk : %s", err.Error())),
		)
		return
	}

	var produksDto []model.ProdukDto

	for _, data := range produk {
		var produkDto model.ProdukDto
		produkDto.FillFromModel(data)
		produksDto = append(produksDto, produkDto)
	}
	c.JSON(http.StatusOK, model.NewSuccessResponse("success", produksDto))

}

func ViewByIdProduk(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.NewFailedResponse(fmt.Sprintf("invalid id : %s", err.Error())),
		)
		return
	}

	var produk model.Produk

	err = repository.Db.First(&produk, id).Error

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			model.NewFailedResponse(fmt.Sprintf("failed to get : %s", err.Error())),
		)
		return
	}

	var produkDto model.ProdukDto
	produkDto.FillFromModel(produk)

	c.JSON(
		http.StatusOK,
		model.NewSuccessResponse("success", produkDto),
	)

}

var produkUploadDir = "uploads/produk"

func UploadGambarProduk(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.NewFailedResponse(fmt.Sprintf("failed to bind request : %s", err.Error())),
		)
		return
	}

	var existProduk model.Produk

	err = repository.Db.First(&existProduk, id).Error

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			model.NewFailedResponse(fmt.Sprintf("not exist id = %s", err.Error())),
		)
		return
	}

	formFile, file, err := c.Request.FormFile("gambar")
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.NewFailedResponse(fmt.Sprintf("failed to get upload file : %s", err.Error())),
		)
		return
	}
	defer formFile.Close()

	maxFileSize := 100 << 10 //10kb
	if file.Size > int64(maxFileSize) {
		c.JSON(
			http.StatusBadRequest,
			model.NewFailedResponse("file size exceed maximum"),
		)
		return

	}

	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" {
		c.JSON(
			http.StatusBadRequest,
			model.NewFailedResponse("please upload jpg or jpeg"),
		)
		return
	}

	buffer := make([]byte, 512)
	_, err = formFile.Read(buffer)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.NewFailedResponse("failed to read file buffer"),
		)
		return
	}

	log.Println("Checking file mime type")

	// Detect the MIME type
	mimeType := http.DetectContentType(buffer)
	if !strings.Contains(mimeType, "image/jpeg") {
		c.JSON(
			http.StatusBadRequest,
			model.NewFailedResponse("file is not a jpeg picture"),
		)
		return
	}

	name := c.PostForm("name")
	path := filepath.Join(produkUploadDir, name)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			model.NewFailedResponse(fmt.Sprintf("failed to save file: %s", err.Error())),
		)
		return
	}

	var produkDto model.ProdukDto
	produkDto.Gambar = &path
	produkDto.ID = existProduk.ID
	produkDto.Nama = existProduk.Nama
	produkDto.Kategori = existProduk.Kategori
	produkDto.CreatedAt = existProduk.CreatedAt
	produkDto.UpdatedAt = existProduk.UpdatedAt

	result := gin.H{
		"id_produk":  produkDto.ID,
		"nama":       produkDto.Nama,
		"kategori":   produkDto.Kategori,
		"created_at": produkDto.CreatedAt,
		"updated_at": produkDto.UpdatedAt,
		"gambar":     produkDto.Gambar,
	}

	err = repository.Db.Model(&existProduk).Where("id = ?", id).Update("gambar", produkDto.Gambar).Error

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			model.NewFailedResponse(fmt.Sprintf("failed to upload gambar: %s", err.Error())),
		)
		return
	}
	c.JSON(http.StatusOK, model.NewSuccessResponse("Success", result))

}

func DownloadImage(c *gin.Context) {
	// Get the filename from the URL parameter
	filename := c.Param("filename")

	// Construct the path to the image file
	imagePath := filepath.Join("uploads", "produk", filename)

	// Check if the file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		c.JSON(404, gin.H{"error": "Image not found"})
		return
	}

	// Serve the image file
	c.File(imagePath)
}
