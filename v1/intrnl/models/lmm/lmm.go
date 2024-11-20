package lmm

type ContentRequest struct {
	ContentID string `json:"content_id" validate:"omitempty"`
}

type ContentResponse struct {
	S3URL string `json:"s3_url" validate:"omitempty"`
}

type LMMContent struct {
	ContentID   string `json:"content_id" validate:"omitempty"`
	ContentType string `json:"content_type" validate:"omitempty"`
}
