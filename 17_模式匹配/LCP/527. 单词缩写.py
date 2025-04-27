# https://leetcode.cn/problems/word-abbreviation/description/

# 思路可以分成三大步：

# 1. **按首尾字母和长度分组**
#    - 两个单词如果想要默认缩写相同（`a…z` 这种形式），它们必须：
#      1. 单词长度相同；
#      2. 首字母相同；
#      3. 末字母相同。
#    - 所以先把所有单词用 `(首字母, 末字母, 长度)` 作为 key 分到同一个桶里。不同桶之间永远不存在冲突，可以各自独立处理。

# 2. **组内消除冲突——最小前缀区分**
#    - 如果一个组里只有一个单词，那么它的最小缩写就是标准的 `首字母 + (长度-2) + 末字母`，若这个缩写并不比原词短，就保留原词。
#    - 如果组内有多于一个单词，就需要找“最短前缀长度”让每个单词和同组的其它单词都不冲突。
#      1. 先把这一组的单词按字典序排好；
#      2. 只需要比较相邻两词的 **最长公共前缀（LCP）**，任何两个不相邻的单词，它们的 LCP 只会更小；
#      3. 用一个数组 `lcp[i]` 存储第 `i` 个词与左右相邻词的最大公共前缀长度；
#      4. 那么第 `i` 个词只要多保留 `lcp[i] + 1` 个字符，就能保证与任何同组词都不一致（前面 `lcp[i]` 位都相同，第 `lcp[i]+1` 位必不同）。

# 3. **根据决定的前缀生成最终缩写**
#    - 对第 `i` 个单词，设保留前缀长度为 `P=lcp[i]+1`，省略中间的字符数为 `num = len(word) - P - 1`；
#    - 如果 `num <= 1`（即省略 0 或 1 个字符），缩写并不会比原词短，这时就直接输出原词；
#    - 否则输出 `word[:P] + str(num) + word[-1]`。

# 这样：
# - **分组** 保证我们只对可能冲突的词进行聚合处理；
# - **字典序 + LCP** 让我们只花组内 `O(k)` 次比较就能为每个词算出最小区分前缀；
# - **条件保留原词** 则保证不会输出更长的不划算缩写。

# 总体时间复杂度主要在：
# - 按 key 分组：\(O(n)\)；
# - 每组排序：\(O(k\log k)\)；
# - 组内线性扫一遍算 LCP 并生成缩写：\(O(k)\)。

# 因此对于所有单词，总体约为 \(O\bigl(\sum k_i \log k_i\bigr)\le O(n\log n)\) 的最优水平。


from collections import defaultdict
from typing import List


def cal_lcp(s1: str, s2: str) -> int:
    i = 0
    while i < len(s1) and i < len(s2) and s1[i] == s2[i]:
        i += 1
    return i


def make_abbr(word: str) -> str:
    n = len(word)
    if n <= 3:
        return word
    cand = f"{word[0]}{n-2}{word[-1]}"
    return cand if len(cand) < n else word


class Solution:
    def wordsAbbreviation(self, words: List[str]) -> List[str]:
        n = len(words)
        res = [""] * n

        groups = defaultdict(list)
        for i, w in enumerate(words):
            groups[(w[0], w[-1], len(w))].append((w, i))

        for (_, _, length), group in groups.items():
            if len(group) == 1:
                w, idx = group[0]
                res[idx] = make_abbr(w)
                continue

            group.sort(key=lambda x: x[0])

            k = len(group)
            lcp = [0] * k
            for j in range(k - 1):
                w1, _ = group[j]
                w2, _ = group[j + 1]
                lcp_ = cal_lcp(w1, w2)
                lcp[j] = max(lcp[j], lcp_)
                lcp[j + 1] = max(lcp[j + 1], lcp_)

            for j, (w, idx) in enumerate(group):
                prefix_len = lcp[j] + 1
                num = length - prefix_len - 1
                if num <= 1:
                    res[idx] = w
                else:
                    abbr = w[:prefix_len] + str(num) + w[-1]
                    res[idx] = abbr if len(abbr) < length else w

        return res
