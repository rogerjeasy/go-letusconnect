package models

import "time"

// Message represents a chat message
type Message struct {
	ID         string    `json:"id" firestore:"id"`
	SenderID   string    `json:"senderId" firestore:"sender_id"`
	ReceiverID string    `json:"receiverId" firestore:"receiver_id"`
	Content    string    `json:"content" firestore:"content"`
	CreatedAt  time.Time `json:"createdAt" firestore:"created_at"`
}

// BaseMessage contains common fields for all messages
type BaseMessage struct {
	ID          string          `json:"id" firestore:"id"`
	SenderID    string          `json:"senderId" firestore:"sender_id"`
	SenderName  string          `json:"senderName" firestore:"sender_name"`
	Content     string          `json:"content" firestore:"content"`
	CreatedAt   string          `json:"createdAt" firestore:"created_at"`
	UpdatedAt   string          `json:"updatedAt,omitempty" firestore:"updated_at,omitempty"`
	ReadStatus  map[string]bool `json:"readStatus" firestore:"read_status"`
	IsDeleted   bool            `json:"isDeleted" firestore:"is_deleted"`
	Attachments []string        `json:"attachments,omitempty" firestore:"attachments,omitempty"`
	Reactions   map[string]int  `json:"reactions,omitempty" firestore:"reactions,omitempty"`
	MessageType string          `json:"messageType" firestore:"message_type"`
	ReplyToID   *string         `json:"replyToId,omitempty" firestore:"reply_to_id,omitempty"`
	IsPinned    bool            `json:"isPinned" firestore:"is_pinned"`
	Priority    string          `json:"priority,omitempty" firestore:"priority,omitempty"`
}

// DirectMessage for one-to-one messaging
type DirectMessage struct {
	BaseMessage
	ReceiverID   string `json:"receiverId" firestore:"receiver_id"`
	ReceiverName string `json:"receiverName" firestore:"receiver_name"`
}

// GroupMessage for group messaging
type GroupMessage struct {
	BaseMessage
	ProjectID string  `json:"projectId" firestore:"project_id"`
	GroupID   *string `json:"groupId,omitempty" firestore:"group_id,omitempty"`
}

type Messages struct {
	ChannelID      string          `json:"channelId" firestore:"channel_id"`
	DirectMessages []DirectMessage `json:"directMessages" firestore:"direct_messages"`
}

type GroupChat struct {
	ID             string          `json:"id" firestore:"id"`
	ProjectID      string          `json:"projectId" firestore:"project_id"`
	CreatedByUID   string          `json:"createdByUid" firestore:"created_by_uid"`
	CreatedByName  string          `json:"createdByName" firestore:"created_by_name"`
	Name           string          `json:"name" firestore:"name"`
	Description    string          `json:"description,omitempty" firestore:"description,omitempty"`
	Participants   []Participant   `json:"participants" firestore:"participants"`
	Messages       []BaseMessage   `json:"messages" firestore:"messages"`
	PinnedMessages []string        `json:"pinnedMessages,omitempty" firestore:"pinned_messages,omitempty"`
	IsArchived     bool            `json:"isArchived" firestore:"is_archived"`
	Notifications  map[string]bool `json:"notifications" firestore:"notifications"`
	CreatedAt      time.Time       `json:"createdAt" firestore:"created_at"`
	UpdatedAt      time.Time       `json:"updatedAt" firestore:"updated_at"`
	ReadStatus     map[string]bool `json:"readStatus" firestore:"read_status"`
	GroupSettings  GroupSettings   `json:"groupSettings" firestore:"group_settings"`
}

type GroupSettings struct {
	AllowFileSharing  bool `json:"allowFileSharing" firestore:"allow_file_sharing"`
	AllowPinning      bool `json:"allowPinning" firestore:"allow_pinning"`
	AllowReactions    bool `json:"allowReactions" firestore:"allow_reactions"`
	AllowReplies      bool `json:"allowReplies" firestore:"allow_replies"`
	MuteNotifications bool `json:"muteNotifications" firestore:"mute_notifications"`
	OnlyAdminsCanPost bool `json:"onlyAdminsCanPost" firestore:"only_admins_can_post"`
}
