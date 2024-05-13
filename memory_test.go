package dataprovider

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryProvider(t *testing.T) {
	ctx := context.Background()

	cfg := &ConfigModule{
		Driver: MemoryDataProviderName,
	}

	provider, err := NewProvider(ctx, cfg)
	assert.NoError(t, err)

	providerStatus := provider.GetProviderStatus()
	assert.Equal(t, MemoryDataProviderName, providerStatus.Driver)
}
