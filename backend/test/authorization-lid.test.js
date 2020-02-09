import AWS from "aws-sdk";
import * as authLib from "../libs/authorization-lib"; 
import * as dbLib from "../libs/database-lib"; 

AWS.config.update({region:'ca-central-1'});
var db = new AWS.DynamoDB({endpoint: 'http://localhost:8000'});

const tableName = "authorizations";

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

describe('authorization lib', () => {
        beforeAll(async () => {
                await createTable();
                await sleep(1000);
            }
        );

        afterAll(async () => {
                await deleteTable();
            }
        );

        test("add authorization", async () => {
                await authLib.addAuthorization("USER-1", "PAGE-1", 1);
                const res = await dbLib.listItems(tableName, {userid: "USER-1"});

                expect(res.Count).toEqual(1);
                expect(res.Items[0].userid).toEqual('USER-1');
                expect(res.Items[0].pageid).toEqual('PAGE-1');
                expect(res.Items[0].level).toEqual(1);
            }
        );

        test("remove authorization", async () => {
                await dbLib.addItem(tableName, {userid: "USER-2", pageid: "PAGE-2", level: 0});
                await authLib.removeAuthorization("USER-2", "PAGE-2");

                const res = await dbLib.listItems(tableName, {userid: "USER-2"});

                expect(res.Count).toEqual(0);
            }
        );

        test("retrieve authorization", async () => {
                await dbLib.addItem(tableName, {userid: "USER-4", pageid: "PAGE-4", level: 0});
                const res = await authLib.retrieveAuthorization("USER-4", "PAGE-4");

                expect(res.userid).toEqual('USER-4');
                expect(res.level).toEqual(0);
                expect(res.pageid).toEqual('PAGE-4');
            }
        );

        test("list authorizations", async () => {
                await dbLib.addItem(tableName, {userid: "USER-7", pageid: "PAGE-8", level: 0});
                await dbLib.addItem(tableName, {userid: "USER-7", pageid: "PAGE-9", level: 1});

                const res = await authLib.listAuthorizations("USER-7");

                expect(res.Count).toEqual(2);
                expect(res.Items[0].userid).toEqual('USER-7');
                expect(res.Items[0].pageid).toEqual('PAGE-8');
                expect(res.Items[0].level).toEqual(0);
                expect(res.Items[1].userid).toEqual('USER-7');
                expect(res.Items[1].pageid).toEqual('PAGE-9');
                expect(res.Items[1].level).toEqual(1);
            }
        );
    }
);