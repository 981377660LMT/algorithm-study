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

---

https://www.cnblogs.com/alex-wei/p/Common_String_Theory_Theory.html
