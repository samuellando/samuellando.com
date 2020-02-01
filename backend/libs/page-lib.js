/**
 * This library provides an API to the page table for the handlers.
 */

import * as DbLib from "./database-lib";
import * as AuthLib from "./authorization-lib";
import uuid from "uuid";

export function addPage(page) {
    page.pageid = uuid.v1();
    return DbLib.addItem(process.env.pageTableName, page);
}

export function removePage(pageid) {
    return DbLib.removeItem(process.env.pageTableName, {"pageid": pageid});
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