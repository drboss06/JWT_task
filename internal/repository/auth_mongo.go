package repository

import (
	JWTServiceObjects "JWTService"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthPostgres struct {
	db *mongo.Client
}

func NewAuthPostgres(db *mongo.Client) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) SetSession(guid string, session JWTServiceObjects.Session) error {
	collection := r.db.Database("test").Collection("users")
	filter := bson.M{"guid": guid}
	update := bson.M{"$set": bson.M{"guid": guid, "refresh_token": session.RefreshToken, "live_time": session.LiveTime}}

	opt := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(context.TODO(), filter, update, opt)

	if err != nil {
		return err
	}
	return nil
}

func (r *AuthPostgres) GetSession(guid string) (JWTServiceObjects.Session, error) {
	collection := r.db.Database("test").Collection("users")

	filter := bson.M{"guid": guid}

	var resault JWTServiceObjects.Session
	err := collection.FindOne(context.TODO(), filter).Decode(&resault)
	if err != nil {
		return JWTServiceObjects.Session{}, err
	}

	return resault, nil
}

func (r *AuthPostgres) SetRefreshToken(refreshToken []byte, session JWTServiceObjects.Session) error {
	collection := r.db.Database("test").Collection("users")

	filter := bson.M{"refresh_token": refreshToken}
	update := bson.D{
		{"$set", bson.D{
			{"refresh_token", session.RefreshToken},
			{"live_time", session.LiveTime},
		}},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	return nil
}
