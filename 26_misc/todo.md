1. Count-Min Sketch / HyperLogLog
   https://github.com/axiomhq/hyperloglog
   https://github.com/axiomhq/hyperminhash
   https://github.com/seiflotfy/hyperbitbit
2. Cuckoo Filter
3. B+树及其各种变种/优化

- 工业数据库与文件系统常见的磁盘索引结构，如 B+Tree, B\*Tree, Bw-tree (LLAMA), FD-tree 等等。
- 在传统数据结构课本中只会讲到基础版 B/B+树，但实际工业落地包含很多缓存、预取、合并等细节优化：
  - 分层预读
  - Write-Ahead Logging
  - 带 Bloom Filter 的多层 B+树
- 这些技巧往往决定了真正的 IO 复杂度和在大数据场景下的性能表现。

4. FST (Finite State Transducer)
5. DAWG (Directed Acyclic Word Graph)
6. FM Index
   是基于 Burrows–Wheeler Transform (BWT) 的全文检索结构，可在压缩后的字符串中做快速子串查找。
7. Gap Buffer、Piece Table、PieceTree、Rope + SumTree
8. Elias-Fano 编码
9. 变长编码（Variable-length Encoding）：将小整数用更短的编码表示，大整数用更长的编码表示
   VByte（Variable Byte）：把一个字节最高位当作“是否继续”的标记；
   Gamma Code / Delta Code：基于二进制前缀和偏移量的编码；
   FST（Finite State Transducer）：适用于压缩大量相似前缀的字符串集合，比如对Term本身的前缀进行压缩。
10. xorfilter
    https://github.com/seiflotfy/xorfilter
