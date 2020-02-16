import * as PageLib from "./libs/page-lib";
import * as AuthLib from "./libs/authorization-lib";

import { success, failure } from "./libs/response-lib";

export async function main(event, context) {
  const userid = event.requestContext.identity.cognitoIdentityId;
  const pageid = event.pathParameters.id;

  const auth = await AuthLib.retrieveAuthorization(userid, pageid);
  if (auth.level <= 1) {
    const res = await PageLib.removePage(userid, pageid);
    if (res) {
      return success({status: "Page deleted."});
    }
  }
  return failure({status: "Failed to delete page."});
}