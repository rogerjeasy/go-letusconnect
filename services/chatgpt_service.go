package services

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/rogerjeasy/go-letusconnect/config"
	openai "github.com/sashabaranov/go-openai"
)

type ChatHistory struct {
	ID        string    `json:"id" firestore:"id"`
	UserID    string    `json:"userId" firestore:"userId"`
	Message   string    `json:"message" firestore:"message"`
	Response  string    `json:"response" firestore:"response"`
	CreatedAt time.Time `json:"createdAt" firestore:"createdAt"`
}

type ChatGPTService struct {
	client          *openai.Client
	firestoreClient *firestore.Client
}

func NewChatGPTService(firestoreClient *firestore.Client) *ChatGPTService {
	return &ChatGPTService{
		client:          openai.NewClient(config.OpenAIKey),
		firestoreClient: firestoreClient,
	}
}

func (s *ChatGPTService) GenerateResponse(ctx context.Context, prompt string, userID string) (*ChatHistory, error) {
	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
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

	chatHistory := &ChatHistory{
		UserID:    userID,
		Message:   prompt,
		Response:  resp.Choices[0].Message.Content,
		CreatedAt: time.Now(),
	}

	// Save to Firestore
	doc, _, err := s.firestoreClient.Collection("chatgpt_conversations").Add(ctx, chatHistory)
	if err != nil {
		return nil, fmt.Errorf("failed to save chat history: %v", err)
	}
	chatHistory.ID = doc.ID

	return chatHistory, nil
}

func (s *ChatGPTService) GetChatHistory(ctx context.Context, userID string) ([]ChatHistory, error) {
	var histories []ChatHistory

	iter := s.firestoreClient.Collection("chatgpt_conversations").
		Where("userId", "==", userID).
		OrderBy("createdAt", firestore.Desc).
		Limit(50).
		Documents(ctx)

	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var history ChatHistory
		if err := doc.DataTo(&history); err != nil {
			return nil, fmt.Errorf("failed to parse chat history: %v", err)
		}
		history.ID = doc.Ref.ID
		histories = append(histories, history)
	}

	return histories, nil
}

func (s *ChatGPTService) DeleteChatHistory(ctx context.Context, historyID string, userID string) error {
	doc := s.firestoreClient.Collection("chatgpt_conversations").Doc(historyID)

	// Verify ownership
	snapshot, err := doc.Get(ctx)
	if err != nil {
		return fmt.Errorf("chat history not found")
	}

	var history ChatHistory
	if err := snapshot.DataTo(&history); err != nil {
		return fmt.Errorf("failed to parse chat history")
	}

	if history.UserID != userID {
		return fmt.Errorf("unauthorized access")
	}

	// if err := doc.Delete(ctx); err != nil {
	//     return fmt.Errorf("failed to delete chat history: %v", err)
	// }

	return nil
}
