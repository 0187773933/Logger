package server

import (
	"fmt"
	"time"
	logrus "github.com/sirupsen/logrus"
	logger "github.com/0187773933/Logger/v1/logger"
	types "github.com/0187773933/Logger/v1/types"
	utils "github.com/0187773933/Logger/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	fiber_cookie "github.com/gofiber/fiber/v2/middleware/encryptcookie"
	fiber_cors "github.com/gofiber/fiber/v2/middleware/cors"
	fiber_favicon "github.com/gofiber/fiber/v2/middleware/favicon"
	// bolt "github.com/boltdb/bolt"
)

type Server struct {
	FiberApp *fiber.App `yaml:"fiber_app"`
	Config *types.ConfigFile `yaml:"config"`
	Location *time.Location `yaml:"-"`
}

var log *logrus.Logger

func ( s *Server ) LogRequest( context *fiber.Ctx ) ( error ) {
	// time_string := s.GetFormattedTimeString()
	ip_address := context.Get( "x-forwarded-for" )
	if ip_address == "" { ip_address = context.IP() }
	log_message := fmt.Sprintf( "%s === %s === %s" , ip_address , context.Method() , context.Path() )
	// fmt.Println( log_message )
	log.Info( log_message )
	// TODO : append to bolt ?
	return context.Next()
}

func ( s *Server ) Start() {
	fmt.Printf( "Admin Login @ http://localhost:%s/%s/%s\n" , s.Config.ServerPort , s.Config.ServerUrlPrefix , s.Config.ServerLoginUrlPrefix )
	fmt.Printf( "Admin Username === %s\n" , s.Config.AdminUsername )
	fmt.Printf( "Admin Password === %s\n" , s.Config.AdminPassword )
	fmt.Printf( "Admin API Key === %s\n" , s.Config.ServerAPIKey )
	local_ip_addresses := utils.GetLocalIPAddresses()
	for _ , ip_address := range local_ip_addresses {
		fmt.Sprintf( "Listening @ http://%s:%s\n" , ip_address , s.Config.ServerPort )
	}
	listen_address := fmt.Sprintf( ":%s" , s.Config.ServerPort )
	log.Info( fmt.Sprintf( "Listening @ %s" , listen_address ) )
	s.FiberApp.Listen( listen_address )
}

func New( config *types.ConfigFile ) ( server Server ) {
	server.Location , _ = time.LoadLocation( config.TimeZone )
	server.FiberApp = fiber.New()
	server.Config = config
	log = logger.New( config )
	server.FiberApp.Use( server.LogRequest )
	server.FiberApp.Use( fiber_favicon.New() )
	server.FiberApp.Use( fiber_cookie.New( fiber_cookie.Config{
		Key: server.Config.ServerCookieSecret ,
	}))
	server.FiberApp.Use( fiber_cors.New( fiber_cors.Config{
		AllowOrigins: config.ServerAllowOriginsString ,
		AllowHeaders:  "Origin, Content-Type, Accept, key, k" ,
	}))
	server.SetupPublicRoutes()
	server.SetupAdminRoutes()
	return
}