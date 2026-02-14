package util

import (
	"fmt"
	"database/sql"
)

func InterfaceToNullString(myInterface interface{}) sql.NullString{
	return sql.NullString{
		String: fmt.Sprintf("%v", myInterface),
		Valid: true,
	}
}

func StringToNullString(myString string) sql.NullString{
	return sql.NullString{
		String: myString,
		Valid: true,
	}
}