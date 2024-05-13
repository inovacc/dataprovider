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

	err := newProvider(ctx, cfg)
	assert.NoError(t, err)

	provider := GetProvider()
}
