package mysql_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tuhochi/langchaingo/tools/sqldatabase"
	_ "github.com/tuhochi/langchaingo/tools/sqldatabase/mysql"
)

func Test(t *testing.T) {
	t.Parallel()

	// export LANGCHAINGO_TEST_MYSQL=user:p@ssw0rd@tcp(localhost:3306)/test
	mysqlURI := os.Getenv("LANGCHAINGO_TEST_MYSQL")
	if mysqlURI == "" {
		t.Skip("LANGCHAINGO_TEST_MYSQL not set")
	}
	db, err := sqldatabase.NewSQLDatabaseWithDSN("mysql", mysqlURI, nil)
	require.NoError(t, err)

	tbs := db.TableNames()
	require.NotEmpty(t, tbs)

	desc, err := db.TableInfo(context.Background(), tbs)
	require.NoError(t, err)

	t.Log(desc)

	for _, tableName := range tbs {
		_, err = db.Query(context.Background(), fmt.Sprintf("SELECT * from %s LIMIT 1", tableName))
		/* exclude no row error,
		since we only need to check if db.Query function can perform query correctly*/
		if errors.Is(err, sql.ErrNoRows) {
			continue
		}
		require.NoError(t, err)
	}
}
