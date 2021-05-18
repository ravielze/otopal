package reminder

import (
	"time"

	"github.com/ravielze/oculi/common"
)

type UpdateRequest struct {
	ReminderType string `json:"reminder_type" binding:"required,lte=6"`
	Last         string `json:"last" binding:"required,eq=10"`
	Next         string `json:"next" binding:"required,eq=10"`
}

type ReminderResponse struct {
	ReminderType string `json:"reminder_type"`
	Last         string `json:"last"`
	Next         string `json:"next"`
}

func (item Reminder) Convert() ReminderResponse {
	return ReminderResponse{
		ReminderType: item.ReminderType,
		Last:         item.Last.Format("02-01-2006"), //DD-MM-YYYY
		Next:         item.Last.Format("02-01-2006"), //DD-MM-YYYY
	}
}

func (item UpdateRequest) Convert(reminderId string, ownerId uint) (Reminder, error) {
	lastTime, err := time.Parse("02-01-2006", item.Last)
	if err != nil {
		return Reminder{}, err
	}
	nextTime, err2 := time.Parse("02-01-2006", item.Next)
	if err2 != nil {
		return Reminder{}, err2
	}
	return Reminder{
		UUIDBase:     common.UUIDBase{ID: reminderId},
		ReminderType: string(ReminderType(item.ReminderType)),
		Last:         lastTime,
		Next:         nextTime,
		OwnerID:      ownerId,
	}, nil
}
