package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ticketing-system/internal/config"
	"ticketing-system/internal/model"

	"github.com/google/uuid"
)

// ======================
// DTO
// ======================

type AttachmentResponse struct {
	ID          uint   `json:"id"`
	TicketID    uint   `json:"ticket_id"`
	UserID      uint   `json:"user_id"`
	FileName    string `json:"file_name"`
	FilePath    string `json:"file_path"` // path relatif untuk preview
	FileSize    int64  `json:"file_size"`
	ContentType string `json:"content_type"`
	CreatedAt   string `json:"created_at"`
}

var allowedTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/gif":       true,
	"image/webp":      true,
	"video/mp4":       true,
	"video/quicktime": true,
	"application/pdf": true,
}

const maxFileSize = 50 * 1024 * 1024 // 50MB

// ======================
// UPLOAD ATTACHMENT
// ======================

func UploadAttachment(ticketID, userID uint, file *multipart.FileHeader) (*AttachmentResponse, error) {
	var ticket model.Ticket
	if err := config.DB.First(&ticket, ticketID).Error; err != nil {
		return nil, errors.New("ticket not found")
	}

	if file.Size > maxFileSize {
		return nil, errors.New("file too large, max 50MB")
	}

	contentType := file.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return nil, errors.New("file type not allowed")
	}

	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, errors.New("failed to create upload directory")
	}

	ext := filepath.Ext(file.Filename)
	safeExt := strings.ToLower(ext)
	uniqueName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), safeExt)
	filePath := filepath.Join(uploadDir, uniqueName)

	src, err := file.Open()
	if err != nil {
		return nil, errors.New("failed to open file")
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, errors.New("failed to save file")
	}
	defer dst.Close()

	buf := make([]byte, 32*1024)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			if _, werr := dst.Write(buf[:n]); werr != nil {
				return nil, errors.New("failed to write file")
			}
		}
		if err != nil {
			break
		}
	}

	attachment := model.Attachment{
		TicketID:    ticketID,
		UserID:      userID,
		FileName:    file.Filename,
		FilePath:    filePath,
		FileSize:    file.Size,
		ContentType: contentType,
	}

	if err := config.DB.Create(&attachment).Error; err != nil {
		os.Remove(filePath)
		return nil, errors.New("failed to save attachment record")
	}

	return &AttachmentResponse{
		ID:          attachment.ID,
		TicketID:    attachment.TicketID,
		UserID:      attachment.UserID,
		FileName:    attachment.FileName,
		FilePath:    "/" + filePath, // untuk static serve
		FileSize:    attachment.FileSize,
		ContentType: attachment.ContentType,
		CreatedAt:   attachment.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// ======================
// GET ATTACHMENTS
// ======================

func GetAttachments(ticketID uint) ([]AttachmentResponse, error) {
	var attachments []model.Attachment
	err := config.DB.Where("ticket_id = ?", ticketID).
		Order("created_at DESC").
		Find(&attachments).Error
	if err != nil {
		return nil, err
	}

	result := make([]AttachmentResponse, 0, len(attachments))
	for _, a := range attachments {
		result = append(result, AttachmentResponse{
			ID:          a.ID,
			TicketID:    a.TicketID,
			UserID:      a.UserID,
			FileName:    a.FileName,
			FilePath:    "/" + a.FilePath,
			FileSize:    a.FileSize,
			ContentType: a.ContentType,
			CreatedAt:   a.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return result, nil
}

// ======================
// GET ATTACHMENT PATH
// ======================

func GetAttachmentPath(attachmentID uint) (string, string, error) {
	var attachment model.Attachment
	if err := config.DB.First(&attachment, attachmentID).Error; err != nil {
		return "", "", errors.New("attachment not found")
	}
	return attachment.FilePath, attachment.FileName, nil
}

// ======================
// DELETE ATTACHMENT
// ======================

func DeleteAttachment(attachmentID, userID uint, role string) error {
	var attachment model.Attachment
	if err := config.DB.First(&attachment, attachmentID).Error; err != nil {
		return errors.New("attachment not found")
	}

	if role != "admin" && attachment.UserID != userID {
		return errors.New("not authorized to delete this attachment")
	}

	os.Remove(attachment.FilePath)
	return config.DB.Delete(&attachment).Error
}
