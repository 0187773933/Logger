package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"path/filepath"
	logger "github.com/0187773933/Logger/v1/logger"
	utils "github.com/0187773933/Logger/v1/utils"
	server "github.com/0187773933/Logger/v1/server"
)

var s server.Server

func SetupCloseHandler() {
	c := make( chan os.Signal )
	signal.Notify( c , os.Interrupt , syscall.SIGTERM , syscall.SIGINT )
	go func() {
		<-c
		logger.Log.Println( "\r- Ctrl+C pressed in Terminal" )
		// fmt.Println( "\r" )
		logger.Log.Printf( "Shutting Down %s Server" , s.Config.ServerName )
		s.FiberApp.Shutdown()
		logger.CloseDB()
		os.Exit( 0 )
	}()
}

func main() {
	// utils.GenerateNewKeys()
	defer utils.SetupStackTraceReport()
	var config_file_path string
	if len( os.Args ) > 1 {
		config_file_path , _ = filepath.Abs( os.Args[ 1 ] )
	} else {
		config_file_path , _ = filepath.Abs( "./config.yaml" )
		if _ , err := os.Stat( config_file_path ); os.IsNotExist( err ) {
			panic( "Config File Not Found" )
		}
	}
	config := utils.ParseConfig( config_file_path )
	fmt.Println( config_file_path , config )
	s = server.New( &config )
	logger.Log.Printf( "Loaded Config File From : %s" , config_file_path )
	SetupCloseHandler()
	s.Start()
}