# Runs（最大周期串）详解

## 1. 基本定义

**Run**（也称为最大周期串或Maximal Periodicity）是字符串中的一个重要概念，它描述了字符串中最大的周期性子串。

### 1.1 形式化定义

一个子串 S[i..j] 是一个 run，如果它满足以下条件：

1. S[i..j] 具有周期 p，其中 p ≤ (j-i+1)/2（即子串长度至少是周期的两倍）
2. 该周期性质不能向左扩展（即 i = 0 或 S[i-1] ≠ S[i-1+p]）
3. 该周期性质不能向右扩展（即 j = |S|-1 或 S[j+1] ≠ S[j+1-p]）
4. p 是 S[i..j] 的最小周期

简而言之，run 是不可扩展的最大周期性子串。

### 1.2 术语说明

- **周期(Period)**：如果对所有有效的 i，S[i] = S[i+p]，则 p 是字符串 S 的一个周期。
- **最小周期(Minimal Period)**：字符串的最小正整数周期。
- **指数(Exponent)**：run 的长度除以其最小周期，表示周期模式重复的次数，可能是分数。

## 2. Runs 的性质

### 2.1 数量性质

**重要定理**：长度为 n 的字符串中，runs 的总数不超过 n。

这个线性上界是字符串算法中的一个重要结果，由 Kolpakov 和 Kucherov 在1999年证明。

### 2.2 结构性质

1. **不重叠性**：不同的 runs 可以重叠，但它们的最小周期或边界必然不同。
2. **包含关系**：runs 之间可能存在包含关系。
3. **周期互斥**：如果两个 runs 重叠且周期不同，那么重叠部分的长度有严格的限制。

## 3. Runs 的检测算法

### 3.1 基于后缀数组的线性时间算法

```python
def compute_lcp(s, sa):
    """计算相邻后缀的最长公共前缀"""
    n = len(s)
    rank = [0] * n
    for i in range(n):
        rank[sa[i]] = i

    lcp = [0] * (n-1)
    h = 0
    for i in range(n):
        if rank[i] > 0:
            j = sa[rank[i] - 1]
            while i + h < n and j + h < n and s[i + h] == s[j + h]:
                h += 1
            lcp[rank[i] - 1] = h
            if h > 0:
                h -= 1
    return lcp

def find_runs(s):
    """查找字符串s中的所有runs"""
    n = len(s)

    # 构建后缀数组和LCP数组
    sa = build_suffix_array(s)  # 假设已实现
    lcp = compute_lcp(s, sa)

    runs = []

    # 使用后缀数组和LCP数组查找潜在周期
    for i in range(n - 1):
        if lcp[i] > 0:
            # 相邻后缀之间的距离可能是周期
            period = sa[i+1] - sa[i]
            if period > 0 and period <= lcp[i] // 2:
                # 检查是否满足run的条件
                left = sa[i]
                right = left + lcp[i] - 1

                # 向左扩展，检查是否可以继续重复
                while left > 0 and s[left-1] == s[left+period-1]:
                    left -= 1

                # 向右扩展，检查是否可以继续重复
                while right < n-1 and s[right+1] == s[right-period+1]:
                    right += 1

                # 确保长度至少是周期的两倍
                if right - left + 1 >= 2 * period:
                    runs.append((left, right, period))

    # 去重，因为一个run可能被多次发现
    runs = list(set(runs))
    return runs
```

### 3.2 基于 Main-Lorentz 算法的变体

Main-Lorentz 算法最初用于检测串联重复，但可以修改为检测 runs：

```python
def find_runs_main_lorentz(s):
    """使用Main-Lorentz算法变体查找runs"""
    n = len(s)
    runs = []

    def extend_run(i, p):
        """从位置i开始，查找周期为p的run"""
        # 向左扩展
        left = i
        while left > 0 and s[left-1] == s[left+p-1]:
            left -= 1

        # 向右扩展
        right = i + p - 1
        while right < n-1 and s[right+1] == s[right-p+1]:
            right += 1

        # 检查是否满足run的条件
        if right - left + 1 >= 2 * p:
            runs.append((left, right, p))

    def find_runs_between(left, mid, right):
        """查找跨越中点的runs"""
        # 构建从中点向左和向右的字符串
        left_str = s[mid-1:left-1:-1] if left <= mid-1 else ""
        right_str = s[mid:right]

        # 计算不同长度的匹配
        for p in range(1, min(len(left_str), len(right_str)) + 1):
            # 检查是否可能形成周期p的run
            if left_str[:p] == right_str[:p]:
                extend_run(mid-p, p)

    def ml_rec(left, right):
        """Main-Lorentz递归查找runs"""
        if right - left <= 1:
            return

        mid = (left + right) // 2

        # 递归处理左右两部分
        ml_rec(left, mid)
        ml_rec(mid, right)

        # 处理跨越中点的runs
        find_runs_between(left, mid, right)

    # 启动递归
    ml_rec(0, n)

    # 去重
    runs = list(set(runs))
    return runs
```

## 4. Runs 的应用

### 4.1 字符串压缩

Runs 可以用于字符串的高效压缩，特别是对于具有大量重复模式的文本：

```python
def compress_using_runs(s):
    """使用runs压缩字符串"""
    runs = find_runs(s)

    # 按起始位置排序
    runs.sort(key=lambda x: x[0])

    result = []
    last_pos = 0

    for left, right, period in runs:
        # 添加runs之间的非重复部分
        if left > last_pos:
            result.append(s[last_pos:left])

        # 添加压缩表示的run
        pattern = s[left:left+period]
        exponent = (right - left + 1) / period
        result.append(f"({pattern}){exponent}")

        last_pos = right + 1

    # 添加最后一部分
    if last_pos < len(s):
        result.append(s[last_pos:])

    return ''.join(result)
```

### 4.2 模式匹配

Runs 可以用于高效的模式匹配算法：

```python
def pattern_matching_with_runs(text, pattern):
    """使用runs进行模式匹配"""
    n, m = len(text), len(pattern)

    # 查找text中的runs
    runs = find_runs(text)

    matches = []

    for left, right, period in runs:
        # 检查pattern是否可能匹配此run
        if right - left + 1 >= m:
            # 周期性匹配检查
            pattern_in_run = True
            for i in range(m):
                if text[(left + i) % period + left] != pattern[i]:
                    pattern_in_run = False
                    break

            if pattern_in_run:
                # 添加所有可能的匹配位置
                for i in range(left, right - m + 2):
                    if text[i:i+m] == pattern:
                        matches.append(i)

    # 检查非run部分的匹配
    # ...

    return sorted(matches)
```

### 4.3 生物信息学应用

在DNA序列分析中，Runs 可以用来识别串联重复序列(Tandem Repeats)，这在基因组研究中非常重要：

```python
def find_dna_repeats(dna_sequence):
    """查找DNA序列中的重复模式"""
    runs = find_runs(dna_sequence)

    significant_repeats = []
    for left, right, period in runs:
        repeat_length = right - left + 1
        exponent = repeat_length / period

        # 筛选生物学上有意义的重复
        if period >= 2 and exponent >= 2.5:  # 例如：至少重复2.5次以上
            significant_repeats.append({
                "position": left,
                "pattern": dna_sequence[left:left+period],
                "length": repeat_length,
                "period": period,
                "exponent": exponent
            })

    return significant_repeats
```

## 5. Runs 与其他字符串概念的比较

### 5.1 Runs vs 串联重复(Tandem Repeats)

- **Runs**：强调最大性，每个run都是不可扩展的
- **串联重复**：任何重复的子串，可能不是最大的
- **数量**：n长字符串中，runs数量是O(n)，串联重复数量是O(n log n)

### 5.2 Runs vs Lyndon分解

- **Runs**：描述字符串的周期性结构
- **Lyndon分解**：将字符串分解为字典序递减的Lyndon词序列
- **关系**：Lyndon词与字符串的最小表示密切相关，可以用于高效查找runs

### 5.3 Runs vs Border

- **Run**：强调周期性子串
- **Border**：既是前缀又是后缀的子串
- **关系**：如果一个字符串有长度为b的border，那么它有周期n-b

## 6. 高级算法和优化

### 6.1 线性时间构建所有Runs

Kolpakov和Kucherov提出了一个真正线性时间的算法来构建所有runs，称为"Main-s-runs"算法：

```python
def compute_all_runs(s):
    """线性时间构建所有runs (伪代码框架)"""
    n = len(s)
    runs = []

    # 1. 计算s的Lempel-Ziv分解
    # 2. 对于每个LZ因子，找到跨越其边界的runs
    # 3. 使用改进的Main-Lorentz算法处理每个LZ因子内的runs

    # 这里是简化的框架，实际实现非常复杂
    # ...

    return runs
```

### 6.2 使用后缀自动机

后缀自动机也可以用来有效地查找runs：

```python
def find_runs_suffix_automaton(s):
    """使用后缀自动机查找runs (概念框架)"""
    # 1. 构建后缀自动机
    # 2. 遍历自动机状态，查找周期性子串
    # 3. 检查是否满足run的条件

    # 这里是简化的框架，实际实现较复杂
    # ...

    return runs
```

## 7. 实际应用中的优化

### 7.1 内存优化

在处理大型文本时，可以使用滑动窗口和流式处理来优化内存使用：

```python
def find_runs_streaming(stream, buffer_size=1000):
    """流式处理查找runs"""
    buffer = []
    runs = []
    global_offset = 0

    while True:
        chunk = stream.read(buffer_size)
        if not chunk:
            break

        buffer.extend(chunk)

        if len(buffer) >= 2 * buffer_size:
            # 处理buffer的前半部分
            local_runs = find_runs(buffer[:buffer_size])

            # 调整偏移量
            for left, right, period in local_runs:
                # 仅保留完全在前半部分的runs
                if right < buffer_size:
                    runs.append((left + global_offset, right + global_offset, period))

            # 滑动窗口
            buffer = buffer[buffer_size:]
            global_offset += buffer_size

    # 处理最后的buffer
    if buffer:
        local_runs = find_runs(buffer)
        for left, right, period in local_runs:
            runs.append((left + global_offset, right + global_offset, period))

    return runs
```

### 7.2 并行化优化

对于大型字符串，可以使用并行处理提高效率：

```python
import multiprocessing

def find_runs_parallel(s, num_processes=4):
    """并行处理查找runs"""
    n = len(s)
    chunk_size = n // num_processes

    with multiprocessing.Pool(num_processes) as pool:
        # 创建重叠的块，确保不会漏掉跨界runs
        chunks = []
        for i in range(num_processes):
            start = i * chunk_size
            end = min(n, (i+1) * chunk_size + chunk_size)  # 重叠一个chunk_size
            chunks.append(s[start:end])

        # 并行处理每个块
        local_runs_lists = pool.map(find_runs, chunks)

        # 合并结果
        all_runs = []
        for i, local_runs in enumerate(local_runs_lists):
            offset = i * chunk_size

            for left, right, period in local_runs:
                # 调整偏移量
                global_left = left + offset
                global_right = right + offset

                # 过滤重叠区域中的重复runs
                if i < num_processes - 1 and global_right >= (i+1) * chunk_size:
                    if global_left >= (i+1) * chunk_size:
                        continue  # 完全在重叠区且会被下一个进程处理

                all_runs.append((global_left, global_right, period))

    # 去重
    all_runs = list(set(all_runs))
    return all_runs
```

## 8. 结论

Runs（最大周期串）是字符串处理中的一个基础概念，它提供了对字符串中周期性模式的精确描述。理解和高效计算runs对于许多字符串算法和应用都至关重要，包括字符串压缩、模式匹配和生物序列分析。

最重要的结果是，一个长度为n的字符串中，runs的数量是线性的O(n)，而且可以在O(n)时间内计算所有runs。这一结果使得许多基于runs的算法能够实现最优的时间复杂度。
