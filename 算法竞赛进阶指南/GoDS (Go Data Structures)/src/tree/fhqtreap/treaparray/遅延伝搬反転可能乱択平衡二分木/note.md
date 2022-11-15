注意 FHQTreap 在数据`较为随机`时表现很好
但是如果数据较为有序时 会进行很多次 split 和 merge 操作
缺点是区间更新时,不能像线段树一样进行有效的懒标记

```Go
这里Query的效率非常底 因为区间更新范围很大 要进行很多次split和merge操作
for _, num := range nums {
  preMax := treap.Query(num-k, num)   // k=5e4 num从5e4开始到1e5
  treap.Update(num, num+1, preMax+1)
}
```

## 注意 nyann 的库里的 merge 和 split 的细节

https://nyaannyaan.github.io/library/rbst/rbst-base.hpp
不维护节点的 priority 而是 merge 时比较左右子树的 size 取随机值

```C++
  static uint64_t rng() {
    static uint64_t x_ = 88172645463325252ULL;
    return x_ ^= x_ << 7, x_ ^= x_ >> 9, x_ & 0xFFFFFFFFull;
  }


  Ptr merge(Ptr l, Ptr r) {
  if (!l || !r) return l ? l : r;
  // 注意merge这里不用权值merge
  if (int((rng() * (l->cnt + r->cnt)) >> 32) < l->cnt) {
    push(l);
    l->r = merge(l->r, r);
    return update(l);
  } else {
    push(r);
    r->l = merge(l, r->l);
    return update(r);
  }
}
```

## !不用指针快很多

- 力扣这种多个 case 的题目,不适合用指针(指针很慢)
- atcoder 这种单个 case 的题目,不用指针可能 MLE,这时可以用指针+禁用 GC 来解决
