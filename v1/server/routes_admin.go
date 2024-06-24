package server

import (
	"fmt"
	"strconv"
	fiber "github.com/gofiber/fiber/v2"
	logger "github.com/0187773933/Logger/v1/logger"
)

func ( s *Server ) Test( c *fiber.Ctx ) ( error ) {
	log.Debug( "hola" )
	return c.JSON( fiber.Map{
		"url": "/test" ,
		"result": true ,
	})
}

func ( s *Server ) GetLogMessages( c *fiber.Ctx ) ( error ) {
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
	admin.Get( "/test" , s.Test )
	admin.Get( "/log/:count" , s.GetLogMessages )
}