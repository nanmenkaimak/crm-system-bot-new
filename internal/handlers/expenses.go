package handlers

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/nanmenkaimak/crm-system-bot-new/internal/tables"
	"strconv"
	"time"
)

func (m *Repository) Expenses(ctx context.Context, b *bot.Bot, update *models.Update) {
	monthAgo := time.Now().UTC().AddDate(0, -1, 0).Format(layout)
	var allExpenses []tables.Expenses

	m.DB.Where("created_at >= ?", monthAgo).Find(&allExpenses)
	expensesText := ""
	totalExpense := 0
	for i, expense := range allExpenses {
		expensesText += fmt.Sprintf("%d. Расход: %s - %d тг, дата: %s \n", i+1, expense.Name, expense.Money, expense.CreatedAt.Format("2006/01/02"))
		totalExpense += expense.Money
	}
	expensesText += fmt.Sprintf("Общий расход: %d тг", totalExpense)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   expensesText,
	})
}

func (m *Repository) AddExpense(ctx context.Context, b *bot.Bot, update *models.Update) {
	if values["name_expenses"] == "" {
		values["name_expenses"] = update.Message.Text
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Канша кетырдын?",
		})
		return
	}
	var newExpense tables.Expenses
	newExpense.Name = values["name_expenses"]
	money, _ := strconv.Atoi(update.Message.Text)
	newExpense.Money = money
	newExpense.CreatedAt, newExpense.UpdatedAt = time.Now(), time.Now()
	m.DB.Create(&newExpense)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Косылды",
	})
}
