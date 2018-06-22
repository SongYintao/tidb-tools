package dbutil

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/juju/errors"
)

// ShowVersion queries variable 'version' and returns its value.
func ShowVersion(db *sql.DB) (value string, err error) {
	return ShowMySQLVariable(db, "version")
}

// ShowLogBin queries variable 'log_bin' and returns its value.
func ShowLogBin(db *sql.DB) (value string, err error) {
	return ShowMySQLVariable(db, "log_bin")
}

// ShowBinlogFormat queries variable 'binlog_format' and returns its value.
func ShowBinlogFormat(db *sql.DB) (value string, err error) {
	return ShowMySQLVariable(db, "binlog_format")
}

// ShowBinlogRowImage queries variable 'binlog_row_image' and returns its values.
func ShowBinlogRowImage(db *sql.DB) (value string, err error) {
	return ShowMySQLVariable(db, "binlog_row_image")
}

// ShowServerID queries variable 'server_id' and returns its value.
func ShowServerID(db *sql.DB) (serverID uint64, err error) {
	value, err := ShowMySQLVariable(db, "server_id")
	if err != nil {
		return 0, errors.Trace(err)
	}

	serverID, err = strconv.ParseUint(value, 10, 64)
	return serverID, errors.Annotatef(err, "parse server_id %s failed", value)
}

// ShowMySQLVariable queries MySQL variable and returns its value.
func ShowMySQLVariable(db *sql.DB, variable string) (value string, err error) {
	query := fmt.Sprintf("SHOW GLOBAL VARIABLES LIKE '%s';", variable)
	err = db.QueryRow(query).Scan(&variable, &value)
	if err != nil {
		return "", errors.Trace(err)
	}
	return value, nil
}

// ShowGrants queries privileges for a mysql user.
func ShowGrants(db *sql.DB, user, host string) ([]string, error) {
	if host == "" {
		host = "%"
	}

	var query string
	if user == "" {
		// for currrent user.
		query = "SHOW GRANTS"
	} else {
		query = fmt.Sprintf("SHOW GRANTS FOR '%s'@'%s'", user, host)
	}
	rows, err := db.Query(query)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer rows.Close()

	grants := make([]string, 0, 8)
	for rows.Next() {
		var grant string
		err = rows.Scan(&grant)
		if err != nil {
			return nil, errors.Trace(err)
		}
		grants = append(grants, grant)
	}
	if rows.Err() != nil {
		return nil, errors.Trace(err)
	}

	return grants, nil

}
