from typing import List, Tuple


def helper(arr: List[int]) -> Tuple[List[int], List[int]]:
    """寻找每个元素右侧比自己大的里最小的和右侧比自己小的里最大的"""
    n = len(arr)
    ids = list(range(n))
    ids.sort(key=lambda i: (arr[i], i))

    # 右侧比自己大的里最小的（相同大的，取index小的）index
    nextBigger = [-1] * n
    stack = []
    for id in ids:
        while stack and stack[-1] < id:
            nextBigger[stack.pop()] = id
        stack.append(id)

    ids.sort(key=lambda i: (-arr[i], i))
    # 右侧比自己小的里最大的（相同大的，取index小的）index
    nextSmaller = [-1] * n
    stack = []
    for id in ids:
        while stack and stack[-1] < id:
            nextSmaller[stack.pop()] = id
        stack.append(id)

    return nextSmaller, nextBigger


class Solution:
    def oddEvenJumps(self, arr: List[int]) -> int:
        n = len(arr)
        nextSmaller, nextBigger = helper(arr)

        # 能跳到最后的索引，起跳是奇数还是偶数次
        odd = [False] * n
        even = [False] * n
        odd[n - 1] = True
        even[n - 1] = True

        for i in range(n - 2, -1, -1):
            if nextBigger[i] != -1:
                # 奇数跳完，下一次是偶数跳
                odd[i] = even[nextBigger[i]]
            if nextSmaller[i] != -1:
                even[i] = odd[nextSmaller[i]]

        # 题目中第一次起跳是奇数跳
        return sum(odd)


print(Solution().oddEvenJumps([10, 13, 12, 14, 15]))

