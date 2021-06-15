package geecache

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

type PeerGetter interface {
	HttpGet(group string, key string) ([]byte, error)
}