package services

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/rogerjeasy/go-letusconnect/config"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	openai "github.com/sashabaranov/go-openai"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatGPTService struct {
	client          *openai.Client
	firestoreClient *firestore.Client
	pdfService      *PDFService
}

func NewChatGPTService(firestoreClient *firestore.Client, pdfService *PDFService) *ChatGPTService {
	return &ChatGPTService{
		client:          openai.NewClient(config.OpenAIKey),
		firestoreClient: firestoreClient,
		pdfService:      pdfService,
	}
}

func (s *ChatGPTService) GenerateResponse(ctx context.Context, prompt string, userID string, conversationID string) (*models.Conversation, error) {
	// Generate response from ChatGPT
	pdfContext := s.pdfService.GetContext()

	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: fmt.Sprintf("You are an AI assistant helping users navigate and understand our website. Here's the context about our website:\n\n%s\n\nPlease use this information to provide accurate and helpful responses. And in your response, do not put URL path in parentheses, and should always start with the forward slash symbole /", pdfContext),
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
			MaxTokens:   1000,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %v", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response generated")
	}

	now := time.Now()
	newMessage := models.MessageConversation{
		ID:        uuid.New().String(),
		CreatedAt: now,
		Message:   prompt,
		Response:  resp.Choices[0].Message.Content,
		Role:      "user",
	}

	var conversation *models.Conversation
	var docRef *firestore.DocumentRef

	// If conversationID is empty, create new conversation
	if conversationID == "" {
		conversation = &models.Conversation{
			ID:        uuid.New().String(),
			UserID:    userID,
			Title:     prompt[:min(30, len(prompt))] + "...", // Create title from first 30 chars
			CreatedAt: now,
			UpdatedAt: now,
			Messages:  []models.MessageConversation{newMessage},
		}

		// Convert to Firestore format and save
		firestoreData := mappers.MapConversationGoToFirestore(*conversation)
		docRef = s.firestoreClient.Collection("chatgpt_conversations").Doc(conversation.ID)

		_, err = docRef.Set(ctx, firestoreData)
		if err != nil {
			return nil, fmt.Errorf("failed to create conversation: %v", err)
		}
	} else {
		// Get existing conversation
		docRef = s.firestoreClient.Collection("chatgpt_conversations").Doc(conversationID)
		doc, err := docRef.Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get conversation: %v", err)
		}

		// Convert Firestore data to Go struct
		firestoreData := doc.Data()
		existingConversation := mappers.MapConversationFirestoreToGo(firestoreData)

		// Update conversation with new message
		existingConversation.Messages = append(existingConversation.Messages, newMessage)
		existingConversation.UpdatedAt = now

		// Convert back to Firestore format and update
		updatedData := mappers.MapConversationGoToFirestore(existingConversation)
		_, err = docRef.Set(ctx, updatedData)
		if err != nil {
			return nil, fmt.Errorf("failed to update conversation: %v", err)
		}

		conversation = &existingConversation
	}

	return conversation, nil
}

func (s *ChatGPTService) GetUserConversations(ctx context.Context, userID string) ([]models.Conversation, error) {
	iter := s.firestoreClient.Collection("chatgpt_conversations").
		Where("user_id", "==", userID).
		Documents(ctx)

	var conversations []models.Conversation
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		conversation := mappers.MapConversationFirestoreToGo(doc.Data())
		conversations = append(conversations, conversation)
	}

	return conversations, nil
}

func (s *ChatGPTService) GetConversation(ctx context.Context, conversationID string) (*models.Conversation, error) {
	doc, err := s.firestoreClient.Collection("chatgpt_conversations").Doc(conversationID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %v", err)
	}

	conversation := mappers.MapConversationFirestoreToGo(doc.Data())
	return &conversation, nil
}

func (s *ChatGPTService) DeleteConversation(ctx context.Context, conversationID string, userID string) error {
	docRef := s.firestoreClient.Collection("chatgpt_conversations").Doc(conversationID)
	doc, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return fmt.Errorf("conversation not found")
		}
		return fmt.Errorf("failed to get conversation: %v", err)
	}

	conversation := mappers.MapConversationFirestoreToGo(doc.Data())

	if conversation.UserID != userID {
		return fmt.Errorf("unauthorized: user does not own this conversation")
	}

	_, err = docRef.Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %v", err)
	}

	return nil
}
