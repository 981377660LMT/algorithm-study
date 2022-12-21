# 65. 有效数字
class Solution:
    def isNumber(self, s: str) -> bool:
        if "inf" in s.lower():
            return False
        try:
            float(s)
            return True
        except ValueError:
            return False
