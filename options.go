package diskcache

import "time"

type Options struct {
	folder               string
	compression          compression
	expiration           time.Duration
	evictionPolicy       evictionPolicy
	maxUsagePercent      float64
	minFreeSpace         uint64
	checkOrphansPeriod   time.Duration
	checkEvictionsPeriod time.Duration
	checkTTLPeriod       time.Duration
}

func DefaultOptions() *Options {
	return new(Options).
		WithFolder("/diskcache").
		WithCompression(CompressionGzip).
		WithExpiration(1 * time.Hour).
		WithEvictionPolicy(EvictionpolicyRemoveOldestFirst).
		WithMaxUsagePercent(90).
		WithMinFreeSpace(0).
		WithCheckOrphansPeriod(1 * time.Minute).
		WithCheckEvictionsPeriod(1 * time.Minute).
		WithCheckTTLPeriod(1 * time.Second)
}

func (opts *Options) WithFolder(input string) *Options {
	opts.folder = input
	return opts
}

func (opts *Options) WithExpiration(input time.Duration) *Options {
	opts.expiration = input
	return opts
}

func (opts *Options) WithCompression(input compressionType) *Options {
	opts.compression = getCompression(input)
	return opts
}

func (opts *Options) WithMaxUsagePercent(input float64) *Options {
	opts.maxUsagePercent = input
	return opts
}

func (opts *Options) WithMinFreeSpace(input uint64) *Options {
	opts.minFreeSpace = input
	return opts
}

func (opts *Options) WithEvictionPolicy(input evictionPolicy) *Options {
	opts.evictionPolicy = input
	return opts
}

func (opts *Options) WithCheckOrphansPeriod(input time.Duration) *Options {
	opts.checkOrphansPeriod = input
	return opts
}

func (opts *Options) WithCheckEvictionsPeriod(input time.Duration) *Options {
	opts.checkEvictionsPeriod = input
	return opts
}

func (opts *Options) WithCheckTTLPeriod(input time.Duration) *Options {
	opts.checkTTLPeriod = input
	return opts
}
