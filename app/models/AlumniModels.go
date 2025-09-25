package models

import "time"

type Alumni struct {
    ID          int       `json:"id"`
    NIM         string    `json:"nim"`
    Nama        string    `json:"nama"`
    Jurusan     string    `json:"jurusan"`
    Angkatan    int       `json:"angkatan"`
    TahunLulus  int       `json:"tahun_lulus"`
    Email       string    `json:"email"`
    NoTelepon   *string   `json:"no_telepon"`
    Alamat      *string   `json:"alamat"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type CreateAlumniRequest struct {
    NIM        string  `json:"nim"`
    Nama       string  `json:"nama"`
    Jurusan    string  `json:"jurusan"`
    Angkatan   int     `json:"angkatan"`
    TahunLulus int     `json:"tahun_lulus"`
    Email      string  `json:"email"`
    NoTelepon  *string `json:"no_telepon"`
    Alamat     *string `json:"alamat"`
}

type UpdateAlumniRequest struct {
    Nama       string  `json:"nama"`
    Jurusan    string  `json:"jurusan"`
    Angkatan   int     `json:"angkatan"`
    TahunLulus int     `json:"tahun_lulus"`
    Email      string  `json:"email"`
    NoTelepon  *string `json:"no_telepon"`
    Alamat     *string `json:"alamat"`
}

type FilterAlumni struct {
    Nama            string `json:"nama" query:"nama"`
    Jurusan         string `json:"jurusan" query:"jurusan"`
    Angkatan        int    `json:"angkatan" query:"angkatan"`
    TahunLulus      int    `json:"tahun_lulus" query:"tahun_lulus"`
	NamaPerusahaan  *string `json:"nama_perusahaan" query:"nama_perusahaan"`
	StatusPekerjaan *string `json:"status_pekerjaan" query:"status_pekerjaan"`
}


