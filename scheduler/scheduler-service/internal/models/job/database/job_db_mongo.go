package database

import (
	"context"
	"time"

	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/internal/models"
	jobErrors "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/pkg/job_errors"
	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/pkg/utils"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	jobsDB         = "jobs"
	jobsCollection = "jobs"
)

type jobDBMongo struct {
	mongoDB *mongo.Client
}

func NewJobDBMongo(mongoDB *mongo.Client) *jobDBMongo {
	return &jobDBMongo{mongoDB: mongoDB}
}

//implement models.JobDB interface for jobDBMongo

func (db jobDBMongo) Create(ctx context.Context, job *models.Job) (*models.Job, error) {

	collection := db.mongoDB.Database(jobsDB).Collection(jobsCollection)

	job.CreatedAt = time.Now().UTC()
	job.UpdatedAt = time.Now().UTC()

	result, err := collection.InsertOne(ctx, job, &options.InsertOneOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "InsertOne")
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.Wrap(jobErrors.ErrObjectIDTypeConversion, "result.InsertedID")
	}

	job.Id = objectID

	return job, nil
}

func (db jobDBMongo) Update(ctx context.Context, job *models.Job) (*models.Job, error) {

	collection := db.mongoDB.Database(jobsDB).Collection(jobsCollection)

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(true)

	var _job models.Job
	if err := collection.FindOneAndUpdate(ctx, bson.M{"_id": job.Id}, bson.M{"$set": job}, ops).Decode(&_job); err != nil {
		return nil, errors.Wrap(err, "Decode")
	}

	return &_job, nil
}

func (db jobDBMongo) GetByID(ctx context.Context, id string) (*models.Job, error) {

	jobID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	collection := db.mongoDB.Database(jobsDB).Collection(jobsCollection)

	var _job models.Job
	if err := collection.FindOne(ctx, bson.M{"_id": jobID}).Decode(&_job); err != nil {
		return nil, errors.Wrap(err, "Decode")
	}

	return &_job, nil
}
func (db jobDBMongo) DeleteByID(ctx context.Context, id string) error {

	jobID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	collection := db.mongoDB.Database(jobsDB).Collection(jobsCollection)
	deleteResult, err2 := collection.DeleteOne(ctx, bson.M{"_id": jobID})
	if err2 != nil {
		return errors.Wrap(err2, "Decode")
	}
	if deleteResult.DeletedCount == 0 {
		return errors.New("No document deleted")
	}

	return nil
}

func (db jobDBMongo) GetByScheduledKey(ctx context.Context, jobScheduledKey int) (*models.Job, error) {
	collection := db.mongoDB.Database(jobsDB).Collection(jobsCollection)

	var _job models.Job
	if err := collection.FindOne(ctx, bson.M{"scheduledKey": jobScheduledKey}).Decode(&_job); err != nil {
		return nil, errors.Wrap(err, "Decode")
	}

	return &_job, nil
}

func (db jobDBMongo) DeleteByScheduledKey(ctx context.Context, jobScheduledKey int) error {
	collection := db.mongoDB.Database(jobsDB).Collection(jobsCollection)
	_, err := collection.DeleteOne(ctx, bson.M{"scheduledKey": jobScheduledKey})
	if err != nil {
		return errors.Wrap(err, "Decode")
	}
	return nil
}

func (db jobDBMongo) ListALL(ctx context.Context, size int64, page int64) (*models.JobsList, error) {

	pagination := utils.NewPaginationQuery(int(size), int(page))

	collection := db.mongoDB.Database(jobsDB).Collection(jobsCollection)

	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return &models.JobsList{
			TotalCount: 0,
			TotalPages: 0,
			Page:       0,
			Size:       0,
			HasMore:    false,
			Jobs:       make([]*models.Job, 0),
		}, nil
	}

	limit := int64(pagination.GetLimit())
	skip := int64(pagination.GetOffset())

	cursor, err := collection.Find(ctx, bson.D{}, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor and retrieve all jobs
	var jobs = make([]*models.Job, 0, pagination.GetSize())
	for cursor.Next(ctx) {
		var job models.Job
		if err := cursor.Decode(&job); err != nil {
			return nil, err
		}
		jobs = append(jobs, &job)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &models.JobsList{
		TotalCount: count,
		TotalPages: int64(pagination.GetTotalPages(int(count))),
		Page:       int64(pagination.GetPage()),
		Size:       int64(pagination.GetSize()),
		HasMore:    pagination.GetHasMore(int(count)),
		Jobs:       jobs,
	}, nil
}
