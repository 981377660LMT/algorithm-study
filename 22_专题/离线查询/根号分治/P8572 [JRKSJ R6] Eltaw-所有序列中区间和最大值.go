// # https://www.luogu.com.cn/blog/cyffff/solution-JRKSJ-Eltaw

// # P8572 [JRKSJ R6] Eltaw
// # 又ROW个长为COL的序列nums_1,nums_2,...,nums_k
// # 有q次询问,每次询问给出一个区间[left,right]
// # !求出所有序列在区间[left,right]的和的最大值
// # ROW,COL,q<=5e5 ROW*COL<=5e5

// # !(注意到ROW*COL<=5e5这个奇怪的条件)

// # COL<=ROW时
// # 询问的区间只有O(COL^2)种 `预处理所有查询`一共(ROW*COL*COL) 即O(5e5*sqrt(5e5))

// # COL>=ROW时
// # 每次询问都要(ROW)回答 一共O(q*ROW) 即O(q*sqrt(5e5))
