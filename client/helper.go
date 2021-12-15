package client

import "strings"

// formatQuery provides simple GraphQL code compression.
func (c *client) formatQuery(query string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(query)), " ")
}
