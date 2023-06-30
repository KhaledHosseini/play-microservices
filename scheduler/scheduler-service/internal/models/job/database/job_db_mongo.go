package database

import(
	"context"
	"time"

	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/internal/models"
	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/pkg/job_errors"
	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/pkg/utils"
	
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	job.JobID = objectID

	return job, nil
}

func (db jobDBMongo) Update(ctx context.Context, job *models.Job) (*models.Job, error) {

	collection := db.mongoDB.Database(jobsDB).Collection(jobsCollection)

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(true)

	var _job models.Job
	if err := collection.FindOneAndUpdate(ctx, bson.M{"_id": job.JobID}, bson.M{"$set": job}, ops).Decode(&_job); err != nil {
		return nil, errors.Wrap(err, "Decode")
	}

	return &_job, nil
}

func (db jobDBMongo) GetByID(ctx context.Context, jobID primitive.ObjectID) (*models.Job, error) {

	collection := db.mongoDB.Database(jobsDB).Collection(jobsCollection)

	var _job models.Job
	if err := collection.FindOne(ctx, bson.M{"_id": jobID}).Decode(&_job); err != nil {
		return nil, errors.Wrap(err, "Decode")
	}

	return &_job, nil
}
func (db jobDBMongo) DeleteByID(ctx context.Context, jobID primitive.ObjectID) (error) {
	
	collection := db.mongoDB.Database(jobsDB).Collection(jobsCollection)
	_, err:= collection.DeleteOne(ctx, bson.M{"_id": jobID}) 
	if err != nil{
		return errors.Wrap(err, "Decode")
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

func (db jobDBMongo) DeleteByScheduledKey(ctx context.Context, jobScheduledKey int) (error) {
	collection := db.mongoDB.Database(jobsDB).Collection(jobsCollection)
	_, err:= collection.DeleteOne(ctx, bson.M{"scheduledKey": jobScheduledKey}) 
	if err != nil{
		return errors.Wrap(err, "Decode")
	}
	return nil
}

func (db jobDBMongo) ListALL(ctx context.Context, pagination *utils.Pagination) (*models.JobsList, error) {

	return nil,nil
}