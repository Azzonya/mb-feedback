package notifier

type Notifier interface {
	SendNotification(clientID, channelID, templateID, headerParam, buttonURLParam string) error
}
