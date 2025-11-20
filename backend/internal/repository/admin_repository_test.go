package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestAdminRepository_FindByEmail(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewAdminRepository(db)

	t.Run("Success - Admin found", func(t *testing.T) {
		email := "admin@example.com"
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "email", "password_hash", "display_name", "role", "status", "created_at", "updated_at",
		}).AddRow(
			int64(1), email, "hashed_password", "Admin User", "system_admin", "active", now, now,
		)

		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE email = \$1`).
			WithArgs(email).
			WillReturnRows(rows)

		admin, err := repo.FindByEmail(email)

		assert.NoError(t, err)
		assert.NotNil(t, admin)
		assert.Equal(t, int64(1), admin.ID)
		assert.Equal(t, email, admin.Email)
		assert.Equal(t, "Admin User", admin.DisplayName)
		assert.Equal(t, domain.RoleSystemAdmin, admin.Role)
		assert.Equal(t, domain.StatusActive, admin.Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success - Admin not found", func(t *testing.T) {
		email := "notfound@example.com"

		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE email = \$1`).
			WithArgs(email).
			WillReturnError(gorm.ErrRecordNotFound)

		admin, err := repo.FindByEmail(email)

		assert.NoError(t, err)
		assert.Nil(t, admin)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAdminRepository_FindByID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewAdminRepository(db)

	t.Run("Success - Admin found", func(t *testing.T) {
		adminID := int64(1)
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "email", "password_hash", "display_name", "role", "status", "created_at", "updated_at",
		}).AddRow(
			adminID, "admin@example.com", "hashed_password", "Admin User", "system_admin", "active", now, now,
		)

		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE "admins"\."id" = \$1`).
			WithArgs(adminID).
			WillReturnRows(rows)

		admin, err := repo.FindByID(adminID)

		assert.NoError(t, err)
		assert.NotNil(t, admin)
		assert.Equal(t, adminID, admin.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success - Admin not found", func(t *testing.T) {
		adminID := int64(999)

		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE "admins"\."id" = \$1`).
			WithArgs(adminID).
			WillReturnError(gorm.ErrRecordNotFound)

		admin, err := repo.FindByID(adminID)

		assert.NoError(t, err)
		assert.Nil(t, admin)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAdminRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewAdminRepository(db)

	t.Run("Success", func(t *testing.T) {
		admin := &domain.Admin{
			Email:        "newadmin@example.com",
			PasswordHash: "hashed_password",
			DisplayName:  "New Admin",
			Role:         domain.RoleAuctioneer,
			Status:       domain.StatusActive,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "admins"`).
			WithArgs(
				admin.Email,
				admin.PasswordHash,
				admin.DisplayName,
				admin.Role,
				admin.Status,
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := repo.Create(admin)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAdminRepository_Update(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewAdminRepository(db)

	t.Run("Success", func(t *testing.T) {
		admin := &domain.Admin{
			ID:           1,
			Email:        "admin@example.com",
			PasswordHash: "hashed_password",
			DisplayName:  "Updated Admin",
			Role:         domain.RoleSystemAdmin,
			Status:       domain.StatusActive,
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "admins"`).
			WithArgs(
				admin.Email,
				admin.PasswordHash,
				admin.DisplayName,
				admin.Role,
				admin.Status,
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
				admin.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(admin)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAdminRepository_Delete(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewAdminRepository(db)

	t.Run("Success", func(t *testing.T) {
		adminID := int64(1)

		mock.ExpectBegin()
		// GORM automatically adds updated_at field, so we need to match that in the expectation
		mock.ExpectExec(`UPDATE "admins" SET "status"=\$1,"updated_at"=\$2 WHERE id = \$3`).
			WithArgs(domain.StatusDeleted, sqlmock.AnyArg(), adminID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(adminID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
