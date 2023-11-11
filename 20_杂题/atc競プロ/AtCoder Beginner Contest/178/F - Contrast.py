# https://atcoder.jp/contests/abc178/tasks/abc178_f
# !重排数组使得对应位置元素不相等

# 输入 n(≤2e5) 和两个非降数组 a 和 b，元素范围在 [1,n]。
# 如果可以重排数组 b，使得 a[i] != b[i] 对每个 i 都成立，则输出 Yes 和重排后的 b，否则输出 No。

# !我们可以将数组 B 倒置，逐个与数组 A 比较，
# !如果相同就向后找出不同的元素进行换位，
# 找不到直接输出 No，最后输出 B 即可

from typing import List, Tuple


def contrast(nums1: List[int], nums2: List[int]) -> Tuple[List[int], bool]:
    n = len(nums1)
    nums2 = nums2[::-1]

    pos = 0  # !记录最左边元素不同的位置
    for i in range(n):
        if nums1[i] == nums2[i]:
            # !注意nums[i] == nums[pos]的情况 pos交换回去不能相等
            while pos < n and (nums1[i] == nums2[pos] or nums1[i] == nums1[pos]):
                pos += 1
            if pos == n:
                return [], False
            nums2[i], nums2[pos] = nums2[pos], nums2[i]
    return nums2, True


# 反例:
# 2
# 1 1
# 1 2
if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n = int(input())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))
    res, ok = contrast(nums1, nums2)
    if not ok:
        print("No")
        exit(0)
    print("Yes")
    print(*res)
