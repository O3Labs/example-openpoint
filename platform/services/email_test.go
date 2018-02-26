package services

import (
	"log"
	"testing"

	"github.com/o3labs/openpoint/platform/config"
	"github.com/o3labs/openpoint/platform/models"
)

func TestSendEmail(t *testing.T) {
	config.Load("local")
	email := models.Email{
		From:    "apisit@o3.network",
		To:      "apisit@o3.network",
		Subject: "test",
		HTML:    "hey there",
	}
	err := NewEmailService().SendEmail(email)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
	}
}
