package repository

import (
	"context"

	"github.com/afrizalsebastian/go-common-modules/logger"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/bootstrap"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/domain/models/dao"
	"gorm.io/gorm"
)

type ICvEvaluatorJobRepository interface {
	CreateJobItem(ctx context.Context, job *dao.CvEvaluatorJob) error
	GetByJobId(ctx context.Context, jobId string) (*dao.CvEvaluatorJob, error)
	UpdateJobByJobId(ctx context.Context, jobId string, job *dao.CvEvaluatorJob) error
}

type cvEvaluatorJobRepository struct {
	db *gorm.DB
}

func NewCvEvaluatorJobRepository(app *bootstrap.Application) ICvEvaluatorJobRepository {
	return &cvEvaluatorJobRepository{
		db: app.DB,
	}
}

func (c *cvEvaluatorJobRepository) CreateJobItem(ctx context.Context, job *dao.CvEvaluatorJob) error {
	l := logger.New().WithContext(ctx)

	if err := c.db.WithContext(ctx).Create(&job).Error; err != nil {
		l.Error("failed to create job").Msg()
		return err
	}
	return nil
}

func (c *cvEvaluatorJobRepository) GetByJobId(ctx context.Context, jobId string) (*dao.CvEvaluatorJob, error) {
	l := logger.New().WithContext(ctx)

	var jobItem dao.CvEvaluatorJob
	if err := c.db.WithContext(ctx).Model(&dao.CvEvaluatorJob{}).Where("job_id = ?", jobId).First(&jobItem).Error; err != nil {
		l.Error("failed to get job by job id").Msg()
		return nil, err
	}

	return &jobItem, nil
}

func (c *cvEvaluatorJobRepository) UpdateJobByJobId(ctx context.Context, jobId string, job *dao.CvEvaluatorJob) error {
	l := logger.New().WithContext(ctx)

	if err := c.db.WithContext(ctx).Where("job_id = ?", jobId).Save(job).Error; err != nil {
		l.Error("failed to update job").Msg()
		return err
	}
	return nil
}
