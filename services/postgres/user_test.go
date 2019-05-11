package postgres

func (s *PostgresSuite) TestInsertUser() {
	if err := InsertUser("philip"); err.Error() != "pq: duplicate key value violates unique constraint \"user_username_key\"" {
		s.NoError(err)
	}
}

func (s *PostgresSuite) GetLastUserID() int64 {
	var userID int64
	s.NoError(DB.QueryRow("SELECT MAX(id) FROM \"user\"").Scan(&userID))
	return userID
}
