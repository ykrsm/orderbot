package order

import (
	"fmt"

	"github.com/nlopes/slack"
)

const ()

func MakeApprovalAttachment(
	dialog slack.DialogCallback,
	triggerID string) slack.PostMessageParameters {

	var (
		itemName   = dialog.Submission["item_name"]
		itemURL    = dialog.Submission["item_url"]
		itemReason = dialog.Submission["item_reason"]
		itemCount  = dialog.Submission["item_count"]
		user       = dialog.User
	)

	attachment := slack.Attachment{
		Text:       fmt.Sprintf("@%s submitted order request", user),
		Color:      "36a64f",
		CallbackID: "order_approval",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Item name",
				Value: itemName,
				Short: false,
			},
			slack.AttachmentField{
				Title: "Reason",
				Value: itemReason,
				Short: false,
			},
			slack.AttachmentField{
				Title: "URL",
				Value: itemURL,
				Short: false,
			},
			slack.AttachmentField{
				Title: "How many",
				Value: itemCount,
				Short: false,
			},
		},
		Actions: []slack.AttachmentAction{
			slack.AttachmentAction{
				Name:  orderApprovalApproved,
				Text:  "Approve",
				Type:  "button",
				Style: "primary",
			},
			slack.AttachmentAction{
				Name:  orderApprovalRejected,
				Text:  "Reject",
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

	return params
}
