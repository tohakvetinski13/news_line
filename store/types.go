package store

import "time"

type Subscriptions struct {
	Subscriptions []string `json:"users"`
}

type Users struct {
	Users []string `json:"users"`
}

type Like struct {
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type News struct {
	ID     string    `json:"id"`
	Title  string    `json:"title"`
	Text   string    `json:"text"`
	Date   time.Time `json:"date_create"`
	Likes  string    `json:"likes"`
	UserID string    `json:"user_id"`
}
