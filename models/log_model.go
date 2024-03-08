package models

import (
	"context"
	"errors"
	"github.com/Olprog59/golog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type LogModel struct {
	Store     *Store              `bson:"-"`
	LogId     *primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Level     string              `bson:"level"`
	Message   string              `bson:"message"`
	Timestamp time.Time           `bson:"timestamp"`
	Metadata  string              `bson:"metadata"`
}

func NewLogModel(store *Store, level, message string, timestamp time.Time, metadata string) *LogModel {
	return &LogModel{
		Store:     store,
		Level:     level,
		Message:   message,
		Timestamp: timestamp,
		Metadata:  metadata,
	}
}

func (l *LogModel) InsertLogMongo(appToken string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := l.Store.Mongo.Client.Database(appToken).Collection("logs").InsertOne(ctx, l)
	if err != nil {
		return err
	}
	golog.Info("Inserted a single log: %v", res.InsertedID)
	return nil
}

func (l *LogModel) DeleteCollectionLogMongo(appToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// remove collection
	err := l.Store.Mongo.Client.Database(appToken).Collection("logs").Drop(ctx)
	if err != nil {
		return err
	}
	golog.Info("Collection dropped")
	return nil
}

//func (l *LogModel) GetLogsMongo(appToken string) ([]LogModel, error) {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	cursor, err := l.Store.Mongo.Client.Database(appToken).Collection("logs").Find(ctx, nil)
//	if err != nil {
//		return nil, err
//	}
//	var logs []LogModel
//	if err = cursor.All(ctx, &logs); err != nil {
//		return nil, err
//	}
//	return logs, nil
//}

func (l *LogModel) GetLogByIdMongo(appToken string, logId string) (LogModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var log LogModel
	objectId, err := primitive.ObjectIDFromHex(logId)
	if err != nil {
		return log, err
	}
	err = l.Store.Mongo.Client.Database(appToken).Collection("logs").FindOne(ctx, primitive.M{"log_id": objectId}).Decode(&log)
	if err != nil {
		return log, err
	}
	return log, nil
}

func (l *LogModel) GetLogsByAppToken(appToken string) ([]LogModel, error) {
	golog.Info("App token: %s", appToken)
	collection := l.Store.Mongo.Client.Database(appToken).Collection("logs")
	if collection == nil {
		return nil, errors.New("Collection not found")
	}
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, mongo.ErrNoDocuments
		}
		if errors.Is(err, mongo.ErrNilDocument) {
			return nil, mongo.ErrNilDocument
		}
		return nil, err
	}
	var logs = make([]LogModel, 0)
	if err = cursor.All(context.TODO(), &logs); err != nil {
		golog.Err("Error getting logs: %s", err)
		return nil, err
	}
	return logs, nil
}

func (l *LogModel) GetConnection(appToken string) *mongo.Collection {
	return l.Store.Mongo.Client.Database(appToken).Collection("logs")
}

func (l *LogModel) DeleteLogByIdMongo(appToken string, logId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(logId)
	if err != nil {
		return err
	}
	_, err = l.Store.Mongo.Client.Database(appToken).Collection("logs").DeleteOne(ctx, primitive.M{"log_id": objectId})
	if err != nil {
		return err
	}
	return nil
}
