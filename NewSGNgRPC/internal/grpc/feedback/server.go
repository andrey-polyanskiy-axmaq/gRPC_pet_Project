package feedback

import (
	"context"
	ssov1 "github.com/andrey-polyanskiy-axmaq/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const empty_string = ""

type Feedback interface {
	Sugggestion(ctx context.Context, TXT string) (HashID []byte, err error)
	Complaint(ctx context.Context, TXT string) (HashID []byte, err error)
	SugComAnswer(ctx context.Context, HashID []byte) (TXT string, err error)
}

type serverAPI struct {
	ssov1.UnimplementedFeedbackServer
	feedback Feedback
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterFeedbackServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Suggestion(ctx context.Context, req *ssov1.SuggestionRequest) (*ssov1.SuggestionResponse, error) {
	//Валидация входящих данных
	txt := req.GetSuggestionTXT()
	if txt == empty_string {
		//TODO: обработка ошибки
		return nil, status.Error(codes.InvalidArgument, "Empty suggestion")
	}

	//Cервисный слой
	id, err := s.feedback.Sugggestion(ctx, txt)

	if err != nil {
		//TODO: обработка ошибки
		return nil, status.Error(codes.Internal, "Internal error")
	}
	return &ssov1.SuggestionResponse{
		HashID: id,
	}, nil
}

func (s *serverAPI) Complaint(ctx context.Context, req *ssov1.ComplaintRequest) (*ssov1.ComplaintResponse, error) {
	//Валидация входящих данных
	txt := req.GetComplaintTXT()
	if txt == empty_string {
		//TODO: обработка ошибки
		return nil, status.Error(codes.InvalidArgument, "Empty complaint")
	}
	//Cервисный слой
	id, err := s.feedback.Complaint(ctx, txt)
	if err != nil {
		//TODO: обработка ошибки
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.ComplaintResponse{
		HashID: id,
	}, nil
}

func (s *serverAPI) SugComAnswer(ctx context.Context, req *ssov1.AnswerRequest) (*ssov1.AnswerResponse, error) {
	//Валидация входящих данных
	id := req.GetHashID()
	// TODO: Декодировка
	var idint int = 0
	if idint == 0 { // TODO: условие отсутствия id в БД
		//TODO: обработка ошибки
		return nil, status.Error(codes.InvalidArgument, "ID doesn't exist")
	}
	//Cервисный слой
	txt, err := s.feedback.SugComAnswer(ctx, id)
	if err != nil {
		//TODO: обработка ошибки
		return nil, status.Error(codes.Internal, "Internal error")
	}
	return &ssov1.AnswerResponse{
		AnswerTXT: txt,
	}, nil
}
