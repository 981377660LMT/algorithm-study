# 975. 奇偶跳.py
from typing import List


# 对每个数，寻找右侧比自己大的数里最小的那个=>
# 遍历前就应该排序，保证遍历总是不断变大的，只需要比较索引是不是满足i<j
class Solution:
    def oddEvenJumps(self, arr: List[int]) -> None:
        n = len(arr)
        ids = list(range(n))
        ids.sort(key=lambda i: (arr[i], i))
        # 右侧比自己大的里最小的（相同大的，取index小的）index
        next_bigger = [-1] * n
        stack = []
        for id in ids:
            while stack and stack[-1] < id:
                next_bigger[stack.pop()] = id
            stack.append(id)

        ids.sort(key=lambda i: (-arr[i], i))
        # 右侧比自己小的里最大的（相同大的，取index小的）index
        next_smaller = [-1] * n
        stack = []
        for id in ids:
            while stack and stack[-1] < id:
                next_smaller[stack.pop()] = id
            stack.append(id)

        print(next_smaller, next_bigger)


print(Solution().oddEvenJumps([10, 13, 12, 14, 15]))
