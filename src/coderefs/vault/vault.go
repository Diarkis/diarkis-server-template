package vault

import (
	"sync"
)

/*
Storage storage of each vault
*/
type Storage struct{ sync.RWMutex }

/*
Setup Sets up vault
*/
func Setup() {
}

/*
SetShutdownWaitTime sets shutdown wait time on process termination in seconds
*/
func SetShutdownWaitTime(waitTime int) {
}

/*
OnEmpty registers a callback to be invoked when all vaults are empty before shutting down
*/
func OnEmpty(cb func()) {
}

/*
Stop Stops vault
*/
func Stop(next func(error)) {
}

/*
NewVault Creates a new vault with a name

[NOTE] Uses mutex lock internal.

	IMPORTANT: create all required vaults BEFORE process start.
	interval is in seconds
*/
func NewVault(name string, interval int64) *Storage {
	return nil
}

/*
OnDel Registers a callback function to a vault on delete.

[NOTE] Uses mutex lock internal.

The callback will be passed the key and the value that has been deleted
*/
func (vs *Storage) OnDel(callback func(string, interface{})) {
}

/*
OnDelAndSet is used internally.
*/
func (vs *Storage) OnDelAndSet(cb func(delKey, setKey string, delValues, setValues interface{})) {
}

/*
Set Sets a value of a key by vault name - ttlDuration is in seconds

[NOTE] Uses mutex lock internal.
*/
func (vs *Storage) Set(key string, values interface{}, ttlDuration int64) bool {
	return false
}

/*
SetIfNotExists sets the given key and the return value of the callback as its value to the storage.

The condition for setting the key and the value is to have the same key not exist in the storage,
and the callback must return an interface{} as its value.

It returns false if the key and the value are not stored.
*/
func (vs *Storage) SetIfNotExists(key string, ttlDuration int64, cb func() interface{}) bool {
	return false
}

/*
SetAsStatic sets a static key that will NEVER be discarded

[NOTE] Uses mutex lock internal.

# IMPORTANT A static key will NEVER be discarded from the server memory

IMPORTANT A static key may sill be deleted by using Del()
*/
func (vs *Storage) SetAsStatic(key string, values interface{}) bool {
	return false
}

/*
EnableStatic makes the value of the given key static (static keys do not expire)

[NOTE] Uses mutex lock internal.
*/
func (vs *Storage) EnableStatic(key string) bool {
	return false
}

/*
DisableStatic makes the value of the given key not static (non-static keys will expire)

[NOTE] Uses mutex lock internal.
*/
func (vs *Storage) DisableStatic(key string) bool {
	return false
}

/*
SetOrUpdate checks to see if the key given exists.

[NOTE] Uses mutex lock internal.

	If the key does not exist or has expired,
	   it creates a new entry by using the returned value interface{} from operation function.
	If key exists and it is valid, it allows operation function to update the entry passed to it.
	If isStatic is true, the item will NEVER be discarded
*/
func (vs *Storage) SetOrUpdate(key string, operation func(interface{}) interface{}, ttlDuration int64, isStatic bool) bool {
	return false
}

/*
Extend Extends TTL of a value data in a vault

[NOTE] Uses mutex lock internal.

it allows expired vault to be extended also. TTL is in seconds.
*/
func (vs *Storage) Extend(key string, ttl int64) bool {
	return false
}

/*
IsEmpty returns true if the vault is empty - it checks TTL on each item thus very slow and expensive

[NOTE] Uses mutex lock internal.
*/
func (vs *Storage) IsEmpty() bool {
	return false
}

/*
Exists returns true if the given key exists - this does NOT update TTL of the key

[NOTE] Uses mutex lock internal.
*/
func (vs *Storage) Exists(key string) bool {
	return false
}

/*
AsyncGet passes the vault data to operation callback function if the key exists.

[NOTE] Uses mutex lock internal.

This function LOCKS the vault data UNTIL operation callback calls the end callback.
*/
func (vs *Storage) AsyncGet(key string, operation func(interface{}, func())) {
}

/*
Get Returns a value of a key by vault name

[NOTE] Uses mutex lock internal.
*/
func (vs *Storage) Get(key string) interface{} {
	return nil
}

/*
Peek returns the stored value of the key if exists.

[NOTE] Uses mutex lock internal.

This does NOT update the key's TTL unlike Get.
*/
func (vs *Storage) Peek(key string) interface{} {
	return nil
}

/*
GetAllKeys returns all keys in the storage as an array

[NOTE] Uses mutex lock internal.
*/
func (vs *Storage) GetAllKeys() []string {
	return nil
}

/*
Length returns the number of keys in the Storage

[NOTE] Uses mutex lock internal.
*/
func (vs *Storage) Length() int {
	return 0
}

/*
GetAll returns all items in the storage as an array

[NOTE] Uses mutex lock internal.
*/
func (vs *Storage) GetAll() []interface{} {
	return nil
}

/*
Update Updates a value of a key by vault name and returns the updated data

[NOTE] Uses mutex lock internal.
*/
func (vs *Storage) Update(key string, operation func(interface{})) interface{} {
	return nil
}

/*
UpdateWithoutTTLChange updates the value of the key given without updating the key's TTL.

[NOTE] Uses mutex lock internal.
*/
func (vs *Storage) UpdateWithoutTTLChange(key string, operation func(interface{})) interface{} {
	return nil
}

/*
DelIf deletes the value of the key given if the conditions implemented by condition function are met

[NOTE] Uses mutex lock internal.
*/
func (vs *Storage) DelIf(key string, condition func(val interface{}) bool) bool {
	return false
}

/*
Del Deletes a key in a vault by name

[NOTE] Uses mutex lock internal.
*/
func (vs *Storage) Del(key string) bool {
	return false
}

/*
DelAndSet is used internally.
*/
func (vs *Storage) DelAndSet(delKey, setKey string, values interface{}, ttlDuration int64) bool {
	return false
}

/*
Clear Deletes all keys and values in the vault

[NOTE] Uses mutex lock internal.

Does trigger OnDel
*/
func (vs *Storage) Clear(vaultName string) {
}
