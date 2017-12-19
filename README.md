# BuildSignal Raspberry Pi Client

AWS Lambda function for managing the state of build lights as part of a [buildsignal](https://buildsignal.github.io) implementation  

## Building
This code utilizes the [aws-lambda-go-shim](https://github.com/eawsy/aws-lambda-go-shim) to run Go code as a lambda.  
You therefore need to setup the shim which enables compiling against a docker image simulating the lambda environment/OS.  
1. docker pull eawsy/aws-lambda-go-shim:latest
2. go get -u -d github.com/eawsy/aws-lambda-go-core/...

To build the project, simply run `make`

## Deployment
TODO: create cloud formation template to perform deployment

Upload the `handler.zip` file created from the build process as a lambda function.  Set the runtime to `Python 2.7` and the handler to `handler.Handle`  

Add an AWS API Gateway as a trigger to the lambda.  Map the root of the API as a `{proxy+}` endpoint

Create two DynamoDB tables  
`buildlights`:  Primary partition key (String): `ClientID`, Primary sort key (String): `LightID`  
`buildstatus`:  Primary partition key (String): `ClientID`, Primary sort key (String): `JobID`

# Configuration
Each light must be configured, either by calling the api with a `PUT` request to `/buildstatus/:id` or by adding a record directly to the `buildlights` table

The configuration of each light looks like:  
```
{  
  "ClientID": "{Any Unique ID}",  
  "Description": "Light 1",  
  "JobIDRegEx": "TEST.*",  
  "LightID": "1"  
}  
```

The `ClientID` is a unique ID that must match the configuration in the CI server and the raspberry pi client  
The `JobIDRegEx` is a regex expression that will be matched against the Job IDs in the status table to find the corresponding builds to include for this light  
The `Description` is only to document the purpose of this light  
The `LightID` is a unique ID within the context of the ClientID that represents a single light  