# https://atcoder.jp/contests/arc120/tasks/arc120_c

# 输入 n(≤2e5) 和两个长为 n 的数组 a b，元素范围在 [0,1e9]。
# 每次操作你可以选择 a 中的两个相邻数字，设 x=a[i], y=a[i+1]，更新 a[i]=y+1, a[i+1]=x-1。
# 输出把 a 变成 b 的最小操作次数，如果无法做到则输出 -1。

# https://atcoder.jp/contests/arc120/submissions/36718795

# !手玩一下可以发现，a[i] 左移 i 位就 +i，右移 i 位就 -i。
# 设 a[i] 最终和 b[j] 匹配，则有 a[i]+i-j=b[j]。
# !移项得 a[i]+i = b[j]+j。
# 设 a'[i] = a[i]+i，b'[i] = b[i]+i。
# 问题变成把 a' 通过邻项交换变成数组 b'，需要的最小操作次数。
# 这可以用树状数组解决，具体见代码。

from typing import List
from minAdjacentSwap import minAdjacentSwap1


def swap2(nums1: List[int], nums2: List[int]) -> int:
    nums1 = [num + index for index, num in enumerate(nums1)]
    nums2 = [num + index for index, num in enumerate(nums2)]
    return minAdjacentSwap1(nums1, nums2)


if __name__ == "__main__":
    n = int(input())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))
    print(swap2(nums1, nums2))
