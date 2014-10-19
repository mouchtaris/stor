package cache

import (
    "math"
)

// A Cache models the server's back-end.
//
// It provides all requested caching functionality while
// abstracting away implementation details.
//
type Cache struct {
    entries map[string] string
    stats Stats
}

type Stats struct {
    Gets, Sets, GetHits, GetMisses, DeleteHits, DeleteMisses uint32
    CurrentItems, Limit uint32
}

//
// Make a cache get without updating statistics.
func (db *Cache) untrackedGet (key string) (item string, present bool) {
    item, present = db.entries[key]
    return
}

//
// Construct and initialise a Cache.
func NewCache () *Cache {
    return &Cache {
        entries: map[string] string { },
    }
}

//
// Returns the current statistics for this cache.
func (db *Cache) Stats () Stats {
    db.stats.CurrentItems = uint32(len(db.entries))
    db.stats.Limit = math.MaxUint16
    return db.stats
}

//
// Set a value in the cache. Also update relative statistics.
func (db *Cache) Set (key, data string) {
    db.stats.Sets++;
    db.entries[key] = data
}

//
// Retrieve a value from the cache, if it is stored.
// Also update relative statistics.
func (db *Cache) Get (key string) (item string, present bool) {
    db.stats.Gets++;
    item, present = db.untrackedGet(key)
    if present {
        db.stats.GetHits++
    } else {
        db.stats.GetMisses++
    }
    return
}

//
// Delete an entry from the cache, if it is stored, and
// return its value. Also update relative statistics.
func (db *Cache) Delete (key string) (deletedItem string, present bool){
    deletedItem, present = db.untrackedGet(key)
    if present {
        db.stats.DeleteHits++
    } else {
        db.stats.DeleteMisses++
    }
    return
}
