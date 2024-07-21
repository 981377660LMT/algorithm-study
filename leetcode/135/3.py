from bisect import bisect_left, bisect_right
from itertools import accumulate
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个长度为 n 的整数数组 nums ，n 是 偶数 ，同时给你一个整数 k 。

# 你可以对数组进行一些操作。每次操作中，你可以将数组中 任一 元素替换为 0 到 k 之间的 任一 整数。

# 执行完所有操作以后，你需要确保最后得到的数组满足以下条件：


# 存在一个整数 X ，满足对于所有的 (0 <= i < n) 都有 abs(a[i] - a[n - i - 1]) = X 。
# 请你返回满足以上条件 最少 修改次数。
#
# 有的pair只需要改一个数，有的pair需要改两个数，有的pair不需要改


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


# TODO: 除开差分有没有更容易理解的set方法
# class Solution {
# public:
#     int minChanges(vector<int>& nums, int K) {
#         int n = nums.size();
#         // 差分数组
#         int f[K + 2];
#         memset(f, 0, sizeof(f));
#         // 枚举每个数对
#         for (int i = 0, j = n - 1; i < j; i++, j--) {
#             int d = abs(nums[i] - nums[j]);
#             int mx = max({nums[i], K - nums[i], nums[j], K - nums[j]});
#             // 0 <= x < d 时需要一次操作
#             f[0]++; f[d]--;
#             // d < x <= mx 时需要一次操作
#             f[d + 1]++; f[mx + 1]--;
#             // x > mx 时需要两次操作
#             f[mx + 1] += 2;
#         }

#         int ans = n;
#         // 枚举 x 的取值，看最少需要几次操作
#         for (int i = 0, now = 0; i <= K + 1; i++) {
#             now += f[i];
#             ans = min(ans, now);
#         }
#         return ans;
#     }
# };


class Solution:
    def minChanges(self, nums: List[int], k: int) -> int:
        diff = [0] * (k + 2)  # diff为x时，需要的操作数

        def add(left: int, right: int, x: int):
            diff[left] += x
            diff[right + 1] -= x

        n = len(nums)
        for i in range(n // 2):
            if nums[i] > nums[~i]:
                nums[i], nums[~i] = nums[~i], nums[i]
            a, b = nums[i], nums[~i]
            add(0, k, 2)
            add(0, max2(k - a, b), -1)
            add(b - a, b - a, -1)

        diff = list(accumulate(diff))
        return min(diff[:-1])


# nums = [1,0,1,2,4,3], k = 4

print(Solution().minChanges([1, 0, 1, 2, 4, 3], 4))
# [9,2,7,7,8,9,1,5,1,9,4,9,4,7]
# 9

print(Solution().minChanges([9, 2, 7, 7, 8, 9, 1, 5, 1, 9, 4, 9, 4, 7], 9))
# 4
