# 273. 整数转换英文表示
# https://leetcode.cn/problems/integer-to-english-words/solution/zheng-shu-zhuan-huan-ying-wen-biao-shi-by-powcai/


from typing import List


class Solution:
    def numberToWords(self, num: int) -> str:
        to19 = (
            "One Two Three Four Five Six Seven Eight Nine Ten Eleven Twelve "
            "Thirteen Fourteen Fifteen Sixteen Seventeen Eighteen Nineteen".split()
        )
        tens = "Twenty Thirty Forty Fifty Sixty Seventy Eighty Ninety".split()

        def trans(num: int) -> List[str]:
            if num < 20:
                return to19[num - 1 : num]
            if num < 100:
                return [tens[num // 10 - 2]] + trans(num % 10)
            if num < 1000:
                return [to19[num // 100 - 1]] + ["Hundred"] + trans(num % 100)
            for p, w in enumerate(["Thousand", "Million", "Billion"], 1):
                if num < 1000 ** (p + 1):
                    return trans(num // 1000**p) + [w] + trans(num % 1000**p)
            return []

        return " ".join(trans(num)) or "Zero"
