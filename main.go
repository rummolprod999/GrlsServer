package main

func init() {
	CreateEnv()
}
func main() {
	defer SaveStack()
	Logging("start")
	server := ServerGrls{":8181"}
	server.run()
	Logging("end")
}
