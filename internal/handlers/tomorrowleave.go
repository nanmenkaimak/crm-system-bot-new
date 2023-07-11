package handlers

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/nanmenkaimak/crm-system-bot-new/internal/tables"
	"sync"
	"time"
)

func (m *Repository) TomorrowLeave(ctx context.Context, b *bot.Bot, mutex *sync.Mutex) {
	mutex.Lock()
	layout := "2006-01-02 15:04:05"
	for {
		today := time.Now().UTC().Format("2006-01-02") + " 00:00:00"
		todayTime, err := time.Parse("2006-01-02 15:04:05", today)
		if err != nil {
			fmt.Println("Error parsing first timestamp:", err)
			return
		}
		location, err := time.LoadLocation("Asia/Almaty")
		if err != nil {
			fmt.Println("Error:", err)
			break
		}

		t := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 12, 0, 0, 0, location)

		time.Sleep(1 * time.Second)
		if time.Now().UTC().Format(layout) == t.UTC().Format(layout) {

			var allCostumers []tables.Customer

			m.DB.Where("departure_day = ? and is_here = ?", todayTime, true).Find(&allCostumers)

			if len(allCostumers) == 0 {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: 1168238274,
					Text:   "Бугын ешкым кетпейды \nБары Норм",
				})
				continue
			}
			textLeave := fmt.Sprintf("Мыналар бугын кетеды \n")

			for i, information := range allCostumers {
				var rooms tables.Room
				m.DB.Table("rooms").
					Joins("JOIN beds ON rooms.id = beds.room_id").
					Joins("JOIN customers ON beds.id = customers.bed_id").
					Where("customers.bed_id = ?", information.BedID).Find(&rooms)
				if information.ID > 0 {
					textLeave += fmt.Sprintf("%d. Аты: %s\nТурган жеры: %s\nИнфо: %s\nКелген куны: %s\nКететын куны: %s\n",
						i+1, information.FullName, rooms.Name, information.Info, information.ArrivalDay.Format("2006-01-02"),
						information.DepartureDay.Format("2006-01-02"))
				}
			}

			var allUsers []tables.User
			m.DB.Find(&allUsers)
			for _, user := range allUsers {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: user.Username,
					Text:   textLeave,
				})
			}
		}
	}
	mutex.Unlock()
}
