/**
 * This library provides an API to the authorization and page libraries 
 * for the underlying database.
 */

import * as dynamoDbLib from "./dynamodb-lib";

export async function addItem(table, item, call) {
    if (!call) {
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

export async function removeItem(table, key, call) {
    if (!call) {
        call = dynamoDbLib.call;
    }

    const params = {
        TableName: table,
        Key: key
    };
    try {
        await call("delete", params);
        return true;
    } catch (e) {
        console.log(e);
        return false;
    }
}

export async function editItem(table, key, item, call) {
    if (call === undefined) {
        call = dynamoDbLib.call;
    }

    var updateExpression = "SET";
    var expressionAttributeValues  = {};
    var expressionAttributeNames = {};
    for (var i = 0; i < Object.keys(item).length; i++) {
        updateExpression += " #"+Object.keys(item)[i]+
            " = :"+Object.keys(item)[i]+",";

        expressionAttributeValues[":"+Object.keys(item)[i]] = 
            item[Object.keys(item)[i]];

        expressionAttributeNames["#"+Object.keys(item)[i]] = 
            Object.keys(item)[i];
    }
    const params = {
        TableName: table,
        Key: key,
        UpdateExpression: updateExpression.substr(0, updateExpression.length - 1),
        ExpressionAttributeValues: expressionAttributeValues,
        ExpressionAttributeNames: expressionAttributeNames,
        ReturnValues: "ALL_NEW"
    };
    try {
        await call("update", params);
        return true;
    } catch (e) {
        return false;
    }
}

export async function retrieveItem(table, key, call) {
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

export async function listItems(table, key, call) {
    if (call === undefined) {
        call = dynamoDbLib.call;
    }

    var keyConditionExpression = "";
    var expressionAttributeValues = {};
    var expressionAttributeNames = {};
    for (var i = 0; i < Object.keys(key).length; i++) {
        keyConditionExpression += " #"+Object.keys(key)[i]+
            " = :"+Object.keys(key)[i]+",";

        expressionAttributeValues[":"+Object.keys(key)[i]] = 
            key[Object.keys(key)[i]];

        expressionAttributeNames["#"+Object.keys(key)[i]] = 
            Object.keys(key)[i];
    }
    const params = {
        TableName: table,

        KeyConditionExpression: keyConditionExpression.substr(0, keyConditionExpression.length - 1),
        ExpressionAttributeValues: expressionAttributeValues,
        ExpressionAttributeNames: expressionAttributeNames,
    };

    try {
        var result = await call("query", params);
        return result;
    } catch (e) {
        return false;
  }
}