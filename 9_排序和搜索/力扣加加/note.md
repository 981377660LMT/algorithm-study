1. 折半解空间才是二分法的灵魂
   难点就是上面提到的两点：什么条件 和 舍弃哪部分。

两种类型(解空间已经明确出来了，如何用代码找出具体的解。)

- 最左插入 bisect.bisect_left
- 最右插入 bisect.bisect_right

四大应用(如何构造解空间。更多的情况则是如何构建有序序列。)

- 能力检测二分

```Python
def ability_test_bs(nums):
  def possible(mid):
    pass
  l, r = 0, len(A) - 1
  while l <= r:
      mid = (l + r) // 2
      # 只有这里和最左二分不一样
      if possible(mid): l = mid + 1
      else: r = mid - 1
  return l
```

- 前缀和二分
- 插入排序二分（不是你理解的插入排序哦）
- 计数二分
