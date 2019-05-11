package postgres

func (s *PostgresSuite) TestPostVoting() {
	s.TestInsertPost()
	s.TestInsertUser()
	postID := s.GetLastPostID()
	userID := s.GetLastUserID()
	s.NoError(PostUnUpvote(userID, postID))
	s.NoError(PostUnDownvote(userID, postID))

	post, err := GetPost(postID)
	if !s.NoError(err) {
		s.FailNow("error getting post")
	}
	s.NotNil(post)
	upvotes := post.Upvotes
	downvotes := post.Downvotes

	s.NoError(PostUpvote(userID, postID))
	s.Error(PostUpvote(userID, postID))
	s.Error(PostDownvote(userID, postID))

	s.NoError(PostUnDownvote(userID, postID))
	s.Error(PostDownvote(userID, postID))
	s.NoError(PostUnUpvote(userID, postID))
	s.NoError(PostDownvote(userID, postID))
	s.Error(PostDownvote(userID, postID))
	s.Error(PostDownvote(userID, postID))
	s.NoError(PostUnDownvote(userID, postID))

	post, err = GetPost(postID)
	s.Equal(upvotes, post.Upvotes)
	s.Equal(downvotes, post.Downvotes)
}
