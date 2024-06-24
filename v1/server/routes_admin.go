package server

import (
	"fmt"
	"strconv"
	net_url "net/url"
	fiber "github.com/gofiber/fiber/v2"
	logger "github.com/0187773933/Logger/v1/logger"
)

func ( s *Server ) LogMessage( c *fiber.Ctx ) ( error ) {
	message := c.Params( "message" )
	if message == "" {
		return c.JSON( fiber.Map{
			"result": false ,
			"error": "No Message Provided" ,
		})
	}
	unescaped_message , _ := net_url.QueryUnescape( message )
	log.Info( unescaped_message )
	return c.JSON( fiber.Map{
		"url": "/log" ,
		"message": unescaped_message ,
		"result": true ,
	})
}

func ( s *Server ) GetLogMessages( c *fiber.Ctx ) ( error ) {
	key := c.Query( "key" )
	if key == "" { key = s.Config.LogKey }
	count := c.Params( "count" )
	count_int , _ := strconv.Atoi( count )
	messages := logger.GetMessages( count_int )
	return c.JSON( fiber.Map{
		"result": true ,
		"url": "/log/:count" ,
		"count": count ,
		"messages": messages ,
	})
}

func ( s *Server ) SetupAdminRoutes() {
	admin := s.FiberApp.Group( fmt.Sprintf( "/%s" , s.Config.ServerUrlPrefix ) )
	admin.Use( s.ValidateAdminMW )
	admin.Get( "/log/:message" , s.LogMessage )
	admin.Get( "/log/view/:count" , s.GetLogMessages )
}