/**
 * This library provides an API to the page table for the handlers.
 */

import * as DbLib from "./database-lib";
//import * as AuthLib from "./authorization-lib";
import uuid from "uuid";

var tableName = process.env.pageTableName || "pages";

export function addPage(userid, title, text)  {
    var page = {
        userid: userid,
        title: title,
        text: text,
    };
    page.pageid = uuid.v1();
    return DbLib.addItem(tableName, page);
}

export function removePage(userid, pageid) {
    return DbLib.removeItem(tableName, {userid: userid, pageid: pageid});
}

export function editPage(userid, pageid, title, text) {
    return DbLib.editItem(tableName, {userid: userid, pageid: pageid}, {title: title, text: text});
}

export function retrievePage(userid, pageid) {
    return DbLib.retrieveItem(tableName, {userid: userid,pageid: pageid});
}

export function listPublicPages() {
    return listPages("PUBLIC");
}

export function listPages(userid) {
    // TODO get authorized pages too.
    return DbLib.listItems(tableName, {userid: userid});
}