package cache

import (
    "errors"
)

var ErrFull = errors.New("the cache is full")

// A Cache models the server's back-end.
//
// It provides all requested caching functionality while
// abstracting away implementation details.
//
type Cache struct {
    entries map[string] string
    stats Stats
    limit uint32
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
func NewCache (limit uint32) *Cache {
    return &Cache {
        entries: map[string] string { },
        limit: limit,
    }
}

//
// Returns the current statistics for this cache.
func (db *Cache) Stats () Stats {
    db.stats.CurrentItems = uint32(len(db.entries))
    db.stats.Limit = db.limit
    return db.stats
}

//
// Set a value in the cache. Also update relative statistics.
// If this would exceed db.Stats().Limit elements, then
// no change is made in the cache and ErrFull is returned.
func (db *Cache) Set (key, data string) error {
    if uint32(len(db.entries)) == db.limit {
        return ErrFull
    }
    db.stats.Sets++;
    db.entries[key] = data
    return nil
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
    delete(db.entries, key)
    return
}
