import { Formal } from "../model/Formal";

function createFormal(data: Partial<Formal>): Formal {
  const template: Formal = {
    id: 0,
    title: "",
    menu: "This is the menu!",
    price: 18.2,
    guestPrice: 23.5,
    options: ["Normal", "Vegan", "Vegetarian", "Pescetarian"],
    guestLimit: 0,
    guestTickets: 0,
    guestTicketsRemaining: 0,
    tickets: 0,
    ticketsRemaining: 0,
    saleStart: new Date("2020/01/01"),
    saleEnd: new Date(),
    dateTime: new Date("2022/05/01 19:30")
  };
  return Object.assign(template, data);
}

export function useFormals(): Formal[] {
  const formals: Formal[] = (<Partial<Formal>[]>[
    {
      id: 1,
      title: "Example Formal",
      guestLimit: 2,
      tickets: 100,
      ticketsRemaining: 50,
      saleStart: new Date("2021/12/25"),
      saleEnd: new Date("2022/01/01"),
    },
    {
      id: 2,
      title: "Example Superformal",
      guestLimit: 0,
      saleEnd: new Date("2022/01/01"),
    },
    {
      id: 3,
      guestLimit: 1,
      title: "One More Formal",
    },
  ]).map(createFormal);
  return formals;
}
