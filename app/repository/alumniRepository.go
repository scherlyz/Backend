package repository

import (
	"backendgo/app/model"
	"backendgo/database"
	"context"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ===================================================
// 🔹 Get All Alumni
// ===================================================
func GetAllAlumni() ([]model.Alumni, error) {
	rows, err := database.DB.Query(`
		SELECT id, user_id, nim, nama, jurusan, angkatan, tahun_lulus, email, 
		       no_telepon, alamat, status_kematian, created_at, updated_at
		FROM alumni 
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumniList []model.Alumni
	for rows.Next() {
		var a model.Alumni
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.NIM, &a.Nama, &a.Jurusan,
			&a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon,
			&a.Alamat, &a.StatusKematian, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		alumniList = append(alumniList, a)
	}
	return alumniList, nil
}

// ===================================================
// 🔹 Get Alumni by ID
// ===================================================
func GetAlumniByID(id int) (model.Alumni, error) {
	var a model.Alumni
	err := database.DB.QueryRow(`
		SELECT id, user_id, nim, nama, jurusan, angkatan, tahun_lulus, email, 
		       no_telepon, alamat, status_kematian, created_at, updated_at
		FROM alumni 
		WHERE id=$1
	`, id).Scan(
		&a.ID, &a.UserID, &a.NIM, &a.Nama, &a.Jurusan,
		&a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon,
		&a.Alamat, &a.StatusKematian, &a.CreatedAt, &a.UpdatedAt,
	)
	return a, err
}

// ===================================================
// 🔹 Create Alumni
// ===================================================
func CreateAlumni(a model.CreateAlumniRequest) (model.Alumni, error) {
	ctx := context.Background()
	tx, err := database.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return model.Alumni{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 1️⃣ Buat user otomatis
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return model.Alumni{}, err
	}

	var userID int
	err = tx.QueryRowContext(ctx, `
		INSERT INTO users (username, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id
	`, a.Nama, a.Email, string(passwordHash), "user").Scan(&userID)
	if err != nil {
		return model.Alumni{}, err
	}

	// 2️⃣ Insert ke tabel alumni
	var newAlumni model.Alumni
	err = tx.QueryRowContext(ctx, `
		INSERT INTO alumni (
			user_id, nim, nama, jurusan, angkatan, tahun_lulus, email,
			no_telepon, alamat, status_kematian, created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,NOW(),NOW()
		) RETURNING id, user_id, nim, nama, jurusan, angkatan, tahun_lulus,
		          email, no_telepon, alamat, status_kematian, created_at, updated_at
	`,
		userID, a.NIM, a.Nama, a.Jurusan, a.Angkatan,
		a.TahunLulus, a.Email, a.NoTelepon, a.Alamat, a.StatusKematian,
	).Scan(
		&newAlumni.ID, &newAlumni.UserID, &newAlumni.NIM, &newAlumni.Nama,
		&newAlumni.Jurusan, &newAlumni.Angkatan, &newAlumni.TahunLulus,
		&newAlumni.Email, &newAlumni.NoTelepon, &newAlumni.Alamat,
		&newAlumni.StatusKematian, &newAlumni.CreatedAt, &newAlumni.UpdatedAt,
	)
	if err != nil {
		return model.Alumni{}, err
	}

	err = tx.Commit()
	if err != nil {
		return model.Alumni{}, err
	}

	return newAlumni, nil
}

// ===================================================
// 🔹 Update Alumni
// ===================================================
func UpdateAlumni(a model.UpdateAlumniRequest) (model.Alumni, error) {
	now := time.Now()

	_, err := database.DB.Exec(`
		UPDATE alumni 
		SET nim=$1, nama=$2, jurusan=$3, angkatan=$4, tahun_lulus=$5,
		    email=$6, no_telepon=$7, alamat=$8, status_kematian=$9, updated_at=$10
		WHERE id=$11
	`,
		a.NIM, a.Nama, a.Jurusan, a.Angkatan, a.TahunLulus,
		a.Email, a.NoTelepon, a.Alamat, a.StatusKematian, now, a.ID,
	)
	if err != nil {
		return model.Alumni{}, err
	}

	// ambil lagi data terbaru
	updated, err := GetAlumniByID(a.ID)
	return updated, err
}

// ===================================================
// 🔹 Delete Alumni
// ===================================================
func DeleteAlumni(id int) error {
	_, err := database.DB.Exec("DELETE FROM alumni WHERE id=$1", id)
	return err
}

// ===================================================
// 🔹 Update Status Kematian
// ===================================================
func UpdateStatusKematian(id int, status bool) error {
	_, err := database.DB.Exec(`
        UPDATE alumni 
        SET status_kematian=$1, updated_at=NOW() 
        WHERE id=$2
    `, status, id)
	return err
}

// ===================================================
// 🔹 Pagination with Search
// ===================================================
func GetAlumniRepo(search, sortBy, order string, limit, offset int) ([]model.Alumni, error) {
	query := fmt.Sprintf(`
        SELECT id, user_id, nim, nama, jurusan, angkatan, tahun_lulus, email, 
               no_telepon, alamat, status_kematian, created_at, updated_at
        FROM alumni
        WHERE nama ILIKE $1 OR email ILIKE $1 OR jurusan ILIKE $1
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `, sortBy, order)

	rows, err := database.DB.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Alumni
	for rows.Next() {
		var a model.Alumni
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.NIM, &a.Nama, &a.Jurusan,
			&a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon,
			&a.Alamat, &a.StatusKematian, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

// ===================================================
// 🔹 Count total alumni (for pagination)
// ===================================================
func CountAlumniRepo(search string) (int, error) {
	var total int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) 
		FROM alumni 
		WHERE nama ILIKE $1 OR email ILIKE $1 OR jurusan ILIKE $1
	`, "%"+search+"%").Scan(&total)
	return total, err
}
