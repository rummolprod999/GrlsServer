package main

import "fmt"

func init() {
	CreateEnv()
}
func main() {
	defer SaveStack()
	Logging("start")
	server := ServerGrls{fmt.Sprintf(":%s", Port)}
	server.run()
	Logging("end")
}
