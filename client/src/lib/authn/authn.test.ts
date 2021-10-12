import { AuthNService } from "./client_gen";

describe("an unauthenticated user", () => {
  const authN = new AuthNService("http://localhost:8001", "");

  beforeAll(async () => {
    const { token, ok } = await authN.signup({
      email: "authn@test.com",
      password: "test",
    });
    expect(ok).toEqual(true);
    expect(token).not.toEqual("");
  });

  it("can signup with an email that doesn't exist in the system", async () => {
    const { token, ok } = await authN.signup({
      email: "authn2@test.com",
      password: "test",
    });
    expect(ok).toEqual(true);
    expect(token).not.toEqual("");
  });

  it("cannot signup with an email that does exist in the system", async () => {
    const { token, ok } = await authN.signup({
      email: "authn@test.com",
      password: "test",
    });
    expect(ok).toEqual(false);
  });

  it("cannot login with an email that doesn't exist in the system", async () => {
    const { token, ok } = await authN.login({
      email: "authn3@test.com",
      password: "test",
    });
    expect(ok).toEqual(false);
  });

  it("cannot login with an incorrect password", async () => {
    const { token, ok } = await authN.login({
      email: "authn@test.com",
      password: "wrong",
    });
    expect(ok).toEqual(false);
  });

  it("can login with a valid email and passowrd", async () => {
    const { token, ok } = await authN.login({
      email: "authn@test.com",
      password: "test",
    });
    expect(ok).toEqual(true);
    expect(token).not.toEqual("");
  });
});
