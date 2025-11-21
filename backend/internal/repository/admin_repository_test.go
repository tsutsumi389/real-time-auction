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

func TestAdminRepository_FindAdminsWithFilters(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewAdminRepository(db)
	now := time.Now()

	t.Run("Success - Default filters (no deleted)", func(t *testing.T) {
		req := &domain.AdminListRequest{
			Page:  1,
			Limit: 20,
		}

		rows := sqlmock.NewRows([]string{
			"id", "email", "password_hash", "display_name", "role", "status", "created_at", "updated_at",
		}).
			AddRow(int64(1), "admin1@example.com", "hash1", "Admin 1", "system_admin", "active", now, now).
			AddRow(int64(2), "admin2@example.com", "hash2", "Admin 2", "auctioneer", "active", now, now)

		// GORM may use LIMIT with hard-coded value instead of placeholder
		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE status IN \(\$1,\$2\) ORDER BY id ASC LIMIT`).
			WithArgs(domain.StatusActive, domain.StatusSuspended).
			WillReturnRows(rows)

		admins, err := repo.FindAdminsWithFilters(req)

		assert.NoError(t, err)
		assert.Len(t, admins, 2)
		assert.Equal(t, int64(1), admins[0].ID)
		assert.Equal(t, "admin1@example.com", admins[0].Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success - With keyword filter", func(t *testing.T) {
		req := &domain.AdminListRequest{
			Page:    1,
			Limit:   20,
			Keyword: "admin1",
		}

		rows := sqlmock.NewRows([]string{
			"id", "email", "password_hash", "display_name", "role", "status", "created_at", "updated_at",
		}).
			AddRow(int64(1), "admin1@example.com", "hash1", "Admin 1", "system_admin", "active", now, now)

		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE email ILIKE \$1 AND status IN \(\$2,\$3\) ORDER BY id ASC LIMIT`).
			WithArgs("%admin1%", domain.StatusActive, domain.StatusSuspended).
			WillReturnRows(rows)

		admins, err := repo.FindAdminsWithFilters(req)

		assert.NoError(t, err)
		assert.Len(t, admins, 1)
		assert.Equal(t, "admin1@example.com", admins[0].Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success - With role filter", func(t *testing.T) {
		req := &domain.AdminListRequest{
			Page:  1,
			Limit: 20,
			Role:  domain.RoleSystemAdmin,
		}

		rows := sqlmock.NewRows([]string{
			"id", "email", "password_hash", "display_name", "role", "status", "created_at", "updated_at",
		}).
			AddRow(int64(1), "admin1@example.com", "hash1", "Admin 1", "system_admin", "active", now, now)

		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE role = \$1 AND status IN \(\$2,\$3\) ORDER BY id ASC LIMIT`).
			WithArgs(domain.RoleSystemAdmin, domain.StatusActive, domain.StatusSuspended).
			WillReturnRows(rows)

		admins, err := repo.FindAdminsWithFilters(req)

		assert.NoError(t, err)
		assert.Len(t, admins, 1)
		assert.Equal(t, domain.RoleSystemAdmin, admins[0].Role)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success - With status filter", func(t *testing.T) {
		req := &domain.AdminListRequest{
			Page:   1,
			Limit:  20,
			Status: []domain.AdminStatus{domain.StatusActive},
		}

		rows := sqlmock.NewRows([]string{
			"id", "email", "password_hash", "display_name", "role", "status", "created_at", "updated_at",
		}).
			AddRow(int64(1), "admin1@example.com", "hash1", "Admin 1", "system_admin", "active", now, now)

		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE status IN \(\$1\) ORDER BY id ASC LIMIT`).
			WithArgs(domain.StatusActive).
			WillReturnRows(rows)

		admins, err := repo.FindAdminsWithFilters(req)

		assert.NoError(t, err)
		assert.Len(t, admins, 1)
		assert.Equal(t, domain.StatusActive, admins[0].Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success - With pagination offset", func(t *testing.T) {
		req := &domain.AdminListRequest{
			Page:  2,
			Limit: 20,
		}

		rows := sqlmock.NewRows([]string{
			"id", "email", "password_hash", "display_name", "role", "status", "created_at", "updated_at",
		}).
			AddRow(int64(21), "admin21@example.com", "hash21", "Admin 21", "auctioneer", "active", now, now)

		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE status IN \(\$1,\$2\) ORDER BY id ASC LIMIT .* OFFSET`).
			WithArgs(domain.StatusActive, domain.StatusSuspended).
			WillReturnRows(rows)

		admins, err := repo.FindAdminsWithFilters(req)

		assert.NoError(t, err)
		assert.Len(t, admins, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success - With sort by email desc", func(t *testing.T) {
		req := &domain.AdminListRequest{
			Page:  1,
			Limit: 20,
			Sort:  "email_desc",
		}

		rows := sqlmock.NewRows([]string{
			"id", "email", "password_hash", "display_name", "role", "status", "created_at", "updated_at",
		}).
			AddRow(int64(2), "admin2@example.com", "hash2", "Admin 2", "auctioneer", "active", now, now).
			AddRow(int64(1), "admin1@example.com", "hash1", "Admin 1", "system_admin", "active", now, now)

		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE status IN \(\$1,\$2\) ORDER BY email DESC LIMIT`).
			WithArgs(domain.StatusActive, domain.StatusSuspended).
			WillReturnRows(rows)

		admins, err := repo.FindAdminsWithFilters(req)

		assert.NoError(t, err)
		assert.Len(t, admins, 2)
		assert.Equal(t, "admin2@example.com", admins[0].Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success - Combined filters", func(t *testing.T) {
		req := &domain.AdminListRequest{
			Page:    1,
			Limit:   10,
			Keyword: "admin",
			Role:    domain.RoleSystemAdmin,
			Status:  []domain.AdminStatus{domain.StatusActive},
			Sort:    "created_at_desc",
		}

		rows := sqlmock.NewRows([]string{
			"id", "email", "password_hash", "display_name", "role", "status", "created_at", "updated_at",
		}).
			AddRow(int64(1), "admin@example.com", "hash", "Admin", "system_admin", "active", now, now)

		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE email ILIKE \$1 AND role = \$2 AND status IN \(\$3\) ORDER BY created_at DESC LIMIT`).
			WithArgs("%admin%", domain.RoleSystemAdmin, domain.StatusActive).
			WillReturnRows(rows)

		admins, err := repo.FindAdminsWithFilters(req)

		assert.NoError(t, err)
		assert.Len(t, admins, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAdminRepository_CountAdminsWithFilters(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewAdminRepository(db)

	t.Run("Success - Default filters", func(t *testing.T) {
		req := &domain.AdminListRequest{
			Page:  1,
			Limit: 20,
		}

		rows := sqlmock.NewRows([]string{"count"}).AddRow(int64(50))

		mock.ExpectQuery(`SELECT count\(\*\) FROM "admins" WHERE status IN \(\$1,\$2\)`).
			WithArgs(domain.StatusActive, domain.StatusSuspended).
			WillReturnRows(rows)

		count, err := repo.CountAdminsWithFilters(req)

		assert.NoError(t, err)
		assert.Equal(t, int64(50), count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success - With keyword filter", func(t *testing.T) {
		req := &domain.AdminListRequest{
			Page:    1,
			Limit:   20,
			Keyword: "test",
		}

		rows := sqlmock.NewRows([]string{"count"}).AddRow(int64(10))

		mock.ExpectQuery(`SELECT count\(\*\) FROM "admins" WHERE email ILIKE \$1 AND status IN \(\$2,\$3\)`).
			WithArgs("%test%", domain.StatusActive, domain.StatusSuspended).
			WillReturnRows(rows)

		count, err := repo.CountAdminsWithFilters(req)

		assert.NoError(t, err)
		assert.Equal(t, int64(10), count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success - With role and status filters", func(t *testing.T) {
		req := &domain.AdminListRequest{
			Page:   1,
			Limit:  20,
			Role:   domain.RoleAuctioneer,
			Status: []domain.AdminStatus{domain.StatusActive, domain.StatusSuspended},
		}

		rows := sqlmock.NewRows([]string{"count"}).AddRow(int64(15))

		mock.ExpectQuery(`SELECT count\(\*\) FROM "admins" WHERE role = \$1 AND status IN \(\$2,\$3\)`).
			WithArgs(domain.RoleAuctioneer, domain.StatusActive, domain.StatusSuspended).
			WillReturnRows(rows)

		count, err := repo.CountAdminsWithFilters(req)

		assert.NoError(t, err)
		assert.Equal(t, int64(15), count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAdminRepository_UpdateAdminStatus(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewAdminRepository(db)

	t.Run("Success - Update to suspended", func(t *testing.T) {
		adminID := int64(1)
		newStatus := domain.StatusSuspended

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "admins" SET "status"=\$1,"updated_at"=\$2 WHERE id = \$3`).
			WithArgs(newStatus, sqlmock.AnyArg(), adminID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.UpdateAdminStatus(adminID, newStatus)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success - Update to active", func(t *testing.T) {
		adminID := int64(2)
		newStatus := domain.StatusActive

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "admins" SET "status"=\$1,"updated_at"=\$2 WHERE id = \$3`).
			WithArgs(newStatus, sqlmock.AnyArg(), adminID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.UpdateAdminStatus(adminID, newStatus)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
