class Solution:
    def firstUniqChar(self, s: str) -> str:
        dic = {}
        for c in s:
            # 这句话保证Once的逻辑 很巧妙
            dic[c] = not c in dic
        for k, v in dic.items():
            if v:
                return k
        return ' '

