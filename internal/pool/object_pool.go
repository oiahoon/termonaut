package pool

import (
	"strings"
	"sync"
	"time"

	"github.com/oiahoon/termonaut/pkg/models"
)

// CommandPool manages a pool of Command objects to reduce allocations
type CommandPool struct {
	pool sync.Pool
}

// NewCommandPool creates a new command pool
func NewCommandPool() *CommandPool {
	return &CommandPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &models.Command{}
			},
		},
	}
}

// Get retrieves a command from the pool
func (p *CommandPool) Get() *models.Command {
	cmd := p.pool.Get().(*models.Command)
	// Reset the command to ensure clean state
	p.resetCommand(cmd)
	return cmd
}

// Put returns a command to the pool
func (p *CommandPool) Put(cmd *models.Command) {
	if cmd != nil {
		p.pool.Put(cmd)
	}
}

// resetCommand resets a command to its zero state
func (p *CommandPool) resetCommand(cmd *models.Command) {
	cmd.ID = 0
	cmd.Timestamp = time.Time{}
	cmd.SessionID = 0
	cmd.Command = ""
	cmd.ExitCode = 0
	cmd.CWD = ""
	cmd.DurationMs = 0
	cmd.Category = ""
}

// StringBuilderPool manages a pool of string builders for efficient string operations
type StringBuilderPool struct {
	pool sync.Pool
}

// NewStringBuilderPool creates a new string builder pool
func NewStringBuilderPool() *StringBuilderPool {
	return &StringBuilderPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &strings.Builder{}
			},
		},
	}
}

// Get retrieves a string builder from the pool
func (p *StringBuilderPool) Get() *strings.Builder {
	sb := p.pool.Get().(*strings.Builder)
	sb.Reset() // Ensure clean state
	return sb
}

// Put returns a string builder to the pool
func (p *StringBuilderPool) Put(sb *strings.Builder) {
	if sb != nil {
		// Only pool builders that aren't too large to avoid memory bloat
		if sb.Cap() < 64*1024 { // 64KB limit
			p.pool.Put(sb)
		}
	}
}

// ByteSlicePool manages a pool of byte slices for I/O operations
type ByteSlicePool struct {
	pool sync.Pool
	size int
}

// NewByteSlicePool creates a new byte slice pool with specified size
func NewByteSlicePool(size int) *ByteSlicePool {
	return &ByteSlicePool{
		size: size,
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, size)
			},
		},
	}
}

// Get retrieves a byte slice from the pool
func (p *ByteSlicePool) Get() []byte {
	return p.pool.Get().([]byte)
}

// Put returns a byte slice to the pool
func (p *ByteSlicePool) Put(b []byte) {
	if b != nil && len(b) == p.size {
		// Clear the slice before returning to pool
		for i := range b {
			b[i] = 0
		}
		p.pool.Put(b)
	}
}

// MapPool manages a pool of string maps for temporary operations
type MapPool struct {
	pool sync.Pool
}

// NewMapPool creates a new map pool
func NewMapPool() *MapPool {
	return &MapPool{
		pool: sync.Pool{
			New: func() interface{} {
				return make(map[string]interface{})
			},
		},
	}
}

// Get retrieves a map from the pool
func (p *MapPool) Get() map[string]interface{} {
	m := p.pool.Get().(map[string]interface{})
	// Clear the map
	for k := range m {
		delete(m, k)
	}
	return m
}

// Put returns a map to the pool
func (p *MapPool) Put(m map[string]interface{}) {
	if m != nil {
		// Only pool maps that aren't too large
		if len(m) < 100 {
			p.pool.Put(m)
		}
	}
}

// SlicePool manages a pool of string slices
type SlicePool struct {
	pool sync.Pool
	cap  int
}

// NewSlicePool creates a new slice pool with specified capacity
func NewSlicePool(capacity int) *SlicePool {
	return &SlicePool{
		cap: capacity,
		pool: sync.Pool{
			New: func() interface{} {
				return make([]string, 0, capacity)
			},
		},
	}
}

// Get retrieves a slice from the pool
func (p *SlicePool) Get() []string {
	slice := p.pool.Get().([]string)
	return slice[:0] // Reset length but keep capacity
}

// Put returns a slice to the pool
func (p *SlicePool) Put(slice []string) {
	if slice != nil && cap(slice) == p.cap {
		p.pool.Put(slice)
	}
}

// GlobalPools provides access to commonly used object pools
type GlobalPools struct {
	Commands      *CommandPool
	StringBuilder *StringBuilderPool
	SmallBytes    *ByteSlicePool // 1KB buffers
	LargeBytes    *ByteSlicePool // 64KB buffers
	Maps          *MapPool
	Slices        *SlicePool
}

// NewGlobalPools creates a new set of global pools
func NewGlobalPools() *GlobalPools {
	return &GlobalPools{
		Commands:      NewCommandPool(),
		StringBuilder: NewStringBuilderPool(),
		SmallBytes:    NewByteSlicePool(1024),      // 1KB
		LargeBytes:    NewByteSlicePool(64 * 1024), // 64KB
		Maps:          NewMapPool(),
		Slices:        NewSlicePool(50), // 50 string capacity
	}
}

// Default global pools instance
var DefaultPools = NewGlobalPools()

// Convenience functions for default pools

// GetCommand gets a command from the default pool
func GetCommand() *models.Command {
	return DefaultPools.Commands.Get()
}

// PutCommand returns a command to the default pool
func PutCommand(cmd *models.Command) {
	DefaultPools.Commands.Put(cmd)
}

// GetStringBuilder gets a string builder from the default pool
func GetStringBuilder() *strings.Builder {
	return DefaultPools.StringBuilder.Get()
}

// PutStringBuilder returns a string builder to the default pool
func PutStringBuilder(sb *strings.Builder) {
	DefaultPools.StringBuilder.Put(sb)
}

// GetSmallBuffer gets a small byte buffer from the default pool
func GetSmallBuffer() []byte {
	return DefaultPools.SmallBytes.Get()
}

// PutSmallBuffer returns a small byte buffer to the default pool
func PutSmallBuffer(b []byte) {
	DefaultPools.SmallBytes.Put(b)
}

// GetLargeBuffer gets a large byte buffer from the default pool
func GetLargeBuffer() []byte {
	return DefaultPools.LargeBytes.Get()
}

// PutLargeBuffer returns a large byte buffer to the default pool
func PutLargeBuffer(b []byte) {
	DefaultPools.LargeBytes.Put(b)
}

// GetMap gets a map from the default pool
func GetMap() map[string]interface{} {
	return DefaultPools.Maps.Get()
}

// PutMap returns a map to the default pool
func PutMap(m map[string]interface{}) {
	DefaultPools.Maps.Put(m)
}

// GetSlice gets a slice from the default pool
func GetSlice() []string {
	return DefaultPools.Slices.Get()
}

// PutSlice returns a slice to the default pool
func PutSlice(slice []string) {
	DefaultPools.Slices.Put(slice)
}
