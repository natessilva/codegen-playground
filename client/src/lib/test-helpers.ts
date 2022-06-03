import fetch from "node-fetch";
import { v4 } from "uuid";
import { AuthNService } from "./authn/client_gen";
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
