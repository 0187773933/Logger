package server

import (
	"fmt"
	fiber "github.com/gofiber/fiber/v2"
)

func ( s *Server ) Test( c *fiber.Ctx ) ( error ) {
	log.Debug( "hola" )
	return c.JSON( fiber.Map{
		"url": "/admin/test" ,
		"result": true ,
	})
}

func ( s *Server ) SetupAdminRoutes() {
	admin := s.FiberApp.Group( fmt.Sprintf( "/%s" , s.Config.ServerUrlPrefix ) )
	admin.Use( s.ValidateAdminMW )
	admin.Get( "/test" , s.Test )
}