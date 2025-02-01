package logger

import (
	"os"
	"fmt"
	"embed"
	"strings"
	"io"
	"time"
	// "encoding/json"
	binary "encoding/binary"
	bolt_api "github.com/boltdb/bolt"
	types "github.com/0187773933/Logger/v1/types"
	encryption "github.com/0187773933/encryption/v1/encryption"
	// utils "github.com/0187773933/Logger/v1/utils"
	logrus "github.com/sirupsen/logrus"
	// ulid "github.com/oklog/ulid/v2"
)

//go:embed zoneinfo
var ZoneInfoFS embed.FS

// var Log *logrus.Logger
var Log *Wrapper
var Config *types.ConfigFile
var Location *time.Location
var DB *bolt_api.DB
var LogKeyBytes []byte
var Encrypting bool

type Wrapper struct {
	*logrus.Logger
}

func Get_location( name string ) ( *time.Location ) {
	return Location
}

func SetLocation( location_string string ) {
	if location_string == "" { location_string = "America/New_York" }
	bs , err := ZoneInfoFS.ReadFile( "zoneinfo/" + location_string )
	if err != nil { panic( err ) }
	loc , err := time.LoadLocationFromTZData( location_string , bs )
	if err != nil { panic( err ) }
	Location = loc
	fmt.Println( "location set to: " , location_string )
}

func FormatTime( input_time *time.Time ) ( result string ) {
	time_object := input_time.In( Location )
	month_name := strings.ToUpper( time_object.Format( "Jan" ) )
	milliseconds := time_object.Format( ".000" )
	date_part := fmt.Sprintf( "%02d%s%d" , time_object.Day() , month_name , time_object.Year() )
	time_part := fmt.Sprintf( "%02d:%02d:%02d%s" , time_object.Hour() , time_object.Minute() , time_object.Second() , milliseconds )
	result = fmt.Sprintf( "%s === %s" , date_part , time_part )
	return
}

// https://github.com/boltdb/bolt#autoincrementing-integer-for-the-bucket
// itob returns an 8-byte big endian representation of v.
func ItoB( v uint64 ) []byte {
    b := make( []byte , 8 )
    binary.BigEndian.PutUint64( b , uint64( v ) )
    return b
}

type CustomTextFormatter struct {
	logrus.TextFormatter
}

type CustomLogrusWriter struct {
	io.Writer
}

type CustomJSONFormatter struct {
	logrus.JSONFormatter
}

func ( f *CustomJSONFormatter ) Format( entry *logrus.Entry ) ( []byte , error ) {
	time_string := FormatTime( &entry.Time )
	fmt.Println( time_string )
	fmt.Println( entry )
	return f.JSONFormatter.Format( entry )
}

// https://github.com/sirupsen/logrus/blob/v1.9.3/entry.go#L44
// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
func ( f *CustomTextFormatter ) Format( entry *logrus.Entry ) ( result_bytes []byte , result_error error ) {
	time_string := FormatTime( &entry.Time )
	// result_bytes , result_error = f.TextFormatter.Format( entry )

	var result_string string
	if entry.Caller != nil {
		var caller_function string
		// test_parts := strings.Split( entry.Caller.Function , "github.com/0187773933/FireC2Server/v1/" )
		// if len( test_parts ) > 1 {
		// 	caller_function = test_parts[ 1 ]
		// } else {
		// 	caller_function = entry.Caller.Function
		// }
        function_parts := strings.Split( entry.Caller.Function , "/" )
        if len( function_parts ) > 0 {
            caller_function = function_parts[ len( function_parts ) - 1 ]
        } else {
            caller_function = entry.Caller.Function
        }
		result_string = fmt.Sprintf( "%s === %s():%d === %s" , time_string , caller_function , entry.Caller.Line , entry.Message )
	} else {
		result_string = fmt.Sprintf( "%s === %s" , time_string , entry.Message )
	}
	result_bytes = []byte( result_string )
	result_error = nil

	// DB.Update( func( tx *bolt_api.Tx ) error {
	// 	b_logs := tx.Bucket( []byte( "logs" ) )
	// 	b_today , _ := b_logs.CreateBucketIfNotExists( []byte( db_log_prefix ) )
	// 	b_today.Put( []byte( ulid_prefix ) , message_bytes )
	// 	return nil
	// })

	// message := &CustomLogMessage{
	// 	Message: result_string ,
	// 	Fields: entry.Data ,
	// 	Time: time_string ,
	// 	Level: entry.Level.String() ,
	// }
	// if entry.Caller != nil {
	// 	message.Frame = CustomLogMessageFrame{
	// 		// Function: entry.Caller.Function ,
	// 		File: entry.Caller.File ,
	// 		Line: entry.Caller.Line ,
	// 	}
	// }
	// db_log_prefix := utils.FormatDBLogPrefix( &entry.Time )
	// ulid_prefix := ulid.Make().String()
	// message_bytes , _ := json.Marshal( message )
	// DB.Update( func( tx *bolt_api.Tx ) error {
	// 	b_logs := tx.Bucket( []byte( "logs" ) )
	// 	b_today , _ := b_logs.CreateBucketIfNotExists( []byte( db_log_prefix ) )
	// 	b_today.Put( []byte( ulid_prefix ) , message_bytes )
	// 	return nil
	// })

	return result_bytes , result_error
}

func ( w *CustomLogrusWriter ) Write( p []byte ) ( n int , err error ) {
	message := string( p )
	// db_result := DB.Update( func( tx *bolt_api.Tx ) error {
	DB.Update( func( tx *bolt_api.Tx ) error {
		uuid_bucket := tx.Bucket( LogKeyBytes )
		sequence_id  , _ := uuid_bucket.NextSequence()
		sequence_id_b := ItoB( sequence_id )
		// fmt.Println( "next sequence id" , sequence_id )
		if Encrypting {
			p_e := encryption.ChaChaEncryptBytes( Config.EncryptionKey , p )
			uuid_bucket.Put( sequence_id_b , p_e )
		} else {
			uuid_bucket.Put( sequence_id_b , p )
		}
		return nil
	})
	// fmt.Println( db_result )
	n_message := message + "\n"
	n , err = fmt.Fprint( os.Stdout , n_message )
	return n , err
}

func ( w *Wrapper ) GetFormattedTimeString() ( result string ) {
	time_object := time.Now().In( Location )
	month_name := strings.ToUpper( time_object.Format( "Jan" ) )
	milliseconds := time_object.Format( ".000" )
	date_part := fmt.Sprintf( "%02d%s%d" , time_object.Day() , month_name , time_object.Year() )
	time_part := fmt.Sprintf( "%02d:%02d:%02d%s" , time_object.Hour() , time_object.Minute() , time_object.Second() , milliseconds )
	result = fmt.Sprintf( "%s === %s" , date_part , time_part )
	return
}

func ( w *Wrapper ) GetFormattedTimeStringOBJ() ( result_string string , result_time time.Time ) {
	result_time = time.Now().In( Location )
	month_name := strings.ToUpper( result_time.Format( "Jan" ) )
	milliseconds := result_time.Format( ".000" )
	date_part := fmt.Sprintf( "%02d%s%d" , result_time.Day() , month_name , result_time.Year() )
	time_part := fmt.Sprintf( "%02d:%02d:%02d%s" , result_time.Hour() , result_time.Minute() , result_time.Second() , milliseconds )
	result_string = fmt.Sprintf( "%s === %s" , date_part , time_part )
	return
}

func ( w *Wrapper ) FormatTime( input_time *time.Time ) ( result string ) {
	time_object := input_time.In( Location )
	month_name := strings.ToUpper( time_object.Format( "Jan" ) )
	milliseconds := time_object.Format( ".000" )
	date_part := fmt.Sprintf( "%02d%s%d" , time_object.Day() , month_name , time_object.Year() )
	time_part := fmt.Sprintf( "%02d:%02d:%02d%s" , time_object.Hour() , time_object.Minute() , time_object.Second() , milliseconds )
	result = fmt.Sprintf( "%s === %s" , date_part , time_part )
	return
}

// so apparently The limitation arises due to the Go language's initialization order:
// Package-level variables are initialized before main() is called.
// Functions in main() execute after package-level initializations.
// something something , singleton
// func GetLogger( config *types.ConfigFile ) *logrus.Logger {
func New( config *types.ConfigFile ) *Wrapper {
	Config = config
	if Log == nil { Init() }
	return Log
}

func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			Log.Fatalf( "Failed to close database: %v" , err )
		}
	}
}

func GetDB() ( *bolt_api.DB ) {
	if DB != nil {
		return DB
	} else {
		return nil
	}
}

func ( w *Wrapper ) GetMessages( count int ) ( messages []string ) {
	DB.View( func( tx *bolt_api.Tx ) error {
		uuid_bucket := tx.Bucket( LogKeyBytes )
		if uuid_bucket == nil {
			return fmt.Errorf( "Bucket not found" )
		}
		c := uuid_bucket.Cursor()
		// Iterate in reverse order
		for k , v := c.Last(); k != nil && ( count != 0 ); k, v = c.Prev() {
			if Encrypting {
				decrytped_message := encryption.ChaChaDecryptBytes( Config.EncryptionKey , v )
				messages = append( messages , string( decrytped_message ) )
			} else {
				messages = append( messages , string( v ) )
			}
			if count > 0 {
				count--
			}
		}
		return nil
	})
	// Reverse the slice to return messages in the original order
	// for i, j := 0, len( messages ) - 1; i < j; i , j = i + 1, j - 1 {
	// 	messages[ i ] , messages[ j ] = messages[ j ] , messages[ i ]
	// }
	return
}

func Init() {
	if Log != nil { return }
	// Log = logrus.New()
	Log = &Wrapper{logrus.New()}
	Location , _ = time.LoadLocation( Config.TimeZone )
	LogKeyBytes = []byte( Config.LogKey )
	db , db_open_error := bolt_api.Open( Config.BoltDBPath , 0600 , &bolt_api.Options{ Timeout: ( 3 * time.Second ) } )
	if db_open_error != nil { Log.Fatalf( "Failed to open database: %v" , db_open_error ); }
	DB = db
	DB.Update( func( tx *bolt_api.Tx ) error {
		tx.CreateBucketIfNotExists( LogKeyBytes )
		return nil
	})
	if Config.EncryptionKey != "" {
		Encrypting = true
	} else {
		Encrypting = false
	}
	// log_level := os.Getenv( "LOG_LEVEL" )
	// fmt.Printf( "LOG_LEVEL=%s\n" , Config.LogLevel )
	switch Config.LogLevel {
		case "debug":
			Log.SetReportCaller( true )
			Log.SetLevel( logrus.DebugLevel )
		default:
			Log.SetReportCaller( false )
			Log.SetLevel( logrus.InfoLevel )
	}
	Log.SetFormatter( &CustomTextFormatter{
		TextFormatter: logrus.TextFormatter{
			DisableColors: false ,
		} ,
	})
	// log.SetFormatter( &CustomJSONFormatter{
	// 	JSONFormatter: logrus.JSONFormatter{} ,
	// })

	// log.SetOutput( os.Stdout )
	Log.SetOutput( &CustomLogrusWriter{} )
}