from typing import List


class Solution:
    def getFolderNames(self, names: List[str]) -> List[str]:
        res = []
        dic = dict()

        for name in names:
            if name not in dic:
                dic[name] = 1
                res.append(name)
            else:
                # 有了就改名，还有的话就继续改名
                newName = f'{name}({dic[name]})'
                while newName in dic:
                    dic[name] += 1
                    newName = f'{name}({dic[name]})'
                dic[newName] = 1
                res.append(newName)

        return res


print(Solution().getFolderNames(names=["gta", "gta(1)", "gta", "avalon"]))
# 输出：["gta","gta(1)","gta(2)","avalon"]
# 解释：文件系统将会这样创建文件名：
# "gta" --> 之前未分配，仍为 "gta"
# "gta(1)" --> 之前未分配，仍为 "gta(1)"
# "gta" --> 文件名被占用，系统为该名称添加后缀 (k)，由于 "gta(1)" 也被占用，所以 k = 2 。实际创建的文件名为 "gta(2)" 。
# "avalon" --> 之前未分配，仍为 "avalon"

