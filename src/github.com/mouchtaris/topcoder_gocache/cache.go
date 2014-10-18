package topcoder_gocache

// A Cache models the server's back-end.
//
// It provides all requested caching functionality while
// abstracting away implementation details.
//
type Cache struct {
    entries map[string] string
    gets, sets, getHits, getMisses, deleteHits, deleteMisses, size, limit uint32
}

//
// Make a cache get without updating statistics.
func (db *Cache) untrackedGet (key string) (item string, present bool) {
    item, present = db.entries[key]
    return
}

//
// Construct and initialise a Cache.
func MakeCache () Cache {
    return Cache {
        map[string] string { },
        0, 0, 0, 0, 0, 0, 0, 0,
    }
}

//
// Set a value in the cache. Also update relative statistics.
func (db *Cache) Set (key, data string) {
    db.sets++;
    db.entries[key] = data
}

//
// Retrieve a value from the cache, if it is stored.
// Also update relative statistics.
func (db *Cache) Get (key string) (item string, present bool) {
    db.gets++;
    item, present = db.untrackedGet(key)
    if present {
        db.getHits++
    } else {
        db.getMisses++
    }
    return
}

//
// Delete an entry from the cache, if it is stored, and
// return its value. Also update relative statistics.
func (db *Cache) Delete (key string) (deletedItem string, present bool){
    deletedItem, present = db.untrackedGet(key)
    if present {
        db.deleteHits++
    } else {
        db.deleteMisses++
    }
    return
}
