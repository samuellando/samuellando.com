import * as dynamoDbLib from "./libs/dynamodb-lib";
import * as authorizationLib from "./libs/authorization-lib";
import { success, failure } from "./libs/response-lib";

export async function main(event, context) {
  const params = {
    TableName: process.env.tableName,
    Key: {
      pageid: event.pathParameters.id
    }
  };

  try {
    if (authorizationLib.get(event) != 0) {
      return false;
    }
    await dynamoDbLib.call("delete", params);
    return success({ status: true });
  } catch (e) {
    return failure({ status: false });
  }
}