- 字符串哈希:O(n)预处理字符串 O(1)获取区间字符串哈希值

  取一固定值 P(131 或 13331)，`把字符串看作 Р 进制数`，并分配一个大于 0 的数值，代表母 l 子符。一般来说，我们分配的数值都远小于 P。例如，对于小写字母构成的字符串，可以令 a = 1,b = 2,…,z= 26。取一固定值 M(2^64)，求出该 Р 进制数对 M 的余数，作为该字符的哈希值。
  S 的哈希值为 H(S) 那么 S 后加一个 c
  `H(S+c)=(H(S)*P+value[c]) mod M`
  S 的哈希值为 H(S) S+T 的哈希值为 H(S+T) 那么 T 的哈希值为
  `H(T)=(H(S+T)-H(S)*p**T.length) mod M` 即相当于在 S 后补零 对齐相减

  根据上面两种操作，可以`O(N)预处理字符串所有前缀的哈希值` 从而在 O(1)时间查询字符串任一个字串的 Hash 值
  哈希值相等 则认为字符串相等

  例题:DNA 序列/最长重复子串(二分+哈希)/最长回文子串(nlogn)

  **Rabin-Karp 字符串哈希算法:P 取 26**,编码对应 0 到 25

- 字符串哈希的「构造 pp 数组」和「计算哈希」的过程，不会溢出吗？

- 如果我们期望做到严格 O(n)，进行计数的「哈希表」就不能是以 String 作为 key，只能使用 Integer（也就是 hash 结果本身）作为 key。`因为 Java 中的 String 的 hashCode 实现是会对字符串进行遍历的`，这样哈希计数过程仍与长度有关，而 Integer 的 hashCode 就是该值本身，这是与长度无关的。

- 注意校验:
  在解决算法题时，我们只要判断两个编码是否相同，就表示它们对应的字符串是否相同。但在实际的应用场景中，`会出现字符串不同但编码相同的情况，因此在实际场景中使用 Rabin-Karp 字符串编码时，推荐在编码相同时再对字符串进行比较，防止出现错误。`

`双哈希，比较快的语言可以使用`
https://github.com/harttle/contest.js/blob/master/src/rolling-hash.ts

```py
hash1 = hasher1(2,5)
hash2 = hasher2(2,5)
newHash = f'{hash1},{hash2}'
```

---

自然溢出哈希冲突
https://hos.ac/blog/#blog0003

---

update : 2024.3.12
https://zhuanlan.zhihu.com/p/25855753

- SDBM HASH
- BKDR HASH

---

hack 哈希

anti-hash
https://heltion.github.io/anti-hash/
