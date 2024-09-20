package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Transaction struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Payer     string             `json:"payer"`
	Points    int                `json:"points"`
	Timestamp time.Time          `json:"timestamp"`
}

func GetTransactions(ctx context.Context, collection *mongo.Collection) ([]Transaction, error) {
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var transactions []Transaction
	if err := cur.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func SpendPoints(ctx context.Context, collection *mongo.Collection, transactions []Transaction, pointsToSpend int) ([]map[string]interface{}, error) {
	spendResults := []map[string]interface{}{}
	for _, transaction := range transactions {
		if pointsToSpend <= 0 {
			break
		}

		pointsToDeduct := 0
		if transaction.Points <= pointsToSpend {
			pointsToDeduct = transaction.Points
		} else {
			pointsToDeduct = pointsToSpend
		}

		pointsToSpend -= pointsToDeduct
		newPoints := transaction.Points - pointsToDeduct

		_, err := collection.UpdateOne(ctx, bson.M{"_id": transaction.ID}, bson.M{
			"$set": bson.M{"points": newPoints},
		})
		if err != nil {
			return nil, err
		}

		spendResults = append(spendResults, map[string]interface{}{
			"payer":  transaction.Payer,
			"points": -pointsToDeduct,
		})
	}

	return spendResults, nil
}

func GetBalance(ctx context.Context, collection *mongo.Collection) (map[string]int, error) {
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var transactions []Transaction
	if err := cur.All(ctx, &transactions); err != nil {
		return nil, err
	}

	balance := make(map[string]int)
	for _, transaction := range transactions {
		balance[transaction.Payer] += transaction.Points
	}

	return balance, nil
}
