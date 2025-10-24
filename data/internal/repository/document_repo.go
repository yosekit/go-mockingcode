package repository

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/go-mockingcode/data/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DocumentRepository struct {
	client *mongo.Client
	dbName string
}

func NewDocumentRepository(client *mongo.Client, dbName string) *DocumentRepository {
	return &DocumentRepository{
		client: client,
		dbName: dbName,
	}
}

func (r *DocumentRepository) GetCollection(collectionName string) *mongo.Collection {
	return r.client.Database(r.dbName).Collection(collectionName)
}

// CreateDocument создает новый документ
func (r *DocumentRepository) CreateDocument(projectID int64, collectionName string, data map[string]any) (*model.MockDocument, error) {
	collection := r.GetCollection(collectionName)
	ctx := context.Background()

	// Генерируем автоинкрементный ID если его нет в data
	if _, hasID := data["id"]; !hasID {
		nextID, err := r.getNextID(projectID, collectionName)
		if err != nil {
			return nil, fmt.Errorf("failed to generate ID: %v", err)
		}
		data["id"] = nextID
	}

	doc := &model.MockDocument{
		ProjectID:      projectID,
		CollectionName: collectionName,
		Data:           data,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %v", err)
	}

	doc.ID = result.InsertedID.(bson.ObjectID)
	return doc, nil
}

// getNextID получает следующий автоинкрементный ID для коллекции проекта
func (r *DocumentRepository) getNextID(projectID int64, collectionName string) (int, error) {
	counterCollection := r.client.Database(r.dbName).Collection("counters")
	ctx := context.Background()

	filter := bson.M{
		"project_id":      projectID,
		"collection_name": collectionName,
	}

	update := bson.M{
		"$inc": bson.M{"seq": 1},
	}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var counter struct {
		Seq int `bson:"seq"`
	}

	err := counterCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&counter)
	if err != nil {
		return 0, fmt.Errorf("failed to get next ID: %v", err)
	}

	return counter.Seq, nil
}

// GetDocuments возвращает документы коллекции
func (r *DocumentRepository) GetDocuments(projectID int64, collectionName string, opts model.QueryOptions) (*model.DocumentsResponse, error) {
	collection := r.GetCollection(collectionName)
	ctx := context.Background()

	// Фильтр по project_id
	filter := bson.M{"project_id": projectID}

	// Настройки пагинации
	findOptions := options.Find()

	if opts.Limit != nil {
		findOptions.SetLimit(*opts.Limit)
	}
	if opts.Offset != nil {
		findOptions.SetSkip(*opts.Offset)
	}

	// Сортировка
	if opts.Sort != "" {
		order := int32(1) // asc
		if opts.Order == "desc" {
			order = int32(-1)
		}
		findOptions.SetSort(bson.D{{Key: opts.Sort, Value: order}})
	} else {
		findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	}

	// Получаем документы
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %v", err)
	}
	defer cursor.Close(ctx)

	var documents []*model.MockDocument
	if err := cursor.All(ctx, &documents); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %v", err)
	}

	// Получаем общее количество
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count documents: %v", err)
	}

	return &model.DocumentsResponse{
		Documents: documents,
		Total:     total,
		Limit:     getInt64Value(opts.Limit, 0),
		Offset:    getInt64Value(opts.Offset, 0),
	}, nil
}

// GetDocumentByID возвращает документ по ID (ищет по data.id)
func (r *DocumentRepository) GetDocumentByID(projectID int64, collectionName, documentID string) (*model.MockDocument, error) {
	collection := r.GetCollection(collectionName)
	ctx := context.Background()

	// Пробуем преобразовать в число
	var filter bson.M
	if numID, err := strconv.Atoi(documentID); err == nil {
		// Числовой ID - ищем по data.id
		filter = bson.M{"project_id": projectID, "data.id": numID}
	} else {
		// Строковый ID - может быть MongoDB ObjectID (для обратной совместимости)
		objID, err := bson.ObjectIDFromHex(documentID)
		if err != nil {
			return nil, fmt.Errorf("invalid document ID: %v", err)
		}
		filter = bson.M{"_id": objID, "project_id": projectID}
	}

	var doc model.MockDocument
	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find document: %v", err)
	}

	return &doc, nil
}

// UpdateDocument обновляет документ (ищет по data.id)
func (r *DocumentRepository) UpdateDocument(projectID int64, collectionName, documentID string, data map[string]interface{}) (*model.MockDocument, error) {
	collection := r.GetCollection(collectionName)
	ctx := context.Background()

	// Пробуем преобразовать в число
	var filter bson.M
	if numID, err := strconv.Atoi(documentID); err == nil {
		// Числовой ID - ищем по data.id и сохраняем id в data
		filter = bson.M{"project_id": projectID, "data.id": numID}
		// Убеждаемся что id не изменяется
		data["id"] = numID
	} else {
		// Строковый ID - может быть MongoDB ObjectID
		objID, err := bson.ObjectIDFromHex(documentID)
		if err != nil {
			return nil, fmt.Errorf("invalid document ID: %v", err)
		}
		filter = bson.M{"_id": objID, "project_id": projectID}
	}

	update := bson.M{
		"$set": bson.M{
			"data":       data,
			"updated_at": time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	var doc model.MockDocument
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %v", err)
	}

	return &doc, nil
}

// ResetCounter сбрасывает автоинкрементный счетчик для коллекции проекта
func (r *DocumentRepository) ResetCounter(projectID int64, collectionName string) error {
	counterCollection := r.client.Database(r.dbName).Collection("counters")
	ctx := context.Background()

	filter := bson.M{
		"project_id":      projectID,
		"collection_name": collectionName,
	}

	// Просто устанавливаем seq в 0
	update := bson.M{
		"$set": bson.M{"seq": 0},
	}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var result struct {
		Seq int `bson:"seq"`
	}

	err := counterCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return fmt.Errorf("failed to reset counter: %v", err)
	}

	slog.Info("Successfully reset counter", 
		slog.Int64("project_id", projectID), 
		slog.String("collection", collectionName), 
		slog.Int("seq", result.Seq))
	return nil
}

// DeleteDocument удаляет документ (ищет по data.id)
func (r *DocumentRepository) DeleteDocument(projectID int64, collectionName, documentID string) error {
	collection := r.GetCollection(collectionName)
	ctx := context.Background()

	// Пробуем преобразовать в число
	var filter bson.M
	if numID, err := strconv.Atoi(documentID); err == nil {
		// Числовой ID - ищем по data.id
		filter = bson.M{"project_id": projectID, "data.id": numID}
	} else {
		// Строковый ID - может быть MongoDB ObjectID
		objID, err := bson.ObjectIDFromHex(documentID)
		if err != nil {
			return fmt.Errorf("invalid document ID: %v", err)
		}
		filter = bson.M{"_id": objID, "project_id": projectID}
	}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete document: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("document not found")
	}

	// Проверяем, стала ли коллекция пустой после удаления
	count, err := r.CountDocuments(projectID, collectionName)
	if err != nil {
		return fmt.Errorf("failed to count documents after deletion: %v", err)
	}

	// Если коллекция стала пустой, сбрасываем счетчик
	if count == 0 {
		if err := r.ResetCounter(projectID, collectionName); err != nil {
			// Логируем ошибку, но не прерываем операцию удаления
			slog.Warn("Failed to reset counter", 
				slog.String("collection", collectionName), 
				slog.String("error", err.Error()))
		}
	}

	return nil
}

// DeleteAllDocuments удаляет все документы коллекции
func (r *DocumentRepository) DeleteAllDocuments(projectID int64, collectionName string) (int64, error) {
	collection := r.GetCollection(collectionName)
	ctx := context.Background()

	filter := bson.M{"project_id": projectID}

	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to delete documents: %v", err)
	}

	// Если удалили документы, сбрасываем счетчик
	if result.DeletedCount > 0 {
		slog.Info("Deleted documents, resetting counter", 
			slog.Int64("deleted_count", result.DeletedCount), 
			slog.String("collection", collectionName), 
			slog.Int64("project_id", projectID))
		if err := r.ResetCounter(projectID, collectionName); err != nil {
			// Логируем ошибку, но не прерываем операцию удаления
			slog.Warn("Failed to reset counter", 
				slog.String("collection", collectionName), 
				slog.String("error", err.Error()))
		} else {
			slog.Info("Successfully reset counter", 
				slog.String("collection", collectionName))
		}
	}

	return result.DeletedCount, nil
}

// CountDocuments возвращает количество документов в коллекции
func (r *DocumentRepository) CountDocuments(projectID int64, collectionName string) (int64, error) {
	collection := r.GetCollection(collectionName)
	ctx := context.Background()

	filter := bson.M{"project_id": projectID}
	return collection.CountDocuments(ctx, filter)
}

// Вспомогательная функция
func getInt64Value(ptr *int64, defaultValue int64) int64 {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}
