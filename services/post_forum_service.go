package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ForumService struct {
	firestoreClient FirestoreClient
	userService     *UserService
}

func NewForumService(fClient FirestoreClient, uService *UserService) *ForumService {
	return &ForumService{
		firestoreClient: fClient,
		userService:     uService,
	}
}

// CreateForum creates a new forum
func (s *ForumService) CreateForum(ctx context.Context, input models.Forum, userID string) (*models.Forum, error) {
	if input.Name == "" {
		return nil, errors.New("forum name is required")
	}

	if input.GroupID == "" {
		return nil, errors.New("group ID is required")
	}

	user, err := s.userService.GetUserByUIDinGoStruct(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	if input.ID == "" {
		input.ID = uuid.New().String()
	}

	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	input.Moderators = []*models.User{user}

	forumData := mappers.MapForumGoToFirestore(input)

	_, err = s.firestoreClient.Collection("forums").Doc(input.ID).Set(ctx, forumData)
	if err != nil {
		return nil, fmt.Errorf("failed to create forum: %v", err)
	}

	return &input, nil
}

// GetForum retrieves a forum by ID
func (s *ForumService) GetForum(ctx context.Context, forumID string) (*models.Forum, error) {
	doc, err := s.firestoreClient.Collection("forums").Doc(forumID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get forum: %v", err)
	}

	forum := mappers.MapForumFirestoreToGo(doc.Data())
	return &forum, nil
}

// UpdateForum updates an existing forum
func (s *ForumService) UpdateForum(ctx context.Context, forumID string, updates models.Forum) error {
	updates.UpdatedAt = time.Now()
	updateData := mappers.MapForumGoToFirestore(updates)

	_, err := s.firestoreClient.Collection("forums").Doc(forumID).Set(ctx, updateData, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("failed to update forum: %v", err)
	}

	return nil
}

// DeleteForum deletes a forum
func (s *ForumService) DeleteForum(ctx context.Context, forumID string, userID string) error {
	forum, err := s.GetForum(ctx, forumID)
	if err != nil {
		return fmt.Errorf("failed to get forum: %v", err)
	}

	isModerator := false
	for _, mod := range forum.Moderators {
		if mod.UID == userID {
			isModerator = true
			break
		}
	}

	if !isModerator {
		return errors.New("unauthorized: only forum moderators can delete the forum")
	}

	_, err = s.firestoreClient.Collection("forums").Doc(forumID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete forum: %v", err)
	}

	return nil
}

// CreatePost creates a new post in a forum
func (s *ForumService) CreatePost(ctx context.Context, forumID string, input models.Post, userID string) (*models.Post, error) {
	forum, err := s.GetForum(ctx, forumID)
	if err != nil {
		return nil, err
	}

	// user, err := s.userService.GetUserByUIDinGoStruct(userID)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to get user: %v", err)
	// }

	input.ID = uuid.New().String()
	input.ForumID = forumID
	input.UserID = userID
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()
	input.Status = "active"

	posts := append(forum.Posts, input)
	forum.Posts = posts
	forum.UpdatedAt = time.Now()

	err = s.UpdateForum(ctx, forumID, *forum)
	if err != nil {
		return nil, err
	}

	return &input, nil
}

// CreateComment adds a comment to a post
func (s *ForumService) CreateComment(ctx context.Context, forumID string, postID string, input models.Comment, userID string) (*models.Comment, error) {
	forum, err := s.GetForum(ctx, forumID)
	if err != nil {
		return nil, err
	}

	var targetPost *models.Post
	for i, post := range forum.Posts {
		if post.ID == postID {
			targetPost = &forum.Posts[i]
			break
		}
	}

	if targetPost == nil {
		return nil, errors.New("post not found")
	}

	input.ID = uuid.New().String()
	input.PostID = postID
	input.UserID = userID
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	targetPost.Comments = append(targetPost.Comments, input)
	forum.UpdatedAt = time.Now()

	err = s.UpdateForum(ctx, forumID, *forum)
	if err != nil {
		return nil, err
	}

	return &input, nil
}

// AddReaction adds a reaction to a post or comment
func (s *ForumService) AddReaction(ctx context.Context, forumID string, reaction models.Reaction, userID string) error {
	forum, err := s.GetForum(ctx, forumID)
	if err != nil {
		return err
	}

	reaction.ID = uuid.New().String()
	reaction.UserID = userID
	reaction.CreatedAt = time.Now()

	if reaction.PostID != nil {
		// Add reaction to post
		for i, post := range forum.Posts {
			if post.ID == *reaction.PostID {
				forum.Posts[i].Reactions = append(forum.Posts[i].Reactions, reaction)
				break
			}
		}
	} else if reaction.CommentID != nil {
		// Add reaction to comment
		for i, post := range forum.Posts {
			for j, comment := range post.Comments {
				if comment.ID == *reaction.CommentID {
					forum.Posts[i].Comments[j].Reactions = append(forum.Posts[i].Comments[j].Reactions, reaction)
					break
				}
			}
		}
	}

	forum.UpdatedAt = time.Now()
	return s.UpdateForum(ctx, forumID, *forum)
}

// ListForumsByGroup retrieves all forums for a specific group
func (s *ForumService) ListForumsByGroup(ctx context.Context, groupID string) ([]models.Forum, error) {
	query := s.firestoreClient.Collection("forums").Where("group_id", "==", groupID)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil && status.Code(err) != codes.NotFound {
		return nil, fmt.Errorf("failed to get forums: %v", err)
	}

	var forums []models.Forum
	for _, doc := range docs {
		forum := mappers.MapForumFirestoreToGo(doc.Data())
		forums = append(forums, forum)
	}

	return forums, nil
}

// SearchPosts searches for posts within a forum
func (s *ForumService) SearchPosts(ctx context.Context, forumID string, query string) ([]models.Post, error) {
	forum, err := s.GetForum(ctx, forumID)
	if err != nil {
		return nil, err
	}

	var matchingPosts []models.Post
	for _, post := range forum.Posts {
		if containsCaseInsensitive(post.Title, query) || containsCaseInsensitive(post.Content, query) {
			matchingPosts = append(matchingPosts, post)
		}
	}

	return matchingPosts, nil
}

// AddModerator adds a moderator to a forum
func (s *ForumService) AddModerator(ctx context.Context, forumID string, userID string) error {
	forum, err := s.GetForum(ctx, forumID)
	if err != nil {
		return err
	}

	user, err := s.userService.GetUserByUIDinGoStruct(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	// Check if already a moderator
	for _, mod := range forum.Moderators {
		if mod.UID == userID {
			return errors.New("user is already a moderator")
		}
	}

	forum.Moderators = append(forum.Moderators, user)
	forum.UpdatedAt = time.Now()

	return s.UpdateForum(ctx, forumID, *forum)
}

// RemoveModerator removes a moderator from a forum
func (s *ForumService) RemoveModerator(ctx context.Context, forumID string, userID string) error {
	forum, err := s.GetForum(ctx, forumID)
	if err != nil {
		return err
	}

	var newModerators []*models.User
	found := false
	for _, mod := range forum.Moderators {
		if mod.UID != userID {
			newModerators = append(newModerators, mod)
		} else {
			found = true
		}
	}

	if !found {
		return errors.New("user is not a moderator")
	}

	if len(newModerators) == 0 {
		return errors.New("cannot remove last moderator")
	}

	forum.Moderators = newModerators
	forum.UpdatedAt = time.Now()

	return s.UpdateForum(ctx, forumID, *forum)
}
