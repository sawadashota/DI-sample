package notification

import (
	"errors"
	"time"
)

type Notification struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	IsDraft   bool      `json:"is_draft"`
	PublishAt time.Time `json:"publish_at"`
}

func NewNotification(id, title, body string, isDraft bool, publishAt time.Time) (*Notification, error) {
	n := &Notification{
		ID:        id,
		Title:     title,
		Body:      body,
		IsDraft:   isDraft,
		PublishAt: publishAt,
	}

	if err := n.validate(); err != nil {
		return nil, err
	}

	return n, nil
}

func (n *Notification) validate() error {
	if len(n.ID) == 0 {
		return errors.New("Notification.ID is required")
	}
	if len(n.Title) == 0 {
		return errors.New("Notification.Title is required")
	}
	if len(n.Body) == 0 {
		return errors.New("Notification.Body is required")
	}
	if n.PublishAt.Equal(time.Time{}) {
		return errors.New("Notification.PublishAt is required")
	}
	return nil
}

type ListNotificationResponse struct {
	Notifications []Notification `json:"notifications"`
}

type AddNotificationRequest struct {
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	IsDraft   bool      `json:"is_draft"`
	PublishAt time.Time `json:"publish_at"`
}

type UpdateNotificationRequest struct {
	Notification
}
