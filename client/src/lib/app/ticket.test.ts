import { getUserToken } from "../test-helpers";
import { TicketService } from "./client_gen";

const url = "http://localhost:8001";
describe("a workspace user", () => {
  let ticket: TicketService;

  beforeAll(async () => {
    const { token } = await getUserToken();
    ticket = new TicketService(url, token);
  });

  it("can crud ticktets", async () => {
    const { id } = await ticket.create({
      subject: "New ticket",
      description: "just a ticket",
    });
    expect(id).not.toBe(null);

    const { ok, ticket: t } = await ticket.get({ id });
    expect(ok).toBe(true);
    expect(t).toEqual({
      description: "just a ticket",
      id,
      status: "open",
      subject: "New ticket",
    });

    const { list } = await ticket.list({});
    expect(list.length).toBe(1);

    const { ok: updateOk } = await ticket.update({
      ...list[0],
      status: "archived",
    });
    expect(updateOk).toBe(true);

    const {
      ok: getOk,
      ticket: { status },
    } = await ticket.get({ id });
    expect(getOk).toBe(true);
    expect(status).toBe("archived");

    const { ok: notOk } = await ticket.get({ id: 0 });
    expect(notOk).toBe(false);
  });

  it("cannot get tickets outside of workspace", async () => {
    const { token } = await getUserToken();
    const newTicket = new TicketService(url, token);

    const { id } = await newTicket.create({
      subject: "a new ticket",
      description: "can't get me!",
    });
    expect(id).not.toBe(null);

    const { ok, ticket: t } = await newTicket.get({ id });
    expect(ok).toBe(true);
    const { ok: alsoOk } = await newTicket.update({
      ...t,
      description: "I can do this!",
    });
    expect(alsoOk).toBe(true);

    const { ok: notOk } = await ticket.get({ id });
    expect(notOk).toBe(false);

    const { ok: alsoNotOk } = await ticket.update({
      ...t,
      status: "closed",
      description: "can't do this!",
    });
    expect(alsoNotOk).toBe(false);
  });
});
