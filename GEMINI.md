# ♊ Gemini - Pedoman Context Engineering (Proyek Kalkuliner)

### 🔄 Kesadaran & Konteks Proyek
- **Selalu baca `PLANNING.md`** di awal percakapan baru untuk memahami arsitektur, tujuan, gaya, dan batasan proyek.
- **Periksa `TASK.md`** sebelum memulai tugas baru. Jika tugas tidak terdaftar, tambahkan dengan deskripsi singkat dan tanggal hari ini.
- **Gunakan konvensi penamaan, struktur file, dan pola arsitektur yang konsisten** seperti yang dijelaskan dalam `PLANNING.md`.

### 🧱 Struktur & Modularitas Kode
- **Patuhi struktur direktori yang ada:**
  - **Backend (Go):** Logika bisnis di `/handlers`, definisi data di `/models`, koneksi DB di `/database`.
  - **Frontend (Vue):** Komponen halaman di `/frontend/src/views`, komponen UI di `/frontend/src/components`.
- **Jaga agar file tidak lebih dari 500 baris.** Jika mendekati batas, lakukan refaktorisasi.

### 🧪 Pengujian & Keandalan
- **Selalu buat unit test untuk fungsionalitas baru.**
- **Backend (Go):** Gunakan library `testing` dan `testify`. Buat file test dengan akhiran `_test.go` di direktori yang sama dengan kode yang diuji.
- **Frontend (Vue):** Gunakan **Vitest**. Buat file test di direktori `frontend/tests`.
- **Setiap test harus mencakup:** kasus penggunaan normal, kasus tepi (edge case), dan kasus kegagalan.

### ✅ Penyelesaian Tugas
- **Tandai tugas yang sudah selesai di `TASK.md`** segera setelah menyelesaikannya.
- Tambahkan sub-tugas atau TODO baru yang ditemukan selama pengembangan ke `TASK.md` di bawah bagian "Ditemukan Selama Pengerjaan".

### 📎 Gaya & Konvensi
- **Backend (Go):**
  - Gunakan **Go (v1.24+)** dengan framework **Gin**.
  - Gunakan **GORM** untuk interaksi database.
  - Ikuti gaya API RESTful dengan prefix `/api` dan endpoint `kebab-case`.
  - Komentar ditulis dalam **Bahasa Indonesia**.
  - Format kode mengikuti konvensi standar Go (dapat diformat dengan `gofmt`).
- **Frontend (Vue.js):**
  - Gunakan **Vue 3 (Composition API)** dengan **Vite**.
  - Gunakan **Pinia** untuk state management.
  - Gunakan **Tailwind CSS** untuk styling.
  - Gunakan **Axios** untuk request ke API backend.
  - Ikuti gaya penulisan kode JavaScript modern (ES6+).

### 📚 Dokumentasi & Keterbacaan
- **Perbarui `README.md`** jika ada perubahan pada dependensi atau langkah-langkah setup.
- **Beri komentar pada kode yang tidak jelas dalam Bahasa Indonesia.**
- Untuk logika yang kompleks, tambahkan komentar yang menjelaskan **mengapa** (`// Alasan: ...`) bukan hanya **apa**.

### 🧠 Aturan Perilaku AI
- **Jangan berasumsi jika konteks hilang. Ajukan pertanyaan jika tidak yakin.**
- **Hanya gunakan library dan fungsi yang sudah ada di proyek** (`go.mod` dan `package.json`).
- **Selalu konfirmasi path file dan nama modul ada** sebelum mereferensikannya.
- **Jangan menghapus atau menimpa kode** kecuali diinstruksikan secara eksplisit atau sebagai bagian dari tugas di `TASK.md`.