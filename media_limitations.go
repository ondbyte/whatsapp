package whatsapp

import "github.com/ondbyte/whatsapp/message"

// mediaLimitations struct represents a media type with its supported MIME types and size limit.
type mediaLimitations[T any] struct {
	_type     string
	mimeTypes map[string]string
	sizeLimit int64
	wrap      func(id string) *T
}

// Constants for different media types
var (
	// audioLimitations
	audioLimitations = mediaLimitations[message.Audio]{
		_type: "audio",
		mimeTypes: map[string]string{
			".aac":  "audio/aac",
			".mp4":  "audio/mp4",
			".mpeg": "audio/mpeg",
			".amr":  "audio/amr",
			".ogg":  "audio/ogg",
		},
		sizeLimit: 16 * 1024 * 1024, // 16MB
		wrap:      message.NewAudio,
	}

	// documentLimitations
	documentLimitations = mediaLimitations[message.Document]{
		_type: "document",
		mimeTypes: map[string]string{
			".txt":  "text/plain",
			".pdf":  "application/pdf",
			".ppt":  "application/vnd.ms-powerpoint",
			".doc":  "application/msword",
			".xls":  "application/vnd.ms-excel",
			".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
			".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		},
		sizeLimit: 100 * 1024 * 1024, // 100MB
		wrap:      message.NewDocument,
	}

	// imageLimitations
	imageLimitations = mediaLimitations[message.Image]{
		_type: "image",
		mimeTypes: map[string]string{
			".jpeg": "image/jpeg",
			".jpg":  "image/jpeg",
			".png":  "image/png",
			".webp": "image/webp",
		},
		sizeLimit: 5 * 1024 * 1024, // 5MB
		wrap:      message.NewImage,
	}

	// videoLimitations
	videoLimitations = mediaLimitations[message.Video]{
		_type: "video",
		mimeTypes: map[string]string{
			".mp4": "video/mp4",
			".3gp": "video/3gp",
		},
		sizeLimit: 16 * 1024 * 1024, // 16MB
		wrap:      message.NewVideo,
	}

	// stickerLimitations
	stickerLimitations = mediaLimitations[message.Sticker]{
		_type: "sticker",
		mimeTypes: map[string]string{
			".webp": "image/webp",
		},
		sizeLimit: 500 * 1024, // 500KB
		wrap:      message.NewSticker,
	}
)
