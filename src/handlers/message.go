package handlers

import (
	"fmt"
	"strings"
	"time"
	"tomoribot-geminiai-version/client"
	command_types "tomoribot-geminiai-version/src/commands/types"
	constants "tomoribot-geminiai-version/src/defaults"
	"tomoribot-geminiai-version/src/handlers/actions"
	infra_whatsmeow_utils "tomoribot-geminiai-version/src/infra/whatsapp/whatsmeow/utils"

	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
)

func IsMentionedBot(QuotedMsgContextInfo *waE2E.ContextInfo, botJid types.JID) bool {
	if QuotedMsgContextInfo == nil {
		return false
	}
	if QuotedMsgContextInfo.Participant == nil {
		return false
	}
	return *QuotedMsgContextInfo.Participant == botJid.String()
}

func MessageHandler(client *client.Client, message *events.Message) {
	processmentStartedTime := time.Now()

	if message.Info.ID == "" || message.Info.IsFromMe {
		return
	}
	userJid := message.Info.Sender.ToNonAD()
	botJid := client.Client.Store.ID.ToNonAD()
	chatJid := message.Info.Chat.ToNonAD()
	if chatJid.String() == constants.STATUS_BROADCAST {
		return
	}
	if strings.HasSuffix(chatJid.String(), constants.GROUP_COMMUNITY) {
		return
	}
	if userJid.String() == botJid.String() {
		return
	}
	if message.IsEphemeral {
		ephemeralMsg := message.Message.GetEphemeralMessage()
		if ephemeralMsg != nil && ephemeralMsg.Message != nil {
			message.Message = ephemeralMsg.Message
		}
	}

	body := infra_whatsmeow_utils.GetMessageBody(message.Message)
	quotedMsgInfo := infra_whatsmeow_utils.GetQuotedMessageContextInfo(message.Message)

	if message.Info.IsGroup && !strings.HasPrefix(strings.ToLower(body), "tomori,") && !IsMentionedBot(quotedMsgInfo, botJid) {
		return
	}
	if len(body) == 0 {
		return
	}

	messageType := infra_whatsmeow_utils.GetMessageType(message.Message)
	quotedMsg := infra_whatsmeow_utils.GetQuotedMessage(message.Message)
	args := strings.Split(body, " ")
	arg := strings.Join(args, " ")
	commandProps := &command_types.CommandProps{
		Client:               client,
		Args:                 args,
		Message:              message,
		QuotedMsgContextInfo: quotedMsgInfo,
		Arg:                  arg,
		Timestamp:            processmentStartedTime,
		MessageType:          messageType,
		QuotedMsg: 					  quotedMsg,
	}
	fmt.Println("🔍 Command received: "+arg)
	actions.ProcessorGeminiAI(commandProps)
}
