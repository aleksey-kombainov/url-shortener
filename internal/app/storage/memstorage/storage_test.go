package memstorage

import (
	"context"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/model"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

const (
	userIDMultiplyTests = "u1"
	errPrefix           = "err"
)

type urlTest struct {
	id       int
	origURL  string
	shortURL string
	userID   string
}

var (
	tests = []urlTest{
		{
			id:       11,
			userID:   userIDMultiplyTests,
			origURL:  `https://practicum.yandex.ru/`,
			shortURL: `sh1`,
		},
		{
			id:       22,
			userID:   userIDMultiplyTests,
			origURL:  `http://ya.ru`,
			shortURL: `sh2`,
		},
		{
			id:       33,
			userID:   "u3",
			origURL:  `http://brbrbr.com`,
			shortURL: `sh3`,
		},
	}
)

func TestStorage(t *testing.T) {
	ctx := context.Background()

	storage := New()

	createRecord(ctx, t, storage)
	deleteRecord(ctx, t, storage)
}

func createRecord(ctx context.Context, t *testing.T, storage *Storage) {

	for i, test := range tests {
		t.Run(`Test CreateRecord #`+strconv.Itoa(i), func(t *testing.T) {
			_, err := storage.CreateRecord(ctx, test.origURL, test.shortURL, test.userID)
			assert.NoError(t, err, "")
		})
	}

	for i, test := range tests {
		t.Run(`Test GetOriginalURLByShortcut #`+strconv.Itoa(i), func(t *testing.T) {
			sh, err := storage.GetOriginalURLByShortcut(ctx, test.shortURL)
			assert.NoError(t, err, "")
			assert.Equal(t, test.origURL, sh.OriginalURL)
			assert.Equal(t, test.userID, sh.UserID)
			assert.False(t, sh.DeletedFlag)
		})
	}

	for i, test := range tests {
		t.Run(`Test gets GetShortcutByOriginalURL #`+strconv.Itoa(i), func(t *testing.T) {
			sh, err := storage.GetShortcutByOriginalURL(ctx, test.origURL)
			assert.NoError(t, err, "")
			assert.Equal(t, test.shortURL, sh.ShortURL)
			assert.Equal(t, test.userID, sh.UserID)
			assert.False(t, sh.DeletedFlag)
		})
	}

	for i, test := range tests {
		t.Run(`Test gets GetShortcutsByUser #`+strconv.Itoa(i), func(t *testing.T) {
			shs, err := storage.GetShortcutsByUser(ctx, test.userID)
			assert.NoError(t, err, "")

			userTests := getTestsByUser(test.userID)

			assert.Len(t, shs, len(userTests))

			for _, sh := range shs {
				assert.Equal(t, test.userID, sh.UserID)
				assert.False(t, sh.DeletedFlag)
			}
		})
	}
}

func deleteRecord(ctx context.Context, t *testing.T, storage *Storage) {
	// DELETE
	delTask := model.DeleteTask{
		UserID: userIDMultiplyTests,
	}
	for _, temp := range getTestsByUser("u1") {
		delTask.ShortURLs = append(delTask.ShortURLs, temp.shortURL)
	}
	err := storage.DeleteByShortcutsAndUser(ctx, append([]model.DeleteTask{}, delTask))
	assert.NoError(t, err, "")

	shs, err := storage.GetShortcutsByUser(ctx, userIDMultiplyTests)
	assert.NoError(t, err, "")
	assert.Len(t, shs, 0)
}

func batch(ctx context.Context, t *testing.T, storage *Storage) {

}

func getTestsByUser(userID string) []urlTest {
	var ret []urlTest
	for _, test := range tests {
		if test.userID == userID {
			ret = append(ret, test)
		}
	}
	return ret
}
