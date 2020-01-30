import * as authorizationLib from "./libs/authorization-lib";
import { success, failure } from "./libs/response-lib";

export async function main(event, context) {
  if (authorizationLib.delete(event)) {
    return success({status: true});
  } else {
    return failure({status: false});
  }
}