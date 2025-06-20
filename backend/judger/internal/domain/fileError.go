package domain

type FileErrorType string

const (
	FileErrorTypeCopyInOpenFile                 FileErrorType = "CopyInOpenFile"
	FileErrorTypeCopyInCreateFile               FileErrorType = "CopyInCreateFile"
	FileErrorTypeCopyInCopyContent              FileErrorType = "CopyInCopyContent"
	FileErrorTypeCopyOutOpen                    FileErrorType = "CopyOutOpen"
	FileErrorTypeCopyOutNotRegularFile          FileErrorType = "CopyOutNotRegularFile"
	FileErrorTypeCopyOutSizeExceeded            FileErrorType = "CopyOutSizeExceeded"
	FileErrorTypeFileErrorTypeCopyOutCreateFile FileErrorType = "CopyOutCreateFile"
	FileErrorTypeCopyOutCopyContent             FileErrorType = "CopyOutCopyContent"
	FileErrorTypeCollectSizeExceeded            FileErrorType = "CollectSizeExceeded"
)

type FileError struct {
	Name    string        `json:"name"`
	Type    FileErrorType `json:"type"`
	Message string        `json:"message"`
}
