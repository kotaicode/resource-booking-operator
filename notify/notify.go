// Package notify holds all different types of ways to notify users about various events.
package notify

import (
	"errors"

	managerv1 "github.com/kotaicode/resource-booking-operator/api/v1"
)

// Notifier is an interface that each type of notifier must implement.
type Notifier interface {
	Prepare(booking managerv1.Booking) Notifier
	Send() error
}

// NewNotifier returns a new Notifier based on the type of notification.
func NewNotifier(notification managerv1.Notification) (Notifier, error) {
	switch notification.Type {
	case "email":
		return &Email{Recipient: notification.Recipient}, nil
	default:
		return nil, errors.New("Notifier type not found")
	}
}
