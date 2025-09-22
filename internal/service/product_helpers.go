package service

import (
	"database/sql"
	"encoding/json"
)

func parseStringArray(raw sql.NullString) []string {
	if !raw.Valid || raw.String == "" {
		return nil
	}

	var parsed []string
	if err := json.Unmarshal([]byte(raw.String), &parsed); err != nil {
		return nil
	}

	return parsed
}

func stringSliceToJSONArg(values []string) (any, error) {
	if values == nil {
		return nil, nil
	}

	data, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	return string(data), nil
}
