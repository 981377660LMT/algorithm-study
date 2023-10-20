# 线段树维护哈希/线段树维护字符串哈希

两个区间的哈希值是可以合并的，所以线段树可以维护区间的哈希值

- 区间的哈希值合并(op/combine):

```go
func (r *RollingHash) Combine(h1, h2 uint, h2len int) uint {
	return h1*r.basePow[h2len] + h2
}
```

- 区间哈希值

```go
func (r *RollingHash) Query(hashes []uint, start, end int) uint {
	return hashes[end] - hashes[start]*r.power[end-start]
}
```

## 例题

1. 单点修改+询问区间回文
   https://blog.csdn.net/qq_21433411/article/details/90719116
   回文=`正串哈希==逆串哈希`
   两个线段树维护正逆串区间哈希值即可
2. 单点修改+子串查询
   https://www.acwing.com/file_system/file/content/whole/index/content/532198/
3. 区间修改字符+询问字符串是否周期
   https://www.luogu.com.cn/problem/CF580E
   **当区间[l+d,r]的哈希值与[l,r-d]的哈希值相等时，那么该区间[l,r]是以 d 为循环节的**
   `卡自然溢出(ull)`

   解法 2:因为每个数字是 0-9,考虑使用切片法
   更新直接 `memcpy`，查询 `memcmp`

4. 区间加+查询区间哈希值
   https://yukicoder.me/problems/no/469

5. https://www.cnblogs.com/zltzlt-blog/p/16797435.html

---

线段树维护多项式哈希

https://www.luogu.com.cn/blog/220037/solution-p3792
哈希函数 f(nums) 将一个序列转为一个整数，且满足：

- 如果 a 为 b 的一个重排列，那么 f(a)=f(b)
- 如果 a 不为 b 的一个重排列，那么大概率 f(a)≠f(b)
- f([i,i+1,...,i+l])可以快速计算
- 在一个带修改的序列中,可以快速计算 f([l,r]) 的值

这样的函数很多，一个通用函数类型就是

`f(k,p,nums) = sum(nums[i]*k^i) % p`

取 k=2,p=1e9+7,就是 rangePow2Sum
