https://leetcode.cn/problems/sum-of-matrix-after-queries/solution/oqlogqzai-xian-cha-xun-by-424479543-q45z/

二维矩阵的两种离线查询:

- 修改和查询交叉的,逆序可能没什么用
- 给若干个修改(覆盖),最后只查询一次,逆序可以简化问题(`因为查到之后那个点就没用了`)
  **染色问题的倒序离线查询**
  二维矩阵很大,不能实际开数组,只能维护行/列的信息
