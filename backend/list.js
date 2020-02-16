import * as PageLib from "./libs/page-lib";

import { success, failure } from "./libs/response-lib";

export async function main(event, context) {
  var userid = event.requestContext.identity.cognitoIdentityId;

  if (userid == undefined) {
    userid = "PUBLIC";
  }

  const res = await PageLib.listPages(userid);
  if (res) {
    return success(res);
  }
  return failure({status: "Failed to list pages."});
}