package mongo

import (
	"context"

	"github.com/cloyop/veetro/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *MongoStorage) GetOffer(offerId string) (*storage.Offer, bool, error) {
	filter := parseOffersIdsFilter(offerId, "")
	o := &storage.Offer{}
	is, err := getOne(s.mongoDB.Collection("offers"), filter, o, nil)
	return o, is, err
}
func (s *MongoStorage) GetAllOffers(keyword, role, location string) (*[]storage.Offer, error) {
	f := bson.D{{Key: "open", Value: true}}
	insensitiveK := bson.E{Key: "$options", Value: "i"}
	if keyword != "" {
		f = append(f, bson.E{Key: "description", Value: bson.D{{Key: "$regex", Value: keyword}, insensitiveK}})
	}
	if role != "" {
		f = append(f, bson.E{Key: "role", Value: bson.D{{Key: "$regex", Value: role}, insensitiveK}})
	}
	if location != "" {
		f = append(f, bson.E{Key: "location", Value: bson.D{{Key: "$regex", Value: location}, insensitiveK}})
	}
	offers := &[]storage.Offer{}
	return offers, getAllOffers(s.mongoDB.Collection("offers"), f, offers, nil)
}
func (s *MongoStorage) GetAllOffersWithApplications(userId string) (*[]storage.Offer, error) {
	ofs := &[]storage.Offer{}
	return ofs, join(s.mongoDB.Collection("offers"), "applications", "id", "offer_id", "applications", bson.D{{Key: "owner_id", Value: userId}}, ofs)
}
func (s *MongoStorage) GetOfferWithApplications(offerId, userId string) (*storage.Offer, bool, error) {
	filters := parseOffersIdsFilter(offerId, userId)
	ofs := []storage.Offer{}
	if err := join(s.mongoDB.Collection("offers"), "applications", "id", "offer_id", "applications", filters, &ofs); err != nil {
		return nil, false, err
	}
	if len(ofs) == 0 {
		return nil, false, nil
	}
	return &ofs[0], true, nil
}
func (s *MongoStorage) OfferExist(id string) (bool, error) {
	quantity, err := howManyDocs(s.mongoDB.Collection("offers"), bson.D{{Key: "id", Value: id}})
	return quantity > 0, err
}

func (s *MongoStorage) UpdateOffer(offerId, userId string, fs *map[string]interface{}) (bool, bool, error) {
	fls := bson.D{}
	for key, value := range *fs {
		fls = append(fls, bson.E{Key: key, Value: value})
	}
	return editOne(s.mongoDB.Collection("offers"), parseOffersIdsFilter(offerId, userId), fls)
}
func (s *MongoStorage) CreateOffer(o *storage.Offer) error {
	_, err := s.mongoDB.Collection("offers").InsertOne(context.Background(), o, options.InsertOne())
	return err
}
func (s *MongoStorage) DeleteOffer(offerId, userId string) (bool, error) {
	filters := parseOffersIdsFilter(offerId, userId)
	r, err := s.mongoDB.Collection("offers").DeleteOne(context.Background(), filters, options.Delete())
	dc := r.DeletedCount > 0
	if dc {
		deleteMany(s.mongoDB.Collection("applications"), bson.D{{Key: "offer_id", Value: offerId}})
	}
	return dc, err
}
func (s *MongoStorage) DeleteOffersByIds(ids *[]string, ownerId string) (int64, error) {
	f := bson.D{bson.E{Key: "id", Value: bson.D{{Key: "$in", Value: ids}}}}
	if ownerId != "" {
		f = append(f, bson.E{Key: "owner_id", Value: ownerId})
	}
	dc, err := deleteMany(s.mongoDB.Collection("offers"), f)
	if err != nil {
		return 0, err
	}
	if dc > 0 {
		deleteMany(s.mongoDB.Collection("applications"), bson.D{bson.E{Key: "offer_id", Value: bson.D{{Key: "$in", Value: ids}}}})
	}
	return dc, err
}
func ParseOfferFilters(fs *map[string]string) interface{} {
	filter := bson.D{}
	for key, value := range *fs {
		filter = append(filter, bson.E{Key: key, Value: value})
	}
	return filter
}
func parseOffersIdsFilter(offerId, userId string) bson.D {
	filter := bson.D{}
	if offerId != "" {
		filter = append(filter, bson.E{Key: "id", Value: offerId})
	}
	if userId != "" {
		filter = append(filter, bson.E{Key: "owner_id", Value: userId})
	}
	return filter
}
