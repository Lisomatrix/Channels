package push

import (
	"github.com/lisomatrix/channels/channels/core"
)

type EmptyPushNotificationHandler struct {}

func (handler *EmptyPushNotificationHandler) EnqueueRequest(*core.PushRequestItem) {}