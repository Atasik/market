package repository

import (
	"fmt"
	"market/internal/model"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUserPostgres_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	r := NewUserPostgresqlRepo(sqlxDB)

	tests := []struct {
		name    string
		mock    func()
		input   model.User
		want    int
		wantErr bool
	}{{
		name: "OK",
		mock: func() {
			rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
			mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", usersTable)).
				WithArgs("Test", "User", "password").WillReturnRows(rows)
		},
		input: model.User{
			Username: "Test",
			Role:     "User",
			Password: "password",
		},
		want: 1,
	}, {
		name: "Empty Input",
		mock: func() {
			rows := sqlmock.NewRows([]string{"id"})
			mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", usersTable)).
				WithArgs("Test", "User", "").WillReturnRows(rows)
		},
		input: model.User{
			Username: "Test",
			Role:     "User",
			Password: "",
		},
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateUser(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserPostgres_GetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	r := NewUserPostgresqlRepo(sqlxDB)

	tests := []struct {
		name    string
		mock    func()
		input   string
		want    model.User
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "role", "username", "password"}).
					AddRow(1, "test_role", "test", "password")
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("test").WillReturnRows(rows)
			},
			input: "test",
			want: model.User{
				ID:       1,
				Role:     "test_role",
				Username: "test",
				Password: "password",
			},
		},
		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "role", "username", "password"})
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("notfound").WillReturnRows(rows)
			},
			input:   "notfound",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetUser(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserPostgres_GetUserById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	r := NewUserPostgresqlRepo(sqlxDB)

	tests := []struct {
		name    string
		mock    func()
		input   int
		want    model.User
		wantErr bool
	}{{
		name: "OK",
		mock: func() {
			rows := sqlmock.NewRows([]string{"id", "role", "username", "password"}).
				AddRow(1, "test_role", "test", "password")
			mock.ExpectQuery("SELECT (.+) FROM users").
				WithArgs(1).WillReturnRows(rows)
		},
		input: 1,
		want: model.User{
			ID:       1,
			Role:     "test_role",
			Username: "test",
			Password: "password",
		},
		wantErr: false,
	}, {
		name: "Wrong Id",
		mock: func() {
			rows := sqlmock.NewRows([]string{"id", "role", "username", "password"})
			mock.ExpectQuery("SELECT (.+) FROM users").
				WithArgs(-1).WillReturnRows(rows)
		},
		input:   -1,
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetUserByID(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
