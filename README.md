
install 
- restore db "db_inventaris.sql"
- config username,password,nama db di "main.go"
- jalankan perintah " go run main.go "

ENDPOINT PRODUK 
- {{host}}/produk?filter=Elektronik //method GET get produk berdasarkan kategori
- {{host}}/produk // method GET get produk all
- {{host}}/produk // method POST create produk
  request json
  {
	  "nama": "Laptop Asus",
    "deksripsi": "Laptop Asus core i5",
    "harga": 2000000,
    "kategori": "Elektronik"
 }
- {{host}}/produk/2 // method PUT param id update produk
  request json
  {
	  "nama": "Laptop Asus",
    "deksripsi": "Laptop Asus core i5",
    "harga": 6000000,
    "kategori": "Elektronik"
  }
- {{host}}/produk/1 // method  DELETE  param id delete produk
- {{host}}/produk/1 // method  GET param id get produk by id
- {{host}}/gambar-produk/3 // method POST param id upload image
  request form data
  gambar = file
  name = nama_gambar.jpg

------------------------------------------------------------------------
ENDPOINT INVENTARIS
- {{host}}/inventaris // method POST create inventaris
  request json
  {
    "id_produk" : 4,
    "jumlah" : 50,
    "lokasi" : "rak 3"
  }
- {{host}}/inventaris/4 // method GET param id get stok produk
- {{host}}/inventaris/4 // method PUT param id update stok
  request json
  {
    "jumlah" : 20
  }

--------------------------------------------------------------------------
ENDPOINT PESANAN
- {{host}}/pesanan // method POST create pesanan
  request json
  {
    "id_produk" : 4,
    "jumlah" : 10
  }
- {{host}}/pesanan/6 // method GET param id get detail pesanan

--------------------------------------------------------------------------

PENJELASAN 
- stok otomatis berkurang pada table inventaris ketika pesanan di buat
- stok bisa di tambah pada endpoint "{{host}}/inventaris/4 // method PUT param id update stok"


