import * as databaseLib from "../libs/database-lib";

import AWS from "aws-sdk";
AWS.config.update({region:'ca-central-1'});


var db = new AWS.DynamoDB({ endpoint: 'http://localhost:8000'});

const tableName = "pages";

async function createTable() {
    var status = false;
    await db.createTable({
            AttributeDefinitions: [{AttributeName: 'pageid', AttributeType: 'S'}],
            TableName: "pages", 
            KeySchema: [{AttributeName: 'pageid', KeyType: 'HASH'}],
            ProvisionedThroughput: {
                'ReadCapacityUnits': 5,
                'WriteCapacityUnits': 5
            },
        },
        (err, data) => {
            if (err == null) {
                status = true;
            } else {
                status = false;
            }
        }
    ).promise();
    return status;
}

async function deleteTable() {
    var status = false;
    await db.deleteTable({
            TableName: "pages", 
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

test("connecting to db", async done => {
        await createTable();
        db.listTables(
            (err, data) => {
                console.log(err);
                console.log(data);
                done();
            }
        );
        await deleteTable();
    }
);