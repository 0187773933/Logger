package server

import (
	"fmt"
	"strconv"
	net_url "net/url"
	bolt_api "github.com/boltdb/bolt"
	encryption "github.com/0187773933/encryption/v1/encryption"
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

func ( s *Server ) LogMessageWithKey( c *fiber.Ctx ) ( error ) {
	key := c.Params( "key" )
	message := c.Params( "message" )
	if message == "" {
		return c.JSON( fiber.Map{
			"result": false ,
			"error": "No Message Provided" ,
		})
	}
	unescaped_message , _ := net_url.QueryUnescape( message )
	db := logger.GetDB()
	var sequence_id uint64
	db.Update( func( tx *bolt_api.Tx ) error {
		bucket , _ := tx.CreateBucketIfNotExists( []byte( key ) )
		sequence_id , _ = bucket.NextSequence()
		if s.Config.EncryptionKey != "" {
			encrypted_message_bytes := encryption.ChaChaEncryptBytes( s.Config.EncryptionKey , []byte( unescaped_message ) )
			bucket.Put( logger.ItoB( sequence_id ) , encrypted_message_bytes )
		} else {
			bucket.Put( logger.ItoB( sequence_id ) , []byte( unescaped_message ) )
		}
		return nil
	})
	return c.JSON( fiber.Map{
		"url": "/log" ,
		"key": key ,
		"sequence_id": sequence_id ,
		"message": unescaped_message ,
		"result": true ,
	})
}

func ( s *Server ) GetLogMessagesWithKey( c *fiber.Ctx ) ( error ) {
	key := c.Params( "key" )
	count := c.Params( "count" )
	count_int , _ := strconv.Atoi( count )
	db := logger.GetDB()
	var  messages []string
	db.View( func( tx *bolt_api.Tx ) error {
		bucket := tx.Bucket( []byte( key ) )
		if bucket == nil {
			return fmt.Errorf( "Bucket not found" )
		}
		c := bucket.Cursor()
		for k , v := c.Last(); k != nil && count_int > 0; k, v = c.Prev() {
			if s.Config.EncryptionKey != "" {
				decrytped_message_bytes := encryption.ChaChaDecryptBytes( s.Config.EncryptionKey , v )
				messages = append( messages , string( decrytped_message_bytes ) )
			} else {
				messages = append( messages , string( v ) )
			}
			count_int--
		}
		return nil
	})
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
	admin.Get( "/log/:key/:message" , s.LogMessageWithKey )
	admin.Get( "/log/:key/view/:count" , s.GetLogMessagesWithKey )
}