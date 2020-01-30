import * as dynamoDbLib from "./libs/dynamodb-lib";
import * as authorizationLib from "./libs/authorization-lib";
import { success, failure } from "./libs/response-lib";
import * as get from './get';

export async function main(event, context) {
  const params = {
    TableName: process.env.tableName,

    KeyConditionExpression: "userid = :userid",
    ExpressionAttributeValues: {
      ":userid": event.requestContext.identity.cognitoIdentityId
    }
  };

  try {
    var result = await dynamoDbLib.call("query", params);
    const auth = authorizationLib.list(event);
    for (var i = 0; i < auth.length; i++) {
      event.pathParameters.pageid = auth[i].pageid;
      results.Items = result.Items.push(await get.get(event));
    }
    // Return the matching list of items in response body
    return success(result.Items);
  } catch (e) {
    return failure({ status: false });
  }
}