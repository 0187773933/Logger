package utils

import (
	"fmt"
	"os"
	"net"
	yaml "gopkg.in/yaml.v3"
	ioutil "io/ioutil"
	runtime "runtime"
	types "github.com/0187773933/Logger/v1/types"
	fiber_cookie "github.com/gofiber/fiber/v2/middleware/encryptcookie"
	encryption "github.com/0187773933/encryption/v1/encryption"
)

func SetupStackTraceReport() {
	if r := recover(); r != nil {
		stacktrace := make( []byte , 1024 )
		runtime.Stack( stacktrace , true )
		fmt.Printf( "%s\n" , stacktrace )
	}
}

func GenerateNewKeys() {
	fiber_cookie_key := fiber_cookie.GenerateKey()
	encryption_key := encryption.GenerateRandomString( 32 )
	server_api_key := encryption.GenerateRandomString( 16 )
	admin_username := encryption.GenerateRandomString( 16 )
	admin_password := encryption.GenerateRandomString( 16 )
	login_url := encryption.GenerateRandomString( 16 )
	url_prefix := encryption.GenerateRandomString( 6 )
	fmt.Println( "Generated New Keys :" )
	fmt.Printf( "\tFiber Cookie Key === %s\n" , fiber_cookie_key )
	fmt.Printf( "\tEncryption Key === %s\n" , encryption_key )
	fmt.Printf( "\tServer API Key === %s\n" , server_api_key )
	fmt.Printf( "\tAdmin Username === %s\n" , admin_username )
	fmt.Printf( "\tAdmin Password === %s\n" , admin_password )
	fmt.Printf( "\tLogin URL === %s\n" , login_url )
	fmt.Printf( "\tURL Prefix === %s\n" , url_prefix )
	panic( "Exiting" )
}

func ParseConfig( file_path string ) ( result types.ConfigFile ) {
	config_file , _ := ioutil.ReadFile( file_path )
	error := yaml.Unmarshal( config_file , &result )
	if error != nil { panic( error ) }
	return
}

func GetLocalIPAddresses() ( ip_addresses []string ) {
	host , _ := os.Hostname()
	addrs , _ := net.LookupIP( host )
	encountered := make( map[ string ]bool )
	for _ , addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			ip := ipv4.String()
			if !encountered[ ip ] {
				encountered[ ip ] = true
				ip_addresses = append( ip_addresses , ip )
			}
		}
	}
	return
}