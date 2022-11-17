package queue

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"github.com/kcsu/store/model"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gorm.io/gorm"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func GetSuccesses(tickets []model.Ticket, maxTickets int) []model.Ticket {
	// Shuffle the list
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(tickets), func(i, j int) {
		tickets[i], tickets[j] = tickets[j], tickets[i]
	})
	// Get the IDs of tickets which have been successfully bought
	ticketSuccess := min(maxTickets, len(tickets))
	return tickets[0:ticketSuccess]
}

func GetGuestSuccesses(tickets []model.Ticket, maxTickets int) []model.Ticket {
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
	successes := make([]model.Ticket, 0, maxTickets)
	for _, userId := range userIds {
		attempts := guestAttempts[userId]
		ticketSuccess := min(maxTickets, len(attempts))
		attempts = attempts[0:ticketSuccess]
		successes = append(successes, attempts...)
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
	err := d.Model(&model.Ticket{}).
		Preload("User").
		Where("formal_id = ? AND is_queue AND NOT is_guest", formal.ID).
		Find(&tickets).Error
	return tickets, err
}

func GetGuestQueueByUsers(formal *model.Formal, users []uuid.UUID, d *gorm.DB) ([]model.Ticket, error) {
	var guestTickets []model.Ticket
	err := d.Model(&model.Ticket{}).
		Preload("User").
		Where("formal_id = ? AND is_queue AND is_guest AND user_id IN ?", formal.ID, users).
		Find(&guestTickets).Error
	return guestTickets, err

}

func GetSuccessfulUserIDs(formal *model.Formal, d *gorm.DB) ([]uuid.UUID, error) {
	var successfulUsers []uuid.UUID
	err := d.Model(&model.Ticket{}).
		Where("formal_id = ?", formal.ID).
		Not("is_guest").Not("is_queue").
		Distinct("user_id").
		Pluck("user_id", &successfulUsers).Error
	return successfulUsers, err
}

type templateFormalData struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
}

type templateTicketData struct {
	Type   string `json:"type"`
	Option string `json:"option"`
}

type templateData struct {
	Username string
	Formal   templateFormalData
	Tickets  []templateTicketData
}

func GeneratePersonalisations(formal *model.Formal, successes []model.Ticket, guestSuccesses []model.Ticket) []*mail.Personalization {
	data := map[string]*templateData{}
	for _, t := range successes {
		data[t.User.Email] = &templateData{
			Tickets: []templateTicketData{
				{
					Type:   "King's",
					Option: t.MealOption,
				},
			},
			Formal: templateFormalData{
				Name: formal.Name,
				ID:   formal.ID,
			},
			Username: t.User.Name,
		}
	}
	for _, t := range guestSuccesses {
		if _, ok := data[t.User.Email]; !ok {
			data[t.User.Email] = &templateData{
				Tickets: []templateTicketData{},
				Formal: templateFormalData{
					Name: formal.Name,
					ID:   formal.ID,
				},
				Username: t.User.Name,
			}
		}
		data[t.User.Email].Tickets = append(data[t.User.Email].Tickets, templateTicketData{
			Type:   "Guest",
			Option: t.MealOption,
		})
	}

	personalizations := make([]*mail.Personalization, 0, len(data))
	for user, d := range data {
		person := mail.NewPersonalization()
		person.AddTos(mail.NewEmail(d.Username, user))
		person.SetDynamicTemplateData("formal", d.Formal)
		person.SetDynamicTemplateData("tickets", d.Tickets)
		personalizations = append(personalizations, person)
	}
	return personalizations
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
		successIDs := make([]uuid.UUID, len(successes))
		for i, t := range successes {
			successIDs[i] = t.ID
		}
		if err := UpdateSuccesses(successIDs, d); err != nil {
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
		guestSuccessIDs := make([]uuid.UUID, len(guestSuccesses))
		for i, t := range guestSuccesses {
			guestSuccessIDs[i] = t.ID
		}
		if err := UpdateSuccesses(guestSuccessIDs, d); err != nil {
			return err
		}

		if len(successes) != 0 || len(guestSuccesses) != 0 {
			personalizations := GeneratePersonalisations(
				&formal, successes, guestSuccesses,
			)
			message := mail.NewV3Mail()
			message.SetFrom(mail.NewEmail("KiFoMaSy", c.MailFrom))
			message.SetTemplateID(c.MailTemplateId)
			message.AddPersonalizations(personalizations...)
			client := sendgrid.NewSendClient(c.MailApiKey)
			response, err := client.Send(message)
			if err != nil {
				return err
			} else {
				fmt.Println(response.StatusCode)
				fmt.Println(response.Body)
				fmt.Println(response.Headers)
			}
		}
	}
	return nil
}
