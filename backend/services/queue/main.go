package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"github.com/kcsu/store/model"
	"gorm.io/gorm"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getSuccesses(tickets []model.Ticket, maxTickets int) []int {
	// Shuffle the list
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(tickets), func(i, j int) {
		tickets[i], tickets[j] = tickets[j], tickets[i]
	})
	// Get the IDs of tickets which have been successfully bought
	ticketSuccess := min(maxTickets, len(tickets))
	tickets = tickets[0:ticketSuccess]
	successes := make([]int, 0, ticketSuccess)
	for _, t := range tickets {
		successes = append(successes, int(t.ID))
	}
	return successes
}

func updateSuccesses(tickets []int, d *gorm.DB) error {
	// Change these tickets from queue tickets to successful purchases
	return d.Model(&model.Ticket{}).Where("id IN ?", tickets).Update("is_queue", false).Error
}

func getNonGuestQueue(formal *model.Formal, d *gorm.DB) ([]model.Ticket, error) {
	var tickets []model.Ticket
	err := d.Model(&formal).Association("TicketSales").Find(&tickets, "NOT is_guest AND is_queue")
	return tickets, err
}

func getGuestQueueByUsers(formal *model.Formal, users []int, d *gorm.DB) ([]model.Ticket, error) {
	var guestTickets []model.Ticket
	err := d.Model(&formal).Association("TicketSales").Find(&guestTickets, "is_guest AND is_queue AND user_id IN ?", users)
	return guestTickets, err

}

func getSuccessfulUserIDs(formal *model.Formal, d *gorm.DB) ([]int, error) {
	var successfulUsers []int
	err := d.Model(&model.Ticket{}).Not("is_guest OR is_queue").Distinct("user_id").Pluck("user_id", &successfulUsers).Error
	return successfulUsers, err
}

func main() {
	// Initialise data
	c := config.Init()
	d, err := db.Init(c)
	if err != nil {
		log.Panic(err)
	}
	f := db.NewFormalStore(d)
	// Query formals
	formals, err := f.Get()

	if err != nil {
		log.Panic(err)
	}
	for _, formal := range formals {
		ticketsRemaining := f.TicketsRemaining(&formal, false)
		guestTicketsRemaining := f.TicketsRemaining(&formal, true)
		// Get all queued tickets
		tickets, err := getNonGuestQueue(&formal, d)
		if err != nil {
			log.Panic(err)
		}

		successes := getSuccesses(tickets, int(ticketsRemaining))
		if err := updateSuccesses(successes, d); err != nil {
			log.Panic(err)
		}

		// Get a list of users with King's tickets
		successfulUsers, err := getSuccessfulUserIDs(&formal, d)
		if err != nil {
			log.Panic(err)
		}

		guestTickets, err := getGuestQueueByUsers(&formal, successfulUsers, d)
		if err != nil {
			log.Panic(err)
		}

		guestSuccesses := getSuccesses(guestTickets, int(guestTicketsRemaining))
		if err := updateSuccesses(guestSuccesses, d); err != nil {
			log.Panic(err)
		}
	}
}
