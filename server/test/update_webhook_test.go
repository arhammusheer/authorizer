package test

import (
	"fmt"
	"testing"

	"github.com/authorizerdev/authorizer/server/constants"
	"github.com/authorizerdev/authorizer/server/crypto"
	"github.com/authorizerdev/authorizer/server/db"
	"github.com/authorizerdev/authorizer/server/graph/model"
	"github.com/authorizerdev/authorizer/server/memorystore"
	"github.com/authorizerdev/authorizer/server/resolvers"
	"github.com/authorizerdev/authorizer/server/utils"
	"github.com/stretchr/testify/assert"
)

func updateWebhookTest(t *testing.T, s TestSetup) {
	t.Helper()
	t.Run("should update webhook", func(t *testing.T) {
		req, ctx := createContext(s)
		adminSecret, err := memorystore.Provider.GetStringStoreEnvVariable(constants.EnvKeyAdminSecret)
		assert.NoError(t, err)
		h, err := crypto.EncryptPassword(adminSecret)
		assert.NoError(t, err)
		req.Header.Set("Cookie", fmt.Sprintf("%s=%s", constants.AdminCookieName, h))
		// get webhook
		webhook, err := db.Provider.GetWebhookByEventName(ctx, constants.UserDeletedWebhookEvent)
		assert.NoError(t, err)
		assert.NotNil(t, webhook)
		webhook.Headers["x-new-test"] = "new-test"

		res, err := resolvers.UpdateWebhookResolver(ctx, model.UpdateWebhookRequest{
			ID:      webhook.ID,
			Headers: webhook.Headers,
			Enabled: utils.NewBoolRef(false),
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.NotEmpty(t, res.Message)

		updatedWebhook, err := db.Provider.GetWebhookByEventName(ctx, constants.UserDeletedWebhookEvent)
		assert.NoError(t, err)
		assert.NotNil(t, updatedWebhook)
		assert.Equal(t, webhook.ID, updatedWebhook.ID)
		assert.Equal(t, utils.StringValue(webhook.EventName), utils.StringValue(updatedWebhook.EventName))
		assert.Equal(t, utils.StringValue(webhook.Endpoint), utils.StringValue(updatedWebhook.Endpoint))
		assert.Len(t, updatedWebhook.Headers, 2)
		assert.False(t, utils.BoolValue(updatedWebhook.Enabled))

		res, err = resolvers.UpdateWebhookResolver(ctx, model.UpdateWebhookRequest{
			ID:      webhook.ID,
			Headers: webhook.Headers,
			Enabled: utils.NewBoolRef(true),
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.NotEmpty(t, res.Message)
	})
}