package mongo

import (
	"context"

	"github.com/cloyop/veetro/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *MongoStorage) GetUserByEmail(email string) (*storage.User, bool, error) {
	u := &storage.User{}
	is, err := getOne(s.mongoDB.Collection("users"), bson.D{{Key: "email", Value: email}}, u, nil)
	return u, is, err
}
func (s *MongoStorage) CreateUser(u *storage.User) error {
	_, err := s.mongoDB.Collection("users").InsertOne(context.Background(), u, options.InsertOne())
	return err
}
func (s *MongoStorage) UserExist(email string) (bool, error) {
	n, err := howManyDocs(s.mongoDB.Collection("users"), bson.D{{Key: "email", Value: email}})
	return n > 0, err
}
func (s *MongoStorage) UpdateUser(userId string, updts *map[string]interface{}) (bool, error) {
	updt := bson.D{}
	for key, value := range *updts {
		updt = append(updt, bson.E{Key: key, Value: value})
	}
	_, upt, err := editOne(s.mongoDB.Collection("users"), bson.D{{Key: "id", Value: userId}}, updt)
	if err != nil {
		return false, err
	}
	return upt, err
}

func (s *MongoStorage) DeleteCustomer(userId string) (bool, error) {
	offers, err := s.GetAllOffersWithApplications(userId)
	if err != nil {
		return false, err
	}
	applicationsId := []string{}
	for _, v := range *offers {
		for _, a := range v.Applications {
			applicationsId = append(applicationsId, a.Id)
		}
	}
	if len(applicationsId) > 0 {
		f := &bson.E{Key: "id", Value: bson.D{{Key: "$in", Value: applicationsId}}}
		if _, err := deleteMany(s.mongoDB.Collection("applications"), bson.D{*f}); err != nil {
			return false, err
		}
	}
	if _, err := deleteMany(s.mongoDB.Collection("offers"), bson.D{{Key: "owner_id", Value: userId}}); err != nil {
		return false, err
	}
	return deleteOne(s.mongoDB.Collection("users"), bson.D{{Key: "id", Value: userId}})
}
func (s *MongoStorage) DeleteEmployee(userId string) (bool, error) {
	if _, err := deleteMany(s.mongoDB.Collection("applications"), bson.D{{Key: "employee_id", Value: userId}}); err != nil {
		return false, err
	}
	return deleteOne(s.mongoDB.Collection("users"), bson.D{{Key: "id", Value: userId}})
}
