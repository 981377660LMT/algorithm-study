# 文件重命名问题
# 工作总结(1)(1)(1)(1)(1)(1).docx
#
# 1487. 保证文件名唯一
# https://leetcode.cn/problems/making-file-names-unique/
# https://taodaling.github.io/blog/2020/06/22/%E6%9C%89%E8%B6%A3%E7%9A%84%E9%97%AE%E9%A2%98/
#
# 给你一个长度为 n 的字符串数组 names 。你将会在文件系统中创建 n 个文件夹：在第 i 分钟，新建名为 names[i] 的文件夹。
# 由于两个文件 不能 共享相同的文件名，因此如果新建文件夹使用的文件名已经被占用，
# 系统会以 (k) 的形式为新文件夹的文件名添加后缀，其中 k 是能保证文件名唯一的 最小正整数 。
# 返回长度为 n 的字符串数组，其中 ans[i] 是创建第 i 个文件夹时系统分配给该文件夹的实际名称。
#
# !注意的是，必须保存每个文件名重复的次数times，重复时要从times开始继续查重，否则还是有可能重复
# !软删除可以这样做, 硬删除不行.


from typing import List


class Solution:
    def getFolderNames(self, names: List[str]) -> List[str]:
        res = []
        next_ = dict()  # 已创建的文件夹的下一后缀序号
        for name in names:
            if name not in next_:
                res.append(name)
                next_[name] = 1
            else:
                # 有了就改名，还有的话就继续改名
                k = next_[name]
                while f"{name}({k})" in next_:
                    k += 1
                newName = f"{name}({k})"
                res.append(newName)
                next_[name] = k + 1
                next_[newName] = 1

        return res


print(Solution().getFolderNames(names=["gta", "gta(1)", "gta", "avalon"]))
# 输出：["gta","gta(1)","gta(2)","avalon"]
# 解释：文件系统将会这样创建文件名：
# "gta" --> 之前未分配，仍为 "gta"
# "gta(1)" --> 之前未分配，仍为 "gta(1)"
# "gta" --> 文件名被占用，系统为该名称添加后缀 (k)，由于 "gta(1)" 也被占用，所以 k = 2 。实际创建的文件名为 "gta(2)" 。
# "avalon" --> 之前未分配，仍为 "avalon"

# timeit : names = ["a"]*100000
import timeit

arr = ["a"] * 100000
print(timeit.timeit("Solution().getFolderNames(arr)", globals=globals(), number=20))
