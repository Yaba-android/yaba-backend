/*
*
* See this doc for more information about HTTP protocol
* https://tools.ietf.org/html/rfc7231#section-4.3
*
 */

package main

import (
	"github.com/nasrat_v/maktaba-android-mvp/src/services/database"
	"github.com/nasrat_v/maktaba-android-mvp/src/services/routes"
)

func main() {
	database.InitDbHandlerForControllers()
	routes.InitRouterForControllers()
}
