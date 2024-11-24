package main

import (
	"assignment_day_26/repository"
	"assignment_day_26/service"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:Lpkia@12345@tcp(127.0.0.1:3307)/db_inventaris?charset=utf8mb4&parseTime=true&loc=Local"

	var err error
	repository.Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	r := gin.Default()
	r.POST("/produk", service.CreateProduk)
	r.GET("/produk", service.ViewProduk)
	r.GET("/produk/:id", service.ViewByIdProduk)
	r.PUT("/produk/:id", service.UpdateProduk)
	r.DELETE("/produk/:id", service.DeleteProduk)
	r.POST("/gambar-produk/:id", service.UploadGambarProduk)
	r.GET("/product/image/:filename", service.DownloadImage)

	r.POST("/inventaris", service.CreateInventaris)
	r.GET("/inventaris/:id", service.ViewStok)
	r.PUT("/inventaris/:id", service.UpdateStok)

	r.POST("/pesanan", service.CreatePesanan)
	r.GET("/pesanan/:id", service.ViewByIdPesanan)

	err = r.Run(":8069")

	if err != nil {
		log.Fatalln(err)
		return

	}

}
