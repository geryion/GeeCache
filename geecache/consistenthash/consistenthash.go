package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//type 对函数进行声明 方便以后调用以及修改
type Hash func(data []byte) uint32

//一个Map包含了一致哈希算法的主要数据结构
type Map struct {
	//todo 是否需要增加mutex
	hash 		Hash
	replicas	int
	keys 		[]int //哈希环数组
	hashMap		map[int]string //
}


//todo 是否需要考虑节点删除和节点新增之后的缓存数据迁移？？
//新增一个Map结构 可以自定义虚拟节点个数和hash算法函数
func MapNew(replicas int, fn Hash) *Map{
	m := &Map{
		replicas: replicas,
		hash: fn,
		hashMap: make(map[int]string),
	}
	//如果没有传入自定义的哈希算法函数，则使用默认的crc32哈希算法
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

//根据节点名称keys创建真实节点和虚拟节点,增加虚拟节点与真实节点映射关系
func (m *Map)MapAdd(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			//todo m.hash return int32, but int(int32)??
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			//增加当前虚拟节点到哈希环上
			m.keys = append(m.keys, hash)
			//把当前的虚拟节点和真实节点映射
			m.hashMap[hash] = key
		}
	}
	//哈希环重新进行排序
	sort.Ints(m.keys)
}

//根据二分查找找到哈希环上最近的节点并返回
func (m *Map)MapGet(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	//查询的key进行hash
	hash := int(m.hash([]byte(key)))
	//对哈希环上的节点值和查询的hash值进行比较 并返回距离最近的节点名称
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}