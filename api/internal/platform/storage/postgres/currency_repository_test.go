package postgres

import (
	bole "boletia/api/internal"
	"context"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_CourseRepository_Save_RepositoryError(t *testing.T) {
	var criteria bole.Criteria
	newCriteria, err := criteria.NewCriteria("MXN", "2024-02-28T23:50:56", "2024-03-28T23:50:56")
	require.NoError(t, err)
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	require.NoError(t, err)

	mock.ExpectQuery(
		"SELECT code, value, last_updated_at FROM currency_rates WHERE code = $1 AND last_updated_at >= $2 AND last_updated_at <= $3").
		WithArgs(newCriteria.CurrencyCode, newCriteria.StartDate, newCriteria.EndDate).
		WillReturnRows(pgxmock.NewRows([]string{"code", "value", "last_updated_at"}).
			AddRow("MXN", 1.23, time.Now()).
			AddRow("MXN", 1.11, time.Now()))

	repo := NewDatabaseRepository(mock, 1*time.Millisecond, "currency_rates")

	data, err := repo.Get(context.Background(), newCriteria)
	require.NoError(t, err)
	require.Equal(t, 2, len(data))
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NoError(t, err)
}

//
//func Test_CourseRepository_Save_Succeed(t *testing.T) {
//	courseID, courseName, courseDuration := "37a0f027-15e6-47cc-a5d2-64183281087e", "Test Course", "10 months"
//
//	course, err := bole.NewCourse(courseID, courseName, courseDuration)
//	require.NoError(t, err)
//
//	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
//	require.NoError(t, err)
//
//	sqlMock.ExpectExec(
//		"INSERT INTO currency (id, name, duration) VALUES (?, ?, ?)").
//		WithArgs(courseID, courseName, courseDuration).
//		WillReturnResult(sqlmock.NewResult(0, 1))
//
//	repo := NewCourseRepository(db, 1*time.Millisecond)
//
//	err = repo.Save(context.Background(), course)
//
//	assert.NoError(t, sqlMock.ExpectationsWereMet())
//	assert.NoError(t, err)
//}
