package client

import "strings"

func (c *client) formatQuery(query string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(query)), "")
}
