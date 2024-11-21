package db

import "github.com/flatgrassdotnet/cloudbox/common"

func FetchNewsEntries() ([]common.NewsEntry, error) {
	var entries []common.NewsEntry
	rows, err := handle.Query("SELECT id, title, body, author, time FROM news")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var entry common.NewsEntry
		err := rows.Scan(&entry.ID, &entry.Title, &entry.Body, &entry.Author, &entry.Time)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
