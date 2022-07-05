package diskcache

import (
	"github.com/shirou/gopsutil/v3/disk"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

const (
	NoExpiration      time.Duration = -1
	DefaultExpiration time.Duration = 0
)

type DiskCache struct {
	Options    *Options
	keys       *keys
	closed     bool
	usageStats *disk.UsageStat
}

type keys struct {
	sync.RWMutex
	keysList []*kvdata
	keysMap  map[string]*kvdata
}

type kvdata struct {
	key      string
	fileSize int
	putTime  time.Time
	duration time.Duration
}

func New(options *Options) *DiskCache {
	c := &DiskCache{
		Options: options,
		keys: &keys{
			keysMap: map[string]*kvdata{},
		},
	}

	go c.startCache()
	return c
}

func (dc *DiskCache) Close() {
	dc.closed = true
}

func (dc *DiskCache) Delete(key string) error {
	dc.keys.Lock()
	if kv, ok := dc.keys.keysMap[key]; ok {
		delete(dc.keys.keysMap, key)
		var indexToRemove = -1
		for i := 0; i < len(dc.keys.keysList); i++ {
			if dc.keys.keysList[i] == kv {
				indexToRemove = i
				break
			}
		}
		if indexToRemove != -1 {
			dc.keys.keysList = append(dc.keys.keysList[0:indexToRemove], dc.keys.keysList[indexToRemove+1:]...)
		}
	}
	dc.keys.Unlock()

	filepath := dc.keyToFilepath(key)
	return dc.removeFile(filepath)
}

func (dc *DiskCache) Set(key string, value []byte, d time.Duration) error {
	compressedValue, _ := dc.Options.compression.compress(value)

	dc.keys.Lock()
	kv := &kvdata{key: key, putTime: time.Now(), duration: d, fileSize: len(compressedValue)}
	dc.keys.keysList = append(dc.keys.keysList, kv)
	dc.keys.keysMap[key] = kv
	dc.keys.Unlock()

	filepath := dc.keyToFilepath(key)

	return ioutil.WriteFile(filepath, compressedValue, 0666)
}

func (dc *DiskCache) Get(key string) ([]byte, error) {
	dc.keys.RLock()
	kv := dc.keys.keysMap[key]
	dc.keys.RUnlock()
	if kv != nil {
		filepath := dc.keyToFilepath(key)
		fileContent, err := ioutil.ReadFile(filepath)
		if err != nil {
			return []byte{}, err
		}
		return dc.Options.compression.decompress(fileContent)
	}
	return []byte{}, err_not_in_cache
}

func (dc *DiskCache) startCache() {
	os.MkdirAll(dc.Options.folder, os.ModePerm)

	if (dc.Options.minFreeSpace > 0 || dc.Options.maxUsagePercent > 0) && dc.Options.evictionPolicy != EvictionpolicyNone {
		go dc.checkCacheLimitsForEviction()
	}

	go dc.checkOrphans()

	go dc.ttlCheck()
}

func (dc *DiskCache) checkCacheLimitsForEviction() {
	for !dc.closed {
		var err error
		dc.usageStats, err = disk.Usage(dc.Options.folder)
		for dc.usageStats, err = disk.Usage(dc.Options.folder); err != nil || (dc.usageStats.Free < dc.Options.minFreeSpace && dc.Options.minFreeSpace > 0) || (dc.usageStats.UsedPercent > dc.Options.maxUsagePercent && dc.Options.maxUsagePercent > 0); dc.usageStats, err = disk.Usage(dc.Options.folder) {
			if err != nil {
				log.Println("read disk usage err", err)
				break
			} else {
				var bytesToRelease uint64
				if dc.usageStats.Free < dc.Options.minFreeSpace && dc.Options.minFreeSpace > 0 {
					if nBytesToRelease := dc.Options.minFreeSpace - dc.usageStats.Free; bytesToRelease < nBytesToRelease {
						bytesToRelease = nBytesToRelease
					}
				}
				if dc.usageStats.UsedPercent > dc.Options.maxUsagePercent && dc.Options.maxUsagePercent > 0 {
					expectedMaxUsage := uint64(float64(dc.usageStats.Total) / 100 * dc.Options.maxUsagePercent)

					if nBytesToRelease := dc.usageStats.Used - expectedMaxUsage; bytesToRelease < nBytesToRelease {
						bytesToRelease = nBytesToRelease
					}
				}

				dc.applyEviction(bytesToRelease)
			}
		}
		time.Sleep(dc.Options.checkEvictionsPeriod)
	}
}
