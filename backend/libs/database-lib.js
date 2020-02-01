/**
 * This library provides an API to the authorization and page libraries 
 * for the underlying database.
 */

import * as dynamoDbLib from "./dynamodb-lib";

export function addItem(table, item) {
    const params = {
        TableName: table,
        Item: item,
    };
    try {
        await dynamoDbLib.call("put", params);
        return true;
    } catch (e) {
        return false;
    }
}

export function removeItem(table, key) {
    const params = {
        TableName: table,
        Key: key,
    };
    try {
        await dynamoDbLib.call("delete", params);
        return true;
    } catch (e) {
        return false;
    }
}

export function editItem(table, key, item) {
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
        await dynamoDbLib.call("update", params);
        return true;
    } catch (e) {
        return false;
    }
}

export function retrieveItem(table, key) {
    const params = {
        TableName: table,
        Key: key,
    };
    try {
        const result = await dynamoDbLib.call("get", params);
        return result.Item;
    } catch (e) {
        return false;
    }
}

export function listItems(table, key) {
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
        var result = await dynamoDbLib.call("query", params);
        return result.Items;
    } catch (e) {
        return false;
  }
}