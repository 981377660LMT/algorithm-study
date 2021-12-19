import collections
from typing import List

# A[i] + rev(A[j]) == A[j] + rev(A[i])
# A[i] - rev(A[i]) == A[j] - rev(A[j])
# B[i] = A[i] - rev(A[i])

# Then it becomes an easy question that,
# how many pairs in B with B[i] == B[j]
def countNicePairs(A: List[int]):
    res = 0
    counter = collections.Counter()
    # 保证i<j
    for a in A:
        b = int(str(a)[::-1])
        res += counter[a - b]
        counter[a - b] += 1
    return res % (10 ** 9 + 7)


print(countNicePairs(A=[42, 11, 1, 97]))
# 输出：2
# 解释：两个坐标对为：
#  - (0,3)：42 + rev(97) = 42 + 79 = 121, 97 + rev(42) = 97 + 24 = 121 。
#  - (1,2)：11 + rev(1) = 11 + 1 = 12, 1 + rev(11) = 1 + 11 = 12 。
