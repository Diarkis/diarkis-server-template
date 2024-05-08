package mesh

/*
Set creates a new value to the SharedData and propagates the change to all pods in the cluster.
*/
func (sd *sharedData) Set(key string, value int16) bool {
	return false
}

/*
Remove deletes the value of SharedData and propagates the change to all pods in the cluster.
*/
func (sd *sharedData) Remove(key string) bool {
	return false
}

/*
Flush flushes out the buffered SharedData value.
*/
func (sd *sharedData) Flush() (string, int16, bool) {
	return "", 0, false
}

/*
SyncFromMARS returns a map of new keys and values
*/
func (sd *sharedData) SyncFromMARS(source map[string]interface{}) map[string]int16 {
	return nil
}

/*
Get returns the SharedData value.
*/
func (sd *sharedData) Get(key string) (int16, bool) {
	return 0, false
}

/*
Length returns the number of shared data keys currently stored.
*/
func (sd *sharedData) Length() int {
	return 0
}

/*
IsFull returns true if the number of shared data key is at maximum.
*/
func (sd *sharedData) IsFull() bool {
	return false
}
