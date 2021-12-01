from collections import Counter

# """
# This is FontInfo's API interface.
# You should not implement it, or speculate about its implementation
# """
class FontInfo(object):
    #  Return the width of char ch when fontSize is used.
    def getWidth(self, fontSize, ch) -> int:
        """
       :type fontSize: int
       :type ch: char
       :rtype int
       """

    def getHeight(self, fontSize) -> int:
        """
       :type fontSize: int
       :rtype int
       """


from typing import List

# 给定一个字符串 text。并能够在 宽为 w 高为 h 的屏幕上显示该文本。
# 字体数组中包含按升序排列的可用字号，您可以从该数组中选择任何字体大小。
# 您可以使用FontInfo接口来获取任何可用字体大小的任何字符的宽度和高度。

# 最右能力二分
class Solution:
    def maxFont(self, text: str, w: int, h: int, fonts: List[int], fontInfo: 'FontInfo') -> int:
        n = len(fonts)
        chr_freq = Counter(text)

        # -------- 检查fonts中第index个是否可以放得下
        def check(m: int) -> bool:
            if fontInfo.getHeight(fonts[m]) > h:
                return False
            width_sum = 0
            for c, f in chr_freq.items():
                width_sum += (fontInfo.getWidth(fonts[m], c)) * f
            return width_sum <= w

        # -------------- 二分查找 寻找符合条件的最右端
        l = 0
        r = n - 1
        while l <= r:
            mid = (l + r) >> 1
            if check(mid):
                l = mid + 1
            else:
                r = mid - 1

        # -- 最后要再检查一下，不一定保证存在
        return fonts[r] if check(r) else -1


# 返回可用于在屏幕上显示文本的最大字体大小。如果文本不能以任何字体大小显示，则返回-1。
# 示例 1:
# 输入: text = "helloworld", w = 80, h = 20, fonts = [6,8,10,12,14,16,18,24,36]
# 输出: 6

