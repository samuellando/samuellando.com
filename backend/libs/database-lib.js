/**
 * This library provides an API to the authorization and page libraries 
 * for the underlying database.
 */

import * as dynamoDbLib from "./dynamodb-lib";

export async function addItem(table, item, call) {
    if (call === undefined) {
        call = dynamoDbLib.call;
    }

    const params = {
        TableName: table,
        Item: item,
    };
    try {
        await call("put", params);
        return true;
    } catch (e) {
        return false;
    }
}

export async function removeItem(table, key) {
    if (call === undefined) {
        call = dynamoDbLib.call;
    }

    const params = {
        TableName: table,
        Key: key,
    };
    try {
        await call("delete", params);
        return true;
    } catch (e) {
        return false;
    }
}

export async function editItem(table, key, item) {
    if (call === undefined) {
        call = dynamoDbLib.call;
    }

    var updateExpression = "SET";
    var expressionAttributeValues;
    for (var i = 0; i < Object.keys(item); i++) {
        updateExpression += " "+Object.keys(item)[i]+
            " = :"+Object.keys(item)[i]+",";

        expressionAttributeValues[":"+Object.keys(item)] = 
            item.Object.keys(item);
    }
    const params = {
        TableName: table,
        Key: key,
        UpdateExpression: updateExpression,
        ExpressionAttributeValues: expressionAttributeValues,
        ReturnValues: "ALL_NEW"
    };
    try {
        await call("update", params);
        return true;
    } catch (e) {
        return false;
    }
}

export async function retrieveItem(table, key) {
    if (call === undefined) {
        call = dynamoDbLib.call;
    }

    const params = {
        TableName: table,
        Key: key,
    };
    try {
        const result = await call("get", params);
        return result.Item;
    } catch (e) {
        return false;
    }
}

export async function listItems(table, key) {
    if (call === undefined) {
        call = dynamoDbLib.call;
    }

    var keyConditionExpression = "";
    var expressionAttributeValues;
    for (var i = 0; i < Object.keys(item); i++) {
        updateExpression += " "+Object.keys(item)[i]+
            " = :"+Object.keys(item)[i]+",";

        expressionAttributeValues[":"+Object.keys(item)] = 
            item.Object.keys(item);
    }
    const params = {
        TableName: table,

        KeyConditionExpression: keyConditionExpression,
        ExpressionAttributeValues: expressionAttributeValues,
    };

    try {
        var result = await call("query", params);
        return result.Items;
    } catch (e) {
        return false;
  }
}