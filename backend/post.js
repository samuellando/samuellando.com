import * as PageLib from "./libs/page-lib";
import * as AuthLib from "./libs/authorization-lib";

import { success, failure } from "./libs/response-lib";

export async function main(event, context) {
  const data = JSON.parse(event.body);
  const userid = event.requestContext.identity.cognitoIdentityId;

  var pUserid = userid;

  if (data.public) {
    pUserid = "PUBLIC";
  }

  const pageid = await PageLib.addPage(pUserid, data.title, data.text);
  if (pageid) {
    const res = await AuthLib.addAuthorization(userid, pageid, 0);
    if (res) {
      return success({status: "Page created."});
    }
  }
  return failure({status: "Failed to create page."});
}