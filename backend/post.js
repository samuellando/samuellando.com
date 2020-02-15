import * as PageLib from "./libs/page-lib";

import { success, failure } from "./libs/response-lib";

export async function main(event, context) {
  const data = JSON.parse(event.body);
  const userid = event.requestContext.identity.cognitoIdentityId;

  const res = await PageLib.addPage(userid, data.title, data.text);
  if (res) {
    return success({status: "Page created."});
  } else {
    return failure({status: "Failed to create page."});
  }
}