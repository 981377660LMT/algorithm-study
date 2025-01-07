1. [Compress your JSON with automatic type extraction](https://stevehanov.ca/blog/?id=104)
   [text](<Compress your JSON with automatic type extraction.md>)
2. [Fast and Easy Levenshtein distance using a Trie](https://stevehanov.ca/blog/?id=114)
   `通过trie实现编辑距离搜索。`
   带编辑距离模糊搜索的trie参见 https://github.com/shivamMg/trie
   但是字典树可能会占用大量内存——甚至可能无法适应 32 位操作系统施加的 2 到 4 GB 的限制。
3. [Compressing dictionaries with a DAWG](https://stevehanov.ca/blog/?id=115)
   `为trie压缩空间。`
   如果你有一个大的单词列表，可以通过 gzip 进行处理，从而获得更好的压缩效果。
   以这种方式存储字典的原因是节省空间并保持易于搜索，而无需先解压缩。

4. [Throw away the keys: Easy, Minimal Perfect Hashing](https://stevehanov.ca/blog/?id=119)
   完美哈希是一种构建没有冲突的哈希表的技术。
5. [Succinct Data Structures: Cramming 80,000 words into a Javascript file.](https://stevehanov.ca/blog/?id=120)
   Succinct Trie 和 DAWG 都是 Trie 的优化版本，旨在提高存储和查询的效率，但它们侧重点不同：

- DAWG 通过共享后缀和最小化自动机，显著减少节点数，适合高空间效率的静态词典存储和查询。
- Succinct Trie 则进一步通过简洁数据结构（如简洁位向量、rank/select 索引等）将 Trie 压缩到接近信息熵下限，同时保持高效的查询性能，适用于内存受限且需要高效模糊匹配的场景。

5. [O(n) Delta Compression With a Suffix Array](https://stevehanov.ca/blog/?id=146)
   `Delta Compression（差分压缩）`是一种高效的数据压缩技术，旨在通过存储和传输两个数据版本之间的差异（delta），而非整个数据集，从而显著减少所需的存储空间和传输带宽。

   在版本控制、文件同步和数据压缩等应用中，常需要高效地存储两个序列（如文件内容、文本字符串）之间的差异。`传统的 INSERT/DELETE 算法`通过记录将序列 A 转换为序列 B 所需的插入和删除操作来实现这一点。然而，这种方法在某些情况下可能会产生冗余操作，尤其是在存在大量相似子串时。
   `COPY/INSERT 替代方案`通过引入 COPY 操作（即从序列 A 的某个位置复制一段子串到序列 B）来减少操作数目，提高压缩效率。为了实现高效的 COPY 操作，本文将介绍如何利用 后缀数组 来快速找到序列 A 中与序列 B 当前未匹配部分的最长匹配子串。

   - INSERT/DELETE ALGORITHMS (插入/删除算法)
     DELETE (p, l)：从序列 A 的位置 p 开始，删除长度为 l 的子串。
     INSERT (p, s)：在序列 A 的位置 p 处，插入字符串 s

     任何基于最长公共子序列的算法在最坏情况下的运行时间都是 O(N^2)。
     现代差异工具通常使用 Meyer 的 O(ND) 算法，该算法与字符串的长度和它们之间的差异数量成正比。
     当两个文本几乎没有相似性时，这将需要很长时间。

   - COPY/INSERT ALGORITHMS (复制/插入算法)
     COPY (p, l)：从序列 A 的位置 p 开始，复制长度为 l 的子串到序列 B。
     INSERT (p, s)：在序列 A 的位置 p 处，插入字符串 s。

     COPY/INSERT 算法的时间复杂度为 O(N)。

---

- https://github.com/ftbe/dawg
- https://github.com/smhanov/dawg
- https://www.wutka.com/dawg.html
