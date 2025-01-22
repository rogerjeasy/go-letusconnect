package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
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
		return s.CreateSchoolExperience(ctx, uid)
	}
	if err != nil {
		return nil, err
	}

	firestoreData := doc.Data()
	experienceGo := mappers.MapUserSchoolExperienceFromFirestoreToGo(firestoreData)
	return experienceGo, nil
}

func (s *UserSchoolExperienceService) UpdateUniversity(ctx context.Context, uid string, universityID string, updateData map[string]interface{}) (*models.UserSchoolExperience, error) {
	err := s.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		docRef := s.firestoreClient.Collection("user_school_experiences").Doc(uid)
		docSnap, err := tx.Get(docRef)
		if err != nil {
			return err
		}

		experience := mappers.MapUserSchoolExperienceFromFirestoreToGo(docSnap.Data())
		if experience == nil {
			return errors.New("failed to map school experience data")
		}

		updated := false
		for i, university := range experience.Universities {
			if university.ID == universityID {
				experience.Universities[i] = mappers.MapUniversityFromFrontendToGo(updateData)
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

		return tx.Set(docRef, firestoreData, firestore.MergeAll)
	})

	if err != nil {
		return nil, err
	}

	return s.GetSchoolExperience(ctx, uid)
}

func (s *UserSchoolExperienceService) AddUniversity(ctx context.Context, uid string, universityData map[string]interface{}) (*models.UserSchoolExperience, error) {
	err := s.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		docRef := s.firestoreClient.Collection("user_school_experiences").Doc(uid)
		docSnap, err := tx.Get(docRef)
		if err != nil {
			return err
		}

		experience := mappers.MapUserSchoolExperienceFromFirestoreToGo(docSnap.Data())
		if experience == nil {
			return errors.New("failed to map school experience data")
		}

		university := mappers.MapUniversityFromFrontendToGo(universityData)
		experience.Universities = append(experience.Universities, university)
		experience.UpdatedAt = time.Now()

		firestoreData := mappers.MapUserSchoolExperienceFromGoToFirestore(experience)
		if firestoreData == nil {
			return errors.New("failed to map experience to firestore format")
		}

		return tx.Set(docRef, firestoreData, firestore.MergeAll)
	})

	if err != nil {
		return nil, err
	}

	return s.GetSchoolExperience(ctx, uid)
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
	err := s.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		docRef := s.firestoreClient.Collection("user_school_experiences").Doc(uid)
		docSnap, err := tx.Get(docRef)
		if err != nil {
			return err
		}

		experience := mappers.MapUserSchoolExperienceFromFirestoreToGo(docSnap.Data())
		if experience == nil {
			return errors.New("failed to map school experience data")
		}

		for _, uniData := range universitiesData {
			university := mappers.MapUniversityFromFrontendToGo(uniData)
			experience.Universities = append(experience.Universities, university)
		}

		experience.UpdatedAt = time.Now()

		firestoreData := mappers.MapUserSchoolExperienceFromGoToFirestore(experience)
		return tx.Set(docRef, firestoreData, firestore.MergeAll)
	})

	if err != nil {
		return nil, err
	}

	return s.GetSchoolExperience(ctx, uid)
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

	firestoreData := mappers.MapUserSchoolExperienceFromGoToFirestore(experience)
	if firestoreData == nil {
		return errors.New("failed to map experience to firestore format")
	}

	docRef := s.firestoreClient.Collection("user_school_experiences").Doc(experience.UID)
	_, err := docRef.Set(ctx, firestoreData, firestore.MergeAll)
	return err
}
