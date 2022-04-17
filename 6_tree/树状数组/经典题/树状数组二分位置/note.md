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
