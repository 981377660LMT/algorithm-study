1. 树状数组维护一个 01 序列的前缀和；
2. 从后往前考虑,因为后面的去除后不影响前面的
3. 二分找到`第 k 个 1` 所在的位置

```Python
    moreThan = nums[i]
    left, right = 1, n
    while left <= right:
        mid = (left + right) >> 1
        if bit.query(mid) >= moreThan + 1:
            right = mid - 1
        else:
            left = mid + 1
    res[i] = left
    bit.add(left, -1)
```

结论:

1. tree[i] = a[i-lowbit(i)+1] + ... + a[i]
2. 每个结点 i 维护区间长度为 lowbit(i) 的前缀和

- 树状数组树上二分的思路
  `从左向右走,看 1<<i 是否可取`

```Python
    def bisectLeft(self, k: int) -> int:
        """返回第一个前缀和大于等于k的位置pos

        1 <= pos <= self.size + 1
        """
        curSum, pos = 0, 0
        for i in range(self.bit, -1, -1):
            nextPos = pos + (1 << i)
            if nextPos <= self.size and curSum + self.tree.get(nextPos, 0) < k:
                pos = nextPos
                curSum += self.tree.get(pos, 0)
        return pos + 1


    def longest(self) -> int:
        """求最长上传前缀,即以1为起点的最长连续1区间的长度

        如果当前结点管理的lowbit(i)长度的和为lowbit(i),则可以向右取这个结点
        """
        res = 0
        for i in range(self.bit, -1, -1):
            nextPos = res + (1 << i)
            if nextPos <= self.size and self.tree.get(nextPos, 0) == nextPos & -nextPos:
                res = nextPos
        return res
```
