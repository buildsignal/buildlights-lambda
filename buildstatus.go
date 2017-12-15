package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func saveBuildStatus(ctx *gin.Context) {

	clientId := ctx.GetHeader("X_CLIENT_ID")
	jobId := ctx.Param("id")

	// TODO validate inputs for security
	var status BuildStatus
	ctx.BindJSON(&status)

	status.ClientID = clientId
	status.JobID = jobId

	createItem(status)

	// TODO update return code for error messages
	ctx.Status(http.StatusOK)
}

func createItem(status BuildStatus) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		// TODO handle exception
		log.Fatalln(err, err.Error())
	}

	svc := dynamodb.New(sess)

	dynamoDocument, err := dynamodbattribute.MarshalMap(status)
	if err != nil {
		// TODO handle exception
		log.Fatalln(err, err.Error())
	}

	putRequest := &dynamodb.PutItemInput{
		Item:                dynamoDocument,
		TableName:           aws.String("buildstatus"),
		ConditionExpression: aws.String("attribute_not_exists (ClientID) or BuildNumber <= :buildno"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":buildno": {
				N: aws.String(strconv.Itoa(status.BuildNumber)),
			},
		},
	}

	_, err = svc.PutItem(putRequest)
	if err != nil {
		// TODO: Handle Error, but we can ignore ConditionalCheckFailedException
		fmt.Println(err.Error())
	}
}

func getFailedBuilds(clientId string) []BuildStatus {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		// TODO handle exception
		log.Fatalln(err, err.Error())
	}

	svc := dynamodb.New(sess)

	input := &dynamodb.ScanInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":clientId": {
				S: aws.String(clientId),
			},
			":buildStatus": {
				BOOL: aws.Bool(false),
			},
		},
		FilterExpression: aws.String("ClientID = :clientId and BuildStatus = :buildStatus"),
		TableName:        aws.String("buildstatus"),
	}

	// TODO proper error handling
	result, err := svc.Scan(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	buildStats := []BuildStatus{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &buildStats)
	if err != nil {
		// TODO handle error
	}

	return buildStats
}
