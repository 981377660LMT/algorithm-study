from typing import List

# 1.将 nums 中每个数字 裁剪 到剩下 最右边 trimi 个数位。
# 2.在裁剪过后的数字中，找到 nums 中第 ki 小数字对应的 下标 。
# 如果两个裁剪后数字一样大，那么下标 更小 的数字视为更小的数字。
# 3.将 nums 中每个数字恢复到原本字符串。

# 请你返回一个长度与 queries 相等的数组 answer，其中 answer[i]是第 i 次查询的结果。
# 1 <= nums.length <= 100
# 1 <= queries.length <= 100

# !基数排序，复杂度 O(nm)，按照每一个位排序
# !把第 i - 1 轮的结果，根据 nums 中右数第 i 位数，依次放入桶中
# !把每个桶的结果连接起来，成为第 i 轮的结果


class Solution:
    def smallestTrimmedNumbers(self, nums: List[str], queries: List[List[int]]) -> List[int]:
        ROW, COL = len(nums), len(nums[0])
        rank = [[] for _ in range(COL + 5)]  # !第i轮排序第j小的数对应的下标
        rank[0] = list(range(ROW))

        # !把第 i - 1 轮的结果，根据 nums 中右数第 i 位数，依次放入桶中
        # !把每个桶的结果连接起来，成为第 i 轮的结果
        for col in range(1, COL + 1):
            buckets = [[] for _ in range(10)]
            for index in rank[col - 1]:
                buckets[int(nums[index][-col])].append(index)
            for bucket in buckets:
                rank[col].extend(bucket)

        res = []
        for k, trim in queries:
            res.append(rank[trim][k - 1])
        return res


print(
    Solution().smallestTrimmedNumbers(
        nums=["102", "473", "251", "814"], queries=[[1, 1], [2, 3], [4, 2], [1, 2]]
    )
)
