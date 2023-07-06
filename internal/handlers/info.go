package handlers

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/datepicker"
	"github.com/nanmenkaimak/crm-system-bot-new/internal/tables"
	"strconv"
	"time"
)

var layout = "2006-01-02"

func (m *Repository) Info(ctx context.Context, b *bot.Bot, update *models.Update) {
	if values["full_name"] == "" {
		values["full_name"] = update.Message.Text
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Фотосын жыбер",
		})
		return
	}
	if values["photo"] == "" {
		if update.Message.Text != "" {
			values["photo"] = aliPhotoFu
		} else {
			values["photo"] = update.Message.Photo[0].FileID
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Телефон номеры кандай",
		})
		return
	}
	if values["phone"] == "" {
		values["phone"] = update.Message.Text
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Енды доп информация жаз",
		})
		return
	}
	if values["info"] == "" {
		values["info"] = update.Message.Text
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Канша акша толеды",
		})
		return
	}
	if values["money"] == "" {
		values["money"] = update.Message.Text
		kb := datepicker.New(b, arrivalDayPicker)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "Келген кунды жаз",
			ReplyMarkup: kb,
		})
		kb = datepicker.New(b, departureDayPicker)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "Кететын кунды жаз",
			ReplyMarkup: kb,
		})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Добавить ету ушын любой сообщения жаз",
		})
		return
	}
	var newCostumer tables.Customer

	bedID, _ := strconv.Atoi(values["bed_id"])
	newCostumer.BedID = bedID

	newCostumer.FullName = values["full_name"]
	newCostumer.Photo = values["photo"]
	newCostumer.Phone = values["phone"]
	newCostumer.Info = values["info"]
	newCostumer.IsHere = true

	money, _ := strconv.Atoi(values["money"])
	newCostumer.Money = money

	arrivalDay, err := time.Parse(layout, values["arrival_day"])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	newCostumer.ArrivalDay = arrivalDay

	departureDay, err := time.Parse(layout, values["departure_day"])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	newCostumer.DepartureDay = departureDay
	newCostumer.CreatedAt, newCostumer.UpdatedAt = time.Now(), time.Now()

	m.DB.Create(&newCostumer)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Косылды",
	})
}

func arrivalDayPicker(ctx context.Context, b *bot.Bot, mes *models.Message, date time.Time) {
	values["arrival_day"] = date.Format(layout)
}

func departureDayPicker(ctx context.Context, b *bot.Bot, mes *models.Message, date time.Time) {
	values["departure_day"] = date.Format(layout)
}
