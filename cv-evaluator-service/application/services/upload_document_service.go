package services

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/afrizalsebastian/go-common-modules/logger"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/api"
	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/domain/models/dto"
	"github.com/google/uuid"
)

var (
	ErrFileNotFound = errors.New("file not found")
	ErrSaveFile     = errors.New("file error to save")
	ErrDeleteFile   = errors.New("file failed to delete")
)

type IUploadDocumentService interface {
	SaveUploadedDocument(context.Context, *dto.UploadDocumentRequest) api.WebResponse
}

type uploadDocumentService struct {
	basePath string
}

func NewUploadDocumentService(basePath string) IUploadDocumentService {
	return &uploadDocumentService{
		basePath: basePath,
	}
}

func (u *uploadDocumentService) SaveUploadedDocument(ctx context.Context, req *dto.UploadDocumentRequest) api.WebResponse {
	l := logger.New().WithContext(ctx)

	folderId := uuid.New().String()
	errChan := make(chan error, 2)
	// save cv file

	go func() {
		l.Info("upload cv").Msg()
		req.CvFileHeader.Filename = "cv_file.pdf"
		if err := u.saveToPath(ctx, folderId, req.CvFile, req.CvFileHeader); err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	go func() {
		l.Info("upload report").Msg()
		req.ReportFileHeader.Filename = "report_file.pdf"
		if err := u.saveToPath(ctx, folderId, req.ReportFile, req.ReportFileHeader); err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	l.Info("upload done").Msg()

	err := <-errChan
	if err != nil {
		l.Error("error when save document").Msg()
		return api.CreateWebResponse("Error when save user document", http.StatusInternalServerError, nil, nil)
	}

	resp := &dto.UploadedDocumentResponse{
		FileId: folderId,
	}

	return api.CreateWebResponse("Success", http.StatusOK, resp, nil)

}

func (u *uploadDocumentService) saveToPath(ctx context.Context, folderId string, file multipart.File, header *multipart.FileHeader) error {
	l := logger.New().WithContext(ctx)
	filename := header.Filename

	safeFilename := filepath.Clean(filename)
	if safeFilename == "." || safeFilename == "/" {
		safeFilename = filename
	}

	_, err := os.Stat(filepath.Join(u.basePath, folderId))
	if os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(u.basePath, folderId), 0o755); err != nil {
			l.Error("failed to create directory").Msg()
			return err
		}
	}

	dest := filepath.Join(u.basePath, folderId, safeFilename)

	out, err := os.Create(dest)
	if err != nil {
		return ErrSaveFile
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		_ = os.Remove(dest)
		return ErrSaveFile
	}

	if err := out.Sync(); err != nil {
		return ErrSaveFile
	}

	return nil
}
