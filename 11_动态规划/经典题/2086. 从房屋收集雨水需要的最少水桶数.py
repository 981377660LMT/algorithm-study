# 你可以在 空位 放置水桶，从相邻的房屋收集雨水。位置在 i - 1 或者 i + 1 的水桶可以收集位置为 i 处房屋的雨水。一个水桶如果相邻两个位置都有房屋，那么它可以收集 两个 房屋的雨水。
# 在确保 每个 房屋旁边都 至少 有一个水桶的前提下，请你返回需要的 最少 水桶数。如果无解请返回 -1 。

# 1 <= street.length <= 105
# If s == 'H', return -1
# If s starts with HH', return -1
# If s ends with HH', return -1
# If s has 'HHH', return -1

# Each house H needs one bucket,
# that's s.count('H')
# Each 'H.H' can save one bucket by sharing one in the middle,
# that's s.count('H.H') (greedy count without overlap)
# So return s.count('H') - s.count('H.H')
class Solution:
    def minimumBuckets(self, s: str) -> int:
        return (
            -1
            if 'HHH' in s or s[:2] == 'HH' or s[-2:] == 'HH' or s == 'H'
            else s.count('H') - s.count('H.H')
        )


print(Solution().minimumBuckets(s="H..H"))
# 输入：street = "H..H"
# 输出：2
# 解释：
# 我们可以在下标为 1 和 2 处放水桶。
# "H..H" -> "HBBH"（'B' 表示放置水桶）。
# 下标为 0 处的房屋右边有水桶，下标为 3 处的房屋左边有水桶。
# 所以每个房屋旁边都至少有一个水桶收集雨水。
