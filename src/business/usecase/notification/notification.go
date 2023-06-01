package notification

import (
	"context"
	notificationDom "go-clean/src/business/domain/notification"
	"go-clean/src/business/entity"
	"go-clean/src/lib/auth"
)

type Interface interface {
	GetList(ctx context.Context) ([]entity.Notification, error)
}

type notification struct {
	notification notificationDom.Interface
	auth         auth.Interface
}

func Init(nd notificationDom.Interface, auth auth.Interface) Interface {
	n := &notification{
		notification: nd,
		auth:         auth,
	}

	return n
}

func (n *notification) GetList(ctx context.Context) ([]entity.Notification, error) {
	var (
		notifications []entity.Notification
		err           error
	)

	user, err := n.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return notifications, err
	}

	notifications, err = n.notification.GetList(entity.NotificationParam{
		UserID: user.User.ID,
	})

	if err != nil {
		return notifications, err
	}

	return notifications, nil
}
