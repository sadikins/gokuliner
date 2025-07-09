# 📝 PLANNING.md - Proyek Kalkuliner

Dokumen ini menjelaskan arsitektur, tujuan, gaya, dan batasan teknis untuk proyek Kalkuliner. Semua kontribusi harus mengikuti pedoman ini.

## 1. Gambaran Umum Proyek

**Kalkuliner** adalah aplikasi web full-stack yang dirancang untuk membantu pemilik usaha kuliner dalam mengelola biaya dan harga jual produk. 

**Fitur Utama:**
- Manajemen Bahan Baku, Tenaga Kerja, dan Biaya Operasional.
- Pembuatan Resep dan Kalkulasi Harga Pokok Penjualan (HPP) otomatis.
- Perhitungan dan penyimpanan Harga Jual berdasarkan HPP dan margin.
- Pembuatan Program Promosi (diskon/potongan).
- Simulasi dampak promosi dan komisi terhadap harga akhir.
- Dasbor untuk ringkasan data.

## 2. Arsitektur

Aplikasi ini menggunakan arsitektur client-server yang terpisah:

- **Backend:** API RESTful yang dibuat menggunakan **Go** dengan framework **Gin**. Bertanggung jawab atas semua logika bisnis, interaksi database, dan kalkulasi.
- **Frontend:** Aplikasi Single Page Application (SPA) yang dibuat menggunakan **Vue.js 3** dengan **Vite**. Bertanggung jawab atas antarmuka pengguna dan interaksi dengan backend API.

## 3. Konvensi Backend (Go)

- **Bahasa:** Go (versi 1.24+)
- **Framework Web:** Gin
- **ORM:** GORM
- **Database:** PostgreSQL untuk produksi, SQLite untuk pengembangan dan pengujian.
- **Struktur Direktori:**
  - `main.go`: Titik masuk aplikasi.
  - `/config`: Logika untuk memuat konfigurasi dari file `.env`.
  - `/database`: Inisialisasi koneksi database dan migrasi otomatis (AutoMigrate).
  - `/models`: Definisi struct GORM yang merepresentasikan tabel database.
  - `/handlers`: Logika untuk menangani request HTTP (controller), dipisahkan per modul (misal: `bahan_baku_handler.go`).
  - `/utils`: Fungsi bantuan umum (misal: pembulatan angka).
- **Gaya API:**
  - Semua endpoint berada di bawah prefix `/api`.
  - Menggunakan format JSON untuk request dan response body.
  - Penamaan endpoint menggunakan format kebab-case (misal: `/bahan-bakus`, `/program-promos`).
- **Manajemen Dependensi:** Menggunakan Go Modules (`go.mod` dan `go.sum`).
- **Variabel Lingkungan:** Dikelola melalui file `.env` dan di-load menggunakan `github.com/joho/godotenv`.

## 4. Konvensi Frontend (Vue.js)

- **Bahasa:** JavaScript (ES6+)
- **Framework:** Vue.js 3 (Composition API)
- **Build Tool:** Vite
- **State Management:** Pinia, untuk mengelola state global.
- **Routing:** Vue Router, untuk navigasi antar halaman.
- **Styling:** Tailwind CSS, untuk styling berbasis utilitas.
- **Komunikasi API:** Axios, untuk melakukan request HTTP ke backend.
- **Struktur Direktori (`/frontend/src`):
  - `/assets`: File statis seperti CSS global dan gambar.
  - `/components`: Komponen Vue yang dapat digunakan kembali (misal: `InputText.vue`).
  - `/views`: Komponen Vue yang merepresentasikan halaman penuh (misal: `BahanBakuView.vue`).
  - `/router`: Konfigurasi rute aplikasi.
  - `/store`: (Direktori potensial untuk modul Pinia jika state menjadi kompleks).
  - `/utils`: Fungsi bantuan JavaScript (misal: `formatters.js`).
- **Manajemen Dependensi:** Menggunakan NPM (`package.json` dan `package-lock.json`).

## 5. Pengujian (Testing)

- **Backend:** Menggunakan library standar `testing` Go dan `github.com/stretchr/testify` untuk assertion. File test diberi nama dengan akhiran `_test.go` dan ditempatkan di direktori yang sama dengan kode yang diuji.
- **Frontend:** (Rekomendasi) Menggunakan **Vitest** untuk unit testing karena integrasinya yang erat dengan Vite. Test dapat ditempatkan di dalam direktori `frontend/tests`.

## 6. Tujuan & Batasan

- **Pemisahan yang Jelas:** Jaga agar logika bisnis tetap di backend dan logika presentasi di frontend.
- **Akurasi Perhitungan:** Pastikan semua kalkulasi (HPP, promo, dll.) akurat dan dapat diandalkan.
- **Konsistensi Kode:** Ikuti gaya dan konvensi yang sudah ada. Kode dan komentar di backend ditulis dalam Bahasa Indonesia.
- **Modularitas:** Buat komponen dan fungsi yang dapat digunakan kembali untuk menghindari duplikasi kode.
