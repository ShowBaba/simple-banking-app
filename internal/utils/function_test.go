package utils

import (
	"github.com/gofiber/fiber/v2/utils"
	"testing"
)

func TestGenerateAccountID(t *testing.T) {
	t.Run("successfully generates account id", func(t *testing.T) {
		accountId := GenerateAccountID()

		utils.AssertEqual(t, "", accountId, "is not empty")
	})

	t.Run("length of generated id", func(t *testing.T) {
		accountId := GenerateAccountID()
		utils.AssertEqual(t, 19, len(accountId), "length is valid")
	})
}
