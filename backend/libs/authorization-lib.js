/**
 * This library provides an API to the authorization table for the handlers.
 */

import * as DbLib from "./database-lib";
import * as PageLib from "./page-lib";

var tableName = process.env.authTableName || "authorizations";

export function addAuthorization(userid, pageid, level) {
    var auth = {
        userid: userid,
        pageid: pageid,
        level: level,
    };
    return DbLib.addItem(tableName, auth);
}

export function removeAuthorization(userid, pageid) {
    var auth = {
        userid: userid, 
        pageid: pageid
    };
    return DbLib.removeItem(tableName, auth);
}

export function retrieveAuthorization(userid, pageid) {
    var auth = {
        userid: userid, 
        pageid: pageid
    };
    return DbLib.retrieveItem(tableName, auth);
}

export function listAuthorizations(userid) {
    var auth = {
        userid: userid, 
    };
    return DbLib.listItems(tableName, auth);
}