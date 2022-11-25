package SCA

import (
	"container/list"
	"fmt"
)

type SCA struct {
	maxpage int64                    //最多存储数据的字节个数，超过此数量便会触发数据的淘汰
	curPage int64                    //目前存储的字节个数
	list    *list.List               //使用go语言内置的双向链表存储节点
	cache   map[string]*list.Element //通过节点的key快速定位到属于哪个节点，不需遍历双向链表
}

type Entry struct {
	Key   string //每个节点的唯一标识，作为key储存到lru的cache里
	Value []byte //携带的数据
}

// New Constructor of LRU
func New(maxPage int64) *SCA {
	return &SCA{
		maxpage: maxPage,
		curPage: 0,
		list:    list.New(),
		cache:   make(map[string]*list.Element),
	}
}

// Get look up a key`s value
func (l *SCA) Get(key string) ([]byte, bool) {
	if ele, exist := l.cache[key]; exist {
		//l.list.MoveToFront(ele)
		if entry, ok := ele.Value.(*Entry); ok {
			return entry.Value, true
		}
	}
	return nil, false
}

// RemoveOldest remove entry from linklist back
func (l *SCA) RemoveOldest() {
	ele := l.list.Front() //取出链表头部节点
	if ele != nil {
		l.list.Remove(ele)
		if entry, ok := ele.Value.(*Entry); ok {
			delete(l.cache, entry.Key)
			l.curPage = l.curPage - 1
		}
	}
}

// Add a value to lru
func (l *SCA) Add(key string, value []byte) {
	// 两种情况
	if ele, ok := l.cache[key]; ok {
		//l.list.MoveToFront(ele)
		if entry, ok := ele.Value.(*Entry); ok {
			entry.Value = value
		}
	} else {
		ele := l.list.PushBack(&Entry{Key: key, Value: value})
		l.cache[key] = ele
		l.curPage = l.curPage + 1
	}
	// 已使用字节数大于最大字节数时，需要移除链表尾部节点，知到已使用字节数小于最大字节数
	for l.maxpage > 0 && l.maxpage < l.curPage {
		l.RemoveOldest()
	}
}

// Len the number of cache entries
func (l *SCA) Len() int {
	return l.list.Len()
}

// Len the number of cache entries
func (l *SCA) Schedule(actualCall ...string) int {
	pagefailNum := 0
	for _, v := range actualCall {
		_, IsExist := l.Get(v)
		fmt.Println("Is", IsExist)
		if IsExist == true {
			//fmt.Println("存在")
		} else {
			pagefailNum = pagefailNum + 1
			if int64(l.Len()) >= l.maxpage {
				l.RemoveOldest()
				l.Add(v, []byte("0"))
			} else {
				l.Add(v, []byte("0"))
			}
		}
	}
	return pagefailNum

}
