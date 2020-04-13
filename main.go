package yaba

/********************************
*
* 		MAIN
*
********************************/

func main() {
	client := redisStartClient()

	startRouter(client)
}
