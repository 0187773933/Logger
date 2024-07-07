# Logger

```
package main

import (
	logger "github.com/0187773933/Logger/v1/logger"
	logger_types "github.com/0187773933/Logger/v1/types"
)

func main() {
	log_config := logger_types.ConfigFile{
		LogLevel: "debug" ,
		TimeZone: "America/New_York" ,
		EncryptionKey: "asdf" ,
		BoltDBPath: "bolt.db" ,
	}
	log := logger.New( &log_config )
	log.Debug( "asdf" )
}
```

```
webpack --config webpack.config.js
```