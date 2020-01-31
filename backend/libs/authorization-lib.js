/**
 * This library provides an API to the authorization table for the handlers.
 */

import * as DbLib from "./database-lib";

export function addAuthorization(authorization) {
    return DbLib.addItem(process.env.authorizationTableName, authorization);
}

export function removeAuthorization(userid, pageid) {
    return DbLib.removeItem(process.env.authorizationTableName, 
        {"userid": userid, "pageid": pageid});
}

export function retrieveAuthorization(userid, pageid) {
    return DbLib.addItem(process.env.authorizationTableName,
        {"userid": userid, "pageid": pageid});
}