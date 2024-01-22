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

    let { ticket, ok: getOK } = await ticketService.get({ id });
    expect(getOK).toBeTruthy();
    expect(ticket).toEqual({
      body: "a new ticket",
      subject: "",
    });

    const { ok: updateOK } = await ticketService.update({
      id,
      body: "a newer ticket",
      subject: "a subject",
    });
    expect(updateOK).toBeTruthy();

    ({ ticket, ok: getOK } = await ticketService.get({ id }));
    expect(getOK).toBeTruthy();
    expect(ticket).toEqual({
      body: "a newer ticket",
      subject: "a subject",
    });
  });

  it("tickets in other spaces can't be accessed", async () => {
    const t1 = await getUserToken();
    const t2 = await getUserToken();
    const ticketService1 = new TicketService(url, t1.token);
    const ticketService2 = new TicketService(url, t2.token);

    const { id: id1 } = await ticketService1.create({
      subject: "a ticket from space 1",
      body: "",
    });
    const { id: id2 } = await ticketService2.create({
      subject: "a ticket from space 2",
      body: "",
    });
    const { ok: ok1 } = await ticketService1.get({ id: id2 });
    expect(ok1).toBeFalsy();
    const { ok: ok2 } = await ticketService2.get({ id: id1 });
    expect(ok2).toBeFalsy();

    const { ok: updateOk } = await ticketService1.update({
      id: id2,
      subject: "a ticket from space 2",
      body: "",
    });
    expect(updateOk).toBeFalsy();
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
