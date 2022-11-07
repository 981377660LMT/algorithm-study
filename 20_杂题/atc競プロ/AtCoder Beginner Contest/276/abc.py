import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


# (1,…,N) の順列 P=(P
# 1
# ​
#  ,…,P
# N
# ​
#  ) が与えられます。ただし、(P
# 1
# ​
#  ,…,P
# N
# ​
#  )
# 
# =(1,…,N) です。

# (1…,N) の順列を全て辞書順で小さい順に並べたとき、P が K 番目であるとします。辞書順で小さい方から K−1 番目の順列を求めてください。


def prePermutation(nums, inPlace=False):
    """返回前一个字典序的排列,如果不存在,返回本身;时间复杂度O(n)"""
    if not inPlace:
        nums = nums[:]

    left = right = len(nums) - 1

    while left > 0 and nums[left - 1] <= nums[left]:  # 1. 找到最后一个递减位置
        left -= 1
    if left == 0:  # 全部递增
        return False, nums
    last = left - 1  # 最后一个递减位置

    while nums[right] >= nums[last]:  # 2. 找到最大的可交换的right，交换这两个数
        right -= 1
    nums[last], nums[right] = nums[right], nums[last]

    left, right = last + 1, len(nums) - 1  # 3. 翻转后面间这段递增数列
    while left < right:
        nums[left], nums[right] = nums[right], nums[left]
        left += 1
        right -= 1
    return True, nums


if __name__ == "__main__":
    n = int(input())
    perms = list(map(int, input().split()))
    print(*prePermutation(perms)[1], sep=" ")
