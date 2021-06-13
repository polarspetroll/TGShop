package app

import (
	"DB"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"telegram"
)

func BotHandler(w http.ResponseWriter, r *http.Request) {
	j, err := ioutil.ReadAll(r.Body)
	DB.CheckErr(err)
	var status telegram.TGMessage
	err = json.Unmarshal(j, &status)
	DB.CheckErr(err)
	DB.SetList("chatids", fmt.Sprintf("%v", status.Message.Chat.ID))
	message := status.Message.Text
	if len(status.Message.Entities) != 0 {
		if status.Message.Entities[0].Type == "bot_command" {
			message = message[1:]
		} else {
			return
		}
	} else {
		return
	}
	if message == "help" {
		telegram.SendMessage(status.Message.Chat.ID, "available commands :\n/help: help menu\n/list: list all products\n/product_code: get a specific product information\neg: /2532\n\n/contact: Contact Us")
	} else if message == "contact" {
		telegram.SendMessage(status.Message.Chat.ID, os.Getenv("CONTACTUS"))
	} else if message == "list" {
		list := DB.ListQuery()
		var msg string
		for i := 0; i < len(list.Id); i++ {
			if i == 0 {
				msg = fmt.Sprintf("/%v : %v\n", list.Id[i], list.Name[i])
			} else {
				msg += fmt.Sprintf("/%v : %v\n", list.Id[i], list.Name[i])
			}
		}
		telegram.SendMessage(status.Message.Chat.ID, msg)
	} else {
		pid, err := strconv.ParseInt(message, 10, 64)
		if err != nil {
			telegram.SendMessage(status.Message.Chat.ID, "Product not found")
			return
		}
		var product DB.QueryOutput
		cacheres := DB.GetCache(pid)
		if fmt.Sprintf("%v", cacheres[0]) == "<nil>" {
			product = DB.QueryById(pid)
		} else {
			var cachestat string
			if cacheres[3] == "1" {
				cachestat = "In stock"
			} else {
				cachestat = "Out of stock"
			}
			caption := fmt.Sprintf("Name: %v\nPrice: %v\nStatus: %v\n\nSource: Cache", cacheres[0], cacheres[1], cachestat)
			telegram.SendPhoto(cacheres[2].(string), caption, status.Message.Chat.ID)
			return
		}
		if len(product.Name) != 1 {
			telegram.SendMessage(status.Message.Chat.ID, "Product not found")
			return
		}
		product.Id = append(product.Id, pid)
		DB.SetCache(product)
		var availability string
		if product.Stat[0] == 1 {
			availability = "In stock"
		} else {
			availability = "Out of stock"
		}
		caption := fmt.Sprintf("Name: %v\nPrice: %v\nStatus: %v", product.Name[0], product.Price[0], availability)
		telegram.SendPhoto(product.Fname[0], caption, status.Message.Chat.ID)
	}
}
