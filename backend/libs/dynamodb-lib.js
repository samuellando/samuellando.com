/**
 * This library provides and API to the aws-dynamo db for the database
 * lib.
 */

import AWS from "aws-sdk";

var dynamoDb;

if (process.env.NODE_ENV === 'test') {
  AWS.config.update({region:'ca-central-1'});
  dynamoDb = new AWS.DynamoDB.DocumentClient({endpoint: 'http://localhost:8000'});
} else {
  dynamoDb = new AWS.DynamoDB.DocumentClient();
}

export function call(action, params) {
  return dynamoDb[action](params).promise();
}
