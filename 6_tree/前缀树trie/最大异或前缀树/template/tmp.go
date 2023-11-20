package main

func main() {

}

// EXTRA: 可持久化字典树
// 注意为了拷贝一份 trieNode，这里的接收器不是指针
// https://oi-wiki.org/ds/persistent-trie/
// roots := make([]*trieNode, n+1)
// roots[0] = &trieNode{}
// roots[i+1] = roots[i].put(s)
func (o trieNode) put(s []byte) *trieNode {
	if len(s) == 0 {
		o.cnt++
		return &o
	}
	b := s[0] - 'a' //
	if o.son[b] == nil {
		o.son[b] = &trieNode{}
	}
	o.son[b] = o.son[b].put(s[1:])
	//o.maintain()
	return &o
}
