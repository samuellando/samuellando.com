import AWS from "aws-sdk";
import * as dbLib from "../libs/database-lib";
import * as dynamoLib from "../libs/dynamodb-lib";

AWS.config.update({region:'ca-central-1'});
var db = new AWS.DynamoDB({endpoint: 'http://localhost:8000'});

var call = dynamoLib.call;

const tableName = "test-pages";

function sleep(ms) {
  return new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
}


async function createTable() {
    var status = false;
    db.createTable({
            AttributeDefinitions: [{AttributeName: 'userid', AttributeType: 'S'}, {AttributeName: 'pageid', AttributeType: 'S'}],
            TableName: tableName, 
            KeySchema: [{AttributeName: 'userid', KeyType: 'HASH'},
                {AttributeName: 'pageid', KeyType: 'RANGE'}],
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
                await dbLib.addItem(tableName, {userid: "USERID", pageid: "TEST-PAGE", data: "Test Data"});
                const res = await call('query', 
                    {
                        TableName: tableName, 
                        KeyConditionExpression: "userid=:userid AND pageid=:pageid", 
                        ExpressionAttributeValues: {":pageid": "TEST-PAGE", ":userid": "USERID"}
                    }
                );
                expect(res.Count).toEqual(1);
                expect(res.Items[0].userid).toEqual('USERID');
                expect(res.Items[0].pageid).toEqual('TEST-PAGE');
                expect(res.Items[0].data).toEqual('Test Data');
            }
        );

        test("remove item", async () => {
                await call('put', {TableName: tableName, Item: {userid: "USERID", pageid: "TEST-PAGE2", data: "Test Data"}});
                await dbLib.removeItem(tableName, {userid: "USERID", pageid: "TEST-PAGE2"});
                const res = await call('query', 
                    {
                        TableName: tableName, 
                        KeyConditionExpression: "pageid=:pageid AND userid=:userid", 
                        ExpressionAttributeValues: {":pageid": "TEST-PAGE2", ":userid": "USERID"}
                    }
                );
                expect(res.Count).toEqual(0);
            }
        );

        test("edit item", async () => {
                await call('put', {TableName: tableName, Item: {userid: "USERID", pageid: "TEST-PAGE3", data: "Test Data"}});
                await dbLib.editItem(tableName, {userid: "USERID", pageid: "TEST-PAGE3"},
                    {data: "New data"});
                const res = await call('query', 
                    {
                        TableName: tableName, 
                        KeyConditionExpression: "userid=:userid AND pageid=:pageid", 
                        ExpressionAttributeValues: {":pageid": "TEST-PAGE3", ":userid": "USERID"}
                    }
                );
                expect(res.Count).toEqual(1);
                expect(res.Items[0].userid).toEqual('USERID');
                expect(res.Items[0].pageid).toEqual('TEST-PAGE3');
                expect(res.Items[0].data).toEqual('New data');
            }
        );

        test("retrieve item", async () => {
                await call('put', {TableName: tableName, Item: {userid: "USERID", pageid: "TEST-PAGE4", data: "Test Data"}});
                const res = await dbLib.retrieveItem(tableName, {userid: "USERID", pageid: "TEST-PAGE4"});
                expect(res.userid).toEqual('USERID');
                expect(res.pageid).toEqual('TEST-PAGE4');
                expect(res.data).toEqual('Test Data');
            }
        );

        test("list items", async () => {
                await call('put', {TableName: tableName, Item: {userid: "USERID2", pageid: "TEST-PAGE5", data: "Test Data 1"}});
                await call('put', {TableName: tableName, Item: {userid: "USERID2", pageid: "TEST-PAGE6", data: "Test Data 2"}});
                const res = await dbLib.listItems(tableName, {userid: "USERID2"});

                expect(res.Count).toEqual(2);
                expect(res.Items[0].userid).toEqual('USERID2');
                expect(res.Items[0].pageid).toEqual('TEST-PAGE5');
                expect(res.Items[0].data).toEqual('Test Data 1');
                expect(res.Items[1].userid).toEqual('USERID2');
                expect(res.Items[1].pageid).toEqual('TEST-PAGE6');
                expect(res.Items[1].data).toEqual('Test Data 2');
            }
        );
    }
);