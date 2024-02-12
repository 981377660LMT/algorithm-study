1. 优化: 对二的整数幂取模(位运算时) 可以换成与运算

```ts
const id = index >> 5
const mask = index & 31
```

---

当 Set 存储的值为 0-1e9 时,可以将 Set 换成 BitSet,这样可以节省空间,加快查询速度

---

bitset 维护邻接表 在字符串匹配中的应用

由于字符种类少，因此可以 bitset 维护邻接表
当字符集大小可以接受时字符串匹配的一种做法：
**给每个字符集开一个 bitset,存每一个字母出现的位(类似邻接表)**

```
// 给定一个字符串s，有两种操作：
// 1 pos c，将s[pos]改为c
// 2 start end word，求字符串word在s[start:end]中出现的次数
//
// 由于字符种类少，因此可以 bitset 维护邻接表
// 对每次询问匹配t，处理出一个bitset.对每个t[i]，让 res &= indexes[t[i]]>>i 即可.
// 最后查询 [start,end-len(t)+1) 中合法的匹配起点个数。
```
