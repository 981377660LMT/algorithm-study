# https://atcoder.jp/contests/arc133/tasks/arc133_c

# 输入 n m k (≤2e5) 和长度分别为 n 和 m 的数组 a 和 b，元素范围 [0,k-1]。
# 构造一个 n 行 m 列，元素范围在 [0,k-1] 的矩阵，
# 使得第 i 行的元素和 % k = a[i]，第 j 列的元素和 % k = b[j]。
# 你只需要输出这个矩阵的元素和的最大值。
# 如果这个矩阵不存在，输出 -1。

# https://atcoder.jp/contests/arc133/submissions/36747125

# 不妨把矩阵的每个数都置为 k-1，然后慢慢减小。

# 对于第 i 行来说，这一行的元素总共需要减小 ((k-1)*m-a[i])%k，累加所有行的减小量，得到 sa。
# 对于第 j 列来说，这一列的元素总共需要减小 ((k-1)*n-b[j])%k，累加所有列的减小量，得到 sb。

# 如果 sa%k != sb%k，则无解。
# 否则 sa-sb 是 k 的倍数，不妨设 sa > sb，那么可以把第一列的减小量 += sa-sb，使得行列减小量之和相同。
# !然后就可以不断选择行列的减小量均为正数的，把这个元素减一。
# 这样最后会操作 max(sa,sb) 次。
# 所以答案是 n*m*(k-1) - max(sa,sb)。

from typing import List


def rowColumnSums(nums1: List[int], nums2: List[int], k: int) -> int:
    # https://atcoder.jp/contests/arc133/submissions/28659748
    ROW, COL = len(nums1), len(nums2)
    nums1 = [((k - 1) * COL - num) % k for num in nums1]
    nums2 = [((k - 1) * ROW - num) % k for num in nums2]
    sum1, sum2 = sum(nums1), sum(nums2)
    if sum1 % k != sum2 % k:
        return -1
    return ROW * COL * (k - 1) - max(sum1, sum2)


if __name__ == "__main__":
    _, _, k = map(int, input().split())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))
    print(rowColumnSums(nums1, nums2, k))
