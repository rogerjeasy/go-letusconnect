package models

import "time"

type Group struct {
	ID            string        `json:"id" gorm:"primaryKey"`
	Name          string        `json:"name" gorm:"uniqueIndex"`
	Description   string        `json:"description"`
	ImageURL      string        `json:"imageUrl"`
	Category      GroupCategory `json:"category" gorm:"embedded"`
	ActivityLevel string        `json:"activityLevel"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
	Members       []Member      `json:"members" gorm:"many2many:group_members;"`
	Topics        []Topic       `json:"topics" gorm:"many2many:group_topics;"`
	Events        []Event       `json:"events" gorm:"foreignKey:GroupID"`
	Resources     []Resource    `json:"resources" gorm:"foreignKey:GroupID"`
	Admins        []*User       `json:"admins" gorm:"many2many:group_admins;"`
	Privacy       string        `json:"privacy"`
	Rules         []Rule        `json:"rules" gorm:"foreignKey:GroupID"`
	Featured      bool          `json:"featured"`
	Size          string        `json:"size"`
}

type Member struct {
	UserID   string    `json:"userId" gorm:"primaryKey"`
	GroupID  string    `json:"groupId" gorm:"primaryKey"`
	JoinedAt time.Time `json:"joinedAt"`
	Role     string    `json:"role"`
	Status   string    `json:"status"`
}

type Topic struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

type GroupCategory struct {
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Count       int    `json:"count"`
}

type Event struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	GroupID     string    `json:"groupId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Location    string    `json:"location"`
	Type        string    `json:"type"`
	Attendees   []*User   `json:"attendees" gorm:"many2many:event_attendees;"`
}

type Resource struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	GroupID     string    `json:"groupId"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	AddedBy     string    `json:"addedBy"`
	AddedAt     time.Time `json:"addedAt"`
}

type Rule struct {
	ID          string `json:"id" gorm:"primaryKey"`
	GroupID     string `json:"groupId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}
