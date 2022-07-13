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

func webhookTest(t *testing.T, s TestSetup) {
	t.Helper()
	t.Run("should get webhook", func(t *testing.T) {
		req, ctx := createContext(s)
		adminSecret, err := memorystore.Provider.GetStringStoreEnvVariable(constants.EnvKeyAdminSecret)
		assert.NoError(t, err)
		h, err := crypto.EncryptPassword(adminSecret)
		assert.NoError(t, err)
		req.Header.Set("Cookie", fmt.Sprintf("%s=%s", constants.AdminCookieName, h))

		// get webhook by event name
		webhook, err := db.Provider.GetWebhookByEventName(ctx, constants.UserCreatedWebhookEvent)
		assert.NoError(t, err)
		assert.NotNil(t, webhook)

		res, err := resolvers.WebhookResolver(ctx, model.WebhookRequest{
			ID: webhook.ID,
		})
		assert.NoError(t, err)
		assert.Equal(t, res.ID, webhook.ID)
		assert.Equal(t, utils.StringValue(res.Endpoint), utils.StringValue(webhook.Endpoint))
		assert.Equal(t, utils.StringValue(res.EventName), utils.StringValue(webhook.EventName))
		assert.Equal(t, utils.BoolValue(res.Enabled), utils.BoolValue(webhook.Enabled))
		assert.Len(t, res.Headers, len(webhook.Headers))
	})
}