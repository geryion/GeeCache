package geecache

//只读数据类型
type ByteView struct {
	b []byte
}

/*
	使用锁操作进行对LRU淘汰算法进行并发控制处理
*/
func (v ByteView)Len() int {
	return len(v.b)
}

func (v ByteView)ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView)String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
