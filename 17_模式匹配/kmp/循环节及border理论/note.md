1. 暴力求循环节
   枚举 n 的因子检查 `O(n^4/3)`

   ```python
   n = len(s)
   for len_ in range(1, n // 2 + 1):
    # 因子大约为n^(1/3)个
    if n % len_ == 0 and s[:len_] * (n // len_) == s:
        return True, s[:len]
   return False
   ```

2. kmp/z 函数求前缀循环节
   如果 S 具有长为|S|-p 的 border，则说明 S 具有周期 p.
   **注意这里的周期不是完整的周期**, 例如 abcabcab，具有周期 5.
   如果需要完整周期，需要 `|S| % p === 0` 来判断
3. 字符串哈希求子串循环节
   `CF580E - Kefa and Watch-区间赋值+区间查询`
   给定一个字符串 s，和一个闭区间 `[left,right]`
   如果 `s[left,right-len] == s[left+len,right]`，那么 `s[left,right]` 就存在一个长度为 len 的循环节

---

# 结论

1. 最小周期 <=> 最大 border
2. 如果 S 具有长为|S|-p 的 border，则说明 S 具有周期 p
   注意这里的周期不是完整的周期, 例如 abcabcab，具有周期 5.
   如果需要完整周期，需要 `|S| % p === 0` 来判断
3. 字符串 s 的所有 border 长度排序后可分成 O(log |s|) 段, 每段是一个等差数列。
4. 给定两个回文串 a、b，当且仅当这`两个回文串的最短回文整周期串相同`时，a+b 是回文串.“
   具体而言，如果 n%(n-next[n-1]) == 0，则最短回文整周期串为 s[:n-next[n-1]]，否则为 s 本身.

---

https://www.cnblogs.com/alex-wei/p/Common_String_Theory_Theory.html

---

https://taodaling.github.io/blog/2019/06/14/%E5%BA%8F%E5%88%97%E9%97%AE%E9%A2%98/

# 旋转周期

定义：对于字符串 S，如果存在某个前缀 P，从 S 中删除前缀 P 并将 P 追加到结果后面，得到字符串 S′，如果 S=S′，那么称 ∣P∣ 是 S 的一个旋转周期。
比如对于 abab，我们可以发现其拥有旋转周期 2，因为(ab)ab+ab=abab。

命题：

- 字符串的旋转周期一定是字符串的周期
- 如果 a 是字符串 S 的旋转周期，那么 ∣S∣−a 也是 S 的一个旋转周期。
- 如果 a,b 是字符串 S 的旋转周期，那么 gcd(a,b)也是 S 的一个旋转周期
- 长度为 n 的字符串 S 的**最小旋转周期** p，一个数 x 是 S 的旋转周期当且仅当 p∣x。
  我们可以通过枚举 p 的因子得到 S 的最小旋转周期，时间复杂度为 O(nlog2n)
