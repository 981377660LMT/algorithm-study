from typing import List


# 对每个数，寻找右侧比自己大的数里最小的那个=>
# 遍历前就应该排序，保证遍历总是不断变大的，只需要比较索引是不是满足i<j
class Solution:
    def oddEvenJumps(self, arr: List[int]) -> int:
        n = len(arr)
        ids = list(range(n))
        ids.sort(key=lambda i: (arr[i], i))
        # 右侧比自己大的里最小的（相同大的，取index小的）index
        odd_next = [-1] * n
        stack = []
        for id in ids:
            while stack and stack[-1] < id:
                odd_next[stack.pop()] = id
            stack.append(id)

        ids.sort(key=lambda i: (-arr[i], i))
        # 右侧比自己小的里最大的（相同大的，取index小的）index
        even_next = [-1] * n
        stack = []
        for id in ids:
            while stack and stack[-1] < id:
                even_next[stack.pop()] = id
            stack.append(id)

        print(even_next, odd_next)

        # 能跳到最后的索引，起跳是奇数还是偶数次
        odd = [False] * n
        even = [False] * n
        odd[n - 1] = True
        even[n - 1] = True

        for i in range(n - 2, -1, -1):
            if odd_next[i] != -1:
                # 奇数跳完，下一次是偶数跳
                odd[i] = even[odd_next[i]]
            if even_next[i] != -1:
                even[i] = odd[even_next[i]]

        # 题目中第一次起跳是奇数跳
        return sum(odd)


print(Solution().oddEvenJumps([10, 13, 12, 14, 15]))


# 第一次(奇数次)跳最近最小的大
# 偶数次跳最近最大的小
# (尽量省力)
# 单调栈解决的是右侧第一个比当前大/小的idnex
# 本题是右侧比当前“最大”/“最小”的index
