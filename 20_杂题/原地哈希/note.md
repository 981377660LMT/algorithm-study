# 原地哈希（Index as Hash Key）

## 核心思想

对于长度为 n 的数组，元素值域在 `[1, n]`（或 `[0, n-1]`）时，可以利用**数组下标本身作为哈希表**，将每个元素放到它"该在"的位置上，从而实现 O(n) 时间、O(1) 空间。

两种常见手法：

1. **交换归位法**：不断 swap，直到 `nums[i] == i+1`（或 `nums[i] == i`）。
2. **取反标记法**：遍历时将 `nums[|val|-1]` 取反为负数，标记该值出现过。

---

## 题目汇总

| #      | 题目                     | 链接                                                                   | 手法     |
| ------ | ------------------------ | ---------------------------------------------------------------------- | -------- |
| 41     | 缺失的第一个正数         | https://leetcode.cn/problems/first-missing-positive/                   | 交换归位 |
| 442    | 数组中重复的数据         | https://leetcode.cn/problems/find-all-duplicates-in-an-array/          | 取反标记 |
| 448    | 找到所有数组中消失的数字 | https://leetcode.cn/problems/find-all-numbers-disappeared-in-an-array/ | 取反标记 |
| 645    | 错误的集合               | https://leetcode.cn/problems/set-mismatch/                             | 取反标记 |
| 287    | 寻找重复数               | https://leetcode.cn/problems/find-the-duplicate-number/                | 交换归位 |
| 268    | 丢失的数字               | https://leetcode.cn/problems/missing-number/                           | 交换归位 |
| 765    | 情侣牵手                 | https://leetcode.cn/problems/couples-holding-hands/                    | 交换归位 |
| 剑指03 | 数组中重复的数字         | https://leetcode.cn/problems/shu-zu-zhong-zhong-fu-de-shu-zi-lcof/     | 交换归位 |

---

## Python 解法

### 41. 缺失的第一个正数

```python
# https://leetcode.cn/problems/first-missing-positive/
# 把每个正数 x 交换到下标 x-1 的位置
class Solution:
    def firstMissingPositive(self, nums: list[int]) -> int:
        n = len(nums)
        for i in range(n):
            while 1 <= nums[i] <= n and nums[i] != nums[nums[i] - 1]:
                nums[nums[i] - 1], nums[i] = nums[i], nums[nums[i] - 1]
        return next((i + 1 for i, v in enumerate(nums) if v != i + 1), n + 1)
```

### 442. 数组中重复的数据

```python
# https://leetcode.cn/problems/find-all-duplicates-in-an-array/
# 遍历时将 nums[|val|-1] 取反，若已为负说明 val 出现过两次
class Solution:
    def findDuplicates(self, nums: list[int]) -> list[int]:
        res = []
        for x in nums:
            i = abs(x) - 1
            if nums[i] < 0:
                res.append(i + 1)
            else:
                nums[i] = -nums[i]
        return res
```

### 448. 找到所有数组中消失的数字

```python
# https://leetcode.cn/problems/find-all-numbers-disappeared-in-an-array/
# 出现过的值对应下标置负，最后仍为正的下标即为消失的数字
class Solution:
    def findDisappearedNumbers(self, nums: list[int]) -> list[int]:
        for x in nums:
            i = abs(x) - 1
            nums[i] = -abs(nums[i])
        return [i + 1 for i, v in enumerate(nums) if v > 0]
```

### 645. 错误的集合

```python
# https://leetcode.cn/problems/set-mismatch/
# 取反标记：已为负的位置对应重复数，仍为正的位置对应丢失数
class Solution:
    def findErrorNums(self, nums: list[int]) -> list[int]:
        dup = -1
        for x in nums:
            i = abs(x) - 1
            if nums[i] < 0:
                dup = i + 1
            else:
                nums[i] = -nums[i]
        miss = next(i + 1 for i, v in enumerate(nums) if v > 0)
        return [dup, miss]
```

### 287. 寻找重复数

```python
# https://leetcode.cn/problems/find-the-duplicate-number/
# 原地哈希解法（会修改数组；不修改数组可用快慢指针）
class Solution:
    def findDuplicate(self, nums: list[int]) -> int:
        for x in nums:
            i = abs(x)
            if nums[i] < 0:
                return i
            nums[i] = -nums[i]
        return -1
```

### 268. 丢失的数字

```python
# https://leetcode.cn/problems/missing-number/
# 原地哈希：将值 v 交换到下标 v，最后 nums[i] != i 的就是丢失的
class Solution:
    def missingNumber(self, nums: list[int]) -> int:
        n = len(nums)
        for i in range(n):
            while nums[i] < n and nums[i] != i:
                nums[nums[i]], nums[i] = nums[i], nums[nums[i]]
        return next((i for i in range(n) if nums[i] != i), n)
```

### 765. 情侣牵手

```python
# https://leetcode.cn/problems/couples-holding-hands/
# 贪心交换：对于位置 [0,1],[2,3],... 如果不是情侣就把正确的人换过来
class Solution:
    def minSwapsCouples(self, row: list[int]) -> int:
        # pos[v] = v 当前所在下标
        pos = [0] * len(row)
        for i, v in enumerate(row):
            pos[v] = i
        res = 0
        for i in range(0, len(row), 2):
            # row[i] 的伴侣
            partner = row[i] ^ 1
            if row[i + 1] != partner:
                # 把 partner 换到 i+1
                j = pos[partner]
                pos[row[i + 1]] = j
                row[i + 1], row[j] = row[j], row[i + 1]
                pos[partner] = i + 1
                res += 1
        return res
```

### 剑指 Offer 03. 数组中重复的数字

```python
# https://leetcode.cn/problems/shu-zu-zhong-zhong-fu-de-shu-zi-lcof/
# 交换归位：nums[i] 应该放在下标 nums[i]，冲突即为重复
class Solution:
    def findRepeatNumber(self, nums: list[int]) -> int:
        for i in range(len(nums)):
            while nums[i] != i:
                if nums[nums[i]] == nums[i]:
                    return nums[i]
                nums[nums[i]], nums[i] = nums[i], nums[nums[i]]
        return -1
```
