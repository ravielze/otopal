package reminder

import (
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) IRepo {
	return Repository{db: db}
}

func (repo Repository) GetOrCreate(userId uint) ([]Reminder, error) {
	var oil Reminder
	var tuneup Reminder

	if err := repo.db.Where(
		Reminder{
			ReminderType: string(OIL),
			OwnerID:      userId,
		},
	).Attrs(
		Reminder{
			ReminderType: string(OIL),
			OwnerID:      userId,
			Last:         time.Now(),
			Next:         time.Now().AddDate(0, 0, 21),
		},
	).FirstOrCreate(&oil).Error; err != nil {
		return nil, err
	}

	if err2 := repo.db.Where(
		Reminder{
			ReminderType: string(TUNE_UP),
			OwnerID:      userId,
		},
	).Attrs(
		Reminder{
			ReminderType: string(TUNE_UP),
			OwnerID:      userId,
			Last:         time.Now(),
			Next:         time.Now().AddDate(0, 1, 14),
		},
	).FirstOrCreate(&tuneup).Error; err2 != nil {
		return nil, err2
	}
	return []Reminder{oil, tuneup}, nil
}

func (repo Repository) Update(reminder Reminder) error {
	if reminder.Next.Before(reminder.Last) {
		reminder.Next, reminder.Last = reminder.Last, reminder.Next
	}
	if err := repo.db.Model(&Reminder{}).
		Where("reminder_id = ?", reminder.ID).
		Where("owner_id = ?", reminder.OwnerID).
		Where("reminder_type = ?", reminder.ReminderType).
		Select("last", "next").Updates(&reminder).
		Error; err != nil {
		return err
	}
	return nil
}
