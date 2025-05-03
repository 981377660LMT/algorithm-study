# Alias Method O(1) 采样算法详解

Alias Method 是一种高效的随机采样算法，它允许我们以 O(1) 的时间复杂度从任意离散概率分布中进行采样。该方法由Michael Vose改进自Walker的工作，因此也被称为Walker-Vose别名方法。

## 基本原理

Alias Method 的核心思想是：

1. 将一个非均匀概率分布转换成一个均匀分布+二元选择的问题
2. 通过预处理构建两个表（概率表和别名表）
3. 采样时只需常数时间

## 详细流程

### 1. 预处理阶段 (O(n))

1. 创建n个大小相等的"桶"，每个桶包含两个候选项：原始元素和其"别名"
2. 将原始概率乘以n进行缩放
3. 将概率<1的元素与概率>1的元素配对，填充桶
4. 构建概率表和别名表

### 2. 采样阶段 (O(1))

1. 随机选择一个桶（均匀分布）
2. 用一个随机数决定是选择原始元素还是其别名

## Python 实现

```python
import random

class AliasMethod:
    def __init__(self, probabilities):
        """
        初始化Alias Method采样器

        参数:
            probabilities: 概率分布列表，所有概率之和应为1
        """
        n = len(probabilities)
        self.prob = [0] * n  # 概率表
        self.alias = [0] * n  # 别名表

        # 缩放概率（乘以n）
        scaled_prob = [p * n for p in probabilities]

        # 创建两个工作队列
        small = []  # 存储概率小于1的索引
        large = []  # 存储概率大于1的索引

        # 将索引分配到两个队列中
        for i, prob in enumerate(scaled_prob):
            if prob < 1.0:
                small.append(i)
            else:
                large.append(i)

        # 主循环: 配对小概率和大概率元素
        while small and large:
            l = small.pop()  # 小概率元素
            g = large.pop()  # 大概率元素

            self.prob[l] = scaled_prob[l]  # 存储小概率元素的概率
            self.alias[l] = g  # 小概率元素的别名是大概率元素

            # 更新大概率元素剩余的概率
            scaled_prob[g] = (scaled_prob[g] + scaled_prob[l]) - 1.0

            # 检查更新后的大概率元素去向
            if scaled_prob[g] < 1.0:
                small.append(g)
            else:
                large.append(g)

        # 处理剩余元素（数值误差处理）
        while large:
            self.prob[large.pop()] = 1.0

        while small:
            self.prob[small.pop()] = 1.0

    def pick(self):
        """以O(1)时间采样一个元素的索引"""
        # 均匀随机选择一个桶
        i = random.randint(0, len(self.prob) - 1)

        # 根据概率决定返回原始元素还是别名元素
        if random.random() < self.prob[i]:
            return i
        else:
            return self.alias[i]


# 使用示例
def test_alias_method():
    # 示例概率分布
    probs = [0.1, 0.3, 0.5, 0.1]

    # 初始化采样器
    sampler = AliasMethod(probs)

    # 进行大量采样并统计结果
    samples = 1000000
    results = [0] * len(probs)

    for _ in range(samples):
        results[sampler.pick()] += 1

    # 输出采样结果与理论概率对比
    print("元素  理论概率  实际频率")
    for i in range(len(probs)):
        frequency = results[i] / samples
        print(f"{i}    {probs[i]:.2f}      {frequency:.4f}")

# 运行测试
test_alias_method()
```

## 复杂度分析

- **空间复杂度**：O(n)，需要存储两个长度为n的表
- **预处理时间**：O(n)，只需遍历每个元素一次并进行常数时间操作
- **采样时间**：O(1)，每次采样只需生成两个随机数并进行一次比较

## 应用场景

Alias Method 特别适用于需要多次从同一个概率分布中进行采样的场景，例如：

- 蒙特卡洛模拟
- 随机游走算法
- 粒子滤波器
- 随机图生成
- 游戏中的随机掉落系统

这种方法的效率远高于传统的累积分布函数 (CDF) 采样方法，尤其在大量采样的情况下。

---

好的，这里用通俗易懂的方式讲解 Alias Method（别名法）：

---

### 1. 问题背景

假如你有一个装着不同颜色球的袋子，每种颜色的球数量不一样。你想要按照每种颜色的概率，快速地随机抓一个球。  
如果直接用累加概率的方法，每次都要遍历一遍概率数组，效率是 O(n)。  
Alias Method 可以让你每次抓球都只需要 O(1) 时间。

---

### 2. 核心思想

把每种颜色的概率分成 n 份（n 是颜色种类数），让每一份都“装满”。  
如果某一份没装满，就用别的颜色来“补满”。  
这样，每一份都只需要记住两个信息：

- 主要颜色（本来的颜色）
- 别名颜色（用来补满的颜色）

---

### 3. 采样过程

1. 随机选一个“份”（桶），每份被选中的概率一样。
2. 再用一个随机数决定，是选主要颜色，还是选别名颜色。

---

### 4. 举个例子

假设有三种颜色，概率分别是 [0.5, 0.3, 0.2]。

- 先把每个概率乘以3（颜色数），得到 [1.5, 0.9, 0.6]。
- 0.9 和 0.6 都小于1，1.5大于1。
- 用0.6去补1.5，0.6+0.9=1.5，最后每个桶都满了1。

每个桶都记下：

- 主要颜色是谁
- 别名颜色是谁
- 主要颜色的概率是多少

---

### 5. 代码简化版

```python
import random

class AliasMethod:
    def __init__(self, probs):
        n = len(probs)
        self.prob = [0] * n
        self.alias = [0] * n
        scaled = [p * n for p in probs]
        small, large = [], []
        for i, p in enumerate(scaled):
            (small if p < 1 else large).append(i)
        while small and large:
            s, l = small.pop(), large.pop()
            self.prob[s] = scaled[s]
            self.alias[s] = l
            scaled[l] = scaled[l] + scaled[s] - 1
            (small if scaled[l] < 1 else large).append(l)
        for i in large + small:
            self.prob[i] = 1

    def pick(self):
        i = random.randint(0, len(self.prob) - 1)
        return i if random.random() < self.prob[i] else self.alias[i]
```

---

### 6. 总结

- 预处理一次，采样永远 O(1)
- 适合大量采样的场景
- 只需要两个数组就能实现

---

https://leetcode.cn/problems/random-point-in-non-overlapping-rectangles/solutions/1543639/o1-pickde-suan-fa-by-hqztrue-vupd/
https://en.wikipedia.org/wiki/Alias_method

---

原因如下：

### 1. 等概率分桶

我们把所有概率分成 n 份（n 是元素个数），每一份（桶）被选中的概率是一样的（1/n）。这样第一步就保证了每个桶被选中的机会均等。

### 2. 桶内概率再分配

每个桶里有两种可能：

- 主要元素（原本属于这个桶的元素）
- 别名元素（用来补满的元素）

每个桶都记录：

- 主要元素的概率（可能小于1）
- 剩下的用别名元素补满

### 3. 二次随机选择

- 先随机选一个桶（概率1/n）
- 再用一个随机数决定选主要元素还是别名元素（概率分别是桶内的主要概率和补满概率）

### 4. 总概率推导

假设第i个元素的原始概率是p[i]，经过分桶和补满后，

- 主要元素被选中的概率 = 选中这个桶的概率 × 主要概率 = (1/n) × prob[i]
- 别名元素被选中的概率 = 选中这个桶的概率 × (1 - prob[i])，但别名元素会被分到多个桶，所有桶里补满它的概率加起来，正好等于它的原始概率。

**所以，每个元素最终被选中的概率，正好等于它的原始概率。**

---

**总结**：  
别名法的本质是把复杂的概率分布，拆成“等概率选桶+桶内二选一”，这样每次采样都只需O(1)时间，且概率分布完全正确。
