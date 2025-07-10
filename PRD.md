# Product Requirements Document (PRD) - Kalkuliner

**Versi:** 1.0
**Tanggal:** 10 Juli 2025
**Disusun oleh:** [Nama Pembuat PRD / Tim Produk]

---

## 1. Pengantar

**Kalkuliner** adalah aplikasi web *full-stack* yang dirancang untuk membantu pemilik usaha kuliner dalam mengelola biaya produksi (HPP), menetapkan harga jual yang optimal, dan mensimulasikan dampak program promosi terhadap keuntungan. Aplikasi ini bertujuan untuk memberikan visibilitas yang lebih baik terhadap kesehatan finansial produk kuliner, memungkinkan pengambilan keputusan berbasis data.

## 2. Visi Produk


Menjadi platform utama yang memberdayakan pemilik usaha kuliner kecil hingga menengah untuk mengoptimalkan profitabilitas mereka melalui manajemen biaya dan penetapan harga yang cerdas dan transparan.

## 3. Tujuan Produk (Product Goals)

* **Optimalisasi Biaya:** Memungkinkan pengguna untuk melacak dan menghitung biaya bahan baku, tenaga kerja, dan operasional dengan akurat.
* **Penetapan Harga Strategis:** Memfasilitasi penetapan harga jual yang optimal berdasarkan berbagai kriteria profitabilitas dan kondisi pasar.
* **Simulasi Keuntungan:** Memberikan alat simulasi untuk memahami dampak promosi, komisi *channel*, dan pajak terhadap margin keuntungan.
* **Efisiensi Operasional:** Mengurangi waktu manual yang dihabiskan untuk perhitungan HPP dan harga jual.
* **Aksesibilitas Data:** Menyajikan data biaya dan profitabilitas dalam format yang mudah dipahami melalui dasbor dan laporan.

## 4. Persona Pengguna

* **Pemilik Usaha Kuliner Kecil/Menengah:**
    * **Kebutuhan:** Memahami biaya riil produk mereka, menetapkan harga yang kompetitif namun menguntungkan, mengelola inventaris bahan baku sederhana, dan menganalisis dampak promosi.
    * **Masalah:** Kesulitan menghitung HPP secara manual, tidak yakin apakah harga jual sudah optimal, sering mengalami kerugian akibat komisi atau promo tak terduga.
* **Koki/Manajer Produksi:**
    * **Kebutuhan:** Memasukkan dan memperbarui resep, mencatat penggunaan bahan baku, dan melacak biaya tenaga kerja/operasional.

## 5. Alur Pengguna Utama (Key User Flows)

1.  **Manajemen Bahan Baku & Biaya:**
    * Pengguna masuk ke modul "Bahan Baku", "Tenaga Kerja", atau "Operasional".
    * Pengguna dapat menambah, melihat, mengedit, atau menghapus item biaya beserta detail harga, satuan beli, dan konversi satuan pemakaian.
2.  **Pembuatan & Perhitungan Resep (HPP):**
    * Pengguna masuk ke modul "Resep & HPP".
    * Pengguna dapat membuat resep baru dengan menambahkan komponen (bahan baku atau resep lain/sub-resep) beserta kuantitasnya.
    * Pengguna dapat melihat detail resep dan memicu perhitungan HPP untuk resep tertentu.
    * Sistem menghitung HPP per unit resep dan HPP per porsi secara otomatis.
3.  **Kalkulasi & Penyimpanan Harga Jual Optimal:**
    * Pengguna masuk ke modul "Harga Jual".
    * Pengguna memilih resep dan memasukkan biaya operasional terkait produk (pajak, komisi *channel*).
    * Pengguna memilih kriteria optimal yang diinginkan (misalnya, minimal profit X% dari *net sales*, target *net sales* Y Rupiah).
    * Sistem menghitung harga jual kotor yang optimal berdasarkan kriteria yang dipilih.
    * Pengguna dapat menyimpan harga jual yang terhitung ini.
4.  **Manajemen Program Promo:**
    * Pengguna masuk ke modul "Promo".
    * Pengguna dapat menambah, melihat, mengedit, atau menghapus program promo (diskon persentase/nominal, minimal belanja, maks potongan, ditanggung *merchant*).
5.  **Simulasi Keuntungan:**
    * Pengguna masuk ke modul "Simulasi".
    * Pengguna memilih menu/produk, jumlah porsi, dan dapat mengaktifkan promo ongkir, promo *channel*, serta memasukkan persentase komisi dan pajak simulasi.
    * Sistem menghitung proyeksi harga jual konsumen, diskon, biaya, *net sales*, dan *gross profit*.
6.  **Dasbor Ringkasan:**
    * Pengguna melihat ringkasan total bahan baku, total resep, total biaya operasional (placeholder), dan top 5 resep dengan HPP per porsi tertinggi.

## 6. Fitur Fungsional (Functional Requirements)

### 6.1. Modul Manajemen Biaya (Bahan Baku, Tenaga Kerja, Operasional)
* [RF.1.1] Pengguna harus dapat **menambah** item biaya baru (Bahan Baku, Tenaga Kerja, Operasional) dengan detail: Nama, Kategori, Harga Beli, Satuan Beli, Netto per Beli, Satuan Pemakaian, Catatan.
* [RF.1.2] Pengguna harus dapat **melihat daftar** semua item biaya, difilter berdasarkan kategori.
* [RF.1.3] Pengguna harus dapat **mengedit** detail item biaya yang sudah ada.
* [RF.1.4] Pengguna harus dapat **menghapus** item biaya yang sudah ada, dengan validasi jika sedang digunakan di resep lain.
* [RF.1.5] **(Backlog)** Mengembangkan modul Tenaga Kerja dan Operasional agar memiliki model dan logika terpisah dari Bahan Baku.

### 6.2. Modul Resep & HPP
* [RF.2.1] Pengguna harus dapat **menambah resep** baru dengan Nama, Jumlah Porsi, dan menandai apakah itu sub-resep.
* [RF.2.2] Pengguna harus dapat **menambahkan komponen** ke resep, memilih antara Bahan Baku atau Resep lain (sub-resep), dan menentukan kuantitas penggunaan.
* [RF.2.3] Sistem harus dapat **menghitung HPP** secara rekursif untuk setiap resep, mengakomodasi sub-resep.
* [RF.2.4] Pengguna harus dapat **melihat detail resep** dan hasil HPP per unit resep serta HPP per porsi.
* [RF.2.5] Sistem harus dapat **menyimpan hasil HPP** ke database jika ada perubahan signifikan.
* [RF.2.6] Pengguna harus dapat **mengedit** resep dan komponennya yang sudah ada.
* [RF.2.7] Pengguna harus dapat **menduplikasi** resep yang sudah ada.
* [RF.2.8] Pengguna harus dapat **menghapus** resep, dengan validasi jika resep tersebut digunakan sebagai komponen di resep lain.

### 6.3. Modul Harga Jual
* [RF.3.1] Pengguna harus dapat **memilih resep** sebagai dasar produk harga jual.
* [RF.3.2] Pengguna harus dapat memasukkan detail produk seperti Nama Produk, Channel Penjualan, Jumlah Porsi Produk, Pajak (%), dan Komisi Channel (%).
* [RF.3.3] Pengguna harus dapat **memilih satu kriteria optimal** untuk perhitungan harga jual (misalnya, min profit % net sales, min profit Rp dari HPP, target net sales Rp, dll.).
* [RF.3.4] Sistem harus dapat **menghitung harga jual kotor optimal** berdasarkan kriteria yang dipilih dan menyimpan hasilnya.
* [RF.3.5] Pengguna harus dapat memilih opsi untuk **membulatkan harga jual kotor**.
* [RF.3.6] Pengguna harus dapat **melihat daftar harga jual** yang sudah disimpan, beserta detail dan metodenya.
* [RF.3.7] Pengguna harus dapat **mengedit dan memperbarui** harga jual yang sudah ada.
* [RF.3.8] Pengguna harus dapat **menghapus** harga jual yang sudah ada.

### 6.4. Modul Program Promo
* [RF.4.1] Pengguna harus dapat **menambah** program promo baru dengan Nama Promo, Channel, Jenis Diskon (persentase/nominal), Besar Diskon, Min Belanja, Maksimal Potongan, Ditanggung Merchant (%), dan Catatan.
* [RF.4.2] Pengguna harus dapat **melihat daftar** semua program promo.
* [RF.4.3] Pengguna harus dapat **mengedit** program promo yang sudah ada.
* [RF.4.4] Pengguna harus dapat **menghapus** program promo yang sudah ada.

### 6.5. Modul Simulasi
* [RF.5.1] Pengguna harus dapat **memilih menu** (berdasarkan harga jual yang tersimpan) atau memasukkan detail menu secara manual (Nama Menu, HPP, Harga Jual Kotor).
* [RF.5.2] Pengguna harus dapat memasukkan **jumlah porsi pembelian** untuk simulasi.
* [RF.5.3] Pengguna harus dapat **mengaktifkan/menonaktifkan simulasi promo ongkir** dan memasukkan nominal ongkir yang ditanggung *merchant*.
* [RF.5.4] Pengguna harus dapat **mengaktifkan/menonaktifkan simulasi promo *channel*** dan memilih program promo yang tersedia.
* [RF.5.5] Pengguna harus dapat memasukkan **persentase komisi *channel* dan pajak** untuk simulasi.
* [RF.5.6] Sistem harus **menghitung dan menampilkan hasil simulasi** yang detail, termasuk harga akhir konsumen, potongan promo, biaya komisi/pajak/subsidi ongkir, *net sales*, *gross profit*, HPP terhadap *net sales* (%), dan *gross profit* terhadap *net sales* (%).

### 6.6. Modul Dasbor
* [RF.6.1] Pengguna harus dapat melihat **ringkasan jumlah total bahan baku** dan **total resep**.
* [RF.6.2] Pengguna harus dapat melihat **total biaya operasional** (saat ini sebagai *placeholder*).
* [RF.6.3] Pengguna harus dapat melihat **top 5 resep dengan HPP per porsi tertinggi**.

## 7. Persyaratan Non-Fungsional (Non-Functional Requirements)

* **Performa:**
    * [NFR.1.1] Perhitungan HPP untuk resep dengan kedalaman komponen hingga 5 level harus selesai dalam waktu kurang dari 2 detik.
    * [NFR.1.2] Pemuatan daftar (bahan baku, resep, harga jual) harus selesai dalam waktu kurang dari 3 detik untuk 1.000 item.
* **Keamanan:**
    * [NFR.2.1] **(Backlog)** Sistem harus memiliki otentikasi pengguna yang aman (login/register).
    * [NFR.2.2] Semua komunikasi antara frontend dan backend harus melalui HTTPS.
    * [NFR.2.3] Input pengguna harus divalidasi dengan kuat di sisi frontend dan backend untuk mencegah serangan umum (misalnya, *SQL injection*, *XSS*).
* **Skalabilitas:**
    * [NFR.3.1] Arsitektur *client-server* yang terpisah mendukung skalabilitas horizontal untuk backend dan frontend secara independen.
    * [NFR.3.2] Penggunaan *caching* di backend untuk data master HPP mendukung performa pada skala data yang lebih besar.
* **Keterpeliharaan (Maintainability):**
    * [NFR.4.1] Kode harus mengikuti konvensi penamaan dan gaya yang konsisten (Go Modules, Gin, GORM untuk Backend; Vue 3 Composition API, Vite, Pinia, Tailwind CSS untuk Frontend).
    * [NFR.4.2] Dokumentasi (termasuk `PLANNING.md`, `GEMINI.md`, `TASK.md`, API Spec, Database Schema, ADRs, Glosarium) harus terus diperbarui.
    * [NFR.4.3] **(Backlog)** Tersedia unit test yang memadai untuk fungsionalitas kritis di backend dan frontend.
* **Keandalan:**
    * [NFR.5.1] Perhitungan HPP harus menghasilkan nilai yang akurat secara matematis (dengan pembulatan konsisten).
    * [NFR.5.2] Sistem harus menangani dan memberikan pesan kesalahan yang jelas kepada pengguna.
    * [NFR.5.3] Transaksi database harus digunakan untuk operasi yang membutuhkan atomisitas (misalnya, membuat resep dengan komponennya).

## 8. Asumsi & Batasan

* **Asumsi:**
    * Pengguna memiliki pemahaman dasar tentang konsep HPP dan harga jual.
    * Aplikasi akan digunakan pada browser modern.
    * Database PostgreSQL akan digunakan di lingkungan produksi.
* **Batasan:**
    * Sistem tidak mengelola inventaris secara *real-time*. Penggunaan bahan baku dalam resep diasumsikan, tidak mengurangi stok.
    * Fokus pada perhitungan biaya dan harga jual; bukan sistem POS atau akuntansi lengkap.
    * Pengguna diasumsikan hanya memiliki satu lokasi/outlet.

## 9. Metrik Keberhasilan

* **Tingkat Penggunaan Fitur HPP:** Jumlah resep yang memiliki HPP terhitung dan tersimpan per minggu.
* **Adopsi Kalkulator Harga Jual Optimal:** Jumlah harga jual yang disimpan menggunakan kriteria optimal per minggu.
* **Akurasi Perhitungan:** Persentase laporan bug terkait perhitungan yang rendah.
* **Kepuasan Pengguna:** Hasil survei atau *feedback* positif dari pengguna.
* **Peningkatan Efisiensi:** Penurunan waktu yang dihabiskan pengguna untuk perhitungan HPP manual (dapat diukur dengan survei).

---