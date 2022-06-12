# 找出源字符串中能通过串联形成目标字符串的子序列的最小数量。

# source 和 target 两个字符串的长度介于 1 和 1000 之间。

# summary
# 双指判断source的子序列是否target的一个部分
# 指针i专门用来遍历target序列
# 指针j专门用来遍历source序列
# 每当j遍历完时,如果i没有移动,说明无解返回-1
# 每当j遍历完时,如果i继续移动,结果+1,直至指针i遍历完,返回最终结果

# O(n^2)
class Solution:
    def shortestWay(self, source: str, target: str) -> int:
        res = 0
        j = 0
        while j < len(target):
            i = 0
            preHit = j

            while j < len(target) and i < len(source):
                if target[j] == source[i]:
                    j += 1
                i += 1

            if i == preHit:
                return -1

            res += 1

        return res


print(Solution().shortestWay(source="abc", target="abcbc"))
# 输入：source = "abc", target = "abcbc"
# 输出：2
# 解释：目标字符串 "abcbc" 可以由 "abc" 和 "bc" 形成，它们都是源字符串 "abc" 的子序列。

