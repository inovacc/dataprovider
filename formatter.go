package dataprovider

import (
	"fmt"
	"strings"
)

// PlaceholderFormatter is dialect-aware: PostgresSQL uses $1, Oracle uses: p1, default uses?
type PlaceholderFormatter interface {
	ReplacePlaceholders(query string) string
}

type postgresFormatter struct{}

func (f *postgresFormatter) ReplacePlaceholders(query string) string {
	for i := 1; strings.Contains(query, "?"); i++ {
		query = strings.Replace(query, "?", fmt.Sprintf("$%d", i), 1)
	}
	return query
}

type oracleFormatter struct{}

func (f *oracleFormatter) ReplacePlaceholders(query string) string {
	for i := 1; strings.Contains(query, "?"); i++ {
		query = strings.Replace(query, "?", fmt.Sprintf(":p%d", i), 1)
	}
	return query
}

type defaultFormatter struct{}

func (f *defaultFormatter) ReplacePlaceholders(query string) string {
	return query
}

func NewFormatter(driver string) PlaceholderFormatter {
	switch driver {
	case PostgresSQLDatabaseProviderName:
		return &postgresFormatter{}
	case OracleDatabaseProviderName:
		return &oracleFormatter{}
	default:
		return &defaultFormatter{}
	}
}
