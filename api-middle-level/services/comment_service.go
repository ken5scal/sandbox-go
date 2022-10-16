package services

import (
	"github.com/ken5scal/api-go-mid-level/apperrors"
	"github.com/ken5scal/api-go-mid-level/models"
	"github.com/ken5scal/api-go-mid-level/repositories"
)

func (s *MyAppService) PostCommentService(comment models.Comment) (models.Comment, error) {
	newComment, err := repositories.InsertComment(s.db, comment)
	if err != nil {
		err = apperrors.InsertDataFailed.Wrap(err, "fail to record data")

		return models.Comment{}, err
	}

	return newComment, nil
}
