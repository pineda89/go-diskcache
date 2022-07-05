### go-diskcache

go-diskcache is a key-value cache 100% on disk, with TTL expiration implementation, and with disk usage limits

## Usage

```go
package main

import (
	"github.com/pineda89/go-diskcache"
	"log"
	"time"
)

func main() {
	dc := diskcache.New(diskcache.DefaultOptions().
		WithExpiration(1 * time.Hour).
		WithCompression(diskcache.CompressionNone).
		WithMaxUsagePercent(66.65))

	key := "mykey"
	originalData := "data here please"

	dc.Set(key, []byte(originalData), diskcache.DefaultExpiration)

	data, err := dc.Get(key)
	log.Println("result:", string(data), err)

	dc.Delete(key)

	data, err = dc.Get(key)
	log.Println("result after delete:", string(data), err)
}
```

## Options

| Function                 | Default value     | Description                                                                                         |
|--------------------------|-------------------|-----------------------------------------------------------------------------------------------------|
| WithFolder               | /diskcache        | folder where the cache will reside. Can be absolute or relative                                     | 
| WithExpiration           | 1 hour            | when making a Set, we can indicate a default TTL. This will be the TTL                              |
 | WithCompression          | gzip              | compression of files on disk. Posibles values are CompressionNone and CompressionGzip               |
| WithMaxUsagePercent      | 90                | maximum usage percentage. If it is left at 0 it is not used                                         |
| WithMinFreeSpace         | 0                 | minimum number of bytes that must be free. By default it is not used                                |
| WithEvictionPolicy       | RemoveOldestFirst | delete criterion in case of reaching the disk limit                                                 |
| WithCheckOrphansPeriod   | 1 second          | how often it checks for orphaned files. An orphaned file is a file that is not covered by the cache |
| WithCheckEvictionsPeriod | 1 second          | how often it checks if we have exceeded the disk limits                                             |
| WithCheckTTLPeriod       | 1 second          | how often it is checked if the keys exceeded the TTL                                                |