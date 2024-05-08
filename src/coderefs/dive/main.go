package dive

/*
Storage represents the distributed data storage within Diarkis cluster.
*/
type Storage struct{}

/*
Options represents optional configurations
for New(name string, options *Options) to create a new storage without the use of configuration file.
*/
type Options struct {
	RingSize           uint32
	TargetNodeType     string
	Migration          bool
	MigrationBatchSize uint16
	MigrationInterval  uint16
}

/*
KeyValuePair represents internally used data
*/
type KeyValuePair struct {
	Name      string      `json:"n"`
	Key       string      `json:"k"`
	Value     interface{} `json:"v"`
	Timestamp int64       `json:"t"`
	TTL       int64       `json:"l"`
}

/*
RangeData represents internally used data
*/
type RangeData struct {
	Name string `json:"n"`
	Key  string `json:"k"`
	From int    `json:"f"`
	To   int    `json:"t"`
	TTL  int64  `json:"l"`
}

/*
MigrationData represents internally used data
*/
type MigrationData struct {
	Name string                   `json:"n"`
	List []map[string]interface{} `json:"l"`
}

/*
Setup must be invoked before calling diarkis.Start() in order to use dive module.
*/
func Setup(confpath string) {
}

/*
IsReadyToShutdown returns true when Dive module is ready to shutdown.

	[NOTE] This function is used internally.
*/
func IsReadyToShutdown() bool {
	return false
}

/*
GetStorageByName returns a storage by the given name.

You must use this function to get a storage if you define storage names in
the configuration JSON.

The function returns nil if the given name does not match existing storages.
*/
func GetStorageByName(name string) *Storage {
	return nil
}

/*
New creates an instance of Storage struct to interact with the distributed storage if the key does not exist.

If you define storage names in the configuration JSON file, you do not need to use this function.
Use GetStorageByName(name string) instead.

Error Cases

	┌────────────────────────┬──────────────────────────────────────────────────────────────┐
	│ Error                  │ Reason                                                       │
	╞════════════════════════╪══════════════════════════════════════════════════════════════╡
	│ Setup must be invoked  │ In order to use Dive module,                                 │
	│                        │ dive.Setup() must be called before calling diarkis.Start()   │
	├────────────────────────┼──────────────────────────────────────────────────────────────┤
	│ Name must not be empty │ Input given name is an empty string.                         │
	├────────────────────────┼──────────────────────────────────────────────────────────────┤
	│ Storage already exists │ The given storage name already exists.                       │
	└────────────────────────┴──────────────────────────────────────────────────────────────┘

Parameters

	name    - Unique name for the storage.
	options - Optional configurations for the new storage.
*/
func New(name string, options *Options) (*Storage, error) {
	return nil, nil
}

/*
UpdateRingSize changes the ring size.

	[IMPORTANT] Changing the ring size will change the key to address index calculation.
	            It means that the same keys may point to different servers leading to values of keys not being found etc.
*/
func (s *Storage) UpdateRingSize(ringSize uint32) {
}

/*
IsStorage returns true, if the server node type is the configured target node type.

	[NOTE] Default target node type is HTTP.
*/
func (s *Storage) IsStorage() bool {
	return false
}

/*
ResolveKey returns the associated internal server node address of the given key.
*/
func (s *Storage) ResolveKey(key string) string {
	return ""
}

/*
GetLocalKeys returns an array of locally stored keys.
*/
func (s *Storage) GetLocalKeys() []string {
	return nil
}

/*
GetNumberOfNodes returns the number of recognized nodes for Dive.
*/
func (s *Storage) GetNumberOfNodes() int {
	return 0
}

/*
SetOnNodeChange assigns a callback to be invoked when the number of target node server changes.

	[NOTE] The callbacks will be called regardless of migration being enabled or not.
*/
func (s *Storage) SetOnNodeChange(cb func()) {
}

/*
SetChangeRingSizeOnNodeUpdate assigns a callback to change ring size before node update.
*/
func (s *Storage) SetChangeRingSizeOnNodeUpdate(cb func(int) uint32) bool {
	return false
}

/*
RemoveOnNodeChange removes the assigned callback.
*/
func (s *Storage) RemoveOnNodeChange(cb func()) bool {
	return false
}
