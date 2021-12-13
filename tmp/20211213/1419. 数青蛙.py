# 青蛙必须 依序 输出 ‘c’, ’r’, ’o’, ’a’, ’k’ 这 5 个字母
# 请你返回模拟字符串中所有蛙鸣所需不同青蛙的最少数目。
# 如果字符串 croakOfFrogs 不是由若干有效的 "croak" 字符混合而成，请返回 -1 。

# 有效的字符:
# 任意时刻都有c>=r>=o>=a>=k
# 结束时c==r==o==a==k 且所有蛙叫完
# 遍历，c的时候+1蛙k的时候-1蛙
class Solution:
    def minNumberOfFrogs(self, croakOfFrogs: str) -> int:
        c = r = o = a = k = 0
        resMax = curMax = 0
        for char in croakOfFrogs:
            if char == 'c':
                c += 1
                # c gives a signal for a frog
                curMax += 1
            elif char == 'r':
                r += 1
            elif char == 'o':
                o += 1
            elif char == 'a':
                a += 1
            else:
                k += 1
                # frog stop croaking
                curMax -= 1

            if c < r or r < o or o < a or a < k:
                return -1
            resMax = max(resMax, curMax)

        if curMax == 0 and c == r == o == a == k:
            return resMax
        return -1


print(Solution().minNumberOfFrogs(croakOfFrogs="crcoakroak"))
# 输出：2
# 解释：最少需要两只青蛙，“呱呱” 声用黑体标注
# 第一只青蛙 "crcoakroak"
# 第二只青蛙 "crcoakroak"
