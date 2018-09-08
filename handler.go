package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/nlopes/slack"
)

// interactionHandler handles interactive message response.
type interactionHandler struct {
	slackClient       *slack.Client
	verificationToken string
}

func (h interactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("[ERROR] Invalid method: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := url.QueryUnescape(string(buf)[8:])
	if err != nil {
		log.Printf("[ERROR] Failed to unespace request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var message slack.AttachmentActionCallback
	if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
		log.Printf("[ERROR] Failed to decode json message from slack: %s", jsonStr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Only accept message from slack with valid token
	if message.Token != h.verificationToken {
		log.Printf("[ERROR] Invalid token: %s", message.Token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	action := message.Actions[0]
	switch action.Name {
	case actionSelect:
		value := action.SelectedOptions[0].Value

		// Overwrite original drop down message.
		attachment := slack.Attachment{
			Text:       fmt.Sprintf("OK to order %s ?", strings.Title(value)),
			Color:      "#f9a41b",
			CallbackID: "beer",
			Actions: []slack.AttachmentAction{
				{
					Name:  actionStart,
					Text:  "Yes",
					Type:  "button",
					Value: "start",
					Style: "primary",
				},
				{
					Name:  actionDialog,
					Text:  "Open Dialog",
					Type:  "button",
					Style: "warning",
				},
				{
					Name:  actionCancel,
					Text:  "No",
					Type:  "button",
					Style: "danger",
				},
			},
		}

		params := slack.PostMessageParameters{
			Attachments: []slack.Attachment{
				attachment,
			},
		}

		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&params)
		return
	case actionStart:
		title := ":ok: Donezo"
		responseMessage(w, message.OriginalMessage, title, "")
		return
	case actionCancel:
		title := fmt.Sprintf(":x: @%s canceled the request", message.User.Name)
		log.Printf("trigger_id: %s", message.TriggerID)
		responseMessage(w, message.OriginalMessage, title, "")
		return
	case actionDialog:
		title := fmt.Sprintf(":x: @%s dialog is opening", message.User.Name)
		h.responseDialog(w, message.OriginalMessage, title, "", message.TriggerID)
		return
	default:
		log.Printf("[ERROR] ]Invalid action was submitted: %s", action.Name)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h interactionHandler) responseDialog(w http.ResponseWriter, original slack.Message, title, value string, triggerID string) {

	log.Printf("trigger_id: %s", triggerID)
	dialog := slack.Dialog{
		CallbackId:     "dialog_callback_id",
		Title:          "dialog test",
		NotifyOnCancel: true,
		Elements: []slack.DialogElement{
			slack.DialogTextElement{
				Label: "hello",
				Name:  "test name",
				Type:  "text",
			},
		},
	}

	if err := h.slackClient.OpenDialog(triggerID, dialog); err != nil {
		fmt.Printf("test %+v", err)
		log.Printf("[ERROR]: %s", err)
	}
}

// responseMessage response to the original slackbutton enabled message.
// It removes button and replace it with message which indicate how bot will work
func responseMessage(w http.ResponseWriter, original slack.Message, title, value string) {

	attachment := slack.Attachment{
		Color:      "#f9a41b",
		CallbackID: "beer",
		Actions:    []slack.AttachmentAction{},
		Fields: []slack.AttachmentField{
			{
				Title: title,
				Value: value,
				Short: false,
			},
		},
	}

	params := slack.PostMessageParameters{
		Attachments: []slack.Attachment{
			attachment,
		},
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&params)
}
