# 每一次转化时，将会一次性将 str1 中出现的 所有 相同字母变成其他 任何 小写英文字母（见示例）。
# 只有在字符串 str1 能够通过上述方式顺利转化为字符串 str2 时才能返回 True，否则返回 False。​​
# 思路:str1里字符相同时，对应的映射对也相同，会被set去重
# 注意str2种类不能是26 否则中途无法转换
class Solution:
    def canConvert(self, str1: str, str2: str) -> bool:
        if str1 == str2:
            return True
        if len(set(str2)) == 26:
            return False
        return len(set(zip(str1, str2))) == len(set(str1))


# 输入：str1 = "aabcc", str2 = "ccdee"
# 输出：true
# 解释：将 'c' 变成 'e'，然后把 'b' 变成 'd'，接着再把 'a' 变成 'c'。注意，转化的顺序也很重要。

