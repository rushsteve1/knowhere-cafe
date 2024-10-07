package mail

import (
	"container/list"
)

type MailQueue struct {
	queue *list.List
}

func NewMailQueue() MailQueue {
	return MailQueue{
		queue: list.New(),
	}
}
