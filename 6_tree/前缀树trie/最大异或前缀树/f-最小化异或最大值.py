from typing import List

# !最小化异或最大值
# !给定一个数组,选择一个非负整数x,使得数组中的每个数都与x异或后的最大值最小化
# n<=1e5 nums[i]<=1e9
# 求出最小化后的最大值

# 解法:
# !递归,对每一位选1还是0,使得异或最大值最小


def minimizeXor(nums: List[int]) -> int:
    def dfs(curGroup: List[int], pos: int) -> int:
        if pos == -1:
            return 0

        zero, one = [], []
        for num in curGroup:
            if (num >> pos) & 1:
                one.append(num)
            else:
                zero.append(num)

        if not zero:
            return dfs(one, pos - 1)
        if not one:
            return dfs(zero, pos - 1)
        return (1 << pos) + min(dfs(one, pos - 1), dfs(zero, pos - 1))

    return dfs(nums, 32)


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n = int(input())
    nums = list(map(int, input().split()))
    print(minimizeXor(nums))
