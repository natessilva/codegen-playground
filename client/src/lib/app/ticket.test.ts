import { authN, getUserToken, url } from "../test-helpers";
import { TicketService } from "./client_gen";

describe("ticket serice", () => {
  let ticketService: TicketService;
  let email: string;
  let password: string;

  beforeAll(async () => {
    const t = await getUserToken();
    email = t.email;
    password = t.password;
    ticketService = new TicketService(url, t.token);
  });

  it("can CRUD tickets", async () => {
    const { id } = await ticketService.create({
      body: "a new ticket",
      subject: "",
    });
    expect(id).toBeTruthy();

    let { ticket, ok } = await ticketService.get({ id });
    expect(ok).toBeTruthy();
    expect(ticket).toEqual({
      body: "a new ticket",
      subject: "",
    });

    await ticketService.update({
      id,
      body: "a newer ticket",
      subject: "a subject",
    });

    ({ ticket, ok } = await ticketService.get({ id }));
    expect(ok).toBeTruthy();
    expect(ticket).toEqual({
      body: "a newer ticket",
      subject: "a subject",
    });
  });

  it("can't get an id that doesn't exist", async () => {
    const { ok } = await ticketService.get({ id: 0 });
    expect(ok).toBeFalsy();
  });

  describe("assignees", () => {
    it("can't assign to things that don't exist", async () => {
      const { id } = await ticketService.create({
        body: "",
        subject: "",
      });
      // userID 0 doesn't exist
      const { ok } = await ticketService.assign({
        ticketID: id,
        userID: 0,
      });
      expect(ok).toBeFalsy();

      // ticketID 0 doesn't exist
      const { ok: ok2 } = await ticketService.assign({
        ticketID: 0,
        userID: 1,
      });

      expect(ok2).toBeFalsy();
    });

    it("can assign self", async () => {
      const { id } = await ticketService.create({
        body: "",
        subject: "",
      });
      await ticketService.assignSelf({ id });
    });
  });
});
