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
	"google.golang.org/api/iterator"
)

type UserSchoolExperienceService struct {
	firestoreClient *firestore.Client
	userService     *UserService
}

func NewUserSchoolExperienceService(client *firestore.Client, userService *UserService) *UserSchoolExperienceService {
	return &UserSchoolExperienceService{
		firestoreClient: client,
		userService:     userService,
	}
}

func (s *UserSchoolExperienceService) CreateSchoolExperience(ctx context.Context, uid string) (*models.UserSchoolExperience, error) {

	docRef := s.firestoreClient.Collection("user_school_experiences").NewDoc()

	var experience = &models.UserSchoolExperience{
		UID:          uid,
		Universities: []models.University{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	experienceFirestore := mappers.MapUserSchoolExperienceFromGoToFirestore(experience)
	_, err := docRef.Create(ctx, experienceFirestore)
	if err != nil {
		return nil, fmt.Errorf("failed to create school experience: %w", err)
	}

	return experience, nil
}

func (s *UserSchoolExperienceService) GetSchoolExperience(ctx context.Context, uid string) (*models.UserSchoolExperience, error) {
	query := s.firestoreClient.Collection("user_school_experiences").Where("uid", "==", uid).Documents(ctx)
	doc, err := query.Next()
	if err == iterator.Done {
		return nil, fmt.Errorf("no school experience found for user")
	}
	if err != nil {
		return nil, err
	}
	firestoreData := doc.Data()
	experienceGo := mappers.MapUserSchoolExperienceFromFirestoreToGo(firestoreData)
	return experienceGo, nil
}

func (s *UserSchoolExperienceService) UpdateUniversity(ctx context.Context, uid string, universityID string, updateData map[string]interface{}) (*models.UserSchoolExperience, error) {
	var updatedExperience *models.UserSchoolExperience

	err := s.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		query := s.firestoreClient.Collection("user_school_experiences").Where("uid", "==", uid).Limit(1)
		docs, err := tx.Documents(query).GetAll()
		if err != nil {
			return fmt.Errorf("failed to get school experience: %w", err)
		}
		if len(docs) == 0 {
			return errors.New("school experience not found")
		}

		docRef := docs[0].Ref
		docSnap := docs[0]

		experience := mappers.MapUserSchoolExperienceFromFirestoreToGo(docSnap.Data())
		if experience == nil {
			return errors.New("failed to map school experience data")
		}

		updated := false
		for i, university := range experience.Universities {
			if university.ID == universityID {
				updatedUniversity := mappers.MapUniversityFromFrontendToGo(updateData)
				updatedUniversity.ID = universityID
				experience.Universities[i] = updatedUniversity
				updated = true
				break
			}
		}

		if !updated {
			return errors.New("university not found")
		}

		experience.UpdatedAt = time.Now()

		firestoreData := mappers.MapUserSchoolExperienceFromGoToFirestore(experience)
		if firestoreData == nil {
			return errors.New("failed to map experience to firestore format")
		}

		updatedExperience = experience

		return tx.Set(docRef, firestoreData, firestore.MergeAll)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to update university: %w", err)
	}

	return updatedExperience, nil
}

func (s *UserSchoolExperienceService) AddUniversity(ctx context.Context, uid string, universityData map[string]interface{}) (*models.UserSchoolExperience, error) {
	var updatedExperience *models.UserSchoolExperience

	err := s.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		query := s.firestoreClient.Collection("user_school_experiences").Where("uid", "==", uid).Limit(1)
		docs, err := tx.Documents(query).GetAll()
		if err != nil {
			return fmt.Errorf("failed to get school experience: %w", err)
		}
		if len(docs) == 0 {
			newExperience, err := s.CreateSchoolExperience(ctx, uid)
			if err != nil {
				return fmt.Errorf("failed to create school experience: %w", err)
			}
			updatedExperience = newExperience
			return nil
		}

		docRef := docs[0].Ref
		docSnap := docs[0]

		experience := mappers.MapUserSchoolExperienceFromFirestoreToGo(docSnap.Data())
		if experience == nil {
			return errors.New("failed to map school experience data")
		}

		university := mappers.MapUniversityFromFrontendToGo(universityData)
		university.ID = uuid.New().String()
		experience.Universities = append(experience.Universities, university)
		experience.UpdatedAt = time.Now()

		firestoreData := mappers.MapUserSchoolExperienceFromGoToFirestore(experience)
		if firestoreData == nil {
			return errors.New("failed to map experience to firestore format")
		}

		updatedExperience = experience

		return tx.Set(docRef, firestoreData, firestore.MergeAll)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to add university: %w", err)
	}

	return updatedExperience, nil
}

func (s *UserSchoolExperienceService) DeleteUniversity(ctx context.Context, uid string, universityID string) error {
	experience, err := s.GetSchoolExperience(ctx, uid)
	if err != nil {
		return err
	}

	newUniversities := []models.University{}
	for _, university := range experience.Universities {
		if university.ID != universityID {
			newUniversities = append(newUniversities, university)
		}
	}

	if len(newUniversities) == len(experience.Universities) {
		return errors.New("university not found")
	}

	experience.Universities = newUniversities
	experience.UpdatedAt = time.Now()

	return s.updateExperience(ctx, experience)
}

func (s *UserSchoolExperienceService) AddListOfUniversities(ctx context.Context, uid string, universitiesData []map[string]interface{}) (*models.UserSchoolExperience, error) {
	var updatedExperience *models.UserSchoolExperience

	err := s.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		query := s.firestoreClient.Collection("user_school_experiences").Where("uid", "==", uid).Limit(1)
		docs, err := tx.Documents(query).GetAll()
		if err != nil {
			return fmt.Errorf("failed to get school experience: %w", err)
		}

		if len(docs) == 0 {
			newExperience, err := s.CreateSchoolExperience(ctx, uid)
			if err != nil {
				return fmt.Errorf("failed to create school experience: %w", err)
			}
			updatedExperience = newExperience
			return nil
		}

		docRef := docs[0].Ref
		docSnap := docs[0]

		experience := mappers.MapUserSchoolExperienceFromFirestoreToGo(docSnap.Data())
		if experience == nil {
			return errors.New("failed to map school experience data")
		}

		for _, uniData := range universitiesData {
			university := mappers.MapUniversityFromFrontendToGo(uniData)
			university.ID = uuid.New().String() // Generate UUID for each university
			experience.Universities = append(experience.Universities, university)
		}

		experience.UpdatedAt = time.Now()

		firestoreData := mappers.MapUserSchoolExperienceFromGoToFirestore(experience)
		if firestoreData == nil {
			return errors.New("failed to map experience to firestore format")
		}

		updatedExperience = experience

		return tx.Set(docRef, firestoreData, firestore.MergeAll)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to add universities: %w", err)
	}

	return updatedExperience, nil
}

// func (s *UserSchoolExperienceService) checkExistingExperience(ctx context.Context, uid string) (bool, error) {
// 	docRef := s.firestoreClient.Collection("user_school_experiences").Doc(uid)
// 	docSnap, err := docRef.Get(ctx)
// 	if err != nil {
// 		if err == iterator.Done {
// 			return false, nil
// 		}
// 		return false, err
// 	}
// 	return docSnap.Exists(), nil
// }

func (s *UserSchoolExperienceService) updateExperience(ctx context.Context, experience *models.UserSchoolExperience) error {
	if experience == nil {
		return errors.New("experience cannot be nil")
	}

	query := s.firestoreClient.Collection("user_school_experiences").Where("uid", "==", experience.UID).Documents(ctx)
	doc, err := query.Next()
	if err == iterator.Done {
		return fmt.Errorf("no school experience found for user")
	}
	if err != nil {
		return err
	}

	firestoreData := mappers.MapUserSchoolExperienceFromGoToFirestore(experience)
	if firestoreData == nil {
		return errors.New("failed to map experience to firestore format")
	}

	_, err = doc.Ref.Set(ctx, firestoreData, firestore.MergeAll)
	return err
}
