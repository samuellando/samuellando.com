/**
 * This library provides and API to the aws-dynamo db for the database
 * lib.
 */

import AWS from "aws-sdk";

export function call(action, params) {
  const dynamoDb = new AWS.DynamoDB.DocumentClient();

  return dynamoDb[action](params).promise();
}
