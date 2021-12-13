# 交换只能发生在两个不同的字符串之间

# 当 s1s1 = "xx"，s2s2 = "yy" 时，只需交换一次，就可以使两个字符串相等
# 当 s1s1 = "xy"，s2s2 = "yx" 时，需要交换两次才可以使两个字符串相等
class Solution:
    def minimumSwap(self, s1: str, s2: str) -> int:
        xy = yx = 0
        # 抽取两数组同位置不同值（x, y）的个数
        for c1, c2 in zip(s1, s2):
            if c1 == c2:
                continue
            if c1 == 'x':
                xy += 1
            else:
                yx += 1

        # 判断差异字符的个数是否为偶数
        if (xy + yx) & 1:
            return -1

        # 优先交换所有的 "xx" "yy" 和 "yy" "xx"，因为只需一次交换
        # 若还剩一组 "xy" "yx" 则再加 2 即可
        return xy // 2 + yx // 2 + 2 * int(xy & 1)


print(Solution().minimumSwap(s1="xy", s2="yx"))
# 输出：2
# 解释：
# 交换 s1[0] 和 s2[0]，得到 s1 = "yy"，s2 = "xx" 。
# 交换 s1[0] 和 s2[1]，得到 s1 = "xy"，s2 = "xy" 。
# 注意，你不能交换 s1[0] 和 s1[1] 使得 s1 变成 "yx"，因为我们只能交换属于两个不同字符串的字符。
