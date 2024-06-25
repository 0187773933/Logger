package server

import (
	"fmt"
	"time"
	"strings"
	fiber "github.com/gofiber/fiber/v2"
	rate_limiter "github.com/gofiber/fiber/v2/middleware/limiter"
	bcrypt "golang.org/x/crypto/bcrypt"
	encryption "github.com/0187773933/encryption/v1/encryption"
	// try "github.com/manucorporat/try"
)

var CDNLimter = rate_limiter.New( rate_limiter.Config{
	Max: 6 ,
	Expiration: 1 * time.Second ,
	KeyGenerator: func( c *fiber.Ctx ) string {
		return c.Get( "x-forwarded-for" )
	} ,
	LimitReached: func( c *fiber.Ctx ) error {
		ip_address := c.IP()
		log_message := fmt.Sprintf( "%s === %s === %s === PUBLIC RATE LIMIT REACHED !!!" , ip_address , c.Method() , c.Path() );
		log.Info( log_message )
		c.Set( "Content-Type" , "text/html" )
		return c.SendString( "<html><h1>loading ...</h1><script>setTimeout(function(){ window.location.reload(1); }, 6000);</script></html>" )
	} ,
})

var PublicLimter = rate_limiter.New( rate_limiter.Config{
	Max: 3 ,
	Expiration: 1 * time.Second ,
	KeyGenerator: func( c *fiber.Ctx ) string {
		return c.Get( "x-forwarded-for" )
	} ,
	LimitReached: func( c *fiber.Ctx ) error {
		ip_address := c.IP()
		log_message := fmt.Sprintf( "%s === %s === %s === PUBLIC RATE LIMIT REACHED !!!" , ip_address , c.Method() , c.Path() );
		log.Info( log_message )
		c.Set( "Content-Type" , "text/html" )
		return c.SendString( "<html><h1>loading ...</h1><script>setTimeout(function(){ window.location.reload(1); }, 6000);</script></html>" )
	} ,
})

func ( s *Server ) ValidateAdmin( context *fiber.Ctx ) ( result bool ) {
	result = false
	admin_cookie := context.Cookies( s.Config.ServerCookieName )
	if admin_cookie != "" {
		admin_cookie_value := encryption.SecretBoxDecrypt( s.Config.EncryptionKey , admin_cookie )
		if admin_cookie_value == s.Config.ServerCookieAdminSecretMessage {
			result = true
			return
		}
	}
	admin_api_key_header := context.Get( "key" )
	if admin_api_key_header != "" {
		if admin_api_key_header == s.Config.ServerAPIKey {
			result = true
			return
		}
	}
	admin_api_key_query := context.Query( "k" )
	if admin_api_key_query != "" {
		if admin_api_key_query == s.Config.ServerAPIKey {
			result = true
			return
		}
	}
	return
}

func ( s *Server ) ValidateAdminMW( context *fiber.Ctx ) ( error ) {
	admin_cookie := context.Cookies( s.Config.ServerCookieName )
	if admin_cookie != "" {
		admin_cookie_value := encryption.SecretBoxDecrypt( s.Config.EncryptionKey , admin_cookie )
		if admin_cookie_value == s.Config.ServerCookieAdminSecretMessage {
			return context.Next()
		}
	}
	admin_api_key_header := context.Get( "key" )
	if admin_api_key_header != "" {
		if admin_api_key_header == s.Config.ServerAPIKey {
			return context.Next()
		}
	}
	admin_api_key_query := context.Query( "k" )
	if admin_api_key_query != "" {
		if admin_api_key_query == s.Config.ServerAPIKey {
			return context.Next()
		}
	}
	return context.Status( fiber.StatusUnauthorized ).SendString( "why" )
}

func ( s *Server ) ValidateLoginCredentials( context *fiber.Ctx ) ( result bool ) {
	result = false
	uploaded_username := context.FormValue( "username" )
	if uploaded_username == "" { fmt.Println( "username empty" ); return }
	if uploaded_username != s.Config.AdminUsername { fmt.Println( "username not correct" ); return }
	uploaded_password := context.FormValue( "password" )
	if uploaded_password == "" { fmt.Println( "password empty" ); return }
	fmt.Println( "uploaded_username ===" , uploaded_username )
	fmt.Println( "uploaded_password ===" , uploaded_password )
	password_matches := bcrypt.CompareHashAndPassword( []byte( uploaded_password ) , []byte( s.Config.AdminPassword ) )
	if password_matches != nil { fmt.Println( "bcrypted password doesn't match" ); return }
	fmt.Println( "password matched" )
	result = true
	return
}

// POST http://localhost:5950/admin/login
func ( s *Server ) HandleLogin( context *fiber.Ctx ) ( error ) {
	valid_login := s.ValidateLoginCredentials( context )
	if valid_login == false { return s.RenderFailedLogin( context ) }
	host := context.Hostname()
	domain := strings.Split( host , ":" )[ 0 ] // setting this leaks url-prefix and locks to specific domain
	context.Cookie(
		&fiber.Cookie{
			Name: s.Config.ServerCookieName ,
			Value: encryption.SecretBoxEncrypt( s.Config.EncryptionKey , s.Config.ServerCookieAdminSecretMessage ) ,
			Secure: true ,
			Path: "/" ,
			// Domain: "blah.ngrok.io" , // probably should set this for webkit
			Domain: domain ,
			HTTPOnly: true ,
			SameSite: "Lax" ,
			Expires: time.Now().AddDate( 10 , 0 , 0 ) , // aka 10 years from now
		} ,
	)
	return context.Redirect( "/" )
}

func ( s *Server ) HandleLogout( context *fiber.Ctx ) ( error ) {
	context.Cookie( &fiber.Cookie{
		Name: s.Config.ServerCookieName ,
		Value: "" ,
		Expires: time.Now().Add( -time.Hour ) , // set the expiration to the past
		HTTPOnly: true ,
		Secure: true ,
	})
	context.Set( "Content-Type" , "text/html" )
	return context.SendString( "<h1>Logged Out</h1>" )
}

func ( s *Server ) RenderFailedLogin( context *fiber.Ctx ) ( error ) {
	context.Set( "Content-Type" , "text/html" )
	return context.SendString( "<h1>no</h1>" )
}

func ( s *Server ) RenderLoginPage( context *fiber.Ctx ) ( error ) {
	context.Set( "Content-Type" , "text/html" )
	return context.SendFile( "./v1/server/html/login.html" )
}

func ( s *Server ) RenderHomePage( context *fiber.Ctx ) ( error ) {
	context.Set( "Content-Type" , "text/html" )
	admin_logged_in := s.ValidateAdmin( context )
	if admin_logged_in == true {
		// fmt.Println( "RenderHomePage() --> Admin" )
		return context.SendFile( "./v1/server/html/admin.html" )
	}
	return context.SendFile( "./v1/server/html/home.html" )
}

func ( s *Server ) SetupPublicRoutes() {
	cdn_group := s.FiberApp.Group( "/cdn" )
	cdn_group.Use( CDNLimter )
	s.FiberApp.Static( "/cdn" , "./v1/server/cdn" )
	s.FiberApp.Get( "/" , PublicLimter , s.RenderHomePage )
	s.FiberApp.Get( fmt.Sprintf( "/%s/%s" , s.Config.ServerUrlPrefix , s.Config.ServerLoginUrlPrefix ) , PublicLimter , s.RenderLoginPage )
	s.FiberApp.Post( fmt.Sprintf( "/%s/%s" , s.Config.ServerUrlPrefix , s.Config.ServerLoginUrlPrefix ) , PublicLimter , s.HandleLogin )
	s.FiberApp.Get( fmt.Sprintf( "/%s/logout" , s.Config.ServerUrlPrefix ) , PublicLimter , s.HandleLogout )
}