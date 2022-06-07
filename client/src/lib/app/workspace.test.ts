import { getUserToken } from "../test-helpers";
import { WorkspaceService } from "./client_gen";

const url = "http://localhost:8001";
describe("a workspace user", () => {
  let workspace: WorkspaceService;

  beforeAll(async () => {
    const { token } = await getUserToken();
    workspace = new WorkspaceService(url, token);
  });

  it("can crud workspaces", async () => {
    const { name } = await workspace.get({});
    expect(name).toEqual("");

    await workspace.update({ name: "custom name" });
    const { name: updatedName } = await workspace.get({});
    expect(updatedName).toEqual("custom name");

    const { list } = await workspace.list({});
    expect(list.length).toBe(1);

    const { id } = await workspace.create({ name: "new workspace" });
    expect(id).not.toBe(null);

    const { list: newList } = await workspace.list({});
    expect(newList).toMatchObject([
      {
        name: "custom name",
      },
      { name: "new workspace" },
    ]);
  });

  it("can switch workspaces", async () => {
    const { id } = await workspace.create({ name: "switch workspace" });
    expect(id).not.toBe(null);
    const { ok: notOk } = await workspace.switch({ id: 0 });
    expect(notOk).toBe(false);

    const { token, ok } = await workspace.switch({ id });
    expect(ok).toBe(true);
    expect(token).not.toEqual("");

    const newWorkspace = new WorkspaceService(url, token);

    const { name } = await newWorkspace.get({});
    expect(name).toBe("switch workspace");
  });
});
