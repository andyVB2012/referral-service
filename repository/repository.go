package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/andyVB2012/referral-service/model"
)

var (
	ErrCodeNotFound      = errors.New("Referral code not found")
	ErrUserNotFound      = errors.New("User not found")
	ErrAttributionFailed = errors.New("Attribution failed")
	ErrCodeError         = errors.New("Referral code error")
	ErrNoTraderAddr      = errors.New("No trader address")
	ErrAlreadyAdded      = errors.New("Pre exsisting Referral Code")
)

type repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) Repository {
	return &repository{db: db}
}

func (r repository) CreateReferralCode(ctx context.Context, traderAddr string) (model.Referral, error) {
	// Define the code to be inserted
	if r.isAddressInReferralCode(ctx, traderAddr) {
		return model.Referral{}, ErrAlreadyAdded
	}

	refCode := r.generateCode()
	code := model.Referral{
		TraderAddr: traderAddr,
		Code:       refCode,
		// Code:       primitive.NewObjectID().Hex(),
	}

	// Insert the code into the database
	_, err := r.db.Collection("referral-codes").InsertOne(ctx, code)
	if err != nil {
		return model.Referral{}, err
	}
	return code, nil
}
func (r repository) AddAttributor(ctx context.Context, refCode string, traderAddr string) error {
	// Define the update operation to increment the attributors field by 1
	if r.IsTraderAddrInDb(ctx, traderAddr) {
		return ErrNoTraderAddr
	}
	collection := r.db.Collection("AttributionData")
	data := model.AttributorData{
		TraderAddress: traderAddr,
		ReferredCode:  refCode,
	}
	_, err := collection.InsertOne(ctx, data)
	return err
}

func (r repository) isAddressInReferralCode(ctx context.Context, traderAddr string) bool {
	var result bson.M
	filter := bson.M{"traderaddr": traderAddr}
	err := r.db.Collection("referral-codes").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No document found with the given refCode
			return false
		}
		return false
	}
	return true
}

func (r repository) generateCode() string {
	len, err := r.db.Collection("referral-codes").CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return "0x1"
	}
	len = len + 1
	return "0x" + strconv.FormatInt(len, 10)
}

func (r repository) GetCode(ctx context.Context, traderAddr string) (string, error) {
	var out model.Referral

	check := model.Referral{
		TraderAddr: traderAddr,
	}
	err := r.db.
		Collection("referral-codes").
		FindOne(ctx, bson.M{"traderaddr": check.TraderAddr}).Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", ErrCodeNotFound
		}
		return "", err
	}
	return out.Code, nil
}

func (r repository) IsTraderAddrInDb(ctx context.Context, traderAddr string) bool {
	check := model.AttributorData{
		TraderAddress: traderAddr,
	}
	res := r.db.
		Collection("AttributionData").
		FindOne(ctx, bson.M{"traderaddress": check.TraderAddress})

	if err := res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		return false
	}
	return true
}

func (r repository) IsCodeInDb(ctx context.Context, code string) bool {
	var result bson.M
	filter := bson.M{"code": code}
	err := r.db.Collection("referral-code").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No document found with the given refCode
			return false
		}
		return false
	}
	return true
}

func (r repository) IsTraderAccountInDb(ctx context.Context, traderAcct string) bool {
	var out model.Account
	err := r.db.
		Collection("AttributionData").
		FindOne(ctx, bson.M{"traderacct": traderAcct}).Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false
		}
		return false
	}

	return true
}

func (r repository) AddTraderAccount(ctx context.Context, traderAddr string, traderAcct string) error {
	update := bson.M{"$set": bson.M{"traderacct": traderAcct}}
	// Find the document with the given traderAddr and update it
	res := r.db.
		Collection("AttributionData").
		FindOneAndUpdate(ctx, bson.M{"traderaddress": traderAddr}, update)

	fmt.Println("res: ", res.Decode(&model.AttributorData{}))
	// Check for errors, including the case where the document does not exist
	if err := res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoTraderAddr
		}
		return err
	}
	return nil
}

func (r repository) AddTraderAcctIfNotExists(ctx context.Context, traderAcct string, traderAddr string) error {
	// Check if the trader exists in the database
	res, err := r.GetTraderAddrFromTraderAcct(ctx, traderAcct)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return r.AddTraderAccount(ctx, traderAcct, traderAddr) // Custom error for no document found
		}
		return r.AddTraderAccount(ctx, traderAcct, traderAddr)
	}
	if res != "" {
		return r.AddTraderAccount(ctx, traderAcct, traderAddr)
	}
	return nil

}

func (r repository) GetTraderAddrFromTraderAcct(ctx context.Context, traderAcct string) (string, error) {
	var out model.AttributorData
	check := model.AttributorData{
		TraderAcct: traderAcct,
	}
	err := r.db.
		Collection("AttributionData").
		FindOne(ctx, bson.M{"traderacct": check.TraderAcct}).
		Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", ErrUserNotFound
		}
		return "", err
	}
	return out.TraderAddress, nil
}

func (r repository) AddVariable(ctx context.Context, traderAddr string, variable string, value int) error {
	// Increment the specified variable by value
	update := bson.M{"$inc": bson.M{variable: value}}
	result := r.db.Collection("AttributionData").FindOneAndUpdate(ctx, bson.M{"traderaddress": traderAddr}, update)

	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoTraderAddr
		}
		return err
	}

	return nil
}

func (r repository) Unsubscribe(ctx context.Context, traderAddr string) error {
	// Define the update operation to decrement the subscriptions field by 1
	update := bson.M{
		"$inc": bson.M{"subscriptions": -1}, // Decrementing the subscriptions field by 1
	}

	// Perform the FindOneAndUpdate operation
	result := r.db.Collection("AttributionData").FindOneAndUpdate(ctx, bson.M{"address": traderAddr}, update)

	// Check for errors (e.g., if the trader doesn't exist)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoTraderAddr
		}
		return err
	}

	return nil
}

func (r repository) GetAllAttributors(ctx context.Context, refCode string) ([]model.AttributorData, error) {
	var attributors []model.AttributorData
	cursor, err := r.db.Collection("AttributionData").Find(ctx, bson.M{"referredcode": refCode})
	if err != nil {
		return attributors, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var attributor model.AttributorData
		if err = cursor.Decode(&attributor); err != nil {
			return attributors, err
		}
		attributors = append(attributors, attributor)
	}

	if err = cursor.Err(); err != nil {
		return attributors, err
	}

	return attributors, nil
}

func (r repository) GetStats(ctx context.Context, code string) (model.AttributionStats, error) {
	var stats model.AttributionStats
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "referredcode", Value: code}}}}
	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "RefferedAmount", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "TotalDeposited", Value: bson.D{{Key: "$sum", Value: "$deposited"}}},
			{Key: "TotalWithdrawn", Value: bson.D{{Key: "$sum", Value: "$withdrawn"}}},
			{Key: "TotalVaultCreated", Value: bson.D{{Key: "$sum", Value: "$vaultscreated"}}},
			{Key: "TotalManualInvested", Value: bson.D{{Key: "$sum", Value: "$manualinvted"}}},
			{Key: "TotalSubscriptionInvested", Value: bson.D{{Key: "$sum", Value: "$subscriptioninvested"}}},
			{Key: "TotalInvested", Value: bson.D{{Key: "$sum", Value: "$totalinvested"}}},
			{Key: "TotalVaultsInvestedIn", Value: bson.D{{Key: "$sum", Value: "$vaultsinvestedin"}}},
			{Key: "TotalSubscriptions", Value: bson.D{{Key: "$sum", Value: "$subscriptions"}}},
			{Key: "TotalStakedAmount", Value: bson.D{{Key: "$sum", Value: "$stakedamount"}}},
		}},
	}

	cursor, err := r.db.Collection("AttributionData").Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return stats, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(&stats)
		if err != nil {
			return stats, err
		}
	} else {
		return stats, mongo.ErrNoDocuments
	}

	return stats, nil
}
