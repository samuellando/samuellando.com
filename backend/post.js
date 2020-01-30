import uuid from "uuid";
import * as authorizationLib from "./libs/authorization-lib";
import * as dynamoDbLib from "./libs/dynamodb-lib";
import { success, failure } from "./libs/response-lib";

export async function main(event, context) {
  const data = JSON.parse(event.body);
  const params = {
    TableName: process.env.tableName,

    Item: {
      pageid: uuid.v1(),
      title: data.title,
      content: data.content,
      private: data.private,
    }
  };

  try {
    await dynamoDbLib.call("put", params);
    event.pathParameters.pageid = params.Item.pageid;
    if (authorizationLib.post(event, force=true)) {
      return success(params.Item);
    } else {
      return failure({status: false});
    }
  } catch (e) {
    return failure({ status: false });
  }
}