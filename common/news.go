package common

import "time"

type NewsEntry struct {
	ID     int       `json:"id"`
	Title  string    `json:"title"`
	Body   string    `json:"body"`
	Author string    `json:"author"`
	Time   time.Time `json:"time"`
}
