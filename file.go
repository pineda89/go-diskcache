package diskcache

import "os"

func (dc *DiskCache) keyToFilepath(key string) string {
	return dc.Options.folder + string(os.PathSeparator) + key
}

func (dc *DiskCache) filenameToKey(filename string) string {
	return filename
}

func (dc *DiskCache) removeFile(filepath string) error {
	return os.RemoveAll(filepath)
}
