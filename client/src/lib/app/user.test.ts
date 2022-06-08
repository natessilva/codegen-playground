import { authN, getUserToken, url } from "../test-helpers";
import { UserService } from "./client_gen";

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
    const { name } = await user.get({});
    expect(name).toEqual("");
    await user.update({ name: "custom name" });
    const { name: updatedName } = await user.get({});
    expect(updatedName).toEqual("custom name");
  });

  it("can change password", async () => {
    const { ok: notOk } = await user.setPassword({
      oldPassword: "wrong",
      newPassword: "new",
    });
    expect(notOk).toEqual(false);

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
    expect(new UserService(url, badLoginToken).get({})).rejects.not.toBe(null);

    const { token, ok: loginOk } = await authN.login({
      email,
      password: "new",
    });
    expect(loginOk).toEqual(true);
    expect(token).not.toEqual("");
    await new UserService(url, token).get({});
  });
});
