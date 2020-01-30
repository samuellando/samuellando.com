import * as dynamoDbLib from "./libs/dynamodb-lib";

export async function get(event) {
    const params = {
        TableName: process.env.authTableName,

        Key: {
        userid: event.requestContext.identity.cognitoIdentityId,
        pageid: event.pathParameters.id,
        }
    };

  try {
    const result = await dynamoDbLib.call("get", params);
    if (result.Item) {
      // Return the retrieved level
      return result.Item.level;
    } else {
      return false;
    }
  } catch (e) {
    return false;
  }
}

export async function post(event, force) {
    if (force && await get(event) != 0) {
        return false;
    }

  const data = JSON.parse(event.body);
    const params = {
        TableName: process.env.authTableName,

        item: {
            userid: event.userid,
            pageid: event.pathParameters.id,
            level: data.level,
        }
    };

  try {
    await dynamoDbLib.call("put", params);
    return true;
  } catch (e) {
    return false;
  }
}

export async function del(event) {
    if (await get(event) != 0) {
        return false;
    }

    const data = JSON.parse(event.body);

    const params = {
        TableName: process.env.authTableName,

        Key: {
            userid: data.userid,
            pageid: event.pathParameters.id,
        }
    };

    try {
        await dynamoDbLib.call("delete", params);
        return true;
    } catch (e) {
        return false;
    }
}

export async function list(event) {
  const params = {
    TableName: process.env.authTableName,

    KeyConditionExpression: "userid = :userid",
    ExpressionAttributeValues: {
      ":userid": event.requestContext.identity.cognitoIdentityId,
    }
  };

  try {
    const result = await dynamoDbLib.call("query", params);
    // Return the matching list of items in response body
    return result.Items;
  } catch (e) {
    return false;
  }
}