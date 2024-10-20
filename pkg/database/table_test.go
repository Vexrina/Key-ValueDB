package database

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTable(t *testing.T) {
	table := NewTableImpl()

	parsedTime, err := table.ParseTime("25.10.2024")
	require.NoError(t, err)

	// 1. Тест вставки новой записи
	t.Run("Insert", func(t *testing.T) {
		initialValue := Value{
			Val: "hello, world",
			Ttl: parsedTime,
		}

		ok, err := table.Insert(1, initialValue)
		require.NoError(t, err)
		assert.True(t, ok)

		actual, err := table.Get(1)
		require.NoError(t, err)
		assert.Equal(t, initialValue.Val, actual.Val)
	})

	// 2. Тест обновления существующей записи
	t.Run("Update", func(t *testing.T) {
		updatedValue := Value{
			Val: "updated value",
			Ttl: parsedTime,
		}

		expectedValue := Value{
			Val: "updated value",
			Ttl: parsedTime,
		}

		ok, err := table.Update(1, updatedValue)
		require.NoError(t, err)
		assert.True(t, ok)

		assert.Equal(t, expectedValue, updatedValue)
	})

	// 3. Тест получения записи
	t.Run("Get", func(t *testing.T) {
		expectedValue := Value{
			Val: "updated value",
			Ttl: parsedTime,
		}

		actual, err := table.Get(1)
		require.NoError(t, err)
		assert.Equal(t, expectedValue.Val, actual.Val)
	})

	// 4. Тест удаления записи
	t.Run("Delete", func(t *testing.T) {
		ok, err := table.Delete(1)
		require.NoError(t, err)
		assert.True(t, ok)

		assert.Equal(t, 0, table.Size())
	})
}
