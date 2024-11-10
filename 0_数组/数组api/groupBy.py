# 分组循环
# https://leetcode.cn/problems/longest-even-odd-subarray-with-threshold/description/
#
# - 适用场景：
#   按照题目要求，数组会被分割成若干组，且每一组的判断/处理逻辑是一样的
# - 流程：
#   !1. 根据题目某些条件，得到可能的分组结果
#   !2. 根据题目某些条件，对分组结果进行验证/处理

from typing import Callable, Generator, List, Tuple


def groupWhile(
    n: int, predicate: Callable[[int, int], bool], skipFalsySingleValueGroup=False
) -> Generator[Tuple[int, int], None, None]:
    """
    分组循环.
    :param n: 数据流长度.
    :param predicate: `[left, curRight]` 闭区间内的元素是否能分到一组.
    :param skipFalsySingleValueGroup: 是否跳过`predicate`返回`False`的单个元素的分组.
    :yield: `[start, end)` 左闭右开分组结果.
    """
    end = 0
    while end < n:
        start = end
        while end < n and predicate(start, end):
            end += 1
        if end == start:
            end += 1
            if skipFalsySingleValueGroup:
                continue
        yield start, end


def snippet(nums):
    def predicate(left: int, curRight: int) -> bool:
        return False

    res = -1
    end, n = 0, len(nums)
    while end < n:
        start = end  # 开始分组
        while end < n and predicate(start, end):  # [start, end] 闭区间内的元素是否能分到同一组
            end += 1

        if end == start:  # 非法的单个元素的分组
            end += 1
            ...
        else:
            res = max(res, end - start)  # 处理分组结果 [start, end)

    return res


if __name__ == "__main__":

    class Solution:
        # !2760. 最长奇偶子数组
        # https://leetcode.cn/problems/longest-even-odd-subarray-with-threshold/description/
        # 找出最长的子数组，满足子数组中的元素:
        # nums[left] % 2 == 0 .
        # nums[i] % 2 != nums[i - 1] % 2 .
        # nums[i] <= threshold .
        def longestAlternatingSubarray(self, nums: List[int], threshold: int) -> int:
            def predicate(left: int, curRight: int) -> bool:
                if nums[left] % 2 != 0:
                    return False
                if curRight != left and nums[curRight] % 2 == nums[curRight - 1] % 2:
                    return False
                if nums[curRight] > threshold:
                    return False
                return True

            groups = groupWhile(len(nums), predicate, skipFalsySingleValueGroup=True)
            return max((end - start for start, end in groups), default=0)

        # !2765. 最长交替子数组
        # https://leetcode.cn/problems/longest-alternating-subarray/description/
        # 1. 长度 > 1.
        # 2. s[1] = s[0] + 1
        # 3. s[len(s) - 1] = s[len(s) - 2] + (-1)^len(s)，例如 [3,4,3,4]
        # 请你返回 nums 中所有 交替 子数组中，最长的长度，如果不存在交替子数组，请你返回 -1
        def alternatingSubarray(self, nums: List[int]) -> int:
            def predicate(left: int, curRight: int) -> bool:
                if curRight == left:
                    return True
                if curRight == left + 1:
                    return nums[curRight] == nums[left] + 1
                return nums[curRight] == nums[curRight - 2]

            res = -1
            end, n = 0, len(nums)
            while end < n:
                start = end  # 开始分组
                while end < n and predicate(
                    start, end
                ):  # [start, end] 闭区间内的元素是否能分到同一组
                    end += 1

                if end == start:  # 非法的单个元素的分组
                    end += 1
                    ...
                else:
                    # 处理分组结果 [start, end)
                    if end - start >= 2:
                        res = max(res, end - start)
                        end -= 1  # !回退，检查是否有更长的交替子数组

            return res

        # !1839. 所有元音按顺序排布的最长子字符串
        # https://leetcode.cn/problems/longest-substring-of-all-vowels-in-order/description/
        # 1. 所有 5 个英文元音字母（'a' ，'e' ，'i' ，'o' ，'u'）都必须 至少 出现一次
        # 2. 子字符串中的元音字母必须按 非递减 的顺序排布
        def longestBeautifulSubstring(self, word: str) -> int:
            def predicate(left: int, curRight: int) -> bool:
                if curRight == left:
                    return True
                return nums[curRight] >= nums[curRight - 1]

            rank = {c: i for i, c in enumerate("aeiou")}
            nums = [rank[c] for c in word]
            groups = groupWhile(len(word), predicate)
            res = 0
            for start, end in groups:
                if len(set(nums[start:end])) == 5:
                    res = max(res, end - start)
            return res

        # 1446. 连续字符
        # https://leetcode.cn/problems/consecutive-characters/description/
        # 只包含一种字符的最长非空子字符串的长度
        def maxPower(self, s: str) -> int:
            def predicate(left: int, curRight: int) -> bool:
                return s[left] == s[curRight]

            groups = groupWhile(len(s), predicate)
            return max((end - start for start, end in groups), default=0)

        # 2110. 股票平滑下跌阶段的数目
        # https://leetcode.cn/problems/number-of-smooth-descent-periods-of-a-stock/description/
        # 每日股价都比 前一日股价恰好少 1
        def getDescentPeriods(self, prices: List[int]) -> int:
            def predicate(left: int, curRight: int) -> bool:
                if curRight == left:
                    return True
                return prices[curRight] == prices[curRight - 1] - 1

            groups = groupWhile(len(prices), predicate)
            res = 0
            for start, end in groups:
                m = end - start
                res += m * (m + 1) // 2
            return res

        # 228. 汇总区间
        # https://leetcode.cn/problems/summary-ranges/description/
        def summaryRanges(self, nums: List[int]) -> List[str]:
            def predicate(left: int, curRight: int) -> bool:
                if curRight == left:
                    return True
                return nums[curRight] == nums[curRight - 1] + 1

            groups = groupWhile(len(nums), predicate)
            res = []
            for start, end in groups:
                if end - start == 1:
                    res.append(str(nums[start]))
                else:
                    res.append(f"{nums[start]}->{nums[end-1]}")
            return res

        # 3350. 检测相邻递增子数组 II
        # https://leetcode.cn/problems/adjacent-increasing-subarrays-detection-ii/description/
        def maxIncreasingSubarrays(self, nums: List[int]) -> int:
            def check(left: int, right: int) -> bool:
                if left == right:
                    return True
                return nums[right - 1] < nums[right]

            lens = [e - s for s, e in groupWhile(len(nums), check)]
            res = 0
            for pre, cur in zip(lens, lens[1:]):
                res = max(res, min(pre, cur))
            for v in lens:
                res = max(res, v // 2)
            return res

    # [2,3,4,3,4]
    print(Solution().alternatingSubarray([2, 3, 4, 3, 4]))  # 4
    # [21,9,5]
    print(Solution().alternatingSubarray([21, 9, 5]))  # 3
