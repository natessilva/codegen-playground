import { getUserToken, url } from "../test-helpers";
import { TicketService } from "./client_gen";

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

    // id: 0 will never exist
    const { ok: notOk } = await ticket.get({ id: 0 });
    expect(notOk).toBe(false);
  });

  it("cannot crud tickets outside of workspace", async () => {
    // newToken is a token for a new tenant in the system.
    // The original token and this token should have entirely
    // isolated data sets. The original ticket service cannot
    // access anything created by the new ticket service
    // created by this new token and vice versa.
    const { token: newToken } = await getUserToken();
    const newTicket = new TicketService(url, newToken);

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
