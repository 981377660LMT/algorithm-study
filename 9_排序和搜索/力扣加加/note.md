基本概念

- 解空间
- 序列有序
- 极值

1. 折半**解空间**才是二分法的灵魂
   其他（比如序列有序，左右双指针）都是二分法的手和脚，都是表象，并不是本质，而折半才是二分法的灵魂。
   难点就是上面提到的两点：什么条件 和 舍弃哪部分。

两种类型(**解空间已经明确出来了，如何用代码找出具体的解。**)

- 最左插入 bisect.bisect_left **希望多向左移动**

  ```JS
  const bisectLeft = (arr: number[], target: number): number => {
  if (arr.length === 0) return -1

  let l = 0
  let r = arr.length - 1
  // 因此当 left <= right 的时候，解空间都不为空，此时我们都需要继续搜索
  while (l <= r) {
    const mid = (l + r) >> 1
    const midElement = arr[mid]
    if (midElement === target) {
     r--
    } else if (midElement < target) {
      // mid 根本就不是答案，直接更新 l = mid + 1，从而将 mid 从解空间排除
      l = mid + 1
    } else if (midElement > target)  {
      // midElement >= target :将 mid 从解空间排除，继续看看有没有更好的
      r = mid - 1
    }
  }

  return l
  }
  ```

- 最右插入 bisect.bisect_right **希望多向右移动**
  改成

  ```JS
   } else if (midElement < target) {
      l = mid + 1
    } else if (midElement > target)  {
      r = mid - 1
    }
      }

  return l
  }
  ```

  对于最左和最右二分，简单用两句话总结一下：
  最左二分不断收缩右边界，最终返回左边界
  最右二分不断收缩左边界，最终返回右边界

四大应用(**如何构造解空间。更多的情况则是如何构建有序序列。**)

1. 能力检测二分

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

2. 前缀和二分:如果数组全是正的，那么其前缀和就是一个严格递增的数组
3. 插入排序二分:自行维护有序序列
4. 计数二分

```python
def main(nums, k):
  def count(mid):
    pass
  l, r = 0, len(A) - 1
  while l <= r:
      mid = (l + r) // 2
      # 只有这里和最左二分不一样
      if count(mid) > k: r = mid - 1
      else: l = mid + 1
  return l
```

虽然二分法不意味着需要序列有序，但大多数二分题目都有有序这个显著特征。只不过：
有的是题目直接限定了有序。这种题目通常难度不高，也容易让人想到用二分。
有的是需要你自己构造有序序列。这种类型的题目通常难度不低，需要大家有一定的观察能力。

堆的一种很重要的用法是求第 k 大的数，而二分法也可以求第 k 大的数，只不过二者的思路完全不同。

**有无重复元素对二分算法影响很大，我们需要小心对待。**
