from typing import List


class Solution:
    def minAbbreviation(self, target: str, dictionary: List[str]) -> str:
        n = len(target)
        diffs = []
        for word in dictionary:
            if len(word) == n:
                mask = 0
                for i in range(n):
                    if target[i] != word[i]:
                        mask |= 1 << i
                diffs.append(mask)
        if not diffs:
            return str(n)

        abbreviations = []  # 所有可能的缩写

        def generate_abbr(pos=0, count=0, cur_mask=0, abbr=""):
            if pos == n:
                if count > 0:
                    abbr += str(count)
                abbreviations.append((abbr, cur_mask))
                return

            # 跳过当前字符，增加缩写的计数
            generate_abbr(pos + 1, count + 1, cur_mask, abbr)

            # 保留当前字符，终止当前的缩写计数
            generate_abbr(
                pos + 1,
                0,
                cur_mask | (1 << pos),
                abbr + (str(count) if count > 0 else "") + target[pos],
            )

        generate_abbr()
        abbreviations.sort(key=lambda x: len(x[0]))

        # 检查缩写的唯一性，找到最短的非冲突缩写
        for abbr, mask in abbreviations:
            if all(mask & diff for diff in diffs):
                return abbr
        return ""
