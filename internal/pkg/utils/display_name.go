package utils

import (
	"log"
)

func GetDisplayName(nickname, username string) string {
	if nickname != "" {
		log.Println("Ada")
		return nickname
	}
	log.Println("Kosong")
	return username
}
