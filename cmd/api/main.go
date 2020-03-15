package main

import (
	"flag"

	api "github.com/saurabh-arch/go-chat/cmd/api/server"
)

var (
	serviceName = flag.String("name", "API Service", "the name of service")
	serviceType = flag.String("type", "API", "the type of service")
	ip          = flag.String("ip", "127.0.0.1", "API IP address")
	port        = flag.Int("port", 6000, "API Port")
)

func parseArgs() {
	flag.Parse()
}

func main() {
	parseArgs()
	res, errStartServer := api.StartWebServer(*ip, *port)
}
