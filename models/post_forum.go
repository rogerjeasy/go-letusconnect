package models

import (
	"time"
)

// Forum represents a discussion forum within a group
type Forum struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	GroupID     string    `json:"groupId" gorm:"index"` // Reference to parent group
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsArchived  bool      `json:"isArchived" gorm:"default:false"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	// Relationships
	Posts      []Post          `json:"posts" gorm:"foreignKey:ForumID"`
	Moderators []*User         `json:"moderators" gorm:"many2many:forum_moderators;"`
	Categories []ForumCategory `json:"categories" gorm:"foreignKey:ForumID"`

	// Settings
	AllowAnonymousPosts bool  `json:"allowAnonymousPosts" gorm:"default:false"`
	RequireModeration   bool  `json:"requireModeration" gorm:"default:false"`
	AllowFiles          bool  `json:"allowFiles" gorm:"default:true"`
	MaxFileSize         int64 `json:"maxFileSize" gorm:"default:5242880"` // Default 5MB
}

// Post represents a discussion thread in the forum
type Post struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	ForumID     string    `json:"forumId" gorm:"index"`
	UserID      string    `json:"userId" gorm:"index"`
	Title       string    `json:"title"`
	Content     string    `json:"content" gorm:"type:text"`
	Status      string    `json:"status" gorm:"default:'active'"` // active, locked, hidden
	IsSticky    bool      `json:"isSticky" gorm:"default:false"`
	IsAnonymous bool      `json:"isAnonymous" gorm:"default:false"`
	ViewCount   int       `json:"viewCount" gorm:"default:0"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	// Relationships
	Comments  []Comment  `json:"comments" gorm:"foreignKey:PostID"`
	Tags      []Tag      `json:"tags" gorm:"many2many:post_tags;"`
	Reactions []Reaction `json:"reactions" gorm:"foreignKey:PostID"`
	Files     []File     `json:"files" gorm:"foreignKey:PostID"`
}

// Comment represents a response to a post
type Comment struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	PostID    string    `json:"postId" gorm:"index"`
	UserID    string    `json:"userId" gorm:"index"`
	ParentID  *string   `json:"parentId"` // For nested comments
	Content   string    `json:"content" gorm:"type:text"`
	IsEdited  bool      `json:"isEdited" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Relationships
	Reactions []Reaction `json:"reactions" gorm:"foreignKey:CommentID"`
}

// ForumCategory helps organize forum posts
type ForumCategory struct {
	ID          string `json:"id" gorm:"primaryKey"`
	ForumID     string `json:"forumId" gorm:"index"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Order       int    `json:"order"`
}

// Tag helps categorize posts
type Tag struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"uniqueIndex"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

// Reaction represents user reactions to posts and comments
type Reaction struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"userId" gorm:"index"`
	PostID    *string   `json:"postId" gorm:"index"`
	CommentID *string   `json:"commentId" gorm:"index"`
	Type      string    `json:"type"` // like, heart, laugh, etc.
	CreatedAt time.Time `json:"createdAt"`
}

// File represents attached files in posts
type File struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	PostID     string    `json:"postId" gorm:"index"`
	UserID     string    `json:"userId"`
	FileName   string    `json:"fileName"`
	FileType   string    `json:"fileType"`
	FileSize   int64     `json:"fileSize"`
	URL        string    `json:"url"`
	UploadedAt time.Time `json:"uploadedAt"`
}

// ForumModerator represents the moderator permissions
type ForumModerator struct {
	UserID      string    `json:"userId" gorm:"primaryKey"`
	ForumID     string    `json:"forumId" gorm:"primaryKey"`
	Permissions []string  `json:"permissions" gorm:"type:json"`
	AddedAt     time.Time `json:"addedAt"`
}
