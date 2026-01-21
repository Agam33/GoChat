package websocket

import "fmt"

func BuildWSTopic(domain string, scope string, id int64) string {
	return fmt.Sprintf("%s_%s_%d", domain, scope, id)
}
