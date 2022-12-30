MOD = int(1e9 + 7)
INF = int(1e20)

# 2259. 移除指定数字得到的最大结果
# 直接字符串max即可 不需数字


class Solution:
    def removeDigit(self, number: str, digit: str) -> str:
        """从 number 中 恰好 移除 一个 等于 digit 的字符后，找出并返回按 十进制 表示 最大 的结果字符串。

        O(n) 贪心
        你需要在人生中找到自己的另一半，在这个过程中你会遇到各种各样的人。
        如果那个人是你所仰慕的对象(number[i + 1]> number[i] )，
        那么你可能就会与ta结为恋人，这件事当然是越早做完越好的。
        如果你走过一生还没有碰到合适的人，你只能在最后一个时刻孤独终老了。
        """
        n = len(number)
        cand = -1
        for i, num in enumerate(number):
            if num == digit:
                cand = i
                if i + 1 < n and number[i + 1] > num:
                    break
        # 如果后面那个数比digit大，那么删除即可
        # 只需要比较下一位就可，而非子串
        return number[:cand] + number[cand + 1 :]

    def removeDigit2(self, number: str, digit: str) -> str:
        """暴力"""
        return max(
            (number[:i] + number[i + 1 :] for i in range(len(number)) if number[i] == digit),
            key=int,
        )


# 1202125
