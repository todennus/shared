package enumdef

import (
	"fmt"
	"strings"

	"github.com/todennus/proto/gen/service/dto"
)

type PolicySource string

const (
	PolicySourceUserAvatar = "user.avatar"
)

func FilePolicyToken(source PolicySource, token string) string {
	return string(source) + "/" + token
}

func ParseFilePolicyToken(token string) (string, string, bool) {
	return strings.Cut(token, "/")
}

type TemporaryFileCommand string

const (
	// TemporaryFileCommandSaveAsImage represents the action of saving the file
	// as an image to persistent storage and deleting the temporary file.
	TemporaryFileCommandDelete TemporaryFileCommand = "delete"

	// TemporaryFileCommandSaveAsImage represents the action of saving the file
	// as an image to persistent storage and deleting the temporary file.
	TemporaryFileCommandSaveAsImage TemporaryFileCommand = "save_as_image"

	// TemporaryFileCommandChangeImageType represents the action of changing the
	// image file to another type. This command is only valid if the file is an
	// image.
	TemporaryFileCommandChangeImageType TemporaryFileCommand = "change_image_type"

	// TemporaryFileCommandImageMetadata represents the action of retrieving the
	// image metadata. This command is only valid if the file is an image.
	TemporaryFileCommandImageMetadata TemporaryFileCommand = "image_metadata"
)

func TemporaryFileCommandFromGRPC(req dto.FileCommandTemporaryFile) TemporaryFileCommand {
	switch req {
	case dto.FileCommandTemporaryFile_Delete:
		return TemporaryFileCommandDelete
	case dto.FileCommandTemporaryFile_SaveAsImage:
		return TemporaryFileCommandSaveAsImage
	case dto.FileCommandTemporaryFile_ImageMetadata:
		return TemporaryFileCommandImageMetadata
	case dto.FileCommandTemporaryFile_ChangeImageType:
		return TemporaryFileCommandChangeImageType
	default:
		panic(fmt.Sprintf("invalid requirement %d", req))
	}
}

func TemporaryFileCommandToGRPC(req TemporaryFileCommand) dto.FileCommandTemporaryFile {
	switch req {
	case TemporaryFileCommandDelete:
		return dto.FileCommandTemporaryFile_Delete
	case TemporaryFileCommandSaveAsImage:
		return dto.FileCommandTemporaryFile_SaveAsImage
	case TemporaryFileCommandImageMetadata:
		return dto.FileCommandTemporaryFile_ImageMetadata
	case TemporaryFileCommandChangeImageType:
		return dto.FileCommandTemporaryFile_ChangeImageType
	default:
		panic(fmt.Sprintf("invalid requirement %d", req))
	}
}
