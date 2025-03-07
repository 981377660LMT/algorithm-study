哈希树，用来快速校验数据完整性或差异
Merkle Tree：本质上是线段树/分块维护区间哈希。
https://zhuanlan.zhihu.com/p/16161065503
https://atcoder.jp/contests/abc351/tasks/abc351_g

# Merkle Tree / Patricia Tree / Merkle Patricia Trie (MPT)

## 一、Merkle Tree

### 1. 是什么

- **Merkle Tree（Merkle 树）**是由 Ralph Merkle 在 1979 年提出的一种基于哈希函数的树形数据结构，也称“哈希树（Hash Tree）”。
- 它的核心思想是：将叶子节点存储数据或数据块的哈希值，然后逐层往上合并，父节点是左右子节点哈希值的组合再进行哈希，直到根节点（也称 Merkle Root）。
- 常见的组合方式是将两个子节点的哈希值连接（如 `hash_left || hash_right`）后，再取哈希作为父节点的值。叶子节点通常用原数据或原数据哈希做输入。

### 2. 为什么需要它

1. **完整性校验**
   - Merkle 树能在大规模数据传输或存储中高效地进行**数据完整性校验**：只需比较根哈希，就能判断整棵树的数据是否被篡改。
2. **局部校验和**
   - 如果只想验证一个或少数叶子数据的正确性，Merkle 树允许你只获取与该叶子在“到根路径”上的节点哈希值进行验证，而无需下载或检查整棵数据。
   - 这种验证方式在区块链（如 Bitcoin）和分布式文件系统（如 IPFS）中非常常见：一个“Merkle Proof（Merkle 证明）”就可以证明叶子数据是否属于一棵特定的 Merkle 树。
3. **减少冗余、快速对比**
   - 若两棵 Merkle 树的根哈希相同，就说明整棵树代表的数据一致；如果不同，可以通过分块对比找出差异部分。

### 3. 怎么办（实现与使用）

1. **构建**
   - 收集所有叶子数据块，为每块计算哈希值，放在叶子节点；
   - 逐层往上做合并哈希：对每对孩子节点连接其哈希值并再哈希，得到父节点哈希值；
   - 直到只剩一个根节点（Merkle Root）。
2. **验证**
   - 若要验证某个叶子是否在树中：请求该叶子哈希 + 它在路径上的所有兄弟节点哈希，通过逐层计算并与根哈希对比即可。
3. **应用**
   - 区块链（区块头中的 Merkle Root，用于快速验证交易是否包含在区块中）、Git（对象哈希树）、IPFS 等广泛使用。

---

## 二、Patricia Tree / Patricia Trie

### 1. 是什么

- **Patricia** 是 “Practical Algorithm To Retrieve Information Coded In Alphanumeric” 的缩写，有时也称 **Patricia Trie** 或 **Radix Trie**。
- 它是一种前缀树（Trie）的压缩版本：对于原本 Trie 中只有单一路径的节点会进行合并，以减少空间占用。
- 在 Patricia Trie 中，每个节点会存储**关键字分支（edge label）**，而不是一个一个字符地分裂节点，这样可以减少深度和节点数量，特别适合存储稀疏关键字或者较长公共前缀的情况。

### 2. 为什么需要它

1. **空间优化**
   - 普通 Trie 在字符级拆分，若有少量关键字且字符集庞大，会出现大量稀疏节点；Patricia Trie 将单一分支的路径合并，大幅减少节点数量和内存开销。
2. **快速匹配**
   - Patricia Trie 在搜索时可以跳过那些长的共享前缀，达到类似二分搜索或分段搜索的效果，减少匹配过程中的比较次数。
3. **常用于 IP 路由查找**
   - IPv4 / IPv6 地址前缀匹配常使用 Patricia Trie 来快速定位路由表项（LPM：Longest Prefix Match）；
   - 也用于数据库中的字符串索引、图形界面菜单索引、编译器关键字匹配等场景。

### 3. 怎么办（实现与使用）

1. **构建**
   - 与一般 Trie 相似，依次插入关键字；但若插入时发现新插入关键字与已有节点存在一个公共前缀，再到分歧点进行“分裂”或“合并”操作，确保路径只在必要处分叉。
2. **搜索**
   - 从根节点开始，根据关键字当前尚未匹配的部分与节点中记录的分支前缀做对比，若部分前缀匹配成功继续深入，否则停止。
3. **更新（插入/删除）**
   - 需要在分歧点处进行相应操作，如果插入新关键字导致原节点分裂，则更新 Patricia Trie 结构；删除时若出现只剩单一路径的节点，也可以进行合并。
4. **变体**
   - 有些实现将 Patricia Trie 与 bit-level 的操作结合，用于 IP 地址的位匹配（每一位或者每几位为一个分叉），更适合路由查找。

---

## 三、Merkle Patricia Trie (MPT)

### 1. 是什么

- **Merkle Patricia Trie (MPT)** 有时也称 **Trie + Merkle** 或 **Secure Trie**，是一种将 **Patricia Trie** 与 **Merkle Tree** 思想相结合的数据结构。
- 它最知名的应用是在 **Ethereum**（以太坊）中，用于存储账户状态、交易收据等，保证数据可验证并节省空间。
- MPT 核心：
  1. 用 Patricia Trie 进行**紧凑前缀存储**；
  2. 在每个节点再做哈希，将不同节点哈希链接起来形成一个**Merkle 化的 Trie**；
  3. 根哈希（Root Hash）就能唯一标识整棵 Trie 中所有键值对的内容。

### 2. 为什么需要它

1. **可验证性**
   - 通过将 Trie 各节点哈希链接起来，任何对底层数据的修改都会改变根哈希。
   - 区块链节点可以只传递部分路径并附带节点哈希，就能证明某键是否存在以及对应值是什么。
2. **空间与层次优化**
   - 使用 Patricia Trie，避免了普通 Trie 过多稀疏节点（每个节点可能是一个“Nibble path”，存一串十六进制半字节路径标签）。
   - 减少了树深度与节点数，提高查询效率。
3. **对部分更新友好**
   - 当更新某键值对时，只需更新与该键相关的那条路径及其上层节点哈希，不必重哈希全树。
4. **以太坊场景**
   - 以太坊中，每个区块头都存储了三颗 MPT 的根哈希：
     - 状态树（accounts state），
     - 交易树（transactions），
     - 收据树（receipts）。
   - 这允许轻节点（light client）或其他节点在不知道全部数据的情况下也可验证一个交易或账户状态。

### 3. 怎么办（实现与使用）

1. **Trie 结构**
   - MPT 的节点类型分为：**扩展节点（Extension node）**、**分支节点（Branch node）** 和 **叶子节点（Leaf node）**。
   - 根据当前键的剩余路径选择进入分支或直接到叶子；扩展节点通常存储一个共享的路径前缀（若存在），分支节点存储多个子分支指针（如 16 个可能性），叶子节点结束某条路径并存储键值对。
2. **Merkle 化**
   - 每个节点的哈希由该节点的结构编码 (RLP in Ethereum) 加上子节点哈希（或值）的组合计算得出。
   - 若子节点很多（分支节点），则其哈希也包含所有子节点哈希的聚合；叶子和扩展节点就包含“路径 + 下一级哈希”。
3. **插入/更新**
   - 在 MPT 中插入或更新一个键值对时，需要沿 Trie 查找插入位置，途中若产生新节点或修改现有节点，都要更新它们的哈希，并向上递归更新父节点哈希，直到根。
   - 此外，还要注意路径压缩/分裂等操作与 Patricia Trie 逻辑一致。
4. **验证 / Proof**
   - 若要证明某键是否存在、或它的对应值是什么，MPT 可以提供一条从根到叶的节点序列（包含每个节点的哈希和 RLP 编码）。只需在本地重新计算这些节点的哈希并和给定的父节点哈希进行对比，最终与根哈希相比对，若一致则证明该键值对确实存在于这棵树中。

---

## 四、常见应用场景

1. **区块链 / 分布式账本**

   - **Merkle Tree**：比特币区块头用 Merkle Tree 存交易列表，方便轻节点验证交易归属；
   - **Merkle Patricia Trie (MPT)**：以太坊中存储账户状态、合约存储、收据等关键数据结构，保障完整性与可验证性。

2. **分布式存储 / P2P 文件系统**

   - **Merkle Tree**：IPFS（InterPlanetary File System）中使用 Merkle DAG 来唯一标识文件内容，支持内容寻址；
   - 可快速验证文件某些分块的正确性，而无需下载全部。

3. **路由表 / 字符串索引**

   - **Patricia Trie**：用于 IP 路由、字符串前缀匹配，对稀疏场景能减少空间浪费，提高查询效率。

4. **版本管理 / 数据去重**
   - **Merkle Tree**：Git 使用（hash-based）树结构来跟踪文件快照和差异；
   - 类似的差分备份系统也常见哈希树或 Patricia Trie 进行版本化管理。

---

## 五、总结

1. **是什么**

   - **Merkle Tree**：用哈希聚合自底向上构建的哈希树，能做完整性校验和局部证明；
   - **Patricia Trie**：压缩前缀树，通过合并单一路径节点，减少空间占用，常用于快速前缀匹配；
   - **Merkle Patricia Trie (MPT)**：将 Patricia Trie 与 Merkle Hash 结合，用于区块链或分布式系统中可验证、节省空间的键值存储结构。

2. **为什么**

   - **Merkle Tree**：简化并行校验，节省带宽，可做局部验证；
   - **Patricia Trie**：空间压缩、前缀检索效率高；
   - **MPT**：既享有 Patricia Trie 的压缩前缀优势，又具备 Merkle Hash 的可验证性和抗篡改能力，最适合区块链或分布式存储等需要可验证、可索引场景。

3. **怎么办**
   - **Merkle Tree**：构建时自底向上计算哈希，验证只需兄弟节点哈希的路径；
   - **Patricia Trie**：插入/删除时注意路径分裂和合并，节点中存储紧凑前缀；
   - **MPT**：基于 Patricia Trie 做节点划分，每个节点再哈希以确保全局完整性，插入/更新时更新整条路径哈希，验证时提供完整的节点哈希序列（Merkle Proof）。
