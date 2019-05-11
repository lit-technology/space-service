package postgres

func (s *PostgresSuite) TestInsertPost() {
	s.NoError(InsertPostForUser(s.GetLastUserID(), "Test", []string{}, true, SpaceID, 0, 0, 0, 0, 0))
}

func (s *PostgresSuite) TestGetPost() {
	s.TestInsertPost()
	post, err := GetPost(s.GetLastPostID())
	s.NoError(err)
	s.NotNil(post)
}

func (s *PostgresSuite) TestGetPostsForSpace() {
	posts, err := GetPostsForSpace(SpaceID, 0)
	s.NoError(err)
	s.NotEmpty(posts)
}

func (s *PostgresSuite) TestGetPostsFromUser() {
	posts, err := GetPostsFromUser(s.GetLastUserID(), 0)
	s.NoError(err)
	s.NotEmpty(posts)
}

func (s *PostgresSuite) GetLastPostID() int64 {
	var postID int64
	s.NoError(DB.QueryRow("SELECT MAX(id) FROM post").Scan(&postID))
	return postID
}
