/**
 * This library provides an API to the authorization and page libraries
 * for the underlying database.
 */

import * as dynamoDbLib from "./dynamodb-lib";

export async function addItem(table, item) {
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

export async function removeItem(table, key) {
    const params = {
        TableName: table,
        Key: key
    };
    try {
        await dynamoDbLib.call("delete", params);
        return true;
    } catch (e) {
        return false;
    }
}

export async function editItem(table, key, item) {
    var updateExpression = "SET";
    var expressionAttributeValues  = {};
    var expressionAttributeNames = {};
    for (var i = 0; i < Object.keys(item).length; i++) {
        if (item[Object.keys(item)[i]] != undefined && item[Object.keys(item)[i]] != null) {
            updateExpression += " #"+Object.keys(item)[i]+
                " = :"+Object.keys(item)[i]+",";

            expressionAttributeValues[":"+Object.keys(item)[i]] =
                item[Object.keys(item)[i]];

            expressionAttributeNames["#"+Object.keys(item)[i]] =
                Object.keys(item)[i];
        }
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
        await dynamoDbLib.call("update", params);
        return true;
    } catch (e) {
        return false;
    }
}

/**
 * @param {*} table
 * @param {*} key
 * @param {*} call
 * @returns a single item.
 */
export async function retrieveItem(table, key) {
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

/**
 * @param {*} table
 * @param {*} key
 * @param {*} call
 * @returns An object with Count int and Items array.
 */
export async function listItems(table, key) {
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
        var result = await dynamoDbLib.call("query", params);
        return result;
    } catch (e) {
        return false;
  }
}