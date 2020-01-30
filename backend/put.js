import * as dynamoDbLib from "./libs/dynamodb-lib";
import * as authorizationLib from "./libs/authorization-lib";
import { success, failure } from "./libs/response-lib";


export async function main(event, context) {
  const data = JSON.parse(event.body);
  const params = {
    TableName: process.env.tableName,
    Key: {
      pageid: event.pathParameters.id
    },
    UpdateExpression: "SET content = :content, title = :title, private = :private",
    ExpressionAttributeValues: {
      ":content": data.content || null,
      ":title": data.title || null,
      ":private": data.private || null,
    },
    ReturnValues: "ALL_NEW"
  };

  try {
    if (authorizationLib.get(event) != 0) {
      return failure({ status: false });
    }
    await dynamoDbLib.call("update", params);
    return success({ status: true });
  } catch (e) {
    return failure({ status: false });
  }
}