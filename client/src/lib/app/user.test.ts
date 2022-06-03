import { authN, getUserToken } from "../test-helpers";
import { UserService } from "./client_gen";

const url = "http://localhost:8001";
describe("an authenticated user", () => {
  let user: UserService;
  let email: string;
  let password: string;

  beforeAll(async () => {
    const t = await getUserToken();
    email = t.email;
    password = t.password;
    user = new UserService(url, t.token);
  });

  it("can get and update user profile", async () => {
    const { name } = await user.get({
      bar: "baz",
    });
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
      oldPassword: password,
      newPassword: "new",
    });
    expect(ok).toEqual(true);
    const { token: badLoginToken, ok: badLoginOk } = await authN.login({
      email,
      password,
    });
    expect(badLoginOk).toEqual(false);
    try {
      await new UserService(url, badLoginToken).get({});
      fail(new Error("user service should throw here"));
    } catch {}
    const { token, ok: loginOk } = await authN.login({
      email,
      password: "new",
    });
    expect(loginOk).toEqual(true);
    expect(token).not.toEqual("");
    await new UserService(url, token).get({});
  });
});
