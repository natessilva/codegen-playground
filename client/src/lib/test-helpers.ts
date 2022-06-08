import fetch from "node-fetch";
import { v4 } from "uuid";
import { AuthNService } from "./authn/client_gen";

// provide node-fetch globally to our jest tests
// this is because the client code expects to
// run in browser and run the built-in fetch
// method.
globalThis.fetch = fetch as any;

export const url = "http://localhost:8001";
export const authN = new AuthNService(url, "");

// getUserToken provides test isolation by signing
// up a random email address on each call. Each
// service test uses this during a beforeAll call
// to set up their service for testing against
// a unique user in the system.
export async function getUserToken() {
  const email = `${v4()}@test.com`;
  const password = "test";
  const { token, ok } = await authN.signup({
    email,
    password,
  });
  expect(ok).toEqual(true);
  expect(token).not.toEqual("");
  return { token, email, password };
}
