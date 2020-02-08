import AWS from "aws-sdk";
import * as pageLib from "../libs/page-lib"; 
import * as dbLib from "../libs/database-lib"; 

AWS.config.update({region:'ca-central-1'});
var db = new AWS.DynamoDB({endpoint: 'http://localhost:8000'});

const tableName = "pages";

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

describe('pagelib', () => {
        beforeAll(async () => {
                await createTable();
                await sleep(1000);
            }
        );

        afterAll(async () => {
                await deleteTable();
            }
        );

        test("add page", async () => {
                await pageLib.addPage("USER-1", "Test Title", "Test Text");
                const res = await dbLib.listItems(tableName, {userid: "USER-1"});

                expect(res.Count).toEqual(1);
                expect(res.Items[0].userid).toEqual('USER-1');
                expect(res.Items[0].title).toEqual('Test Title');
                expect(res.Items[0].text).toEqual('Test Text');
                expect(res.Items[0].pageid).toBeDefined();
            }
        );

        test("remove page", async () => {
                await dbLib.addItem(tableName, {userid: "USER-2", pageid: "PAGE-2", text: "Page Text", title: "Page Title"});
                await pageLib.removePage("USER-2", "PAGE-2");

                const res = await dbLib.listItems(tableName, {userid: "USER-2"});

                expect(res.Count).toEqual(0);
            }
        );

        test("edit page", async () => {
                await dbLib.addItem(tableName, {userid: "USER-3", pageid: "PAGE-3", text: "Page Text", title: "Page Title"});
                await pageLib.editPage("USER-3", "PAGE-3", "New Title", "New Text");

                const res = await dbLib.listItems(tableName, {userid: "USER-3"});

                expect(res.Count).toEqual(1);
                expect(res.Items[0].userid).toEqual('USER-3');
                expect(res.Items[0].title).toEqual('New Title');
                expect(res.Items[0].text).toEqual('New Text');
                expect(res.Items[0].pageid).toEqual('PAGE-3');
            }
        );

        test("retrieve page", async () => {
                await dbLib.addItem(tableName, {userid: "USER-4", pageid: "PAGE-4", text: "Page Text", title: "Page Title"});
                const res = await pageLib.retrievePage("USER-4", "PAGE-4");

                expect(res.userid).toEqual('USER-4');
                expect(res.title).toEqual('Page Title');
                expect(res.text).toEqual('Page Text');
                expect(res.pageid).toEqual('PAGE-4');
            }
        );

        test("list public pages", async () => {
            }
        );

        test("list pages", async () => {
            }
        );
    }
);