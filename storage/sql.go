package storage

// UpdateOnly updates only the given columns (includeColumns) and excludes all other
func (conn *Connection) UpdateOnly(model interface{}, includeColumns ...string) error {
	xcols, err := getExcludedColumns(model, includeColumns...)
	if err != nil {
		return err
	}
	return conn.Update(model, xcols...)
}
