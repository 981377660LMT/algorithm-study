from itertools import permutations
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def nextPermutation(nums) -> bool:
    """O(n)返回下一个字典序的排列."""
    left = right = len(nums) - 1
    while left > 0 and nums[left - 1] >= nums[left]:  # 1. 找到最后一个递增位置
        left -= 1
    if left == 0:  # 全部递减
        return False
    last = left - 1  # 最后一个递增位置
    while nums[right] <= nums[last]:  # 2. 找到最小的可交换的right，交换这两个数
        right -= 1
    nums[last], nums[right] = nums[right], nums[last]
    left, right = last + 1, len(nums) - 1  # 3. 翻转后面间这段递减数列
    while left < right:
        nums[left], nums[right] = nums[right], nums[left]
        left += 1
        right -= 1
    return True


if __name__ == "__main__":
    N, K = map(int, input().split())
    S = input()
    sb = sorted(S)
    res = 0
    while True:
        flag = True
        for i in range(N - K + 1):
            for j in range((K + 1) // 2):
                if S[i + j] != S[i + K - 1 - j]:
                    flag = False
                    break
            if not flag:
                break
        if flag:
            res += 1
        if not nextPermutation(sb):
            break

    print(res)
