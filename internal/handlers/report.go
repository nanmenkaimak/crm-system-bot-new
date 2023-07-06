package handlers

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/nanmenkaimak/crm-system-bot-new/internal/tables"
	"time"
)

func (m *Repository) Report(ctx context.Context, b *bot.Bot, update *models.Update) {
	monthAgo := time.Now().UTC().AddDate(0, -1, 0).Format("2006-01-02")
	var allCostumers []tables.Customer
	var allExpenses []tables.Expenses

	m.DB.Where("arrival_day >= ?", monthAgo).Find(&allCostumers)
	totalMoney := 0
	totalDay := 0
	for _, costumer := range allCostumers {
		totalMoney += costumer.Money
		if costumer.DepartureDay.After(time.Now().UTC()) {
			totalDay += int(time.Now().UTC().Sub(costumer.ArrivalDay).Hours() / 24)
		} else {
			totalDay += int(costumer.DepartureDay.Sub(costumer.ArrivalDay).Hours() / 24)
		}
	}
	m.DB.Where("created_at >= ?", monthAgo).Find(&allExpenses)
	totalExpense := 0
	for _, expense := range allExpenses {
		totalExpense += expense.Money
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: fmt.Sprintf("Быр айдын отчеты \nОткызылген кундер: %d \nЖиналган акша: %d тг \nРасходтар: %d тг \nОбщий: %d тг \n",
			totalDay, totalMoney, totalExpense, totalMoney-totalExpense),
	})
}

func (m *Repository) ReportTime(ctx context.Context, b *bot.Bot, update *models.Update) {
	arrivalDay, err := time.Parse(layout, values["arrival_day"])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	departureDay, err := time.Parse(layout, values["departure_day"])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var allCostumers []tables.Customer
	var allExpenses []tables.Expenses

	m.DB.Where("arrival_day >= ? and arrival_day <= ?", arrivalDay, departureDay).Find(&allCostumers)
	totalMoney := 0
	totalDay := 0
	for _, costumer := range allCostumers {
		totalMoney += costumer.Money
		if costumer.DepartureDay.After(time.Now().UTC()) {
			totalDay += int(time.Now().UTC().Sub(costumer.ArrivalDay).Hours() / 24)
		} else {
			totalDay += int(costumer.DepartureDay.Sub(costumer.ArrivalDay).Hours() / 24)
		}
	}
	m.DB.Where("created_at >= ? and created_at <= ?", arrivalDay, departureDay).Find(&allExpenses)
	totalExpense := 0
	for _, expense := range allExpenses {
		totalExpense += expense.Money
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: fmt.Sprintf("Тандаган периодтын отчеты \nОткызылген кундер: %d \nЖиналган акша: %d тг \nРасходтар: %d тг \nОбщий: %d тг \n",
			totalDay, totalMoney, totalExpense, totalMoney-totalExpense),
	})
}
