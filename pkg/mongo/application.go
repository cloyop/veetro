package mongo

import (
	"context"

	"github.com/cloyop/veetro/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *MongoStorage) GetAllApplications(offerId, employeeId string) (*[]storage.Application, error) {
	filters := parseApplyIdsFilter("", offerId, employeeId)
	apps := &[]storage.Application{}
	return apps, getAllOffers(s.mongoDB.Collection("applications"), filters, apps, options.Find())
}
func (s *MongoStorage) GetApplication(id, offerId, employeeId string) (*storage.Application, bool, error) {
	filters := parseApplyIdsFilter(id, offerId, employeeId)
	appl := &storage.Application{}
	is, err := getOne(s.mongoDB.Collection("applications"), filters, appl, options.FindOne())
	return appl, is, err
}
func (s *MongoStorage) CreateApplication(a *storage.Application) error {
	_, err := s.mongoDB.Collection("applications").InsertOne(context.Background(), a, options.InsertOne())
	return err
}
func (s *MongoStorage) ApplicationExist(id, offerId, employeeId string) (bool, error) {
	filters := parseApplyIdsFilter(id, offerId, employeeId)
	n, err := howManyDocs(s.mongoDB.Collection("applications"), filters)
	return n > 0, err
}
func (s *MongoStorage) DeleteApplication(id, offerId, employeeId string) (bool, error) {
	return deleteOne(s.mongoDB.Collection("applications"), parseApplyIdsFilter(id, offerId, employeeId))
}
func (s *MongoStorage) DeleteApplicationsByIds(ids *[]string, employeeId string) (int64, error) {
	f := bson.D{{Key: "id", Value: bson.D{{Key: "$in", Value: ids}}}}
	if employeeId != "" {
		f = append(f, bson.E{Key: "employee_id", Value: employeeId})
	}
	return deleteMany(s.mongoDB.Collection("applications"), f)
}

func (s *MongoStorage) UpdateApplication(applyId, employee_id string, updts *map[string]string) (bool, bool, error) {
	f := bson.D{}
	for key, value := range *updts {
		f = append(f, bson.E{Key: key, Value: value})
	}
	return editOne(s.mongoDB.Collection("applications"), parseApplyIdsFilter(applyId, "", employee_id), f)
}

func parseApplyIdsFilter(id, offerId, employeeId string) bson.D {
	filter := bson.D{}
	if id != "" {
		filter = append(filter, bson.E{Key: "id", Value: id})
	}
	if offerId != "" {
		filter = append(filter, bson.E{Key: "offer_id", Value: offerId})
	}
	if employeeId != "" {
		filter = append(filter, bson.E{Key: "employee_id", Value: employeeId})
	}
	return filter
}
