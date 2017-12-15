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
	"regexp"
)

func getLight(ctx *gin.Context) {
	clientId := ctx.GetHeader("X_CLIENT_ID")
	lightId := ctx.Param("id")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		// TODO handle exception
		log.Fatalln(err, err.Error())
	}

	svc := dynamodb.New(sess)

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ClientID": {
				S: aws.String(clientId),
			},
			"LightID": {
				S: aws.String(lightId),
			},
		},
		TableName: aws.String("buildlights"),
	}

	result, err := svc.GetItem(input)
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
		return
	}

	lightConfig := LightConfig{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &lightConfig)
	if err != nil {
		// TODO handle error
	}

	ctx.JSON(http.StatusOK, lightConfig)
}

func getLightConfigs(ctx *gin.Context) {
	clientId := ctx.GetHeader("X_CLIENT_ID")

	lightConfigs := getLightCfgs(clientId)

	ctx.JSON(http.StatusOK, lightConfigs)
}

func getLightCfgs(clientId string) []LightConfig {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		// TODO handle exception
		log.Fatalln(err, err.Error())
	}

	svc := dynamodb.New(sess)

	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":clientId": {
				S: aws.String(clientId),
			},
		},
		KeyConditionExpression: aws.String("ClientID = :clientId"),
		TableName:              aws.String("buildlights"),
	}

	// TODO proper error handling
	result, err := svc.Query(input)
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

	lightConfigs := []LightConfig{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &lightConfigs)
	if err != nil {
		// TODO handle error
	}

	return lightConfigs
}

func putLightConfig(ctx *gin.Context) {
	clientId := ctx.GetHeader("X_CLIENT_ID")
	lightId := ctx.Param("id")

	// TODO validate inputs for security
	var lightConfig LightConfig
	ctx.BindJSON(&lightConfig)

	lightConfig.ClientID = clientId
	lightConfig.LightID = lightId

	createLightConfig(lightConfig)

	// TODO update return code for error messages
	ctx.Status(http.StatusOK)
}

func createLightConfig(config LightConfig) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		// TODO handle exception
		log.Fatalln(err, err.Error())
	}

	svc := dynamodb.New(sess)

	dynamoDocument, err := dynamodbattribute.MarshalMap(config)
	if err != nil {
		// TODO handle exception
		log.Fatalln(err, err.Error())
	}

	// TODO Optomistic locking?
	putRequest := &dynamodb.PutItemInput{
		Item:      dynamoDocument,
		TableName: aws.String("buildlights"),
	}

	_, err = svc.PutItem(putRequest)
	if err != nil {
		// TODO: Handle Error, but we can ignore ConditionalCheckFailedException
		fmt.Println(err.Error())
	}
}

func getLightStatus(ctx *gin.Context) {

	clientId := ctx.GetHeader("X_CLIENT_ID")
	lightConfigs := getLightCfgs(clientId)
	buildStats := getFailedBuilds(clientId)

	currentStatus := make(map[string]bool)

	for _, lightConfig := range lightConfigs {
		lightStatus := true

		for _, buildStat := range buildStats {
			matched, _ := regexp.MatchString(lightConfig.JobIDRegEx, buildStat.JobID)
			if matched {
				if !buildStat.BuildStatus {
					lightStatus = false
					break
				}
			}
		}

		currentStatus[lightConfig.LightID] = lightStatus
	}

	ctx.JSON(http.StatusOK, currentStatus)
}
