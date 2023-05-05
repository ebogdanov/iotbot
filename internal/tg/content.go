package tg

//
//import (
//	"fmt"
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	"su27bot/internal/acl"
//	"su27bot/internal/config"
//	"su27bot/internal/handler"
//	"su27bot/internal/model"
//)
//
//type Render struct {
//	Actions  map[string][]*model.Action
//	stickers []string
//	cfg      *config.Config
//	acl      *acl.Storage
//}
//
//func NewContent(cfg *config.Config, acl *acl.Storage, apiActions []*model.Action) *Render {
//	c := &Render{
//		Actions: filter(cfg, apiActions),
//		cfg:     cfg,
//		acl:     acl,
//	}
//
//	c.stickers = make([]string, 0)
//	if len(cfg.Stickers) != 0 {
//		c.stickers = cfg.Stickers
//	}
//
//	return c
//}
//
//func filter(cfg *config.Config, a []*model.Action) map[string][]*model.Action {
//	actions := make(map[string][]*model.Action, 0)
//
//	// Todo - refactor this shit
//
//	if len(cfg.Acl.Actions.Allow) > 0 {
//		actions["all"] = make([]*model.Action, 0)
//
//		// Filter out what we got from API
//		for _, item := range a {
//			found := false
//
//			for _, id := range cfg.Acl.Actions.Allow {
//				if item.Name == id {
//					found = true
//					break
//				}
//			}
//
//			if found {
//				actions["all"] = append(actions["all"], item)
//			}
//		}
//	}
//
//	for name, members := range cfg.Acl.Actions.Only {
//		actions[name] = make([]*model.Action, 0)
//
//		for _, item := range a {
//			found := false
//
//			for _, id := range members {
//				if item.Name == id {
//					found = true
//					break
//				}
//			}
//
//			if found {
//				actions[name] = append(actions[name], item)
//			}
//		}
//	}
//
//	return actions
//}
//
//// MainMenu returns menu to be shown for user
//func (c *Render) MainMenu(userId string) tgbotapi.InlineKeyboardMarkup {
//	menu := tgbotapi.NewInlineKeyboardMarkup()
//
//	for groupName, groupItems := range c.Actions {
//		if !c.acl.IsMember(groupName, userId) {
//			continue
//		}
//
//		for _, item := range groupItems {
//			id := fmt.Sprintf("DO_%d_%s", item.HomeID, item.ID)
//
//			menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
//				tgbotapi.NewInlineKeyboardButtonData(item.Name, id),
//			))
//		}
//	}
//
//	if c.acl.IsMember("admin", userId) {
//		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
//			// Submenu - admin section for admin only
//			tgbotapi.NewInlineKeyboardButtonData("Admin section", handler.MenuAdmin),
//		))
//	}
//
//	return menu
//}
//
//func (c *Render) GroupMenu(groups []string, prefix string) tgbotapi.InlineKeyboardMarkup {
//	menu := tgbotapi.NewInlineKeyboardMarkup()
//
//	for _, g := range groups {
//		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
//			// Group item
//			tgbotapi.NewInlineKeyboardButtonData(g, fmt.Sprintf("%s:%s", prefix, g))))
//	}
//
//	menu.InlineKeyboard = append(menu.InlineKeyboard, c.AddDeleteButtons(handler.GroupList)...)
//
//	return menu
//}
//
//func (c *Render) AddDeleteButtons(prefix string) [][]tgbotapi.InlineKeyboardButton {
//	menuPrefix := prefix
//	res := [][]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardRow(
//		tgbotapi.NewInlineKeyboardButtonData("Добавить", menuPrefix+":ADD"),
//	), tgbotapi.NewInlineKeyboardRow(
//		tgbotapi.NewInlineKeyboardButtonData("Удалить", menuPrefix+":DELETE"),
//	), tgbotapi.NewInlineKeyboardRow(
//		// В главное меню
//		tgbotapi.NewInlineKeyboardButtonData("Назад", handler.MenuAdmin),
//	)}
//
//	return res
//}
//
//func (c *Render) EnterUsername() {
//	// tgbotapi.NewInlineKeyboardButtonData("Удалить все", "CLEAR_INVITES"),
//}
