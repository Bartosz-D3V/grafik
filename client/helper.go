// Package client contains the code used internally by grafik to prepare & send HTTP requests.
package client

import "strings"

// formatQuery provides simple GraphQL code compression.
func (c *client) formatQuery(query string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(query)), " ")
}
