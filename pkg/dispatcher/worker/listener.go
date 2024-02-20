package dispatcher


// Watcher is the woker that listens to the data change buffer for new coming data change files.
type Watcher struct {
	// The data change buffer to listen to.
	buffer *DataChangeBuffer

	// The listener to listen to the data change buffer.
	listener *Listener

	// The worker to handle the data change files.
	worker *Worker
}