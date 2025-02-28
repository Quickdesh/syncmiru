package sync

import (
	"context"
	"fmt"
	"github.com/Quickdesh/SyncMiru/internal/notification"

	"github.com/Quickdesh/SyncMiru/internal/domain"
	"github.com/Quickdesh/SyncMiru/internal/logger"
	"github.com/rs/zerolog"
)

type Service interface {
	// Get etag of sync data.
	// For avoid memory usage, only the etag will be returnedj
	GetSyncDataETag(ctx context.Context, apiKey string) (*string, error)
	// Get sync data and etag
	GetSyncDataAndETag(ctx context.Context, apiKey string) ([]byte, *string, error)
	// Create or replace sync data, returns the new etag.
	SetSyncData(ctx context.Context, apiKey string, data []byte) (*string, error)
	// Replace sync data only if the etag matches,
	// returns the new etag if updated, or nil if not.
	SetSyncDataIfMatch(ctx context.Context, apiKey string, etag string, data []byte) (*string, error)
}

func NewService(log logger.Logger, repo domain.SyncRepo, notificationSvc notification.Service, apiRepo domain.APIRepo) Service {
	return &service{
		log:                 log.With().Str("module", "sync").Logger(),
		repo:                repo,
		notificationService: notificationSvc,
		apiRepo:             apiRepo,
	}
}

type service struct {
	log                 zerolog.Logger
	repo                domain.SyncRepo
	notificationService notification.Service
	apiRepo             domain.APIRepo
}

// Get etag of sync data.
// For avoid memory usage, only the etag will be returned.
func (s service) GetSyncDataETag(ctx context.Context, apiKey string) (*string, error) {
	return s.repo.GetSyncDataETag(ctx, apiKey)
}

// Get sync data and etag
func (s service) GetSyncDataAndETag(ctx context.Context, apiKey string) ([]byte, *string, error) {
	return s.repo.GetSyncDataAndETag(ctx, apiKey)
}

// Create or replace sync data, returns the new etag.
func (s service) SetSyncData(ctx context.Context, apiKey string, data []byte) (*string, error) {
	return s.repo.SetSyncData(ctx, apiKey, data)
}

// Replace sync data only if the etag matches,
// returns the new etag if updated, or nil if not.
func (s service) SetSyncDataIfMatch(ctx context.Context, apiKey string, etag string, data []byte) (*string, error) {
	return s.repo.SetSyncDataIfMatch(ctx, apiKey, etag, data)
}

func (s service) notifySyncStarted(apiKeyName string) {
	s.notificationService.Send(domain.NotificationEventSyncStarted, domain.NotificationPayload{
		Subject: "Sync Initiated",
		Message: fmt.Sprintf("A sync operation between Animiru and your library has been initiated for user **%s**. "+
			"Please wait for the process to complete.", apiKeyName),
	})
}

func (s service) notifySyncSuccess(apiKeyName string) {
	s.notificationService.Send(domain.NotificationEventSyncSuccess, domain.NotificationPayload{
		Subject: "Sync Completed Successfully",
		Message: fmt.Sprintf("The synchronization with your Animiru library has completed successfully for user **%s**.", apiKeyName),
	})
}

func (s service) notifySyncFailed(apiKeyName string, errMsg string) {
	s.notificationService.Send(domain.NotificationEventSyncFailed, domain.NotificationPayload{
		Subject: "Sync Operation Failed",
		Message: fmt.Sprintf("The synchronization with Animiru failed for user **%s**. Error: %s", apiKeyName, errMsg),
	})
}

func (s service) notifySyncError(apiKeyName string, errMsg string) {
	s.notificationService.Send(domain.NotificationEventSyncError, domain.NotificationPayload{
		Subject: "Error During Sync",
		Message: fmt.Sprintf("An error occurred during synchronization with Animiru for user **%s**. Error: %s", apiKeyName, errMsg),
	})
}
