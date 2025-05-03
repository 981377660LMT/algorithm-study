# 面试题 16.18. 模式匹配
# https://leetcode.cn/problems/pattern-matching-lcci/description/
# 超级复杂的O(N + P)算法——后缀树和O(1) LCA https://leetcode.cn/problems/pattern-matching-lcci/solutions/875256/chao-ji-fu-za-de-on-psuan-fa-hou-zhui-sh-5q5h/


class Solution:
    def patternMatching(self, pattern: str, value: str) -> bool:
        # 若模式为空，则value也应为空才返回True
        if not pattern:
            return not value

        # 计算pattern中'a'和'b'的数量
        count_a = pattern.count("a")
        count_b = len(pattern) - count_a

        # 若value为空，只有一种情况：pattern中只出现某一个字母
        if not value:
            return count_a == 0 or count_b == 0

        # 分情况：若count_a或count_b为0，只需一种子串匹配整个value
        if count_b == 0:
            # 全都是'a'
            # value必须能平均分成count_a段，且a串的重复构成value
            if len(value) % count_a != 0:
                return False
            la = len(value) // count_a
            a_str = value[:la]
            return a_str * count_a == value

        if count_a == 0:
            # 全都是'b'
            if len(value) % count_b != 0:
                return False
            lb = len(value) // count_b
            b_str = value[:lb]
            return b_str * count_b == value

        # 一般情况：枚举a_str长度la，再计算b_str长度lb
        for la in range(len(value) // count_a + 1):
            total_len_b = len(value) - la * count_a
            if count_b != 0 and total_len_b % count_b == 0:
                lb = total_len_b // count_b
                # 根据pattern推导出a_str和b_str
                pos = 0
                a_str, b_str = None, None
                match = True
                for ch in pattern:
                    if ch == "a":
                        sub = value[pos : pos + la]
                        if a_str is None:
                            a_str = sub
                        elif a_str != sub:
                            match = False
                            break
                        pos += la
                    else:  # ch == 'b'
                        sub = value[pos : pos + lb]
                        if b_str is None:
                            b_str = sub
                        elif b_str != sub:
                            match = False
                            break
                        pos += lb
                # 还需确保a串和b串不同
                if match and a_str != b_str:
                    return True
        return False
