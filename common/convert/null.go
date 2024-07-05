package convert

import (
	"database/sql"
	"time"
)

func StrToNullString(str string) (data sql.NullString) {
	if str == "" {
		data.Valid = false
	} else {
		data.Valid = true
	}
	data.String = str
	return
}

func IntToNullInt(value int64) (data sql.NullInt64) {
	data.Valid = true
	data.Int64 = value
	return
}

func TimeToNullTime(value time.Time) (data sql.NullTime) {
	data.Valid = true
	data.Time = value
	return
}
