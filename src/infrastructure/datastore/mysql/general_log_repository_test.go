package mysql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestQueries tests the GetQueries method using a mock database
func TestQueries(t *testing.T) {
	// テスト用にモックを使用する場合のサンプル
	// 実際のDBに接続する場合は環境変数などから設定を読み込む
	t.Skip("Skipping test that requires database connection")

	handler, closer, err := NewMySQLHandler(
		"root",
		"password",
		"localhost",
		3306,
		"mysql")

	require.NoError(t, err)
	defer func() {
		if err := closer(); err != nil {
			t.Logf("Error closing database connection: %v", err)
		}
	}()

	repo := NewGeneralLogRepository(handler)

	// データクリーンアップ
	ctx := context.Background()
	err = repo.Clear(ctx)
	require.NoError(t, err)

	// テストデータを挿入
	db, err := handler.DB()
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO general_log (command_type, argument) VALUES ('Query', 'SELECT * FROM test_table')")
	require.NoError(t, err)

	// テスト実行
	queries, err := repo.GetQueries(ctx)
	require.NoError(t, err)

	// アサーション
	assert.Len(t, queries, 1)
	assert.Equal(t, "SELECT * FROM test_table", queries[0])
}
