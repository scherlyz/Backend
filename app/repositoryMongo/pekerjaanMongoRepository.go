package repositoryMongo

import (
	"backendgo/app/modelmongo"
	"backendgo/database"
	"context"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Ambil koleksi MongoDB
func getPekerjaanCollection() *mongo.Collection {
	return database.MongoDB.Collection("pekerjaan_alumni")
}

// -------------------- CREATE --------------------
func CreatePekerjaanMongo(data modelmongo.PekerjaanAlumni) (*modelmongo.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data.ID = primitive.NewObjectID()
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	data.IsDeleted = false

	_, err := getPekerjaanCollection().InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// -------------------- GET ALL (Non Deleted) --------------------
func GetAllPekerjaanMongo() ([]modelmongo.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"is_deleted": false}
	cursor, err := getPekerjaanCollection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []modelmongo.PekerjaanAlumni
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// -------------------- GET BY ID --------------------
func GetPekerjaanByIDMongo(id string) (*modelmongo.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var pekerjaan modelmongo.PekerjaanAlumni
	err = getPekerjaanCollection().FindOne(ctx, bson.M{"_id": objID}).Decode(&pekerjaan)
	if err != nil {
		return nil, err
	}
	return &pekerjaan, nil
}

// -------------------- UPDATE --------------------
func UpdatePekerjaanMongo(id string, req modelmongo.UpdatePekerjaanRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updateFields := bson.M{}
	if req.NamaPerusahaan != nil {
		updateFields["nama_perusahaan"] = *req.NamaPerusahaan
	}
	if req.PosisiJabatan != nil {
		updateFields["posisi_jabatan"] = *req.PosisiJabatan
	}
	if req.BidangIndustri != nil {
		updateFields["bidang_industri"] = *req.BidangIndustri
	}
	if req.LokasiKerja != nil {
		updateFields["lokasi_kerja"] = *req.LokasiKerja
	}
	if req.GajiRange != nil {
		updateFields["gaji_range"] = *req.GajiRange
	}
	if req.TanggalMulaiKerja != nil {
		updateFields["tanggal_mulai_kerja"] = *req.TanggalMulaiKerja
	}
	if req.TanggalSelesaiKerja != nil {
		updateFields["tanggal_selesai_kerja"] = *req.TanggalSelesaiKerja
	}
	if req.StatusPekerjaan != nil {
		updateFields["status_pekerjaan"] = *req.StatusPekerjaan
	}
	if req.DeskripsiPekerjaan != nil {
		updateFields["deskripsi_pekerjaan"] = *req.DeskripsiPekerjaan
	}
	updateFields["updated_at"] = time.Now()

	update := bson.M{"$set": updateFields}
	_, err = getPekerjaanCollection().UpdateByID(ctx, objID, update)
	return err
}

// -------------------- SOFT DELETE --------------------
func SoftDeletePekerjaanMongo(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"updated_at": time.Now(),
		},
	}
	_, err = getPekerjaanCollection().UpdateByID(ctx, objID, update)
	return err
}

// -------------------- RESTORE --------------------
func RestorePekerjaanMongo(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": false,
			"updated_at": time.Now(),
		},
	}
	_, err = getPekerjaanCollection().UpdateByID(ctx, objID, update)
	return err
}

// -------------------- HARD DELETE --------------------
func HardDeletePekerjaanMongo(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = getPekerjaanCollection().DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// -------------------- GET TRASHED --------------------
func GetTrashedPekerjaanMongo() ([]modelmongo.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"is_deleted": true}
	cursor, err := getPekerjaanCollection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []modelmongo.PekerjaanAlumni
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// -------------------- GET BY ALUMNI --------------------
func GetPekerjaanByAlumniMongo(alumniIDStr string) ([]modelmongo.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	alumniID, err := strconv.Atoi(alumniIDStr)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"alumni_id": alumniID,
		"is_deleted": false,
	}

	cursor, err := getPekerjaanCollection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []modelmongo.PekerjaanAlumni
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
