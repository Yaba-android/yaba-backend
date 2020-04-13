package main

/********************************
*
* 		MAIN
*
********************************/

func main() {
	client := redisStartClient()
	startRouter(client)
}
