package postgres

func (s *PostgresSuite) TestFollowSpace() {
	userID := s.GetLastUserID()
	// UnFollowSpace should be idiomatic without errors
	s.NoError(UnFollowSpace(userID, SpaceID))

	s.NoError(FollowSpace(userID, SpaceID))
	s.Error(FollowSpace(userID, SpaceID))
	s.NoError(UnFollowSpace(userID, SpaceID))
}
