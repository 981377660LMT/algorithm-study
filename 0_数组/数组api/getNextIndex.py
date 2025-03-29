# getPrevIndex/getNextIndex
# 获取数组中相同元素的前一个和后一个元素的位置


from typing import List, Tuple


def getPrevIndex(arr: List[int]) -> List[int]:
    """
    获取数组中相同元素的前一个元素的位置.不存在则返回-1.
    """
    n = len(arr)
    left = [0] * n
    last = dict()
    for i, x in enumerate(arr):
        left[i] = last.get(x, -1)
        last[x] = i
    return left


def getNextIndex(arr: List[int]) -> List[int]:
    """
    获取数组中相同元素的后一个元素的位置.不存在则返回-1.
    """
    n = len(arr)
    right = [0] * n
    last = dict()
    for i in range(n - 1, -1, -1):
        v = arr[i]
        right[i] = last.get(v, -1)
        last[v] = i
    return right


def getPrevAndNextIndex(arr: List[int]) -> Tuple[List[int], List[int]]:
    """
    获取数组中相同元素的前一个和后一个元素的位置.不存在则返回-1.
    """
    n = len(arr)
    left, right = [0] * n, [-1] * n
    last = dict()
    for i, x in enumerate(arr):
        left[i] = j = last.get(x, -1)
        if j >= 0:
            right[j] = i
        last[x] = i
    return left, right


if __name__ == "__main__":
    assert getPrevIndex([1, 2, 3, 4, 5, 1, 2, 3, 4, 5]) == [-1, -1, -1, -1, -1, 0, 1, 2, 3, 4]
    assert getNextIndex([1, 2, 3, 4, 5, 1, 2, 3, 4, 5]) == [5, 6, 7, 8, 9, -1, -1, -1, -1, -1]
    assert getPrevAndNextIndex([1, 2, 3, 4, 5, 1, 2, 3, 4, 5]) == (
        [-1, -1, -1, -1, -1, 0, 1, 2, 3, 4],
        [5, 6, 7, 8, 9, -1, -1, -1, -1, -1],
    )

    # 3488. 距离最小相等元素查询
    # https://leetcode.cn/problems/closest-equal-element-queries/description/
    class Solution:
        def solveQueries(self, nums: List[int], queries: List[int]) -> List[int]:
            INF = int(1e18)
            res = [-1] * len(queries)
            left, right = getPrevAndNextIndex(nums + nums)

            def distToSame(pos: int) -> int:
                if left[pos] == -1 and right[pos] == -1:
                    return INF
                if left[pos] == -1:
                    return right[pos] - pos
                if right[pos] == -1:
                    return pos - left[pos]
                return min(pos - left[pos], right[pos] - pos)

            for i, pos in enumerate(queries):
                d1, d2 = distToSame(pos), distToSame(pos + len(nums))
                min_ = min(d1, d2)
                if min_ < len(nums):
                    res[i] = min_
            return res
