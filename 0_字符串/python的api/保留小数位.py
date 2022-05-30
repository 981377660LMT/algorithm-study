# python保留小数位要用format字符串实现

# 保留两位小数
# f`{variable:.2f}`
# f`{123.22:.2f}`
class Solution:
    def discountPrices(self, sentence: str, discount: int) -> str:
        words = sentence.split()
        res = []

        for w in words:
            if w.startswith('$') and w[1:].isdigit():
                count = float(w[1:])
                discountedCount = count * (100 - discount) / 100
                # 保留两位小数
                res.append('$' + f'{discountedCount:.2f}')
            else:
                res.append(w)

        return ' '.join(res)


print(f'{123.2:.2f}')
