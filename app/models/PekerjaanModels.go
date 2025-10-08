package models

import (
	"database/sql"
	"time"
)

type PekerjaanAlumni struct {
    ID                  int        `json:"id"`
    AlumniID            int        `json:"alumni_id"`
    NamaPerusahaan      string     `json:"nama_perusahaan"`
    PosisiJabatan       string     `json:"posisi_jabatan"`
    BidangIndustri      string     `json:"bidang_industri"`
    LokasiKerja         string     `json:"lokasi_kerja"`
    GajiRange           *string    `json:"gaji_range"`
    TanggalMulaiKerja   time.Time  `json:"tanggal_mulai_kerja"`
    TanggalSelesaiKerja *time.Time `json:"tanggal_selesai_kerja"`
    StatusPekerjaan     string     `json:"status_pekerjaan"`
    DeskripsiPekerjaan  *string    `json:"deskripsi_pekerjaan"`
    CreatedAt           time.Time  `json:"created_at"`
    UpdatedAt           time.Time  `json:"updated_at"`
    AlumniNama *string `json:"alumni_nama,omitempty"`
    UserID     *int    `json:"user_id,omitempty"` 
}

type CreatePekerjaanRequest struct {
    AlumniID            int        `json:"alumni_id"`
    NamaPerusahaan      string     `json:"nama_perusahaan"`
    PosisiJabatan       string     `json:"posisi_jabatan"`
    BidangIndustri      string     `json:"bidang_industri"`
    LokasiKerja         string     `json:"lokasi_kerja"`
    GajiRange           *string    `json:"gaji_range"`
    TanggalMulaiKerja   string     `json:"tanggal_mulai_kerja"` // Format: YYYY-MM-DD
    TanggalSelesaiKerja *string    `json:"tanggal_selesai_kerja"`
    StatusPekerjaan     string     `json:"status_pekerjaan"`
    DeskripsiPekerjaan  *string    `json:"deskripsi_pekerjaan"`
}

type UpdatePekerjaanRequest struct {
    NamaPerusahaan      string     `json:"nama_perusahaan"`
    PosisiJabatan       string     `json:"posisi_jabatan"`
    BidangIndustri      string     `json:"bidang_industri"`
    LokasiKerja         string     `json:"lokasi_kerja"`
    GajiRange           *string    `json:"gaji_range"`
    TanggalMulaiKerja   string     `json:"tanggal_mulai_kerja"`
    TanggalSelesaiKerja *string    `json:"tanggal_selesai_kerja"`
    StatusPekerjaan     string     `json:"status_pekerjaan"`
    DeskripsiPekerjaan  *string    `json:"deskripsi_pekerjaan"`
}

type FilterPekerjaan struct {
	// AlumniID            *int        `json:"alumni_id"`
    Nama                string      `json:"nama"`
	NamaPerusahaan      *string     `json:"nama_perusahaan"`
	PosisiJabatan       *string     `json:"posisi_jabatan"`
	BidangIndustri      *string     `json:"bidang_industri"`
	LokasiKerja         *string     `json:"lokasi_kerja"`
	GajiRange           *string     `json:"gaji_range" query:"gaji_range"`
	StatusPekerjaan     *string     `json:"status_pekerjaan"`
	TahunMulaiKerja     *int        `json:"tahun_mulai"`
	TahunSelesaiKerja   *int        `json:"tahun_selesai"`
}

type PekerjaanDeleted struct {
    ID 			int 		`json:"id"`
    IsDeleted 	sql.NullTime   `json:"is_deleted"`
}
