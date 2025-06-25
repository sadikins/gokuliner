// src/utils/formatters.js

/**
 * Memformat angka kuantitas agar ditampilkan dengan presisi yang tepat.
 * Menghilangkan angka desimal .000 jika nilai adalah bilangan bulat.
 * Mempertahankan hingga 3 angka desimal jika ada.
 * @param {string|number} value - Nilai kuantitas yang akan diformat.
 * @returns {string} - Nilai kuantitas yang sudah diformat.
 */
export function formatQuantity(value) {
    const num = parseFloat(value);
    if (isNaN(num)) {
      return '0'; // Atau sesuaikan dengan string default yang Anda inginkan
    }
    // Jika nilai adalah bilangan bulat (tidak ada desimal atau desimalnya 0)
    if (num === Math.floor(num)) {
      return num.toString();
    }
    // Jika ada desimal, pertahankan hingga 3 angka di belakang koma
    // toFixed(3) akan membulatkan. Jika Anda ingin memangkas, gunakan logika lain.
    return num.toFixed(3).replace(/\.?0+$/, ''); // Menghilangkan trailing zeros setelah titik desimal
  }

  /**
   * Memformat nilai mata uang ke format Indonesia.
   * @param {string|number} value - Nilai mata uang.
   * @returns {string} - Nilai mata uang yang sudah diformat.
   */
  export function formatCurrency(value) {
    const num = parseFloat(value);
    if (isNaN(num)) {
      return '0.00'; // Atau sesuaikan
    }
    return num.toLocaleString('id-ID', { minimumFractionDigits: 0, maximumFractionDigits: 2 });
  }

   /**
   * Memformat tanggal.
   * @param {date} value - Tanggal.
   * @returns {string} - tanggal yang sudah diformat.
   */
  export const formatDateTime = (dateTimeString) => {
    if (!dateTimeString) return '-';
    const date = new Date(dateTimeString);
    return date.toLocaleString('id-ID', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    });
  };

