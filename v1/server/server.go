package server

import (
	"fmt"
	"time"
	"strings"
	// logrus "github.com/sirupsen/logrus"
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

var log *logger.Wrapper

func ( s *Server ) LogRequest( context *fiber.Ctx ) ( error ) {
	ip_address := context.Get( "x-forwarded-for" )
	if ip_address == "" { ip_address = context.IP() }
	c_method := context.Method()
	c_path := context.Path()
	// c_route := context.Route()
	// fmt.Println( c_route.Path )
	// fmt.Println( c_method )
	// fmt.Println( c_path )
	// c_handlers := c_route.Handlers
	if strings.HasPrefix( c_path , fmt.Sprintf( "/%s/favicon" , s.Config.ServerUrlPrefix ) ) {
		return context.Next()
	}
	// avoids double store in global log
	if strings.HasPrefix( c_path , fmt.Sprintf( "/%s/log/c/" , s.Config.ServerUrlPrefix ) ) {
		if strings.Contains( c_path , "/view" ) == false {
			time_string := s.GetFormattedTimeString()
			log_message := fmt.Sprintf( "%s === %s === %s === %s , skipping global log storage" , time_string , ip_address , c_method , c_path )
			fmt.Println( log_message )
			return context.Next()
		}
	}
	log_message := fmt.Sprintf( "%s === %s === %s" , ip_address , c_method , c_path )
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
	test := log.GetFormattedTimeString()
	fmt.Println( "test" , test )
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