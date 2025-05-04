# 串联重复（Tandem repeats）

https://taodaling.github.io/blog/2019/06/14/%E5%BA%8F%E5%88%97%E9%97%AE%E9%A2%98/#heading-%E5%AD%97%E7%AC%A6%E4%B8%B2%E7%9A%84%E4%B8%80%E4%BA%9B%E5%91%A8%E6%9C%9Fborder%E6%80%A7%E8%B4%A8

https://oi-wiki.org/string/main-lorentz/

https://leetcode.cn/problems/count-beautiful-splits-in-an-array/solutions/3023632/main-lorentzcha-zhao-zhong-fu-zi-chuan-f-nwsc/
https://leetcode.cn/problems/count-beautiful-splits-in-an-array/solutions/3020684/onlogn-hou-zhui-shu-zu-jie-fa-by-vclip-suiq/
https://leetcode.cn/problems/distinct-echo-substrings/

---

# 串联重复（Tandem Repeats）及相关算法

## 1. 基本概念

**串联重复（Tandem Repeats）** 是指字符串中连续出现两次或多次的子串模式。形式上，如果一个字符串 S 可以表示为 S = uvu，其中 u 和 v 都是非空字符串，那么 uu 就是 S 中的一个串联重复。

### 1.1 术语定义

- **原始串联重复**：如果重复的子串本身不能再被分解为更小的重复模式，则称为原始串联重复。
- **平方串（Square）**：形如 ww 的字符串，即一个子串连续重复两次。
- **立方串（Cube）**：形如 www 的字符串，即一个子串连续重复三次。
- **串联重复的次数**：子串重复的次数，如 "abab" 中 "ab" 重复了2次。
- **最大串联重复**：在给定位置开始的最长串联重复。

## 2. 串联重复的性质

### 2.1 周期性质

如果 S[i...j] 是一个周期为 p 的串联重复，那么：

- S[i...j] 的长度至少为 2p
- S[i...i+p-1] = S[i+p...i+2p-1] = ... = S[j-p+1...j]

### 2.2 数量界限

- 一个长度为 n 的字符串中，不同的串联重复数量可以达到 Θ(n log n)
- 原始串联重复的数量也是 Θ(n log n)

## 3. 检测算法

### 3.1 朴素算法

最简单的方法是枚举所有可能的子串开始位置和长度，然后检查是否形成重复：

```python
def find_tandem_repeats_naive(s):
    n = len(s)
    results = []

    for i in range(n):  # 起始位置
        for length in range(1, (n-i)//2 + 1):  # 子串长度
            repeat_count = 1
            j = i + length
            while j + length <= n and s[j:j+length] == s[j-length:j]:
                repeat_count += 1
                j += length

            if repeat_count >= 2:
                results.append((i, length, repeat_count))

    return results
```

时间复杂度：O(n³)，因为需要枚举 O(n²) 个子串，每个子串的比较需要 O(n) 时间。

### 3.2 基于 Z 算法的检测

Z 算法可以在线性时间内计算每个位置的 Z 值（与前缀的最长公共子串长度）。利用 Z 值可以高效检测重复：

```python
def z_function(s):
    n = len(s)
    z = [0] * n
    l, r = 0, 0
    for i in range(1, n):
        if i <= r:
            z[i] = min(r - i + 1, z[i - l])
        while i + z[i] < n and s[z[i]] == s[i + z[i]]:
            z[i] += 1
        if i + z[i] - 1 > r:
            l, r = i, i + z[i] - 1
    return z

def find_tandem_repeats_z(s):
    n = len(s)
    results = []

    for i in range(n):
        for length in range(1, (n-i)//2 + 1):
            pattern = s[i:i+length]
            extended = pattern + s[i+length:]
            z = z_function(extended)

            if z[length] >= length:  # 找到重复
                repeat_count = 1 + z[length] // length
                results.append((i, length, repeat_count))

    return results
```

时间复杂度：仍然是 O(n³)，但在实践中通常比朴素算法快。

### 3.3 主-Lorentz 算法（Main-Lorentz Algorithm）

Main-Lorentz 算法是一种高效的串联重复查找算法，能在 O(n log n) 时间内找出所有原始串联重复：

#### 算法步骤

1. **分治策略**：将字符串 S 分成两半，递归处理
2. **处理跨越中点的串联重复**
3. **合并结果**

```python
def find_tandem_repeats_main_lorentz(s):
    def find_crossing_tandems(s, center):
        # 查找跨越中点的串联重复
        n = len(s)
        results = []

        # 左侧扩展序列和右侧扩展序列
        left_ext = s[center-1::-1]
        right_ext = s[center:]

        # 计算 LCP（最长公共前缀）数组
        lcp_l = compute_lcp(left_ext)
        lcp_r = compute_lcp(right_ext)

        for period in range(1, min(center, n-center) + 1):
            # 检查 period 是否形成串联重复
            if lcp_l[period] + lcp_r[period] >= period:
                # 找到串联重复
                start = center - period - lcp_l[period]
                length = period
                repeat_count = (lcp_l[period] + lcp_r[period]) // period + 1
                results.append((start, length, repeat_count))

        return results

    def main_lorentz_rec(s, start, end):
        if end - start <= 1:
            return []

        mid = (start + end) // 2

        # 递归处理左半部分和右半部分
        left_results = main_lorentz_rec(s, start, mid)
        right_results = main_lorentz_rec(s, mid, end)

        # 处理跨越中点的串联重复
        crossing = find_crossing_tandems(s[start:end], mid-start)
        # 调整crossing中的位置索引
        crossing = [(start+pos, length, count) for pos, length, count in crossing]

        # 合并结果
        return left_results + right_results + crossing

    return main_lorentz_rec(s, 0, len(s))

def compute_lcp(s):
    # 计算每个位置与字符串开头的最长公共前缀
    n = len(s)
    lcp = [0] * (n+1)
    for i in range(n):
        j = 0
        while i+j < n and s[j] == s[i+j]:
            j += 1
        lcp[i] = j
    return lcp
```

时间复杂度：O(n log² n)，通过更复杂的实现可以优化到 O(n log n)。

### 3.4 后缀数组方法

后缀数组也可以用来高效找出串联重复：

```python
def find_tandem_repeats_suffix_array(s):
    n = len(s)

    # 构建后缀数组
    sa = build_suffix_array(s)
    # 构建LCP数组（相邻后缀的最长公共前缀）
    lcp = build_lcp_array(s, sa)

    results = []

    # 使用LCP数组查找潜在的串联重复
    for i in range(1, n):
        if lcp[i] > 0:
            pos1, pos2 = sa[i-1], sa[i]
            common_len = lcp[i]

            # 检查是否形成串联重复
            for length in range(1, common_len + 1):
                if common_len >= length and abs(pos1 - pos2) == length:
                    start = min(pos1, pos2)
                    repeat_count = common_len // length + 1
                    results.append((start, length, repeat_count))

    return results
```

时间复杂度：O(n log n)，主要受限于后缀数组的构建。

## 4. 高级应用

### 4.1 串联重复的压缩

串联重复可用于字符串压缩：

```python
def compress_string(s):
    n = len(s)
    result = []
    i = 0

    while i < n:
        # 查找从位置i开始的最大串联重复
        max_repeat = 1
        max_length = 0

        for length in range(1, (n-i)//2 + 1):
            repeat_count = 1
            j = i + length

            while j + length <= n and s[j:j+length] == s[i:i+length]:
                repeat_count += 1
                j += length

            if repeat_count > max_repeat or (repeat_count == max_repeat and length > max_length):
                max_repeat = repeat_count
                max_length = length

        if max_repeat > 1:
            result.append(f"{max_repeat}({s[i:i+max_length]})")
            i += max_length * max_repeat
        else:
            result.append(s[i])
            i += 1

    return ''.join(result)
```

### 4.2 DNA序列分析

在生物信息学中，串联重复在DNA序列分析中有重要应用：

```python
def find_dna_tandem_repeats(dna_sequence, min_length=2, min_repeats=2):
    results = find_tandem_repeats_main_lorentz(dna_sequence)

    # 过滤结果，只保留满足最小长度和最小重复次数的
    filtered = [(pos, length, count) for pos, length, count in results
                if length >= min_length and count >= min_repeats]

    # 转换为更可读的格式
    readable_results = []
    for pos, length, count in filtered:
        pattern = dna_sequence[pos:pos+length]
        readable_results.append({
            "position": pos,
            "pattern": pattern,
            "length": length,
            "repeats": count,
            "sequence": dna_sequence[pos:pos+length*count]
        })

    return readable_results
```

## 5. 扩展与优化

### 5.1 近似串联重复

允许一定数量的错误或变异的串联重复：

```python
def find_approximate_tandems(s, error_threshold=1):
    n = len(s)
    results = []

    for i in range(n):
        for length in range(1, (n-i)//2 + 1):
            j = i + length
            errors = 0

            for k in range(length):
                if s[i+k] != s[j+k]:
                    errors += 1

            if errors <= error_threshold:
                results.append((i, length, 2, errors))

    return results
```

### 5.2 最长串联重复

查找字符串中最长的串联重复：

```python
def longest_tandem_repeat(s):
    all_repeats = find_tandem_repeats_main_lorentz(s)
    if not all_repeats:
        return None

    # 按总长度排序
    all_repeats.sort(key=lambda x: x[1] * x[2], reverse=True)
    pos, length, count = all_repeats[0]

    return {
        "position": pos,
        "pattern": s[pos:pos+length],
        "length": length,
        "repeats": count,
        "total_length": length * count,
        "sequence": s[pos:pos+length*count]
    }
```

## 6. 总结

串联重复是字符串处理中一个重要的概念，有广泛的应用：

1. **数据压缩**：识别重复模式可以提高压缩率
2. **生物信息学**：DNA序列中的串联重复有重要生物学意义
3. **密码学**：某些密码分析中需要识别重复模式
4. **自然语言处理**：帮助识别文本中的重复模式

最有效的串联重复查找算法（如Main-Lorentz算法或基于后缀数组的方法）能够在O(n log n)时间内找出所有串联重复，这在处理大规模数据时非常重要。

在实际应用中，通常需要根据具体问题选择合适的算法，并考虑允许的错误率、最小重复长度等因素来优化结果。

---

# 串联重复（Tandem Repeats）、Lyndon 分解和 Runs 的关系

这三个概念都是字符串算法领域中的重要概念，它们之间有着紧密的联系，同时也有各自特定的性质。下面将详细解释它们的关系和区别。

## 1. 概念定义

### 1.1 串联重复（Tandem Repeats）

- **定义**：字符串中连续出现两次或多次的子串模式，形如 `ww...w`（w重复k次）。
- **例子**：在字符串 "abababc" 中，"ab" 的三次重复 "ababab" 是一个串联重复。

### 1.2 Lyndon 词（Lyndon Word）

- **定义**：一个字符串 w 是 Lyndon 词，如果 w 严格小于它的所有非平凡后缀（按字典序）。
- **例子**："aab" 是 Lyndon 词，因为它小于其后缀 "ab" 和 "b"。

### 1.3 Lyndon 分解（Lyndon Factorization）

- **定义**：将字符串 S 分解为 S = w₁w₂...wₖ，其中每个 wᵢ 是 Lyndon 词，且 w₁ ≥ w₂ ≥ ... ≥ wₖ（字典序）。
- **例子**："aababc" 的 Lyndon 分解是 "a·ab·abc"。

### 1.4 Runs（最大周期串）

- **定义**：如果字符串 S[i..j] 具有周期 p，且这个周期性无法向左右扩展（即最大的周期性子串），且满足 j-i+1 ≥ 2p，则称 S[i..j] 是一个 run。
- **例子**：在 "abababc" 中，"ababab" 是一个 run，周期为 2。

## 2. 关系分析

### 2.1 Runs 与 串联重复的关系

Runs 和串联重复密切相关，但有一些关键区别：

1. **Runs 是最大化的**：

   - 每个 run 都是不可扩展的周期性子串
   - 串联重复可能是 run 的一部分，但不一定是最大的

2. **数量关系**：

   - 一个长度为 n 的字符串中，runs 的数量最多为 O(n)
   - 而串联重复的数量可能达到 O(n log n)

3. **包含关系**：
   - 每个 run 都包含至少一个串联重复
   - 不是所有串联重复都是 run（可能是 run 的一部分）

### 2.2 Lyndon 分解与串联重复的关系

1. **结构性差异**：

   - Lyndon 分解关注的是将字符串分解为字典序递减的 Lyndon 词序列
   - 串联重复关注的是重复模式

2. **周期性检测**：

   - Lyndon 分解可以用来检测字符串的周期性
   - 如果字符串 S 的 Lyndon 分解是 wⁿ（即 n 个相同的 Lyndon 词 w），那么 S 是一个串联重复

3. **算法联系**：
   - Duval 算法（计算 Lyndon 分解的算法）可以用来检测某些类型的串联重复

### 2.3 Runs 与 Lyndon 分解的关系

1. **周期性表示**：

   - Runs 直接表示字符串中的周期性部分
   - Lyndon 分解可以间接反映字符串的周期性结构

2. **不重叠性质**：

   - Lyndon 分解将字符串分成不重叠的部分
   - Runs 可能彼此重叠

3. **最小表示**：
   - 每个 run 的一个最小周期可以表示为 Lyndon 词的重复
   - 如果一个 run 的周期是 p，则它可以看作是某个 Lyndon 词（长度可能小于 p）的重复加上可能的前缀

## 3. 通过实例说明关系

考虑字符串 S = "abababcabcabc"：

### 串联重复（Tandem Repeats）：

- "abab" (位置 0-3)："ab" 重复两次
- "ababab" (位置 0-5)："ab" 重复三次
- "abcabc" (位置 7-12)："abc" 重复两次
- 等等...

### Runs：

- "ababab" (位置 0-5)：周期为 2 的 run
- "abcabc" (位置 7-12)：周期为 3 的 run

### Lyndon 分解：

- S = "a·b·ababc·abc·abc"
  - "a"、"b"、"ababc"、"abc"、"abc" 都是 Lyndon 词
  - 且满足 "a" > "b" > "ababc" > "abc" = "abc" (字典序)

## 4. 算法联系

### 4.1 Main-Lorentz 算法：

- 用于高效查找串联重复
- 时间复杂度：O(n log n)

### 4.2 Duval 算法：

- 用于计算 Lyndon 分解
- 时间复杂度：O(n)
- 可以用于某些串联重复的识别

### 4.3 Runs 算法：

- 基于后缀数组或 Main-Lorentz 算法的变体
- 时间复杂度：O(n)

## 5. 总结

- **串联重复**：关注重复模式，数量为 O(n log n)，是字符串中最基本的重复结构
- **Runs**：关注最大周期性子串，数量为 O(n)，每个 run 包含多个串联重复
- **Lyndon 分解**：关注字典序递减分解，提供字符串的另一种结构视角，与周期性有关
