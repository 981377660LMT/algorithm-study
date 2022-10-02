## Longest Common Prefix

**LCP(i,j)** 表示字符串 s 从下标 i 开始的后缀和从下标 j 开始的`后缀的最长公共前缀`

1. dp 求 LCP O(n^2)
2. 后缀数组+RMQ 求 LCP O(nlogn)
   [`求两个后缀的最长公共前缀` O(nlogn)预处理 O(1)查询](<../%E5%90%8E%E7%BC%80%E6%95%B0%E7%BB%84/golang/%E4%BB%BB%E6%84%8F%E4%B8%A4%E4%B8%AA%E5%90%8E%E7%BC%80(%E5%AD%90%E4%B8%B2)%E7%9A%84LCP.go>)
