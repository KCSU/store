package queue

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
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

func GetSuccesses(tickets []model.Ticket, maxTickets int) []uuid.UUID {
	// Shuffle the list
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(tickets), func(i, j int) {
		tickets[i], tickets[j] = tickets[j], tickets[i]
	})
	// Get the IDs of tickets which have been successfully bought
	ticketSuccess := min(maxTickets, len(tickets))
	tickets = tickets[0:ticketSuccess]
	successes := make([]uuid.UUID, 0, ticketSuccess)
	for _, t := range tickets {
		successes = append(successes, t.ID)
	}
	return successes
}

func GetGuestSuccesses(tickets []model.Ticket, maxTickets int) []uuid.UUID {
	guestAttempts := make(map[uuid.UUID][]model.Ticket, len(tickets))
	for _, guestTicket := range tickets {
		guestAttempts[guestTicket.UserID] = append(guestAttempts[guestTicket.UserID], guestTicket)
	}
	userIds := make([]uuid.UUID, 0, len(guestAttempts))
	for userId := range guestAttempts {
		userIds = append(userIds, userId)
	}
	// Shuffle the list
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(userIds), func(i, j int) {
		userIds[i], userIds[j] = userIds[j], userIds[i]
	})
	// Get the IDs of tickets which have been successfully bought
	successes := make([]uuid.UUID, 0, maxTickets)
	for _, userId := range userIds {
		attempts := guestAttempts[userId]
		ticketSuccess := min(maxTickets, len(attempts))
		attempts = attempts[0:ticketSuccess]
		for _, t := range attempts {
			successes = append(successes, t.ID)
		}
		maxTickets -= ticketSuccess
		if maxTickets == 0 {
			break
		}
	}
	return successes
}

func UpdateSuccesses(tickets []uuid.UUID, d *gorm.DB) error {
	// Change these tickets from queue tickets to successful purchases
	return d.Model(&model.Ticket{}).Where("id IN ?", tickets).Update("is_queue", false).Error
}

func GetNonGuestQueue(formal *model.Formal, d *gorm.DB) ([]model.Ticket, error) {
	var tickets []model.Ticket
	err := d.Model(&formal).Association("TicketSales").Find(&tickets, "NOT is_guest AND is_queue")
	return tickets, err
}

func GetGuestQueueByUsers(formal *model.Formal, users []uuid.UUID, d *gorm.DB) ([]model.Ticket, error) {
	var guestTickets []model.Ticket
	err := d.Model(&formal).Association("TicketSales").Find(&guestTickets, "is_guest AND is_queue AND user_id IN ?", users)
	return guestTickets, err

}

func GetSuccessfulUserIDs(formal *model.Formal, d *gorm.DB) ([]uuid.UUID, error) {
	var successfulUsers []uuid.UUID
	err := d.Model(&model.Ticket{}).Not("is_guest").Not("is_queue").Distinct("user_id").Pluck("user_id", &successfulUsers).Error
	return successfulUsers, err
}

func Run(c *config.Config, d *gorm.DB, f db.FormalStore) error {
	// Query formals
	formals, err := f.GetActive()
	if err != nil {
		return err
	}

	for _, formal := range formals {
		ticketsAvailable := f.TicketsRemaining(&formal, false)
		guestTicketsAvailable := f.TicketsRemaining(&formal, true)

		if time.Now().Before(formal.SecondSaleStart) {
			// Second sale hasn't started yet
			ticketsAvailable -= formal.SecondSaleTickets
			guestTicketsAvailable -= formal.SecondSaleGuestTickets
		}

		// Get all queued tickets
		tickets, err := GetNonGuestQueue(&formal, d)
		if err != nil {
			return err
		}

		successes := GetSuccesses(tickets, int(ticketsAvailable))
		if err := UpdateSuccesses(successes, d); err != nil {
			return err
		}

		// Get a list of users with King's tickets
		successfulUsers, err := GetSuccessfulUserIDs(&formal, d)
		if err != nil {
			return err
		}

		guestTickets, err := GetGuestQueueByUsers(&formal, successfulUsers, d)
		if err != nil {
			return err
		}

		guestSuccesses := GetGuestSuccesses(guestTickets, int(guestTicketsAvailable))
		if err := UpdateSuccesses(guestSuccesses, d); err != nil {
			return err
		}
	}
	return nil
}
