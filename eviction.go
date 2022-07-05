package diskcache

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type evictionPolicy byte

const (
	EvictionpolicyNone evictionPolicy = iota
	EvictionpolicyRemoveOldestFirst
)

func (dc *DiskCache) applyEviction(toFree uint64) {
	switch dc.Options.evictionPolicy {
	case EvictionpolicyRemoveOldestFirst:
		dc.applyEvictionOldestFirst(toFree)
	default:
		panic("not eviction policy found")
	}
}

func (dc *DiskCache) applyEvictionOldestFirst(toFree uint64) {
	resto := int64(toFree)
	for resto > 0 {
		var f *kvdata
		dc.keys.RWMutex.RLock()
		if len(dc.keys.keysList) > 0 {
			f = dc.keys.keysList[0]
		}
		dc.keys.RWMutex.RUnlock()

		if f != nil {
			dc.Delete(f.key)
			resto = resto - int64(f.fileSize)
		} else {
			break
		}
	}
}

func (dc *DiskCache) checkOrphans() {
	for !dc.closed {
		files, err := ioutil.ReadDir(dc.Options.folder)
		if err != nil {
			log.Println("err checking orphans", err)
		} else {
			toDelete := make([]fs.FileInfo, 0)
			dc.keys.RLock()
			for _, f := range files {
				if !f.IsDir() {
					key := dc.filenameToKey(f.Name())
					if _, ok := dc.keys.keysMap[key]; !ok {
						toDelete = append(toDelete, f)
					}
				}
			}
			dc.keys.RUnlock()

			for i := range toDelete {
				dc.removeFile(dc.Options.folder + string(os.PathSeparator) + toDelete[i].Name())
			}
		}
		time.Sleep(dc.Options.checkOrphansPeriod)
	}
}
