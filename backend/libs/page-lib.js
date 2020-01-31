/**
 * This library provides an API to the page table for the handlers.
 */

import * as DbLib from "./database-lib";
import uuid from "uuid";

export function addPage(page) {
    page.pageid = uuid.v1();
    return DbLib.addItem(process.env.pageTableName, page)
}

export function removePage(pageid) {
    return DbLib.removeItem(process.env.pageTableName, pageid)
}

export function editPage(pageid, page) {
    return DbLib.addItem(process.env.pageTableName, pageid, page)
}

export function retrievePage(pageid) {
    return DbLib.addItem(process.env.pageTableName, pageid)
}