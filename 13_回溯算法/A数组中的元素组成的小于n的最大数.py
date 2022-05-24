# 数组A中给定可以使用的1~9的数，返回由A数组中的元素组成的小于n的最大数。
# 例如A={1, 2, 4, 9}，n=2533，返回2499
# https://leetcode.cn/circle/discuss/BlvA0l/
from typing import List


def getLower(target: int, A: List[int]) -> int:
    # A数组中的元素组成的小于target的最大数 如果不存在 输出-1
    def dfs(index: int, isLimit: bool, path: List[int]) -> None:
        """index表示当前数字的位置,isLimit表示是否被原数字限制"""
        nonlocal res2
        if index == n:
            cand = int(''.join(map(str, path)))
            if cand < target:
                res2 = max(res2, cand)
            return
        upper = digits[index] if isLimit else maxA
        for cur in A:
            if cur > upper:
                break
            path.append(cur)
            dfs(index + 1, (isLimit and cur == upper), path)
            path.pop()

    n, A = len(str(target)), sorted(A)
    digits = list(map(int, str(target)))
    maxA = max(A)

    res1, res2 = int((n - 1) * str(maxA) or -1), -1  # n-1位最大数和n位最大数
    dfs(0, True, [])
    return max(res1, res2)


print(getLower(2533, [1, 2, 4, 9]))
print(getLower(12121, [1, 2]))
