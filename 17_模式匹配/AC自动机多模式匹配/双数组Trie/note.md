TODO: zhihu
双数组 Trie 树(DoubleArrayTrie)Java 实现
https://www.hankcs.com/program/java/triedoublearraytriejava.html

- 两个数组: base[] 和 check[]
  base: 数组中的每个元素相当于 trie 树的一个节点
  check: 相当于当前状态的前一状态

- 对于从状态 s 转移到状态 t，有
  `base[s] + c = t`
  `check[base[s] + c] = s`
  其中 c 是输入变量
