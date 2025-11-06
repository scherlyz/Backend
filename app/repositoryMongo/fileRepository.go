package repositoryMongo

import (
	"backendgo/app/modelmongo"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileRepository interface {
	Create(file *modelmongo.File) error
	FindAll() ([]modelmongo.File, error)
	FindByUser(userID int) ([]modelmongo.File, error)
	Delete(id string) error
}

type fileRepository struct {
	collection *mongo.Collection
}

func NewFileRepository(db *mongo.Database) FileRepository {
	if db == nil {
		panic("❌ MongoDB belum terhubung (db == nil) saat membuat FileRepository")
	}
	return &fileRepository{
		collection: db.Collection("files"),
	}
}

func (r *fileRepository) Create(file *modelmongo.File) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	file.UploadedAt = time.Now()
	res, err := r.collection.InsertOne(ctx, file)
	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		file.ID = oid
	}
	return nil
}

func (r *fileRepository) FindAll() ([]modelmongo.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var files []modelmongo.File
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &files); err != nil {
		return nil, err
	}
	return files, nil
}

func (r *fileRepository) FindByUser(userID int) ([]modelmongo.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var files []modelmongo.File
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &files); err != nil {
		return nil, err
	}
	return files, nil
}

func (r *fileRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID tidak valid")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("gagal menghapus file: %v", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("file tidak ditemukan")
	}
	return nil
}
