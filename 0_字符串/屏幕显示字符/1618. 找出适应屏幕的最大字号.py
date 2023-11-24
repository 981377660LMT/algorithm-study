# 1618. 找出适应屏幕的最大字号
# https://leetcode.cn/problems/maximum-font-to-fit-a-sentence-in-a-screen/description/


# """
# This is FontInfo's API interface.
# You should not implement it, or speculate about its implementation
# """
# class FontInfo(object):
#    Return the width of char ch when fontSize is used.
#    def getWidth(self, fontSize, ch):
#        """
#        :type fontSize: int
#        :type ch: char
#        :rtype int
#        """
#
#    def getHeight(self, fontSize):
#        """
#        :type fontSize: int
#        :rtype int
#        """


class Solution:
    def maxFont(self, text: str, w: int, h: int, fonts: List[int], fontInfo: "FontInfo") -> int:
        ...
