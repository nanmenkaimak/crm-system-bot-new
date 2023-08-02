package handlers

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/datepicker"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/go-telegram/ui/slider"
	"github.com/nanmenkaimak/crm-system-bot-new/internal/tables"
	"gorm.io/gorm"
	"strconv"
)

// Repo is repository used by handlers
var Repo *Repository

// Repository is repository type
type Repository struct {
	DB *gorm.DB
}

// NewRepo creates new repository
func NewRepo(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

// NewHandlers sets repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

var values = make(map[string]string)
var slides []slider.Slide
var adminID = 1
var aliPhoto = "AgACAgIAAxkBAAIQr2TJ63_S8PBKh86oLoZMnY9i_y9cAAKYyTEbbSVQShHLDpFLNjuLAQADAgADcwADMAQ"
var aliPhotoFu = "AgACAgIAAxkBAAIQn2TJ6SQflYbZeL1_a40LRlTXwntXAAKUyTEbbSVQSi8R7wvEBouKAQADAgADcwADMAQ"
var allBeds []tables.Bed

func (m *Repository) Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var currentUser tables.User
	m.DB.Where("username = ?", update.Message.Chat.ID).Find(&currentUser)
	fmt.Println(update.Message.Chat.ID)
	if currentUser.ID == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Сен сообщения жаза алмайсын",
		})
		return
	}
	fmt.Println(update.Message.Photo[0].FileID)
	switch update.Message.Text {
	case "/info":
		values = make(map[string]string)
		slides = []slider.Slide{}
		allBeds = []tables.Bed{}
		values["for"] = "info"
		var allRooms []tables.Room
		m.DB.Where("admin_id = ?", adminID).Find(&allRooms)
		if len(allRooms) == 0 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "There is no room that you are admin",
			})
			return
		}
		kb := inline.New(b)
		for _, room := range allRooms {
			roomID := []byte(strconv.Itoa(room.ID))
			kb.Row().Button(room.Name, roomID, m.onInlineKeyboardSelectForRooms)
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "Select the variant",
			ReplyMarkup: kb,
		})
		break
	case "/addexpense":
		values = make(map[string]string)
		slides = []slider.Slide{}
		values["for"] = "addexpense"
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Кандай расход косасын?",
		})
		break
	case "/report":
		values = make(map[string]string)
		slides = []slider.Slide{}
		values["for"] = "report"
		m.Report(ctx, b, update)
		break
	case "/reporttime":
		values = make(map[string]string)
		slides = []slider.Slide{}
		values["for"] = "reporttime"
		kb := datepicker.New(b, arrivalDayPicker)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "Кай куннен бастап",
			ReplyMarkup: kb,
		})
		kb = datepicker.New(b, departureDayPicker)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "Кай кунге дейын",
			ReplyMarkup: kb,
		})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Карау ушын любой сообщения жаз",
		})
		break
	case "/expenses":
		values = make(map[string]string)
		slides = []slider.Slide{}
		values["for"] = "expenses"
		m.Expenses(ctx, b, update)
		break
	default:
		if values["for"] == "info" {
			m.Info(ctx, b, update)
			return
		}
		if values["for"] == "addexpense" {
			m.AddExpense(ctx, b, update)
			return
		}
		if values["for"] == "reporttime" {
			m.ReportTime(ctx, b, update)
			return
		}
		if values["for"] == "change_date" {
			m.UpdateDate(ctx, b, update)
			return
		}
		break
	}
}

func (m *Repository) onInlineKeyboardSelectForRooms(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	str := ""
	for _, b := range data {
		str += string(byte(b))
	}
	values["room_id"] = str
	roomID, _ := strconv.Atoi(str)

	m.DB.Where("room_id = ?", roomID).Find(&allBeds)

	var allBedsWithCustomer = make([]tables.BedWithCustomer, len(allBeds))
	for i, bed := range allBeds {
		if bed.ID == 0 {
			continue
		}
		var customer tables.Customer
		m.DB.Where("bed_id = ? and is_here = ?", bed.ID, true).Find(&customer)
		allBedsWithCustomer[i].Bed = bed
		allBedsWithCustomer[i].Customer = customer
	}

	for i, information := range allBedsWithCustomer {
		var inforOfCustomer slider.Slide
		if information.Customer.ID > 0 {
			inforOfCustomer = slider.Slide{
				Text: fmt.Sprintf("*%d\\.* Full Name: %s\n Phone: %s\n Position: %s\n Info: %s\n Money: %d тг толеды \n Arrival Day: %s\n Departure Day: %s\n", i+1,
					information.Customer.FullName, information.Customer.Phone, information.Bed.Name,
					information.Customer.Info, information.Customer.Money, information.Customer.ArrivalDay.Format("2006/01/02"),
					information.Customer.DepartureDay.Format("2006/01/02")),
				Photo: information.Customer.Photo,
			}
		} else {
			inforOfCustomer = slider.Slide{
				Text:  fmt.Sprintf("*%d\\.* Казыр мында ешкым турып жаткан жок \n %s кровать", i+1, information.Bed.Name),
				Photo: aliPhoto,
			}
		}
		slides = append(slides, inforOfCustomer)
	}

	opts := []slider.Option{
		slider.OnSelect("Update", true, m.sliderOnUpdate),
	}

	sl := slider.New(slides, opts...)

	sl.Show(ctx, b, strconv.FormatInt(mes.Chat.ID, 10))
}

func (m *Repository) sliderOnUpdate(ctx context.Context, b *bot.Bot, message *models.Message, item int) {
	kb := inline.New(b)
	var cnt int64
	m.DB.Model(&tables.Customer{}).Where("bed_id = ? and is_here = ?", allBeds[item].ID, true).Count(&cnt)
	if cnt == 0 {
		kb.Row().Button("Add", []byte(strconv.Itoa(allBeds[item].ID)), m.onInlineKeyboardSelectAfterUpdateAdd)
	} else {
		kb.Row().Button("Delete", []byte(strconv.Itoa(allBeds[item].ID)), m.onInlineKeyboardSelectAfterUpdateDelete)
		kb.Row().Button("Change", []byte(strconv.Itoa(allBeds[item].ID)), m.onInlineKeyboardSelectAfterUpdateChange)
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      message.Chat.ID,
		Text:        "Select the variant",
		ReplyMarkup: kb,
	})
}

func (m *Repository) onInlineKeyboardSelectAfterUpdateAdd(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	values["bed_id"] = string(data)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "Аты кым конактын?",
	})
}

func (m *Repository) onInlineKeyboardSelectAfterUpdateDelete(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	bedID, _ := strconv.Atoi(string(data))
	m.DB.Model(&tables.Customer{}).Where("bed_id = ?", bedID).Update("is_here", false)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "Удалить еттым",
	})
}

func (m *Repository) onInlineKeyboardSelectAfterUpdateChange(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	bedID := string(data)
	values["for"] = "change_date"
	values["bed_id"] = bedID
	kb := datepicker.New(b, arrivalDayPicker)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        "Продлить еткен куны",
		ReplyMarkup: kb,
	})
	kb = datepicker.New(b, departureDayPicker)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        "Кететын кунды жаз",
		ReplyMarkup: kb,
	})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "Канша акша толеды",
	})
}
