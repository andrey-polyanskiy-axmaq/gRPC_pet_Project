package feedback

import (
	"NewSGNgRPC/internal/domain/model"
	"NewSGNgRPC/internal/services/storage"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

type Feedback struct {
	log      slog.Logger
	saver    Saver
	provider Provider
}

type Saver interface {
	SaveSugComText(ctx context.Context, txt string, IDhash []byte) (err error)
}

type Provider interface {
	ProvideSugComText(ctx context.Context, ID int64) (model.Text, error)
	ProvideAnswerText(ctx context.Context, ID int64) (model.Text, error)
	ProvideApp(ctx context.Context, ID int64) (model.App, error)
}

// New возращает новою структуры feedback-сервиса
func New(log slog.Logger, saver Saver, provider Provider) *Feedback {

	return &Feedback{
		log,
		saver,
		provider,
	}
}

var (
	ErrInvalidCredentials    = errors.New("Invalid Credentials")
	ErrIdempotentCredentials = errors.New("Such Credentials already exist")
)

func (f *Feedback) Suggestion(ctx context.Context, txt string, appID int64, sugcomID int64) ([]byte, error) {
	const op = "feedback.CreateSuggestionText"

	// TODO: в БД берем возможный ID

	log := f.log.With(slog.String("current operation", op), slog.Int64("New Suggestion ID", sugcomID))
	log.Info("Creating New Suggestion...")

	IDHash, err := bcrypt.GenerateFromPassword([]byte(string(sugcomID)), bcrypt.MinCost)
	if err != nil {
		log.Error("failed to generate text ID hash", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = f.saver.SaveSugComText(ctx, txt, IDHash)
	if err != nil {

		if errors.Is(err, storage.ErrTextExist) {
			f.log.Warn("Text already exists")
			return nil, fmt.Errorf("%s: %w", op, ErrIdempotentCredentials)
		}

		log.Error("failed to save suggestion text", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("Suggestion is saved")
	return IDHash, nil

}

func (f *Feedback) Complaint(ctx context.Context, txt string, appID int64, sugcomID int64) ([]byte, error) {
	const op = "feedback.CreateComplaintText"

	// TODO: в БД берем возможный ID

	log := f.log.With(slog.String("current operation", op), slog.Int64("New Complaint ID", sugcomID))
	log.Info("Creating New Complaint...")

	IDHash, err := bcrypt.GenerateFromPassword([]byte(string(sugcomID)), bcrypt.MinCost)
	if err != nil {
		log.Error("failed to generate text ID hash", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = f.saver.SaveSugComText(ctx, txt, IDHash)
	if err != nil {

		if errors.Is(err, storage.ErrTextExist) {
			f.log.Warn("Text already exists")
			return nil, fmt.Errorf("%s: %w", op, ErrIdempotentCredentials)
		}

		log.Error("failed to save suggestion text", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("Complaint is saved")

	return IDHash, nil
}

func (f *Feedback) SugComAnswer(ctx context.Context, txt string, appID int64, IDHash []byte) (string, error) {
	const op = "feedback.CreateSugComAnswer"

	log := f.log.With(slog.String("current operation", op))

	//декодирование HashID
	storagelen := 10 // костыли
	var IntID int64
	for i := 0; i < storagelen; i++ {
		if err := bcrypt.CompareHashAndPassword(IDHash, []byte((string(i)))); err == nil {
			IntID = int64(i)
			break
		}
	}

	log.Info("user logged in successfully")

	//Обращение к БД, получение текста ответа
	txtAnswer, err := f.provider.ProvideAnswerText(ctx, IntID)

	if err != nil {

		if errors.Is(err, storage.ErrIDNotFound) {
			f.log.Warn("ID not found")
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to provide answer text", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("Text is provided")

	return txtAnswer.TXT, nil
}
