- [一文搞懂 LSM(Log-structured Merge Tree)](https://www.cnblogs.com/wxiaotong/p/15919650.html)
- [在 Golang 中实现 LSM 树：综合指南](https://dzone.com/articles/implementing-lsm-trees-in-golang)
- [LSMTree 笔记](https://huanglei.rocks/posts/note-on-lsmt/)
- [DDIA 读书笔记（三）：B-Tree 和 LSM-Tree](https://www.qtmuniao.com/2022/04/16/ddia-reading-chapter3-part1/)
- [B+树,B-link 树,LSM 树...一个视频带你了解常用存储引擎数据结构](https://www.bilibili.com/video/BV1se4y1U7Dn)
  https://jasonkayzk.github.io/2022/11/05/BTree%E3%80%81B-Tree%E5%92%8CLSM-Tree%E5%B8%B8%E7%94%A8%E5%AD%98%E5%82%A8%E5%BC%95%E6%93%8E%E6%95%B0%E6%8D%AE%E7%BB%93%E6%9E%84%E6%80%BB%E7%BB%93/
- [【数据结构与算法】LSM-Tree 的实现原理](https://www.bilibili.com/video/BV1oUtKerE5n)
- [初探 rocksdb 之 lsm tree](https://www.bilibili.com/video/BV11u411P7GP)

---

1. LSM 树并没有一种固定死的实现方式，更多的是一种将：
   **“磁盘顺序写” + “多个树(状数据结构)” + “冷热（新老）数据分级” + “定期归并” + “非原地更新”这几种特性统一在一起的思想**

   设计原则：
   先内存再磁盘
   内存原地更新
   磁盘追加更新
   归并保留新值

   `如果说 B/B+树的读写性能基本平衡的话，LSM 树的设计原则通过舍弃部分读性能，换取了无与伦比的写性能。该数据结构适合用于写吞吐量远远大于读吞吐量的场景，得到了 NoSQL 届的喜爱和好评。`

2. LSM 树通过将数据分为内存组件和磁盘组件来工作：

- MemTable (in-memory component): A balanced tree structure that temporarily stores recent writes
  MemTable (内存组件)：一种平衡树结构，临时存储最近的写入

  MemTable is an in-memory data structure holding data before they are flushed to SST files. It serves both read and write - new writes always insert data to memtable, and reads has to query memtable before reading from SST files, because data in memtable is newer. Once a memtable is full, it becomes immutable and replaced by a new memtable. A background thread will flush the content of the memtable into a SST file, after which the memtable can be destroyed.
  MemTable 是一个内存数据结构，用于在数据刷新到 SST 文件之前保存数据。它同时支持读和写——新的写入总是将数据插入到 memtable 中，而读取必须在从 SST 文件读取之前查询 memtable，因为 memtable 中的数据更新。 一旦 memtable 满了，它就变为不可变，并被一个新的 memtable 替换。一个后台线程将把 memtable 的内容刷新到 SST 文件中，之后 memtable 可以被销毁。

- SSTables (on-disk component): Sorted String tables that store data permanently, organized in levels
  SSTables (磁盘组件)：排序字符串表，永久存储数据，按层级组织

The basic operation flow is as follows:
基本操作流程如下：

- Writes are handled by the MemTable.
- When the MemTable exceeds a threshold size, it is flushed to disk as a sorted SSTable.
  当 MemTable 超过阈值大小时，它会被刷新到磁盘作为一个排序的 SSTable 。
- Reads first check the MemTable, and if the key is not found, it searches through the on-disk SSTables.
  读取第一个检查 MemTable ，如果未找到密钥，则在磁盘上的 SSTable 中进行搜索。
- Background processes periodically merge and compact the SSTables to improve performance and manage disk space efficiently.
  后台进程定期合并和压缩 SSTable 以提高性能并有效管理磁盘空间。

3. https://github.com/facebook/rocksdb/wiki/MemTable
   **RocksDB 的数据结构比较：选择跳表的原因是跳表支持并发插入。**
   ![alt text](image.png)
4. 流程图
   ![alt text](image-1.png)
   如果想让一个引擎工程上可用，还会做大量的性能优化。对于 LSM-Tree 来说，包括

   - 优化 SSTable 的查找。
   - 层级化组织 SSTable。

5. 同一个 Key 会存多遍，一般使用 MVCC 进行控制。
6. 小徐先生笔记
   https://github.com/981377660LMT/golsm

---

level compaction 类似二进制分组
