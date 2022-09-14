# https://www.acwing.com/problem/content/description/1287/

# 某人读论文，一篇论文是由许多单词组成的。
# 但他发现一个单词会在论文中出现很多次，
# !现在他想知道每个单词分别在论文中出现多少次。

# !给定k个单词和一段包含n个字符的文章,求有多少个单词在文章里`出现过`。
# 若使用KMP算法,则每个模式串T,都要与主串S进行一次匹配,
# 总时间复杂度为O(n×k+m),其中n为主串S的长度,m为各个模式串的长度之和,k为模式串的个数。
# !而采用AC自动机,时间复杂度只需O(n+m)。


from AutoMaton import AhoCorasick

if __name__ == "__main__":
    n = int(input())
    words = [input() for _ in range(n)]
    ac = AhoCorasick(words)
    match = ac.search("#".join(words))
    res = [0] * n
    for *_, wid in match:
        res[wid] += 1

    print(*res, sep="\n")
