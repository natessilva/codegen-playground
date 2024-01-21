import { getUserToken, url } from "../test-helpers";
import { SpaceService } from "./client_gen";

describe("a user", () => {
  let space: SpaceService;
  let userId: number;

  beforeAll(async () => {
    const { token, id } = await getUserToken();
    userId = id;
    space = new SpaceService(url, token);
  });

  it("can crud spaces", async () => {
    const { name } = await space.get({});
    expect(name).toEqual("");

    await space.update({ name: "custom name" });
    const { name: updatedName } = await space.get({});
    expect(updatedName).toEqual("custom name");
  });

  it("can get users in space", async () => {
    const { users } = await space.getUsers({});
    expect(users.length).toEqual(1);
    expect(users[0].id).toEqual(userId);
  });
});
