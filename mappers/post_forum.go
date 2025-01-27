package mappers

import (
	"time"

	"github.com/rogerjeasy/go-letusconnect/models"
)

// Forum Mappers
func MapForumFrontendToGo(data map[string]interface{}) models.Forum {
	return models.Forum{
		ID:                  getStringValue(data, "id"),
		GroupID:             getStringValue(data, "groupId"),
		Name:                getStringValue(data, "name"),
		Description:         getStringValue(data, "description"),
		IsArchived:          getBoolValue(data, "isArchived"),
		CreatedAt:           getTimeValue(data, "createdAt"),
		UpdatedAt:           getTimeValue(data, "updatedAt"),
		Posts:               MapPostsArrayFrontendToGo(data, "posts"),
		Moderators:          MapModeratorsArrayFrontendToGo(data, "moderators"),
		Categories:          MapForumCategoriesArrayFrontendToGo(data, "categories"),
		AllowAnonymousPosts: getBoolValue(data, "allowAnonymousPosts"),
		RequireModeration:   getBoolValue(data, "requireModeration"),
		AllowFiles:          getBoolValue(data, "allowFiles"),
		MaxFileSize:         getInt64Value(data, "maxFileSize"),
	}
}

func MapForumGoToFirestore(forum models.Forum) map[string]interface{} {
	return map[string]interface{}{
		"id":                    forum.ID,
		"group_id":              forum.GroupID,
		"name":                  forum.Name,
		"description":           forum.Description,
		"is_archived":           forum.IsArchived,
		"created_at":            forum.CreatedAt,
		"updated_at":            forum.UpdatedAt,
		"posts":                 MapPostsArrayGoToFirestore(forum.Posts),
		"moderators":            MapModeratorsArrayGoToFirestore(forum.Moderators),
		"categories":            MapForumCategoriesArrayGoToFirestore(forum.Categories),
		"allow_anonymous_posts": forum.AllowAnonymousPosts,
		"require_moderation":    forum.RequireModeration,
		"allow_files":           forum.AllowFiles,
		"max_file_size":         forum.MaxFileSize,
	}
}

func MapForumFirestoreToGo(data map[string]interface{}) models.Forum {
	return models.Forum{
		ID:                  getStringValue(data, "id"),
		GroupID:             getStringValue(data, "group_id"),
		Name:                getStringValue(data, "name"),
		Description:         getStringValue(data, "description"),
		IsArchived:          getBoolValue(data, "is_archived"),
		CreatedAt:           getFirestoreTimeToGoTime(data["created_at"]),
		UpdatedAt:           getFirestoreTimeToGoTime(data["updated_at"]),
		Posts:               MapPostsArrayFirestoreToGo(data, "posts"),
		Moderators:          MapModeratorsArrayFirestoreToGo(data, "moderators"),
		Categories:          MapForumCategoriesArrayFirestoreToGo(data, "categories"),
		AllowAnonymousPosts: getBoolValue(data, "allow_anonymous_posts"),
		RequireModeration:   getBoolValue(data, "require_moderation"),
		AllowFiles:          getBoolValue(data, "allow_files"),
		MaxFileSize:         getInt64Value(data, "max_file_size"),
	}
}

func MapForumGoToFrontend(forum models.Forum) map[string]interface{} {
	return map[string]interface{}{
		"id":                  forum.ID,
		"groupId":             forum.GroupID,
		"name":                forum.Name,
		"description":         forum.Description,
		"isArchived":          forum.IsArchived,
		"createdAt":           forum.CreatedAt.Format(time.RFC3339),
		"updatedAt":           forum.UpdatedAt.Format(time.RFC3339),
		"posts":               MapPostsArrayGoToFrontend(forum.Posts),
		"moderators":          MapModeratorsArrayGoToFrontend(forum.Moderators),
		"categories":          MapForumCategoriesArrayGoToFrontend(forum.Categories),
		"allowAnonymousPosts": forum.AllowAnonymousPosts,
		"requireModeration":   forum.RequireModeration,
		"allowFiles":          forum.AllowFiles,
		"maxFileSize":         forum.MaxFileSize,
	}
}

func MapPostsArrayGoToFrontend(posts []models.Post) []map[string]interface{} {
	var result []map[string]interface{}
	for _, post := range posts {
		result = append(result, map[string]interface{}{
			"id":          post.ID,
			"forumId":     post.ForumID,
			"userId":      post.UserID,
			"title":       post.Title,
			"content":     post.Content,
			"status":      post.Status,
			"isSticky":    post.IsSticky,
			"isAnonymous": post.IsAnonymous,
			"viewCount":   post.ViewCount,
			"createdAt":   post.CreatedAt.Format(time.RFC3339),
			"updatedAt":   post.UpdatedAt.Format(time.RFC3339),
			"comments":    MapCommentsArrayGoToFrontend(post.Comments),
			"tags":        MapTagsArrayGoToFrontend(post.Tags),
			"reactions":   MapReactionsArrayGoToFrontend(post.Reactions),
			"files":       MapFilesArrayGoToFrontend(post.Files),
		})
	}
	return result
}

// Post Mappers
func MapPostFrontendToGo(data map[string]interface{}) models.Post {
	return models.Post{
		ID:          getStringValue(data, "id"),
		ForumID:     getStringValue(data, "forumId"),
		UserID:      getStringValue(data, "userId"),
		Title:       getStringValue(data, "title"),
		Content:     getStringValue(data, "content"),
		Status:      getStringValue(data, "status"),
		IsSticky:    getBoolValue(data, "isSticky"),
		IsAnonymous: getBoolValue(data, "isAnonymous"),
		ViewCount:   getIntValueSafe(data, "viewCount"),
		CreatedAt:   getTimeValue(data, "createdAt"),
		UpdatedAt:   getTimeValue(data, "updatedAt"),
		Comments:    MapCommentsArrayFrontendToGo(data, "comments"),
		Tags:        MapTagsArrayFrontendToGo(data, "tags"),
		Reactions:   MapReactionsArrayFrontendToGo(data, "reactions"),
		Files:       MapFilesArrayFrontendToGo(data, "files"),
	}
}

func MapPostGoToFirestore(post models.Post) map[string]interface{} {
	return map[string]interface{}{
		"id":           post.ID,
		"forum_id":     post.ForumID,
		"user_id":      post.UserID,
		"title":        post.Title,
		"content":      post.Content,
		"status":       post.Status,
		"is_sticky":    post.IsSticky,
		"is_anonymous": post.IsAnonymous,
		"view_count":   post.ViewCount,
		"created_at":   post.CreatedAt,
		"updated_at":   post.UpdatedAt,
		"comments":     MapCommentsArrayGoToFirestore(post.Comments),
		"tags":         MapTagsArrayGoToFirestore(post.Tags),
		"reactions":    MapReactionsArrayGoToFirestore(post.Reactions),
		"files":        MapFilesArrayGoToFirestore(post.Files),
	}
}

// Comments Array Frontend Mapper
func MapCommentsArrayGoToFrontend(comments []models.Comment) []map[string]interface{} {
	var result []map[string]interface{}
	for _, comment := range comments {
		result = append(result, map[string]interface{}{
			"id":        comment.ID,
			"postId":    comment.PostID,
			"userId":    comment.UserID,
			"parentId":  comment.ParentID,
			"content":   comment.Content,
			"isEdited":  comment.IsEdited,
			"createdAt": comment.CreatedAt.Format(time.RFC3339),
			"updatedAt": comment.UpdatedAt.Format(time.RFC3339),
			"reactions": MapReactionsArrayGoToFrontend(comment.Reactions),
		})
	}
	return result
}

// Tags Array Frontend Mapper
func MapTagsArrayGoToFrontend(tags []models.Tag) []map[string]interface{} {
	var result []map[string]interface{}
	for _, tag := range tags {
		result = append(result, map[string]interface{}{
			"id":          tag.ID,
			"name":        tag.Name,
			"description": tag.Description,
			"color":       tag.Color,
		})
	}
	return result
}

// Reactions Array Frontend Mapper
func MapReactionsArrayGoToFrontend(reactions []models.Reaction) []map[string]interface{} {
	var result []map[string]interface{}
	for _, reaction := range reactions {
		result = append(result, map[string]interface{}{
			"id":        reaction.ID,
			"userId":    reaction.UserID,
			"postId":    reaction.PostID,
			"commentId": reaction.CommentID,
			"type":      reaction.Type,
			"createdAt": reaction.CreatedAt.Format(time.RFC3339),
		})
	}
	return result
}

// Files Array Frontend Mapper
func MapFilesArrayGoToFrontend(files []models.File) []map[string]interface{} {
	var result []map[string]interface{}
	for _, file := range files {
		result = append(result, map[string]interface{}{
			"id":         file.ID,
			"postId":     file.PostID,
			"userId":     file.UserID,
			"fileName":   file.FileName,
			"fileType":   file.FileType,
			"fileSize":   file.FileSize,
			"url":        file.URL,
			"uploadedAt": file.UploadedAt.Format(time.RFC3339),
		})
	}
	return result
}

// Post Firestore to Go Mapper
func MapPostFirestoreToGo(data map[string]interface{}) models.Post {
	return models.Post{
		ID:          getStringValue(data, "id"),
		ForumID:     getStringValue(data, "forum_id"),
		UserID:      getStringValue(data, "user_id"),
		Title:       getStringValue(data, "title"),
		Content:     getStringValue(data, "content"),
		Status:      getStringValue(data, "status"),
		IsSticky:    getBoolValue(data, "is_sticky"),
		IsAnonymous: getBoolValue(data, "is_anonymous"),
		ViewCount:   getIntValueSafe(data, "view_count"),
		CreatedAt:   getFirestoreTimeToGoTime(data["created_at"]),
		UpdatedAt:   getFirestoreTimeToGoTime(data["updated_at"]),
		Comments:    MapCommentsArrayFirestoreToGo(data, "comments"),
		Tags:        MapTagsArrayFirestoreToGo(data, "tags"),
		Reactions:   MapReactionsArrayFirestoreToGo(data, "reactions"),
		Files:       MapFilesArrayFirestoreToGo(data, "files"),
	}
}

// Comments Array Firestore Mapper
func MapCommentsArrayGoToFirestore(comments []models.Comment) []map[string]interface{} {
	var result []map[string]interface{}
	for _, comment := range comments {
		result = append(result, map[string]interface{}{
			"id":         comment.ID,
			"post_id":    comment.PostID,
			"user_id":    comment.UserID,
			"parent_id":  comment.ParentID,
			"content":    comment.Content,
			"is_edited":  comment.IsEdited,
			"created_at": comment.CreatedAt,
			"updated_at": comment.UpdatedAt,
			"reactions":  MapReactionsArrayGoToFirestore(comment.Reactions),
		})
	}
	return result
}

func MapCommentsArrayFirestoreToGo(data map[string]interface{}, key string) []models.Comment {
	var comments []models.Comment
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if commentMap, ok := item.(map[string]interface{}); ok {
				var parentID *string
				if pid := getStringValue(commentMap, "parent_id"); pid != "" {
					parentID = &pid
				}

				comments = append(comments, models.Comment{
					ID:        getStringValue(commentMap, "id"),
					PostID:    getStringValue(commentMap, "post_id"),
					UserID:    getStringValue(commentMap, "user_id"),
					ParentID:  parentID,
					Content:   getStringValue(commentMap, "content"),
					IsEdited:  getBoolValue(commentMap, "is_edited"),
					CreatedAt: getFirestoreTimeToGoTime(commentMap["created_at"]),
					UpdatedAt: getFirestoreTimeToGoTime(commentMap["updated_at"]),
					Reactions: MapReactionsArrayFirestoreToGo(commentMap, "reactions"),
				})
			}
		}
	}
	return comments
}

func MapTagsArrayFirestoreToGo(data map[string]interface{}, key string) []models.Tag {
	var tags []models.Tag
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if tagMap, ok := item.(map[string]interface{}); ok {
				tags = append(tags, models.Tag{
					ID:          getStringValue(tagMap, "id"),
					Name:        getStringValue(tagMap, "name"),
					Description: getStringValue(tagMap, "description"),
					Color:       getStringValue(tagMap, "color"),
				})
			}
		}
	}
	return tags
}

func MapReactionsArrayFirestoreToGo(data map[string]interface{}, key string) []models.Reaction {
	var reactions []models.Reaction
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if reactionMap, ok := item.(map[string]interface{}); ok {
				var postID, commentID *string
				if pid := getStringValue(reactionMap, "post_id"); pid != "" {
					postID = &pid
				}
				if cid := getStringValue(reactionMap, "comment_id"); cid != "" {
					commentID = &cid
				}

				reactions = append(reactions, models.Reaction{
					ID:        getStringValue(reactionMap, "id"),
					UserID:    getStringValue(reactionMap, "user_id"),
					PostID:    postID,
					CommentID: commentID,
					Type:      getStringValue(reactionMap, "type"),
					CreatedAt: getFirestoreTimeToGoTime(reactionMap["created_at"]),
				})
			}
		}
	}
	return reactions
}

func MapFilesArrayFirestoreToGo(data map[string]interface{}, key string) []models.File {
	var files []models.File
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if fileMap, ok := item.(map[string]interface{}); ok {
				files = append(files, models.File{
					ID:         getStringValue(fileMap, "id"),
					PostID:     getStringValue(fileMap, "post_id"),
					UserID:     getStringValue(fileMap, "user_id"),
					FileName:   getStringValue(fileMap, "file_name"),
					FileType:   getStringValue(fileMap, "file_type"),
					FileSize:   getInt64Value(fileMap, "file_size"),
					URL:        getStringValue(fileMap, "url"),
					UploadedAt: getFirestoreTimeToGoTime(fileMap["uploaded_at"]),
				})
			}
		}
	}
	return files
}

// Tags Array Firestore Mapper
func MapTagsArrayGoToFirestore(tags []models.Tag) []map[string]interface{} {
	var result []map[string]interface{}
	for _, tag := range tags {
		result = append(result, map[string]interface{}{
			"id":          tag.ID,
			"name":        tag.Name,
			"description": tag.Description,
			"color":       tag.Color,
		})
	}
	return result
}

// Reactions Array Firestore Mapper
func MapReactionsArrayGoToFirestore(reactions []models.Reaction) []map[string]interface{} {
	var result []map[string]interface{}
	for _, reaction := range reactions {
		result = append(result, map[string]interface{}{
			"id":         reaction.ID,
			"user_id":    reaction.UserID,
			"post_id":    reaction.PostID,
			"comment_id": reaction.CommentID,
			"type":       reaction.Type,
			"created_at": reaction.CreatedAt,
		})
	}
	return result
}

// Files Array Firestore Mapper
func MapFilesArrayGoToFirestore(files []models.File) []map[string]interface{} {
	var result []map[string]interface{}
	for _, file := range files {
		result = append(result, map[string]interface{}{
			"id":          file.ID,
			"post_id":     file.PostID,
			"user_id":     file.UserID,
			"file_name":   file.FileName,
			"file_type":   file.FileType,
			"file_size":   file.FileSize,
			"url":         file.URL,
			"uploaded_at": file.UploadedAt,
		})
	}
	return result
}

func MapModeratorsArrayGoToFirestore(moderators []*models.User) []map[string]interface{} {
	var result []map[string]interface{}
	for _, moderator := range moderators {
		result = append(result, MapUserFrontendToBackend(moderator))
	}
	return result
}

func MapModeratorsArrayGoToFrontend(moderators []*models.User) []map[string]interface{} {
	var result []map[string]interface{}
	for _, moderator := range moderators {
		result = append(result, MapUserToFrontend(moderator))
	}
	return result
}

func MapForumCategoriesArrayGoToFrontend(categories []models.ForumCategory) []map[string]interface{} {
	var result []map[string]interface{}
	for _, category := range categories {
		result = append(result, map[string]interface{}{
			"id":          category.ID,
			"forumId":     category.ForumID,
			"name":        category.Name,
			"description": category.Description,
			"color":       category.Color,
			"order":       category.Order,
		})
	}
	return result
}

func MapModeratorsArrayFirestoreToGo(data map[string]interface{}, key string) []*models.User {
	var moderators []*models.User
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if moderatorMap, ok := item.(map[string]interface{}); ok {
				user := MapBackendToUser(moderatorMap)
				moderators = append(moderators, &user)
			}
		}
	}
	return moderators
}

func MapForumCategoriesArrayFirestoreToGo(data map[string]interface{}, key string) []models.ForumCategory {
	var categories []models.ForumCategory
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if categoryMap, ok := item.(map[string]interface{}); ok {
				categories = append(categories, models.ForumCategory{
					ID:          getStringValue(categoryMap, "id"),
					ForumID:     getStringValue(categoryMap, "forum_id"),
					Name:        getStringValue(categoryMap, "name"),
					Description: getStringValue(categoryMap, "description"),
					Color:       getStringValue(categoryMap, "color"),
					Order:       getIntValueSafe(categoryMap, "order"),
				})
			}
		}
	}
	return categories
}

// Forum Categories Array Mappers
func MapForumCategoriesArrayGoToFirestore(categories []models.ForumCategory) []map[string]interface{} {
	var result []map[string]interface{}
	for _, category := range categories {
		result = append(result, map[string]interface{}{
			"id":          category.ID,
			"forum_id":    category.ForumID,
			"name":        category.Name,
			"description": category.Description,
			"color":       category.Color,
			"order":       category.Order,
		})
	}
	return result
}

func MapModeratorsArrayFrontendToGo(data map[string]interface{}, key string) []*models.User {
	var moderators []*models.User
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if moderatorMap, ok := item.(map[string]interface{}); ok {
				user := MapFrontendToUser(moderatorMap)
				moderators = append(moderators, &user)
			}
		}
	}
	return moderators
}

// Array Mappers for Posts
func MapPostsArrayFrontendToGo(data map[string]interface{}, key string) []models.Post {
	var posts []models.Post
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if postMap, ok := item.(map[string]interface{}); ok {
				posts = append(posts, MapPostFrontendToGo(postMap))
			}
		}
	}
	return posts
}

func MapPostsArrayGoToFirestore(posts []models.Post) []map[string]interface{} {
	var result []map[string]interface{}
	for _, post := range posts {
		result = append(result, MapPostGoToFirestore(post))
	}
	return result
}

func MapPostsArrayFirestoreToGo(data map[string]interface{}, key string) []models.Post {
	var posts []models.Post
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if postMap, ok := item.(map[string]interface{}); ok {
				post := MapPostFirestoreToGo(postMap)
				posts = append(posts, post)
			}
		}
	}
	return posts
}

// Array Mappers for Comments
func MapCommentFrontendToGo(data map[string]interface{}) models.Comment {
	var parentID *string
	if pid := getStringValue(data, "parentId"); pid != "" {
		parentID = &pid
	}

	return models.Comment{
		ID:        getStringValue(data, "id"),
		PostID:    getStringValue(data, "postId"),
		UserID:    getStringValue(data, "userId"),
		ParentID:  parentID,
		Content:   getStringValue(data, "content"),
		IsEdited:  getBoolValue(data, "isEdited"),
		CreatedAt: getTimeValue(data, "createdAt"),
		UpdatedAt: getTimeValue(data, "updatedAt"),
		Reactions: MapReactionsArrayFrontendToGo(data, "reactions"),
	}
}

// Comments Array Frontend Mapper
func MapCommentsArrayFrontendToGo(data map[string]interface{}, key string) []models.Comment {
	var comments []models.Comment
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if commentMap, ok := item.(map[string]interface{}); ok {
				comment := MapCommentFrontendToGo(commentMap)
				comments = append(comments, comment)
			}
		}
	}
	return comments
}

// Reaction Mappers
func MapReactionFrontendToGo(data map[string]interface{}) models.Reaction {
	var postID, commentID *string
	if pid := getStringValue(data, "postId"); pid != "" {
		postID = &pid
	}
	if cid := getStringValue(data, "commentId"); cid != "" {
		commentID = &cid
	}

	return models.Reaction{
		ID:        getStringValue(data, "id"),
		UserID:    getStringValue(data, "userId"),
		PostID:    postID,
		CommentID: commentID,
		Type:      getStringValue(data, "type"),
		CreatedAt: getTimeValue(data, "createdAt"),
	}
}

// Array Mappers for Reactions
func MapReactionsArrayFrontendToGo(data map[string]interface{}, key string) []models.Reaction {
	var reactions []models.Reaction
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if reactionMap, ok := item.(map[string]interface{}); ok {
				reactions = append(reactions, MapReactionFrontendToGo(reactionMap))
			}
		}
	}
	return reactions
}

// File Mappers
func MapFileFrontendToGo(data map[string]interface{}) models.File {
	return models.File{
		ID:         getStringValue(data, "id"),
		PostID:     getStringValue(data, "postId"),
		UserID:     getStringValue(data, "userId"),
		FileName:   getStringValue(data, "fileName"),
		FileType:   getStringValue(data, "fileType"),
		FileSize:   getInt64Value(data, "fileSize"),
		URL:        getStringValue(data, "url"),
		UploadedAt: getTimeValue(data, "uploadedAt"),
	}
}

// Array Mappers for Files
func MapFilesArrayFrontendToGo(data map[string]interface{}, key string) []models.File {
	var files []models.File
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if fileMap, ok := item.(map[string]interface{}); ok {
				files = append(files, MapFileFrontendToGo(fileMap))
			}
		}
	}
	return files
}

// Tag Mappers
func MapTagFrontendToGo(data map[string]interface{}) models.Tag {
	return models.Tag{
		ID:          getStringValue(data, "id"),
		Name:        getStringValue(data, "name"),
		Description: getStringValue(data, "description"),
		Color:       getStringValue(data, "color"),
	}
}

// Array Mappers for Tags
func MapTagsArrayFrontendToGo(data map[string]interface{}, key string) []models.Tag {
	var tags []models.Tag
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if tagMap, ok := item.(map[string]interface{}); ok {
				tags = append(tags, MapTagFrontendToGo(tagMap))
			}
		}
	}
	return tags
}

// Forum Category Mappers
func MapForumCategoryFrontendToGo(data map[string]interface{}) models.ForumCategory {
	return models.ForumCategory{
		ID:          getStringValue(data, "id"),
		ForumID:     getStringValue(data, "forumId"),
		Name:        getStringValue(data, "name"),
		Description: getStringValue(data, "description"),
		Color:       getStringValue(data, "color"),
		Order:       getIntValueSafe(data, "order"),
	}
}

// Array Mappers for Forum Categories
func MapForumCategoriesArrayFrontendToGo(data map[string]interface{}, key string) []models.ForumCategory {
	var categories []models.ForumCategory
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if categoryMap, ok := item.(map[string]interface{}); ok {
				categories = append(categories, MapForumCategoryFrontendToGo(categoryMap))
			}
		}
	}
	return categories
}

// ForumModerator Mappers
func MapForumModeratorFrontendToGo(data map[string]interface{}) models.ForumModerator {
	return models.ForumModerator{
		UserID:      getStringValue(data, "userId"),
		ForumID:     getStringValue(data, "forumId"),
		Permissions: getStringArrayValue(data, "permissions"),
		AddedAt:     getTimeValue(data, "addedAt"),
	}
}
