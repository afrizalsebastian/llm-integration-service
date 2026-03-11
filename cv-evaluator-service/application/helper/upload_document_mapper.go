package helper

import (
	"mime/multipart"

	"github.com/afrizalsebastian/llm-integration-service/cv-evaluator-service/domain/models/dto"
)

func MultipartToUploadDocumentRequest(req *MultipartRequest) *dto.UploadDocumentRequest {
	var cvFile multipart.File
	var cvFileHeader *multipart.FileHeader
	var reportFile multipart.File
	var reportFileHeader *multipart.FileHeader

	if cvInfo, ok := req.Files["cv_file"]; ok && cvInfo != nil {
		cvFile = cvInfo.File
		cvFileHeader = cvInfo.FileHeader
	}

	if reportInfo, ok := req.Files["report_file"]; ok && reportInfo != nil {
		reportFile = reportInfo.File
		reportFileHeader = reportInfo.FileHeader
	}

	return &dto.UploadDocumentRequest{
		CvFile:           cvFile,
		CvFileHeader:     cvFileHeader,
		ReportFile:       reportFile,
		ReportFileHeader: reportFileHeader,
	}
}
