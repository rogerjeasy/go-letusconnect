package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"google.golang.org/api/iterator"
)

type GroupService struct {
	firestoreClient  *firestore.Client
	cloudinaryClient *cloudinary.Cloudinary
	userService      *UserService
}

func NewGroupService(fClient *firestore.Client, cClient *cloudinary.Cloudinary, uUserService *UserService) *GroupService {
	return &GroupService{
		firestoreClient:  fClient,
		cloudinaryClient: cClient,
		userService:      uUserService,
	}
}

// CreateGroup creates a new group with the given input
func (s *GroupService) CreateGroup(ctx context.Context, input models.Group, userId string) (*models.Group, error) {
	if input.Name == "" {
		return nil, errors.New("group name is required")
	}

	if input.ID == "" {
		input.ID = uuid.New().String()
	}

	user, err := s.userService.GetUserByUIDinGoStruct(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	input.Admins = []*models.User{user}

	groupData := mappers.MapGroupGoToFirestore(input)

	_, err = s.firestoreClient.Collection("group_forums").Doc(input.ID).Set(ctx, groupData)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %v", err)
	}

	return &input, nil
}

// GetGroup retrieves a group by ID
func (s *GroupService) GetGroup(ctx context.Context, groupID string) (*models.Group, error) {
	doc, err := s.firestoreClient.Collection("group_forums").Doc(groupID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %v", err)
	}

	group := mappers.MapGroupFirestoreToGo(doc.Data())
	return &group, nil
}

// UpdateGroup updates an existing group
func (s *GroupService) UpdateGroup(ctx context.Context, groupID string, updates models.Group) error {
	updates.UpdatedAt = time.Now()

	updateData := mappers.MapGroupGoToFirestore(updates)

	_, err := s.firestoreClient.Collection("group_forums").Doc(groupID).Set(ctx, updateData, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("failed to update group: %v", err)
	}

	return nil
}

func (s *GroupService) DeleteGroup(ctx context.Context, groupID string, userID string) error {
	group, err := s.GetGroup(ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %v", err)
	}

	isAdmin := false
	for _, admin := range group.Admins {
		if admin.UID == userID {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return errors.New("unauthorized: only group admins can delete the group")
	}

	_, err = s.firestoreClient.Collection("group_forums").Doc(groupID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete group: %v", err)
	}

	return nil
}

// ListGroups retrieves all groups with optional filtering
func (s *GroupService) ListGroups(ctx context.Context, filters map[string]interface{}) ([]models.Group, error) {
	collRef := s.firestoreClient.Collection("group_forums")

	var query firestore.Query
	first := true
	for key, value := range filters {
		if first {
			query = collRef.Where(key, "==", value)
			first = false
		} else {
			query = query.Where(key, "==", value)
		}
	}

	var iter *firestore.DocumentIterator
	if first {
		iter = collRef.Documents(ctx)
	} else {
		iter = query.Documents(ctx)
	}

	var groups []models.Group
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate groups: %v", err)
		}

		group := mappers.MapGroupFirestoreToGo(doc.Data())
		groups = append(groups, group)
	}

	return groups, nil
}

// AddMember adds a member to a group
func (s *GroupService) AddMember(ctx context.Context, groupID string, member models.Member) error {
	group, err := s.GetGroup(ctx, groupID)
	if err != nil {
		return err
	}

	// Check if member already exists
	for _, m := range group.Members {
		if m.UserID == member.UserID {
			return errors.New("member already exists in group")
		}
	}

	member.JoinedAt = time.Now()
	group.Members = append(group.Members, member)
	group.UpdatedAt = time.Now()

	return s.UpdateGroup(ctx, groupID, *group)
}

// RemoveMember removes a member from a group
func (s *GroupService) RemoveMember(ctx context.Context, groupID string, userID string) error {
	group, err := s.GetGroup(ctx, groupID)
	if err != nil {
		return err
	}

	found := false
	newMembers := []models.Member{}
	for _, member := range group.Members {
		if member.UserID != userID {
			newMembers = append(newMembers, member)
		} else {
			found = true
		}
	}

	if !found {
		return errors.New("member not found in group")
	}

	group.Members = newMembers
	group.UpdatedAt = time.Now()

	return s.UpdateGroup(ctx, groupID, *group)
}

// UploadGroupImage uploads a group image to Cloudinary and updates the group
func (s *GroupService) UploadGroupImage(ctx context.Context, groupID string, file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	uploadResult, err := s.cloudinaryClient.Upload.Upload(ctx, src, uploader.UploadParams{
		Folder:   "groups",
		PublicID: fmt.Sprintf("group_%s", groupID),
	})
	if err != nil {
		return fmt.Errorf("failed to upload image: %v", err)
	}

	group, err := s.GetGroup(ctx, groupID)
	if err != nil {
		return err
	}

	group.ImageURL = uploadResult.SecureURL
	group.UpdatedAt = time.Now()

	return s.UpdateGroup(ctx, groupID, *group)
}

// AddEvent adds an event to a group
func (s *GroupService) AddEvent(ctx context.Context, groupID string, event models.Event) error {
	group, err := s.GetGroup(ctx, groupID)
	if err != nil {
		return err
	}

	event.ID = uuid.New().String()
	event.GroupID = groupID
	group.Events = append(group.Events, event)
	group.UpdatedAt = time.Now()

	return s.UpdateGroup(ctx, groupID, *group)
}

// RemoveEvent removes an event from a group
func (s *GroupService) RemoveEvent(ctx context.Context, groupID string, eventID string) error {
	group, err := s.GetGroup(ctx, groupID)
	if err != nil {
		return err
	}

	found := false
	newEvents := []models.Event{}
	for _, event := range group.Events {
		if event.ID != eventID {
			newEvents = append(newEvents, event)
		} else {
			found = true
		}
	}

	if !found {
		return errors.New("event not found in group")
	}

	group.Events = newEvents
	group.UpdatedAt = time.Now()

	return s.UpdateGroup(ctx, groupID, *group)
}

// AddResource adds a resource to a group
func (s *GroupService) AddResource(ctx context.Context, groupID string, resource models.Resource) error {
	group, err := s.GetGroup(ctx, groupID)
	if err != nil {
		return err
	}

	resource.ID = uuid.New().String()
	resource.GroupID = groupID
	resource.AddedAt = time.Now()
	group.Resources = append(group.Resources, resource)
	group.UpdatedAt = time.Now()

	return s.UpdateGroup(ctx, groupID, *group)
}

// RemoveResource removes a resource from a group
func (s *GroupService) RemoveResource(ctx context.Context, groupID string, resourceID string) error {
	group, err := s.GetGroup(ctx, groupID)
	if err != nil {
		return err
	}

	found := false
	newResources := []models.Resource{}
	for _, resource := range group.Resources {
		if resource.ID != resourceID {
			newResources = append(newResources, resource)
		} else {
			found = true
		}
	}

	if !found {
		return errors.New("resource not found in group")
	}

	group.Resources = newResources
	group.UpdatedAt = time.Now()

	return s.UpdateGroup(ctx, groupID, *group)
}

// UpdateGroupSettings updates a group's settings
func (s *GroupService) UpdateGroupSettings(ctx context.Context, groupID string, privacy string, featured bool) error {
	group, err := s.GetGroup(ctx, groupID)
	if err != nil {
		return err
	}

	group.Privacy = privacy
	group.Featured = featured
	group.UpdatedAt = time.Now()

	return s.UpdateGroup(ctx, groupID, *group)
}

// SearchGroups searches for groups based on name or description
func (s *GroupService) SearchGroups(ctx context.Context, query string) ([]models.Group, error) {
	iter := s.firestoreClient.Collection("groups").Documents(ctx)
	var groups []models.Group

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate groups: %v", err)
		}

		group := mappers.MapGroupFirestoreToGo(doc.Data())

		// Simple case-insensitive search
		if containsCaseInsensitive(group.Name, query) || containsCaseInsensitive(group.Description, query) {
			groups = append(groups, group)
		}
	}

	return groups, nil
}

// Helper function for case-insensitive search
func containsCaseInsensitive(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Contains(s, substr)
}
