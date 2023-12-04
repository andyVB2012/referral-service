package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/parsaakbari1209/go-mongo-crud-rest-api/model"
)

var (
	ErrCodeNotFound      = errors.New("Referral code not found")
	ErrUserNotFound      = errors.New("User not found")
	ErrAttributionFailed = errors.New("Attribution failed")
	ErrCodeError         = errors.New("Referral code error")
)

type repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) Repository {
	return &repository{db: db}
}

func (r repository) CreateReferralCode(ctx context.Context, traderAddr string) (model.Code, error) {
	// create unique code
	// check if code exists
	// if exists, create new code
	// if not exists, create code

	code := generateCode()
	// check if code exsists

	doc := bson.M{"traderAddr": traderAddr, "code": code}
	_, err1 := r.db.
		Collection("StfxReferralCodes").
		InsertOne(ctx, doc)
	if err1 != nil {
		return model.Code{}, err1
	}
	return model.Code{TraderAddr: traderAddr, Code: code}, nil

}

func (r repository) GetReferralCode(ctx context.Context, traderAddr string) (model.Code, error) {
	var out Code
	err := r.db.
		Collection("StfxReferralCodes").
		FindOne(ctx, bson.M{"traderAddr": traderAddr}).
		Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Code{}, ErrCodeNotFound
		}
		return model.Code{}, err
	}
	return toModel(out), nil
}

func generateCode() string {
	return "test"
}

func (r repository) CreateAttribution(ctx context.Context, attribution model.Attribution) (model.Attribution, error) {
	var out model.Attribution
	err := r.db.
		Collection("StfxAttributions").
		FindOne(ctx, bson.M{"code": attribution.Code}).
		Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			doc := bson.M{"code": attribution.Code, "traderAddr": attribution.TraderAddr}
			_, err1 := r.db.
				Collection("StfxAttributions").
				InsertOne(ctx, doc)
			if err1 != nil {
				return model.Attribution{}, err1
			}
			return attribution, nil
		}
		return model.Attribution{}, err
	}
	return toModel(out), ErrAttributionFailed
}

func (r repository) GetAttributions(ctx context.Context, traderAddr string) (model.AttributorData, error) {
	var results []model.Attribution
	cur, err := r.db.
		Collection("StfxAttributions").
		Find(ctx, bson.M{"traderAddr": traderAddr})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.AttributorData{}, ErrUserNotFound
		}
		return model.AttributorData{}, err
	}

	for cur.Next(context.TODO()) {
		var elem Attribution
		err := cur.Decode(&elem)
		if err != nil {
			panic(err)
		}
		results = append(results, toModel(elem))
	}

	if err := cur.Err(); err != nil {
		panic(err)
	}

	//Close the cursor once finished
	cur.Close(context.TODO())

	return model.AttributorData{TraderAddr: traderAddr, Attributions: results}, nil
}

func (r repository) GetTotalAttributionAmount(ctx context.Context, traderAddr string) (model.AttributionAmount, error) {
}

func (r repository) GetAttributionByUser(ctx context.Context, traderAddr string, attributer string) (model.AttributionAmount, error) {
}

func toModel(in Attribution) model.Attribution {
	return model.Attribution{
		TraderAddr: in.TraderAddr,
		Code:       in.Code,
	}
}

func fromModel(in model.Attribution) Attribution {
	return Attribution{
		TraderAddr: in.TraderAddr,
		Code:       in.Code,
	}
}

// func (r repository) GetFollow(ctx context.Context, user1 string, user2 string) (model.Follow, error) {
// 	var out Follow
// 	err := r.db.
// 		Collection("StfxFollows").
// 		FindOne(ctx, bson.M{"user1": user1, "user2": user2}).
// 		Decode(&out)
// 	if err != nil {
// 		if errors.Is(err, mongo.ErrNoDocuments) {
// 			return model.Follow{}, ErrUserNotFound
// 		}
// 		return model.Follow{}, err
// 	}
// 	return toModel(out), nil
// }

// func (r repository) GetFollowings(ctx context.Context, user1 string) ([]model.Follow, error) {
// 	var results []model.Follow
// 	cur, err := r.db.
// 		Collection("StfxFollows").
// 		Find(ctx, bson.M{"user1": user1})
// 	if err != nil {
// 		if errors.Is(err, mongo.ErrNoDocuments) {
// 			return []model.Follow{}, ErrUserNotFound
// 		}
// 		return []model.Follow{}, err
// 	}

// 	for cur.Next(context.TODO()) {
// 		var elem Follow
// 		err := cur.Decode(&elem)
// 		if err != nil {
// 			panic(err)
// 		}
// 		results = append(results, toModel(elem))
// 	}

// 	if err := cur.Err(); err != nil {
// 		panic(err)
// 	}

// 	//Close the cursor once finished
// 	cur.Close(context.TODO())

// 	return results, nil
// }

// func (r repository) GetFollowers(ctx context.Context, user2 string) ([]model.Follow, error) {
// 	var results []model.Follow
// 	cur, err := r.db.
// 		Collection("StfxFollows").
// 		Find(ctx, bson.M{"user2": user2})
// 	if err != nil {
// 		if errors.Is(err, mongo.ErrNoDocuments) {
// 			return []model.Follow{}, ErrUserNotFound
// 		}
// 		return []model.Follow{}, err
// 	}

// 	for cur.Next(context.TODO()) {
// 		var elem Follow
// 		err := cur.Decode(&elem)
// 		if err != nil {
// 			panic(err)
// 		}
// 		results = append(results, toModel(elem))
// 	}

// 	if err := cur.Err(); err != nil {
// 		panic(err)
// 	}

// 	//Close the cursor once finished
// 	cur.Close(context.TODO())

// 	return results, nil
// }

// func (r repository) GetAll(ctx context.Context) ([]model.Follow, error) {
// 	var results []model.Follow
// 	cur, err := r.db.
// 		Collection("StfxFollows").
// 		Find(ctx, bson.M{})
// 	if err != nil {
// 		if errors.Is(err, mongo.ErrNoDocuments) {
// 			return []model.Follow{}, ErrUserNotFound
// 		}
// 		return []model.Follow{}, err
// 	}

// 	for cur.Next(context.TODO()) {
// 		var elem Follow
// 		err := cur.Decode(&elem)
// 		if err != nil {
// 			panic(err)
// 		}
// 		results = append(results, toModel(elem))
// 	}

// 	if err := cur.Err(); err != nil {
// 		panic(err)
// 	}

// 	//Close the cursor once finished
// 	cur.Close(context.TODO())

// 	return results, nil
// }

// func (r repository) CreateFollow(ctx context.Context, follow model.Follow) (model.Follow, error) {
// 	var out Follow
// 	err := r.db.
// 		Collection("StfxFollows").
// 		FindOne(ctx, bson.M{"user1": follow.User1, "user2": follow.User2}).
// 		Decode(&out)
// 	if err != nil {
// 		if errors.Is(err, mongo.ErrNoDocuments) {
// 			doc := bson.M{"user1": follow.User1, "user2": follow.User2}
// 			_, err1 := r.db.
// 				Collection("StfxFollows").
// 				InsertOne(ctx, doc)
// 			if err1 != nil {
// 				return model.Follow{}, err1
// 			}
// 			return follow, nil
// 		}
// 		return model.Follow{}, err
// 	}
// 	return toModel(out), ErrUserFound
// }

// func (r repository) CreateFollowBatch(ctx context.Context, follows []model.Follow) ([]model.Follow, error) {
// 	wm := make([]mongo.WriteModel, len(follows))
// 	for i, follow := range follows {
// 		wm[i] = mongo.NewReplaceOneModel().
// 			SetUpsert(true).
// 			SetFilter(bson.M{"user1": follow.User1, "user2": follow.User2}).
// 			SetReplacement(follow)
// 	}
// 	_, err := r.db.
// 		Collection("StfxFollows").
// 		BulkWrite(ctx, wm)
// 	if err != nil {
// 		return []model.Follow{}, err
// 	}
// 	return follows, nil
// }

// func (r repository) DeleteFollow(ctx context.Context, user1 string, user2 string) error {
// 	out, err := r.db.
// 		Collection("StfxFollows").
// 		DeleteOne(ctx, bson.M{"user1": user1, "user2": user2})
// 	if err != nil {
// 		return err
// 	}
// 	if out.DeletedCount == 0 {
// 		return ErrUserNotFound
// 	}
// 	return nil
// }

// type Follow struct {
// 	User1 string `json:"user1"`
// 	User2 string `json:"user2"`
// }

// func fromModel(in model.Follow) Follow {
// 	return Follow{
// 		User1: in.User1,
// 		User2: in.User2,
// 	}
// }

// func toModel(in Follow) model.Follow {
// 	return model.Follow{
// 		User1: in.User1,
// 		User2: in.User2,
// 	}
// }
