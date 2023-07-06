package handlers

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/ui/slider"
	"github.com/nanmenkaimak/crm-system-bot-new/internal/tables"
	"sync"
	"time"
)

func (m *Repository) TomorrowLeave(ctx context.Context, b *bot.Bot, mutex *sync.Mutex) {
	mutex.Lock()
	layout := "2006-01-02 15:04:05"
	for {
		today := time.Now().UTC().Format("2006-01-02")
		location, err := time.LoadLocation("Asia/Almaty")
		if err != nil {
			fmt.Println("Error:", err)
			break
		}

		t := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 12, 0, 0, 0, location)

		time.Sleep(1 * time.Second)
		if time.Now().UTC().Format(layout) == t.UTC().Format(layout) {

			var allCostumers []tables.Customer

			m.DB.Where("departure_day = ? and is_here = ?", today, true).Find(&allCostumers)

			if len(allCostumers) == 0 {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: 1168238274,
					Text:   "Ертен ешкым кетпейды \nБары Норм",
				})
				continue
			}
			for i, information := range allCostumers {
				var inforOfCustomer slider.Slide
				if information.ID > 0 {
					inforOfCustomer = slider.Slide{
						Text: fmt.Sprintf("Мыналар бугын кету керек\n*%d\\.* Full Name: %s\n Phone: %s\n Info: %s\n Money: %d тг толеды \n Arrival Day: %s\n Departure Day: %s\n", i+1,
							information.FullName, information.Phone,
							information.Info, information.Money, information.ArrivalDay.Format("2006/01/02"),
							information.DepartureDay.Format("2006/01/02")),
						Photo: information.Photo,
					}
				}
				slides = append(slides, inforOfCustomer)
			}

			opts := []slider.Option{
				slider.OnSelect("Update", true, m.sliderOnUpdate),
			}

			sl := slider.New(slides, opts...)

			var allUsers []tables.User
			m.DB.Find(&allUsers)
			for _, user := range allUsers {
				sl.Show(ctx, b, user.Username)
			}
		}
	}
	mutex.Unlock()
}
