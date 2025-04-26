from typing import List


class Solution:
    def fullJustify(self, words: List[str], maxWidth: int) -> List[str]:
        res = []
        n = len(words)
        i = 0

        while i < n:
            # 1. 确定这一行能放多少单词
            row_len = len(words[i])
            j = i + 1
            while j < n and row_len + 1 + len(words[j]) <= maxWidth:
                row_len += 1 + len(words[j])  # +1 是预留一个空格
                j += 1

            # words[i:j] 是本行要放的单词
            num_words = j - i
            line = ""

            # 2. 若是最后一行，或本行只有一个单词，则左对齐
            if j == n or num_words == 1:
                line = " ".join(words[i:j])
                # 补足右边空格
                line += " " * (maxWidth - len(line))
            else:
                # 3. 否则，均匀分配空格
                total_spaces = maxWidth - sum(len(w) for w in words[i:j])
                gaps = num_words - 1
                space, extra = divmod(total_spaces, gaps)
                # 前 extra 个缝隙多分配一个空格
                for k in range(gaps):
                    line += words[i + k]
                    line += " " * (space + (k < extra))
                # 最后一个单词
                line += words[j - 1]

            res.append(line)
            i = j

        return res


if __name__ == "__main__":
    words = ["This", "is", "an", "example", "of", "text", "justification."]
    maxWidth = 16
    ans = Solution().fullJustify(words, maxWidth)
    for row in ans:
        print(f"'{row}'")
    # 输出：
    # 'This    is    an'
    # 'example  of text'
    # 'justification.  '
