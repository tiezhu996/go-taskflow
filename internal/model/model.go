package model

type Status string

const (
	StatusPending Status = "pending"
	StatusRunning Status = "running"
	StatusDone    Status = "done"
	StatusFailed  Status = "failed"
)

type Task struct {
	ID      string
	Payload string
	Status  Status
	Result  string
}

// ChunkIDs splits ids into chunks of at most size; each returned chunk is an
// independent slice, so mutating a chunk never affects the caller's input.
func ChunkIDs(ids []string, size int) [][]string {
	if size <= 0 {
		size = 1
	}
	chunks := make([][]string, 0, (len(ids)+size-1)/size)
	for i := 0; i < len(ids); i += size {
		end := i + size
		if end > len(ids) {
			end = len(ids)
		}
		chunk := make([]string, end-i)
		copy(chunk, ids[i:end])
		chunks = append(chunks, chunk)
	}
	return chunks
}
