package database

func (conn *Connection) InsertSubreddit(name string) error {
	stmt := `INSERT INTO subreddits (name) VALUES ($1)`
	_, err := conn.db.Exec(stmt, name)
	return err
}

func (conn *Connection) UpsertSubreddit(name string) error {
	stmt := `
INSERT INTO subreddits (name)
VALUES ($1)
ON CONFLICT ON CONSTRAINT subreddits_name_key
DO NOTHING;
	`
	_, err := conn.db.Exec(stmt, name)
	return err
}

func (conn *Connection) GetAllSubreddits() (subs []string, err error) {
	stmt := `SELECT name FROM subreddits`
	rows, err := conn.db.Query(stmt)
	if err != nil {
		return subs, err
	}
	defer rows.Close()
	var sub string
	for rows.Next() {
		err = rows.Scan(&sub)
		if err != nil {
			return subs, err
		}
		subs = append(subs, sub)
	}

	return subs, err
}
