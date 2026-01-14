package util

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

//URL request: GET /users?limit=20&before=
//URL request: GET /users?limit=20&after=

type Cursor struct {
	Id        string
	CreatedAt time.Time
}

func EncodeCursor(cursor Cursor) string {
	if len(cursor.Id) == 0 {
		return ""
	}

	serializedCursor, err := json.Marshal(cursor)
	if err != nil {
		return ""
	}

	encodedCursor := base64.StdEncoding.EncodeToString(serializedCursor)
	return encodedCursor
}

func DecodeCursor(cursor string) (Cursor, error) {
	decodedCursor, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return Cursor{}, err
	}

	var unmarshaledCrs Cursor
	if err := json.Unmarshal(decodedCursor, &unmarshaledCrs); err != nil {
		return Cursor{}, err
	}

	return unmarshaledCrs, nil
}
