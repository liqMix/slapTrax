package assets

import (
	"sync"
	"time"

	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/types"
)

// SongCache provides caching for parsed songs to improve loading performance
type SongCache struct {
	mu      sync.RWMutex
	cache   map[string]*CacheEntry
	maxSize int
	maxAge  time.Duration
	hits    int64
	misses  int64
}

// CacheEntry represents a cached song with metadata
type CacheEntry struct {
	Song        *types.Song
	Hash        string
	Size        int64
	CreatedAt   time.Time
	AccessedAt  time.Time
	AccessCount int64
}

// NewSongCache creates a new song cache with specified limits
func NewSongCache(maxSize int, maxAge time.Duration) *SongCache {
	return &SongCache{
		cache:   make(map[string]*CacheEntry),
		maxSize: maxSize,
		maxAge:  maxAge,
	}
}

// Get retrieves a song from cache
func (sc *SongCache) Get(key string) (*types.Song, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	entry, exists := sc.cache[key]
	if !exists {
		sc.misses++
		return nil, false
	}

	// Check if entry is expired
	if sc.maxAge > 0 && time.Since(entry.CreatedAt) > sc.maxAge {
		sc.mu.RUnlock()
		sc.mu.Lock()
		delete(sc.cache, key)
		sc.mu.Unlock()
		sc.mu.RLock()
		sc.misses++
		return nil, false
	}

	// Update access information
	entry.AccessedAt = time.Now()
	entry.AccessCount++
	sc.hits++

	return entry.Song, true
}

// Put stores a song in cache
func (sc *SongCache) Put(key string, song *types.Song, hash string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Check if we need to evict entries
	if len(sc.cache) >= sc.maxSize {
		sc.evictLRU()
	}

	entry := &CacheEntry{
		Song:        song,
		Hash:        hash,
		Size:        sc.estimateSize(song),
		CreatedAt:   time.Now(),
		AccessedAt:  time.Now(),
		AccessCount: 1,
	}

	sc.cache[key] = entry
	logger.Debug("Cached song: %s (cache size: %d)", key, len(sc.cache))
}

// evictLRU removes the least recently used entry
func (sc *SongCache) evictLRU() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range sc.cache {
		if oldestKey == "" || entry.AccessedAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.AccessedAt
		}
	}

	if oldestKey != "" {
		delete(sc.cache, oldestKey)
		logger.Debug("Evicted song from cache: %s", oldestKey)
	}
}

// estimateSize estimates memory usage of a song
func (sc *SongCache) estimateSize(song *types.Song) int64 {
	size := int64(len(song.Title) + len(song.Artist) + len(song.Album))

	for _, chart := range song.Charts {
		size += int64(chart.TotalNotes * 128) // Rough estimate per note
		if chart.EventManager != nil {
			size += int64(chart.EventManager.GetEventCount() * 64) // Rough estimate per event
		}
	}

	return size
}

// Clear removes all entries from cache
func (sc *SongCache) Clear() {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.cache = make(map[string]*CacheEntry)
	sc.hits = 0
	sc.misses = 0
	logger.Debug("Cleared song cache")
}

// Stats returns cache statistics
func (sc *SongCache) Stats() CacheStats {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	var totalSize int64
	for _, entry := range sc.cache {
		totalSize += entry.Size
	}

	hitRate := float64(0)
	total := sc.hits + sc.misses
	if total > 0 {
		hitRate = float64(sc.hits) / float64(total) * 100
	}

	return CacheStats{
		Entries:   len(sc.cache),
		MaxSize:   sc.maxSize,
		TotalSize: totalSize,
		Hits:      sc.hits,
		Misses:    sc.misses,
		HitRate:   hitRate,
	}
}

// CacheStats provides cache performance metrics
type CacheStats struct {
	Entries   int
	MaxSize   int
	TotalSize int64
	Hits      int64
	Misses    int64
	HitRate   float64
}

// Global cache instance
var songCache *SongCache

// InitCache initializes the song cache
func InitCache() {
	songCache = NewSongCache(50, 30*time.Minute) // 50 songs, 30 min TTL
	logger.Debug("Initialized song cache")
}

// GetCachedSong retrieves a song from the global cache
func GetCachedSong(key string) (*types.Song, bool) {
	if songCache == nil {
		InitCache()
	}
	return songCache.Get(key)
}

// CacheSong stores a song in the global cache
func CacheSong(key string, song *types.Song, hash string) {
	if songCache == nil {
		InitCache()
	}
	songCache.Put(key, song, hash)
}

// GetCacheStats returns global cache statistics
func GetCacheStats() CacheStats {
	if songCache == nil {
		return CacheStats{}
	}
	return songCache.Stats()
}

// StreamingLoader provides streaming/lazy loading capabilities
type StreamingLoader struct {
	chunkSize  int
	bufferSize int
}

// NewStreamingLoader creates a streaming loader
func NewStreamingLoader(chunkSize, bufferSize int) *StreamingLoader {
	return &StreamingLoader{
		chunkSize:  chunkSize,
		bufferSize: bufferSize,
	}
}

// LoadSongLazy loads a song with lazy initialization
func (sl *StreamingLoader) LoadSongLazy(folderName string) (*types.Song, error) {
	// Check cache first
	if song, found := GetCachedSong(folderName); found {
		logger.Debug("Loaded song from cache: %s", folderName)
		return song, nil
	}

	// Load song normally
	song, err := LoadSongJSON(folderName)
	if err != nil {
		return nil, err
	}

	// Cache the result
	CacheSong(folderName, song, song.Hash)

	return song, nil
}

// CompressedSongData represents compressed song data for storage
type CompressedSongData struct {
	Data         []byte
	OriginalSize int64
	Compressed   bool
	Format       string
}

// CompressionManager handles song data compression
type CompressionManager struct {
	enabled     bool
	threshold   int64  // Compress files larger than this
	compression string // Compression algorithm
}

// NewCompressionManager creates a compression manager
func NewCompressionManager(enabled bool, threshold int64) *CompressionManager {
	return &CompressionManager{
		enabled:     enabled,
		threshold:   threshold,
		compression: "gzip",
	}
}

// Compress compresses song data if beneficial
func (cm *CompressionManager) Compress(data []byte) *CompressedSongData {
	if !cm.enabled || int64(len(data)) < cm.threshold {
		return &CompressedSongData{
			Data:         data,
			OriginalSize: int64(len(data)),
			Compressed:   false,
			Format:       "raw",
		}
	}

	// TODO: Implement actual compression (gzip, etc.)
	// For now, return uncompressed
	return &CompressedSongData{
		Data:         data,
		OriginalSize: int64(len(data)),
		Compressed:   false,
		Format:       "raw",
	}
}

// ObjectPool provides object pooling for frequently allocated objects
type ObjectPool struct {
	notes  sync.Pool
	events sync.Pool
	charts sync.Pool
}

// NewObjectPool creates a new object pool
func NewObjectPool() *ObjectPool {
	return &ObjectPool{
		notes: sync.Pool{
			New: func() interface{} {
				return &types.Note{}
			},
		},
		events: sync.Pool{
			New: func() interface{} {
				return &types.BaseEvent{}
			},
		},
		charts: sync.Pool{
			New: func() interface{} {
				return &types.Chart{
					EventManager: types.NewEventManager(),
				}
			},
		},
	}
}

// GetNote gets a note from the pool
func (op *ObjectPool) GetNote() *types.Note {
	note := op.notes.Get().(*types.Note)
	note.Reset() // Ensure clean state
	return note
}

// PutNote returns a note to the pool
func (op *ObjectPool) PutNote(note *types.Note) {
	op.notes.Put(note)
}

// GetChart gets a chart from the pool
func (op *ObjectPool) GetChart() *types.Chart {
	chart := op.charts.Get().(*types.Chart)
	chart.EventManager.Clear() // Ensure clean state
	return chart
}

// PutChart returns a chart to the pool
func (op *ObjectPool) PutChart(chart *types.Chart) {
	op.charts.Put(chart)
}

// Global object pool
var globalPool *ObjectPool

// InitObjectPool initializes the global object pool
func InitObjectPool() {
	globalPool = NewObjectPool()
	logger.Debug("Initialized object pool")
}

// GetPooledNote gets a note from the global pool
func GetPooledNote() *types.Note {
	if globalPool == nil {
		InitObjectPool()
	}
	return globalPool.GetNote()
}

// ReturnNote returns a note to the global pool
func ReturnNote(note *types.Note) {
	if globalPool != nil {
		globalPool.PutNote(note)
	}
}

// Performance monitoring
type PerformanceMonitor struct {
	mu          sync.RWMutex
	loadTimes   map[string]time.Duration
	parseTimes  map[string]time.Duration
	memorySizes map[string]int64
}

// NewPerformanceMonitor creates a performance monitor
func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{
		loadTimes:   make(map[string]time.Duration),
		parseTimes:  make(map[string]time.Duration),
		memorySizes: make(map[string]int64),
	}
}

// RecordLoadTime records song loading time
func (pm *PerformanceMonitor) RecordLoadTime(songName string, duration time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.loadTimes[songName] = duration
}

// RecordParseTime records song parsing time
func (pm *PerformanceMonitor) RecordParseTime(songName string, duration time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.parseTimes[songName] = duration
}

// RecordMemorySize records song memory usage
func (pm *PerformanceMonitor) RecordMemorySize(songName string, size int64) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.memorySizes[songName] = size
}

// GetStats returns performance statistics
func (pm *PerformanceMonitor) GetStats() PerformanceStats {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var totalLoadTime, totalParseTime time.Duration
	var totalMemory int64
	var maxLoadTime, maxParseTime time.Duration
	var maxMemory int64

	count := len(pm.loadTimes)

	for _, loadTime := range pm.loadTimes {
		totalLoadTime += loadTime
		if loadTime > maxLoadTime {
			maxLoadTime = loadTime
		}
	}

	for _, parseTime := range pm.parseTimes {
		totalParseTime += parseTime
		if parseTime > maxParseTime {
			maxParseTime = parseTime
		}
	}

	for _, memory := range pm.memorySizes {
		totalMemory += memory
		if memory > maxMemory {
			maxMemory = memory
		}
	}

	stats := PerformanceStats{
		TotalSongs:     count,
		TotalLoadTime:  totalLoadTime,
		TotalParseTime: totalParseTime,
		TotalMemory:    totalMemory,
		MaxLoadTime:    maxLoadTime,
		MaxParseTime:   maxParseTime,
		MaxMemory:      maxMemory,
	}

	if count > 0 {
		stats.AvgLoadTime = totalLoadTime / time.Duration(count)
		stats.AvgParseTime = totalParseTime / time.Duration(count)
		stats.AvgMemory = totalMemory / int64(count)
	}

	return stats
}

// PerformanceStats contains performance metrics
type PerformanceStats struct {
	TotalSongs     int
	TotalLoadTime  time.Duration
	TotalParseTime time.Duration
	TotalMemory    int64
	AvgLoadTime    time.Duration
	AvgParseTime   time.Duration
	AvgMemory      int64
	MaxLoadTime    time.Duration
	MaxParseTime   time.Duration
	MaxMemory      int64
}

// Global performance monitor
var perfMonitor *PerformanceMonitor

// InitPerformanceMonitor initializes performance monitoring
func InitPerformanceMonitor() {
	perfMonitor = NewPerformanceMonitor()
	logger.Debug("Initialized performance monitor")
}

// RecordSongLoadTime records loading time for a song
func RecordSongLoadTime(songName string, duration time.Duration) {
	if perfMonitor == nil {
		InitPerformanceMonitor()
	}
	perfMonitor.RecordLoadTime(songName, duration)
}

// GetPerformanceStats returns global performance statistics
func GetPerformanceStats() PerformanceStats {
	if perfMonitor == nil {
		return PerformanceStats{}
	}
	return perfMonitor.GetStats()
}
