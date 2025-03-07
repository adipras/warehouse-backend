# Warehouse Management System (Backend)

Sistem manajemen gudang sederhana menggunakan **Golang Gin** & **MySQL**.

---

## 🚀 1. Instalasi & Menjalankan Server

### **1.1 Clone Repository**
```sh
git clone https://github.com/adipras/warehouse-backend.git
cd warehouse-backend
```

### **1.2 Buat & Konfigurasi Database MySQL**
Buat database MySQL:
```sql
CREATE DATABASE warehouse_db;
```

### **1.3 Buat File `.env`**
Buat file `.env` di root proyek, lalu isi:
```env
DB_USER=root
DB_PASSWORD=password123
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=warehouse_db
JWT_SECRET=supersecretkey
```

### **1.4 Instal Dependensi**
Pastikan Go sudah terinstall, lalu jalankan:
```sh
go mod tidy
```

### **1.5 Jalankan Migrasi Database**
```sh
go run main.go migrate
```
> Perintah ini akan menjalankan fungsi migrasi database tanpa menjalankan server.

### **1.6 Menjalankan Server**
```sh
go run main.go
```
Server akan berjalan di `http://localhost:8080`

---

## 📜 2. API Endpoint

### **2.1 Autentikasi**
| Method | Endpoint          | Deskripsi          |
|--------|------------------|--------------------|
| POST   | `/auth/register` | Register User     |
| POST   | `/auth/login`    | Login & Dapatkan JWT |

### **2.2 Produk**
| Method | Endpoint             | Deskripsi                 |
|--------|----------------------|---------------------------|
| POST   | `/products`          | Tambah Produk             |
| GET    | `/products`          | Ambil Semua Produk        |
| GET    | `/products/:id`      | Ambil Produk Berdasarkan ID |
| PUT    | `/products/:id`      | Update Produk             |
| DELETE | `/products/:id`      | Hapus Produk              |
| POST   | `/products/bulk`     | Tambah Banyak Produk Sekaligus |

### **2.3 Ekspor & Barcode**
| Method | Endpoint             | Deskripsi                 |
|--------|----------------------|---------------------------|
| GET    | `/products/export`   | Ekspor Produk ke CSV      |
| GET    | `/products/barcode/:sku` | Generate Barcode SKU |

---

## 📖 3. Dokumentasi API Swagger
Swagger tersedia di:
```
http://localhost:8080/swagger/index.html
```

---

## 🛠 4. Teknologi yang Digunakan
- Golang Gin
- GORM (ORM untuk MySQL)
- JWT untuk Autentikasi
- Swagger untuk Dokumentasi API

---

