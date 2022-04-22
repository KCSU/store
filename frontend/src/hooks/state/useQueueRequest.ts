import { useReducer } from "react";
import { QueueRequest } from "../../model/QueueRequest";
import { TicketRequest } from "../../model/TicketRequest";

export type QueueRequestAction =
  | {
      type: "ticket";
      value: TicketRequest;
    }
  | {
      type: "guestTickets";
      value: TicketRequest[];
    }
  | {
      type: "option";
      value: string;
    }
  | {
      type: "guestTicket";
      index: number;
      value: string;
    }
  | {
      type: "id";
      value: number;
    }
  | {
      type: "removeGuestTicket";
      index: number;
    }
  | {
      type: "addGuestTicket";
    };

function reducer(
  state: QueueRequest,
  action: QueueRequestAction
): QueueRequest {
  switch (action.type) {
    case "ticket":
      return {
        ...state,
        ticket: action.value,
      };
    case "guestTickets":
      return {
        ...state,
        guestTickets: action.value,
      };
    case "option":
      return {
        ...state,
        ticket: {
          option: action.value,
        },
      };
    case "guestTicket":
      return {
        ...state,
        guestTickets: [
          ...state.guestTickets.slice(0, action.index),
          {
            option: action.value,
          },
          ...state.guestTickets.slice(action.index + 1),
        ],
      };
    case "id":
      return {
        ...state,
        formalId: action.value,
      };
    case "removeGuestTicket":
      return {
        ...state,
        guestTickets: [
          ...state.guestTickets.slice(0, action.index),
          ...state.guestTickets.slice(action.index + 1),
        ],
      };
    case "addGuestTicket":
      return {
        ...state,
        guestTickets: [
          ...state.guestTickets,
          {
            option: "Normal",
          },
        ],
      };
  }
}

export function useQueueRequest(id: string) {
  return useReducer(reducer, <QueueRequest>{
    formalId: id,
    ticket: {
      option: "Normal",
    },
    guestTickets: [],
  });
}
