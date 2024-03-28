package media

import (
	"context"
	"time"
	"wp-media-core/database"
)

func saveMedia(source string, media_type string) (*Media, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var media Media
	err := database.DB.QueryRowContext(ctx, "INSERT wp_media(source, type) VALUES($1,$2) RETURNING id, source, type, thumbnail, created_at, updated_at", source, media_type).Scan(&media.ID, &media.Source, &media.Type, &media.Thumbnail, media.CreatedAt, media.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &media, nil
}
