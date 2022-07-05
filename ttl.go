package diskcache

import "time"

func (dc *DiskCache) ttlCheck() {
	for !dc.closed {
		toRemove := make([]*kvdata, 0)
		dc.keys.RLock()
		for _, v := range dc.keys.keysList {
			if v.duration != NoExpiration {
				expiration := v.putTime
				if v.duration == DefaultExpiration {
					if dc.Options.expiration == NoExpiration {
						continue
					}
					expiration = expiration.Add(dc.Options.expiration)
				} else {
					expiration = expiration.Add(v.duration)
				}

				if expiration.Before(time.Now()) {
					toRemove = append(toRemove, v)
				}
			}
		}
		dc.keys.RUnlock()

		for _, v := range toRemove {
			dc.Delete(v.key)
		}

		time.Sleep(dc.Options.checkTTLPeriod)
	}
}
