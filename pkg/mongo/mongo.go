package mongo

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	mongoDB     *mongo.Database
	mongoClient *mongo.Client
}

func New() *MongoStorage {
	return &MongoStorage{}
}
func (m *MongoStorage) Init() {
	var MONGO_URL = os.Getenv("MONGO_URL")
	bsonOpts := &options.BSONOptions{
		UseJSONStructTags: true,
		NilSliceAsEmpty:   true,
	}
	opts := options.Client().ApplyURI(MONGO_URL).SetBSONOptions(bsonOpts).SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}
	m.mongoClient, m.mongoDB = client, client.Database("veetro")
}
func (s *MongoStorage) Close() {
	if err := s.mongoClient.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

// -
func join(c *mongo.Collection, from, localField, foreignField, as string, match bson.D, obj interface{}) error {
	pipe := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: from},
			{Key: "localField", Value: localField},
			{Key: "foreignField", Value: foreignField},
			{Key: "as", Value: as},
		}},
		}, {{Key: "$match", Value: match}},
	}
	r, err := c.Aggregate(context.Background(), pipe, nil)
	if err != nil {
		return err
	}
	return r.All(context.Background(), obj)
}
func getOne(c *mongo.Collection, f bson.D, obj interface{}, opts *options.FindOneOptions) (bool, error) {
	if opts == nil {
		opts = options.FindOne()
	}
	r := c.FindOne(context.Background(), f, opts)
	if r.Err() != nil {
		if r.Err() == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, r.Err()
	}
	return true, r.Decode(obj)
}
func getAllOffers(c *mongo.Collection, f bson.D, obj interface{}, opts *options.FindOptions) error {
	if opts == nil {
		opts = options.Find()
	}
	r, err := c.Find(context.Background(), f, opts)
	if err != nil {
		return err
	}
	return r.All(context.Background(), obj)
}
func howManyDocs(c *mongo.Collection, f bson.D) (int64, error) {
	return c.CountDocuments(context.Background(), f, options.Count())
}
func deleteOne(c *mongo.Collection, f bson.D) (bool, error) {
	r, err := c.DeleteOne(context.Background(), f, options.Delete())
	if err != nil {
		return false, err
	}
	return r.DeletedCount > 0, nil
}
func deleteMany(c *mongo.Collection, f bson.D) (int64, error) {
	r, err := c.DeleteMany(context.Background(), f, options.Delete())
	return r.DeletedCount, err
}
func editOne(c *mongo.Collection, f bson.D, updts bson.D) (bool, bool, error) {
	r, err := c.UpdateOne(context.Background(), f, bson.D{{Key: "$set", Value: updts}}, options.Update())
	if err != nil {
		return false, false, err
	}
	return r.MatchedCount > 0, r.ModifiedCount > 0, err
}
