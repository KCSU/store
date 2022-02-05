package queue

import (
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

func GetSuccesses(tickets []model.Ticket, maxTickets int) []int {
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

func UpdateSuccesses(tickets []int, d *gorm.DB) error {
	// Change these tickets from queue tickets to successful purchases
	return d.Model(&model.Ticket{}).Where("id IN ?", tickets).Update("is_queue", false).Error
}

func GetNonGuestQueue(formal *model.Formal, d *gorm.DB) ([]model.Ticket, error) {
	var tickets []model.Ticket
	err := d.Model(&formal).Association("TicketSales").Find(&tickets, "NOT is_guest AND is_queue")
	return tickets, err
}

func GetGuestQueueByUsers(formal *model.Formal, users []int, d *gorm.DB) ([]model.Ticket, error) {
	var guestTickets []model.Ticket
	err := d.Model(&formal).Association("TicketSales").Find(&guestTickets, "is_guest AND is_queue AND user_id IN ?", users)
	return guestTickets, err

}

func GetSuccessfulUserIDs(formal *model.Formal, d *gorm.DB) ([]int, error) {
	var successfulUsers []int
	err := d.Model(&model.Ticket{}).Not("is_guest OR is_queue").Distinct("user_id").Pluck("user_id", &successfulUsers).Error
	return successfulUsers, err
}

func Run(c *config.Config, d *gorm.DB, f db.FormalStore) error {
	// Query formals
	// TODO: only get those whose sales have started
	formals, err := f.Get()
	if err != nil {
		return err
	}

	for _, formal := range formals {
		ticketsRemaining := f.TicketsRemaining(&formal, false)
		guestTicketsRemaining := f.TicketsRemaining(&formal, true)
		// Get all queued tickets
		tickets, err := GetNonGuestQueue(&formal, d)
		if err != nil {
			return err
		}

		successes := GetSuccesses(tickets, int(ticketsRemaining))
		if err := UpdateSuccesses(successes, d); err != nil {
			return err
		}

		// Get a list of users with King's tickets
		successfulUsers, err := GetSuccessfulUserIDs(&formal, d)
		if err != nil {
			return err
		}

		// TODO: make this fair to people who only want one!
		guestTickets, err := GetGuestQueueByUsers(&formal, successfulUsers, d)
		if err != nil {
			return err
		}

		guestSuccesses := GetSuccesses(guestTickets, int(guestTicketsRemaining))
		if err := UpdateSuccesses(guestSuccesses, d); err != nil {
			return err
		}
	}
	return nil
}
