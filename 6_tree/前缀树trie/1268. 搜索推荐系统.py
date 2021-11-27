# 在依次输入单词 searchWord 的每一个字母后，
# 推荐 products 数组中前缀与 searchWord 相同的最多三个产品。
# 如果前缀相同的可推荐产品超过三个，请按字典序返回最小的三个。
# products有重复的，考虑不周啊


# console.log(suggestedProducts(['mobile', 'mouse', 'moneypot', 'monitor', 'mousepad'], 'mouse'))
# // 输出：[
# //   ["mobile","moneypot","monitor"],
# //   ["mobile","moneypot","monitor"],
# //   ["mouse","mousepad"],
# //   ["mouse","mousepad"],
# //   ["mouse","mousepad"]
# //   ]
from typing import List
from bisect import bisect_left

# 1. sort input array
# 2. use binary search to find first index
# 3. check the following 3 words


class Solution:
    def suggestedProducts(self, products: List[str], searchWord: str) -> List[List[str]]:
        products.sort()
        res, prefix, index = [], '', 0
        for char in searchWord:
            prefix += char
            index = bisect_left(products, prefix, index)  # 因为每次搜索index都会单增，所以用index限制lower,是用来加速的
            res.append([word for word in products[index : index + 3] if word.startswith(prefix)])

        return res


print(Solution().suggestedProducts(['mobile', 'mouse', 'moneypot', 'monitor', 'mousepad'], 'mouse'))

