package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"gorm.io/gorm"
)

func TestItemMediaRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewItemMediaRepository(db)

	t.Run("Success - Create media", func(t *testing.T) {
		itemID := uuid.New()
		media := &domain.ItemMedia{
			ItemID:       itemID,
			MediaType:    domain.MediaTypeImage,
			URL:          "http://example.com/image.webp",
			DisplayOrder: 0,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "item_media"`).
			WithArgs(
				sqlmock.AnyArg(), // item_id
				sqlmock.AnyArg(), // media_type
				sqlmock.AnyArg(), // url
				sqlmock.AnyArg(), // thumbnail_url
				sqlmock.AnyArg(), // display_order
				sqlmock.AnyArg(), // created_at
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := repo.Create(media)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestItemMediaRepository_FindByItemID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewItemMediaRepository(db)

	t.Run("Success - Find media by item ID", func(t *testing.T) {
		itemID := uuid.New()
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "item_id", "media_type", "url", "thumbnail_url", "display_order", "created_at",
		}).
			AddRow(1, itemID, "image", "http://example.com/image1.webp", "http://example.com/thumb1.webp", 0, now).
			AddRow(2, itemID, "image", "http://example.com/image2.webp", "http://example.com/thumb2.webp", 1, now)

		mock.ExpectQuery(`SELECT \* FROM "item_media" WHERE item_id = \$1 ORDER BY display_order ASC, created_at ASC`).
			WithArgs(itemID).
			WillReturnRows(rows)

		mediaList, err := repo.FindByItemID(itemID)

		assert.NoError(t, err)
		assert.Len(t, mediaList, 2)
		assert.Equal(t, int64(1), mediaList[0].ID)
		assert.Equal(t, 0, mediaList[0].DisplayOrder)
		assert.Equal(t, int64(2), mediaList[1].ID)
		assert.Equal(t, 1, mediaList[1].DisplayOrder)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success - No media found", func(t *testing.T) {
		itemID := uuid.New()

		rows := sqlmock.NewRows([]string{
			"id", "item_id", "media_type", "url", "thumbnail_url", "display_order", "created_at",
		})

		mock.ExpectQuery(`SELECT \* FROM "item_media" WHERE item_id = \$1 ORDER BY display_order ASC, created_at ASC`).
			WithArgs(itemID).
			WillReturnRows(rows)

		mediaList, err := repo.FindByItemID(itemID)

		assert.NoError(t, err)
		assert.Len(t, mediaList, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestItemMediaRepository_FindByID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewItemMediaRepository(db)

	t.Run("Success - Media found", func(t *testing.T) {
		mediaID := int64(1)
		itemID := uuid.New()
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "item_id", "media_type", "url", "thumbnail_url", "display_order", "created_at",
		}).AddRow(mediaID, itemID, "image", "http://example.com/image.webp", "http://example.com/thumb.webp", 0, now)

		mock.ExpectQuery(`SELECT \* FROM "item_media" WHERE "item_media"\."id" = \$1`).
			WithArgs(mediaID).
			WillReturnRows(rows)

		media, err := repo.FindByID(mediaID)

		assert.NoError(t, err)
		assert.NotNil(t, media)
		assert.Equal(t, mediaID, media.ID)
		assert.Equal(t, itemID, media.ItemID)
		assert.Equal(t, domain.MediaTypeImage, media.MediaType)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success - Media not found", func(t *testing.T) {
		mediaID := int64(999)

		mock.ExpectQuery(`SELECT \* FROM "item_media" WHERE "item_media"\."id" = \$1`).
			WithArgs(mediaID).
			WillReturnError(gorm.ErrRecordNotFound)

		media, err := repo.FindByID(mediaID)

		assert.NoError(t, err)
		assert.Nil(t, media)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestItemMediaRepository_FindByIDAndItemID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewItemMediaRepository(db)

	t.Run("Success - Media found", func(t *testing.T) {
		mediaID := int64(1)
		itemID := uuid.New()
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "item_id", "media_type", "url", "thumbnail_url", "display_order", "created_at",
		}).AddRow(mediaID, itemID, "image", "http://example.com/image.webp", "http://example.com/thumb.webp", 0, now)

		mock.ExpectQuery(`SELECT \* FROM "item_media" WHERE id = \$1 AND item_id = \$2`).
			WithArgs(mediaID, itemID).
			WillReturnRows(rows)

		media, err := repo.FindByIDAndItemID(mediaID, itemID)

		assert.NoError(t, err)
		assert.NotNil(t, media)
		assert.Equal(t, mediaID, media.ID)
		assert.Equal(t, itemID, media.ItemID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success - Media not found", func(t *testing.T) {
		mediaID := int64(1)
		itemID := uuid.New()

		mock.ExpectQuery(`SELECT \* FROM "item_media" WHERE id = \$1 AND item_id = \$2`).
			WithArgs(mediaID, itemID).
			WillReturnError(gorm.ErrRecordNotFound)

		media, err := repo.FindByIDAndItemID(mediaID, itemID)

		assert.NoError(t, err)
		assert.Nil(t, media)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestItemMediaRepository_Delete(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewItemMediaRepository(db)

	t.Run("Success - Delete media", func(t *testing.T) {
		mediaID := int64(1)
		itemID := uuid.New()
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "item_id", "media_type", "url", "thumbnail_url", "display_order", "created_at",
		}).AddRow(mediaID, itemID, "image", "http://example.com/image.webp", "http://example.com/thumb.webp", 0, now)

		mock.ExpectQuery(`SELECT \* FROM "item_media" WHERE "item_media"\."id" = \$1`).
			WithArgs(mediaID).
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "item_media" WHERE "item_media"\."id" = \$1`).
			WithArgs(mediaID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		media, err := repo.Delete(mediaID)

		assert.NoError(t, err)
		assert.NotNil(t, media)
		assert.Equal(t, mediaID, media.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success - Media not found", func(t *testing.T) {
		mediaID := int64(999)

		mock.ExpectQuery(`SELECT \* FROM "item_media" WHERE "item_media"\."id" = \$1`).
			WithArgs(mediaID).
			WillReturnError(gorm.ErrRecordNotFound)

		media, err := repo.Delete(mediaID)

		assert.NoError(t, err)
		assert.Nil(t, media)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestItemMediaRepository_CountByItemIDAndType(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewItemMediaRepository(db)

	t.Run("Success - Count images", func(t *testing.T) {
		itemID := uuid.New()

		rows := sqlmock.NewRows([]string{"count"}).AddRow(5)

		mock.ExpectQuery(`SELECT count\(\*\) FROM "item_media" WHERE item_id = \$1 AND media_type = \$2`).
			WithArgs(itemID, domain.MediaTypeImage).
			WillReturnRows(rows)

		count, err := repo.CountByItemIDAndType(itemID, domain.MediaTypeImage)

		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestItemMediaRepository_CountByItemID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewItemMediaRepository(db)

	t.Run("Success - Count all media", func(t *testing.T) {
		itemID := uuid.New()

		rows := sqlmock.NewRows([]string{"count"}).AddRow(8)

		mock.ExpectQuery(`SELECT count\(\*\) FROM "item_media" WHERE item_id = \$1`).
			WithArgs(itemID).
			WillReturnRows(rows)

		count, err := repo.CountByItemID(itemID)

		assert.NoError(t, err)
		assert.Equal(t, int64(8), count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestItemMediaRepository_GetNextDisplayOrder(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewItemMediaRepository(db)

	t.Run("Success - Has existing media", func(t *testing.T) {
		itemID := uuid.New()

		rows := sqlmock.NewRows([]string{"coalesce"}).AddRow(3)

		mock.ExpectQuery(`SELECT COALESCE\(MAX\(display_order\), -1\) FROM "item_media" WHERE item_id = \$1`).
			WithArgs(itemID).
			WillReturnRows(rows)

		order, err := repo.GetNextDisplayOrder(itemID)

		assert.NoError(t, err)
		assert.Equal(t, 4, order)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success - No existing media", func(t *testing.T) {
		itemID := uuid.New()

		rows := sqlmock.NewRows([]string{"coalesce"}).AddRow(-1)

		mock.ExpectQuery(`SELECT COALESCE\(MAX\(display_order\), -1\) FROM "item_media" WHERE item_id = \$1`).
			WithArgs(itemID).
			WillReturnRows(rows)

		order, err := repo.GetNextDisplayOrder(itemID)

		assert.NoError(t, err)
		assert.Equal(t, 0, order)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestItemMediaRepository_FindByIDsAndItemID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewItemMediaRepository(db)

	t.Run("Success - Find multiple media", func(t *testing.T) {
		itemID := uuid.New()
		ids := []int64{1, 2, 3}
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "item_id", "media_type", "url", "thumbnail_url", "display_order", "created_at",
		}).
			AddRow(1, itemID, "image", "http://example.com/image1.webp", nil, 0, now).
			AddRow(2, itemID, "image", "http://example.com/image2.webp", nil, 1, now).
			AddRow(3, itemID, "image", "http://example.com/image3.webp", nil, 2, now)

		mock.ExpectQuery(`SELECT \* FROM "item_media" WHERE id IN \(\$1,\$2,\$3\) AND item_id = \$4`).
			WithArgs(int64(1), int64(2), int64(3), itemID).
			WillReturnRows(rows)

		mediaList, err := repo.FindByIDsAndItemID(ids, itemID)

		assert.NoError(t, err)
		assert.Len(t, mediaList, 3)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestItemMediaRepository_DeleteByItemID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewItemMediaRepository(db)

	t.Run("Success - Delete all media for item", func(t *testing.T) {
		itemID := uuid.New()
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "item_id", "media_type", "url", "thumbnail_url", "display_order", "created_at",
		}).
			AddRow(1, itemID, "image", "http://example.com/image1.webp", nil, 0, now).
			AddRow(2, itemID, "image", "http://example.com/image2.webp", nil, 1, now)

		mock.ExpectQuery(`SELECT \* FROM "item_media" WHERE item_id = \$1`).
			WithArgs(itemID).
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "item_media" WHERE item_id = \$1`).
			WithArgs(itemID).
			WillReturnResult(sqlmock.NewResult(0, 2))
		mock.ExpectCommit()

		mediaList, err := repo.DeleteByItemID(itemID)

		assert.NoError(t, err)
		assert.Len(t, mediaList, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success - No media to delete", func(t *testing.T) {
		itemID := uuid.New()

		rows := sqlmock.NewRows([]string{
			"id", "item_id", "media_type", "url", "thumbnail_url", "display_order", "created_at",
		})

		mock.ExpectQuery(`SELECT \* FROM "item_media" WHERE item_id = \$1`).
			WithArgs(itemID).
			WillReturnRows(rows)

		mediaList, err := repo.DeleteByItemID(itemID)

		assert.NoError(t, err)
		assert.Len(t, mediaList, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
