# abc334c-删除一个元素后的最小相邻差之和
# https://atcoder.jp/contests/abc334/tasks/abc334_c


from itertools import accumulate
from typing import List


def minAdjacentDiffSum(nums: List[int]) -> int:
    """
    给定一个奇数长度的数组.
    从数组中删除一个元素后，求最小相邻差之和.
    相邻差之和为:sum(abs(nums[i] - nums[i + 1]) for i in range(0, len(nums) - 1, 2)
    """

    def makePreSum(seq: List[int]) -> List[int]:
        n = len(seq)
        diff = [abs(seq[2 * i + 1] - seq[2 * i]) for i in range(n // 2)]
        preSum = [0] + list(accumulate(diff))
        return preSum

    sum1 = makePreSum(nums)
    sum2 = makePreSum(nums[::-1])[::-1]
    return min(sum1[i] + sum2[i] for i in range(len(sum1)))


if __name__ == "__main__":
    N, K = map(int, input().split())
    indexes = list(map(int, input().split()))
    counter = [2] * N
    for v in indexes:
        counter[v - 1] -= 1
    newNums = []
    for i, v in enumerate(counter):
        newNums.extend([i] * v)

    if len(newNums) % 2 == 0:
        print(sum(abs(newNums[i] - newNums[i + 1]) for i in range(0, len(newNums) - 1, 2)))
    else:
        print(minAdjacentDiffSum(newNums))
