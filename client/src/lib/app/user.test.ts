import { UserService } from "./client_gen";
import { AuthNService } from "../authn/client_gen";
const url = "http://localhost:8001";
describe("an authenticated user", () => {
  const authN = new AuthNService(url, "");
  let user = new UserService(url, "");

  beforeAll(async () => {
    const { token, ok } = await authN.signup({
      email: "user@test.com",
      password: "test",
    });
    expect(ok).toEqual(true);
    expect(token).not.toEqual("");
    user = new UserService(url, token);
  });

  it("can get and update user profile", async () => {
    const { name } = await user.get({});
    expect(name).toEqual("");
    await user.update({ name: "custom name" });
    const { name: updatedName } = await user.get({});
    expect(updatedName).toEqual("custom name");
  });

  it("cannot change password without correct old passwrod", async () => {
    const { ok } = await user.setPassword({
      oldPassword: "wrong",
      newPassword: "new",
    });
    expect(ok).toEqual(false);
  });

  it("can change password with correct old passwrod", async () => {
    const { ok } = await user.setPassword({
      oldPassword: "test",
      newPassword: "new",
    });
    expect(ok).toEqual(true);
    const { token: badLoginToken, ok: badLoginOk } = await authN.login({
      email: "user@test.com",
      password: "test",
    });
    expect(badLoginOk).toEqual(false);
    try {
      await new UserService(url, badLoginToken).get({});
      fail(new Error("user service should throw here"));
    } catch {}
    const { token, ok: loginOk } = await authN.login({
      email: "user@test.com",
      password: "new",
    });
    expect(loginOk).toEqual(true);
    expect(token).not.toEqual("");
    await new UserService(url, token).get({});
  });
});
