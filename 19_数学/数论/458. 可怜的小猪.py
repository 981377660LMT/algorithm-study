# 458. 可怜的小猪
# https://leetcode-cn.com/problems/poor-pigs/comments/107570
# 有 buckets 桶液体，其中 正好 有一桶含有毒药，其余装的都是水
# 小猪喝完后，必须有 minutesToDie 分钟的冷却时间。
# 你只有 minutesToTest 分钟时间来确定哪桶液体是有毒的
# 返回在规定时间内判断哪个桶有毒所需的 最小 猪数


class Solution:
    def poorPigs(self, buckets: int, minutesToDie: int, minutesToTest: int) -> int:
        radix = (minutesToTest // minutesToDie) + 1  # 可以尝试的次数
        x = 0
        while radix**x < buckets:  # radix^x>=buckets
            x += 1
        return x


if __name__ == "__main__":
    print(Solution().poorPigs(1000, 15, 60))
