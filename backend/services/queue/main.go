package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"github.com/kcsu/store/model"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// FIXME: this doesn't work.
// We need to ensure all guest tickets correspond to
// successful "normal" tickets
func main() {
	// Initialise data
	c := config.Init()
	d, err := db.Init(c)
	if err != nil {
		log.Fatal(err)
	}
	fs := db.NewFormalStore(d)
	// Query formals
	formals, err := fs.GetWithQueue()
	if err != nil {
		log.Fatal(err)
	}
	for _, formal := range formals {
		ticketsRemaining := fs.TicketsRemaining(&formal, false)
		guestTicketsRemaining := fs.TicketsRemaining(&formal, true)
		// Sort tickets into guest/king's tickets
		tickets := []model.Ticket{}
		guestTickets := []model.Ticket{}
		for _, ticket := range formal.TicketSales {
			if ticket.IsGuest {
				guestTickets = append(guestTickets, ticket)
			} else {
				tickets = append(tickets, ticket)
			}
		}
		// Shuffle each list
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(tickets), func(i, j int) {
			tickets[i], tickets[j] = tickets[j], tickets[i]
		})
		rand.Shuffle(len(guestTickets), func(i, j int) {
			guestTickets[i], guestTickets[j] = guestTickets[j], guestTickets[i]
		})
		// Get the IDs of tickets which have been successfully bought
		ticketSuccess := min(int(ticketsRemaining), len(tickets))
		guestSuccess := min(int(guestTicketsRemaining), len(guestTickets))
		tickets = tickets[0:ticketSuccess]
		guestTickets = guestTickets[0:guestSuccess]
		successes := make([]uint, 0, ticketSuccess+guestSuccess)
		for _, t := range tickets {
			successes = append(successes, t.ID)
		}
		for _, t := range guestTickets {
			successes = append(successes, t.ID)
		}
		// Change these tickets from queue tickets to successful purchases
		err := d.Model(&model.Ticket{}).Where("id IN ?", successes).Update("is_queue", false).Error
		if err != nil {
			log.Fatal(err)
		}
	}
}
