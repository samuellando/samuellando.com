import * as dbLib from "../libs/database-lib";

import AWS from "aws-sdk";
AWS.config.update({region:'ca-central-1'});

const db = new AWS.DynamoDB({endpoint: 'http://localhost:8000'});

const tableName = "pages";

function sleep(ms) {
  return new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
}

function call(action, params) {
  const callDb = new AWS.DynamoDB.DocumentClient({endpoint: 'http://localhost:8000'});
  return callDb[action](params).promise();
}

async function createTable() {
    var status = false;
    db.createTable({
            AttributeDefinitions: [{AttributeName: 'pageid', AttributeType: 'S'}],
            TableName: tableName, 
            KeySchema: [{AttributeName: 'pageid', KeyType: 'HASH'}],
            ProvisionedThroughput: {
                'ReadCapacityUnits': 5,
                'WriteCapacityUnits': 5
            },
        }, (err, data) => {
            if (err == null) {
                return true;
            } else {
                return false;
            }
        }
    );
    return status;
}

async function deleteTable() {
    var status = false;
    await db.deleteTable({
            TableName: tableName, 
        },
        (err, data) => {
            if (err === null) {
                status = true;
            } else {
                status = false;
            }
        }
    ).promise();
    return status;
}

describe('dblib', () => {
        beforeAll(async () => {
                await createTable();
                await sleep(1000);
            }
        );

        afterAll(async () => {
                await deleteTable();
            }
        );

        test("add item", async () => {
                await dbLib.addItem(tableName, {pageid: "TEST-PAGE", data: "Test Data"}, call);
                const res = await call('query', 
                    {
                        TableName: tableName, 
                        KeyConditionExpression: "pageid=:pageid", 
                        ExpressionAttributeValues: {":pageid": "TEST-PAGE"}
                    }
                );
                expect(res.Count).toEqual(1);
                expect(res.Items[0].pageid).toEqual('TEST-PAGE');
                expect(res.Items[0].data).toEqual('Test Data');
            }
        );

        test("remove item", async () => {
                await call('put', {TableName: tableName, Item: {pageid: "TEST-PAGE2", data: "Test Data"}});
                await dbLib.removeItem(tableName, {pageid: "TEST-PAGE2"}, call);
                const res = await call('query', 
                    {
                        TableName: tableName, 
                        KeyConditionExpression: "pageid=:pageid", 
                        ExpressionAttributeValues: {":pageid": "TEST-PAGE2"}
                    }
                );
                expect(res.Count).toEqual(0);
            }
        );

        test("edit item", async () => {
                await call('put', {TableName: tableName, Item: {pageid: "TEST-PAGE3", data: "Test Data"}});
                await dbLib.editItem(tableName, {pageid: "TEST-PAGE3"},
                    {data: "New data"}, call);
                const res = await call('query', 
                    {
                        TableName: tableName, 
                        KeyConditionExpression: "pageid=:pageid", 
                        ExpressionAttributeValues: {":pageid": "TEST-PAGE3"}
                    }
                );
                expect(res.Count).toEqual(1);
                expect(res.Items[0].pageid).toEqual('TEST-PAGE3');
                expect(res.Items[0].data).toEqual('New data');
            }
        );
    }
);