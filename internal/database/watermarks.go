package database

import (
	"time"
)

func (conn *Connection) GenerateRunID() (int, error) {
	stmt :=
		`
INSERT INTO watermarks (run_start)
VALUES ($1);
	`
	_, err := conn.db.Exec(stmt, time.Now())
	if err != nil {
		return 0, err
	}

	// get the last id inserted
	stmt =
		`
SELECT 
	MAX(id)
FROM 
	watermarks;
	`
	res, err := conn.db.Query(stmt)
	if err != nil {
		return 0, err
	}
	defer res.Close()

	var id int
	for res.Next() {
		err = res.Scan(&id)
	}

	return id, err
}

type UpdateWatermarkInput struct {
	Id         int
	RunEnd     time.Time
	Successful bool
}

func (conn *Connection) UpdateWatermark(input UpdateWatermarkInput) error {
	stmt :=
		`
UPDATE
	watermarks
SET
	run_end=$1,
	successful=$2
WHERE
	id=$3;
	`
	_, err := conn.db.Exec(stmt, input.RunEnd, input.Successful, input.Id)
	return err
}
