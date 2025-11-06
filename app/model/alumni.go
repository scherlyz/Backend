package model

import "time"

type Alumni struct {
	ID         int       `json:"id"`
	UserID     int   	`json:"user_id"`
	NIM        string    `json:"nim"`
	Nama       string    `json:"nama"`
	Jurusan    string    `json:"jurusan"`
	Angkatan   int       `json:"angkatan"`
	TahunLulus int       `json:"tahun_lulus"`
	Email      string    `json:"email"`
	NoTelepon  string    `json:"no_telepon"`
	Alamat     string    `json:"alamat"`
	StatusKematian bool `json:"status_kematian"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateAlumniRequest struct {
	NIM            string `json:"nim" example:"12345678"`
	Nama           string `json:"nama" example:"John Doe"`
	Jurusan        string `json:"jurusan" example:"Teknik Informatika"`
	Angkatan       int    `json:"angkatan" example:"2019"`
	TahunLulus     int    `json:"tahun_lulus" example:"2023"`
	Email          string `json:"email" example:"john@example.com"`
	NoTelepon      string `json:"no_telepon" example:"08123456789"`
	Alamat         string `json:"alamat" example:"Jl. Merdeka No. 1"`
	StatusKematian bool   `json:"status_kematian" example:"false"`
}

type UpdateAlumniRequest struct {
	ID             int    `json:"id"` // biar bisa dipakai di repository
	UserID         int    `json:"user_id,omitempty"`
	NIM            string `json:"nim" example:"12345678"`
	Nama           string `json:"nama" example:"John Doe"`
	Jurusan        string `json:"jurusan" example:"Teknik Informatika"`
	Angkatan       int    `json:"angkatan" example:"2019"`
	TahunLulus     int    `json:"tahun_lulus" example:"2023"`
	Email          string `json:"email" example:"john@example.com"`
	NoTelepon      string `json:"no_telepon" example:"08123456789"`
	Alamat         string `json:"alamat" example:"Jl. Merdeka No. 1"`
	StatusKematian bool   `json:"status_kematian" example:"false"`
}

type UpdateStatusKematianRequest struct {
	StatusKematian bool `json:"status_kematian" example:"true"`
}


