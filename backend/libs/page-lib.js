/**
 * This library provides an API to the page table for the handlers.
 */

import * as DbLib from "./database-lib";
import * as AuthLib from "./authorization-lib";
import uuid from "uuid";

var tableName = process.env.pageTableName || "pages";

export function addPage(userid, title, text, isPrivate)  {
    var page = {
        userid: userid,
        title: title,
        text: text,
        private: isPrivate,
    }
    page.pageid = uuid.v1();
    return DbLib.addItem(tableName, page);
}

export function removePage(userid, pageid) {
    return DbLib.removeItem(tableName, {userid: userid, pageid: pageid});
}

export function editPage(pageid, page) {
    return DbLib.addItem(process.env.pageTableName, {"pageid": pageid}, page);
}

export function retrievePage(pageid) {
    return DbLib.addItem(process.env.pageTableName, {"pageid": pageid});
}

export function listPublicPages() {
    return listPages("PUBLIC");
}

export function listPages(userid) {
    const auths = AuthLib.listAuthorizations(userid);
    var items = [];
    for (var auth in auths) {
        items.push(retrievePage(auth.pageid));
    }
    items.push(...listPublicPages());
    // REMOVE DUPLICATES.
    return items;
}