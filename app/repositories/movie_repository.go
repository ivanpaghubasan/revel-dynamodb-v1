package repositories

import (
	"errors"
	"log"
	"revel-dynamodb-v1/app/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type MovieRepository struct {
	Client    *dynamodb.DynamoDB
	TableName string
}

func New(client *dynamodb.DynamoDB, tableName string) *MovieRepository {
	return &MovieRepository{
		Client:    client,
		TableName: tableName,
	}
}

func (s *MovieRepository) GetMovies() ([]*models.Movie, error) {
	var movies []*models.Movie

	result, err := s.Client.Scan(&dynamodb.ScanInput{
		TableName: aws.String(s.TableName),
	})

	if err != nil {
		log.Fatalf("Got an error calling scan: %s", err)
		return nil, errors.New(err.Error())
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &movies)
	if err != nil {
		log.Fatalf("Got an error calling UnmarshalListOfMaps: %s", err)
		return nil, errors.New(err.Error())
	}

	return movies, nil
}
