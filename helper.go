package goteleport

import "strings"

func getHost(hostWithPort string)string{
	return strings.Split(hostWithPort, ":")[0]
}

