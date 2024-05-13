package dataprovider

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProvider(t *testing.T) {
	cfg := &ConfigModule{
		Driver: MemoryDataProviderName,
	}

	provider, err := NewProvider(context.Background(), cfg)
	assert.NoError(t, err)

	providerStatus := provider.GetProviderStatus()
	assert.Equal(t, MemoryDataProviderName, providerStatus.Driver)
}
