# Approximate Membership Data Structure (概率数据结构/近似成员数据结构/AMQ)

在“近似成员查询（Approximate Membership Query, AMQ）”领域里，没有一个放之四海而皆准的“最好”数据结构，因为它们在**空间占用**、**查询速度**、**插入删除操作支持**、**实现复杂度**、**负载压力**等维度各有取舍。不同场景中，适合的选择也会不一样。

如果一定要从常见的几种数据结构中进行比较，可以大致分为以下几类，并结合各自的优劣势来判断在某个场景下“哪一种最好”：

---

## 1. 布隆过滤器（Bloom Filter）

- **特点**：最经典、实现简单、查询时间和插入时间都为 \(O(k)\)（\(k\) 通常是哈希次数）。
- **优点**：实现最为简单，库和资料丰富；结构固定、查询也很快。
- **缺点**：
  1. 不支持删除（除非用计数 Bloom Filter，但那会增加空间占用和实现复杂度）；
  2. 需要事先预估元素数量并分配好位数组，否则若后续元素远超预期，误报率会激增或者需要整体扩容重建。

在静态或插入后很少变动的场景（如批量构建、只查不删），Bloom Filter 依旧是一个非常好用的选择，尤其是对实现复杂度要求低的系统。

---

## 2. Cuckoo Filter（布谷鸟过滤器）

- **特点**：基于布谷鸟哈希（Cuckoo Hashing）的“踢出（kick-out）”思想，通过存储元素短指纹（fingerprint）来节省空间。
- **优点**：
  1. 原生支持删除；
  2. 在很多负载下比同等误报率的 Bloom Filter 空间效率更好；
  3. 对插入也相对灵活。
- **缺点**：
  1. 在高负载因子时，可能会频繁发生踢出操作，导致插入不稳定；
  2. 实现复杂度较 Bloom Filter 稍高；
  3. 空间利用率受桶大小、指纹长度、负载因子等参数影响较大，需调参。

当需要**动态插入和删除**，或对空间效率有更高要求时，Cuckoo Filter 往往是一个很好的平衡点。

---

## 3. XOR Filter

- **特点**：由 Daniel Lemire 等人提出的一种结构，针对**静态集合**（即插入完成后集合不再变化）的场景而优化；通过 XOR-based 的构造方式在保证低假阳性率的同时可获得极高的查询速度和非常紧凑的空间占用。
- **优点**：
  1. 在静态场景下往往拥有比传统 Bloom、Cuckoo Filter 更好的空间效率以及更快的查询速度；
  2. 构造完成后查询只需要极少的 CPU 指令，适合对查询性能要求苛刻的场景。
- **缺点**：
  1. 不支持删除，且不支持后续插入（或者需要重建）；
  2. 构造过程相比 Bloom Filter 要复杂一些。

如果你有一批**静态数据**，构建完过滤器后只做查询，不怎么改动，那么 XOR Filter 是非常“紧凑且快速”的选择。

---

## 4. Quotient Filter（商过滤器）及其变体（比如 CQF、层级 QF 等）

- **特点**：使用商（quotient）和余数（remainder）的概念将哈希值进行分段存储，在插入和查询上都有不错的复杂度。
- **优点**：
  1. 空间利用率高；
  2. 支持一定程度的插入和删除。
- **缺点**：
  1. 具体实现较为复杂；
  2. 在高负载时也会遇到冲突处理、重建等问题；
  3. 根据实现的不同，对并发或局部操作的支持也可能比较复杂。

Quotient Filter 更像是对 Bloom Filter、Cuckoo Filter 之间的一种折衷，在一些系统里也有应用。但其并没有像 Cuckoo Filter 或 XOR Filter 那样知名。

---

## 5. Morton Filter

- **特点**：Morton Filter（VLDB 2019 论文中提出）在实践中具有与 Cuckoo Filter 类似的插入、删除、查询操作特性，但是在高并发、大规模数据集的场景下往往有更好的性能和更优的缓存局部性。
- **优点**：
  1. 在高负载时比传统的 Cuckoo Filter 更稳定，减少踢出操作的开销；
  2. 仍可支持插入和删除操作，假阳性率可控。
- **缺点**：
  1. 数据结构及实现更复杂；
  2. 在某些读多写少的场景下，性能优势不一定明显。

如果在高并发、高负载的 OLTP 场景或实时系统里，需要一个更可扩展、更稳定的滤器，Morton Filter 可能是一个好选择。

---

## 6. 小结：哪个是“最好”？

1. **静态集合，极致空间效率 + 极快查询**：

   - **XOR Filter** 或者类似的 **Golomb-coded sets（GCS）**、**Ranked-based structures** 一般是最佳，假阳性率和空间占用都可做得相当好。

2. **动态集合，需要频繁插入和删除**：

   - **Cuckoo Filter** 或 **Morton Filter** 更具优势，原生支持删除，且在保持较好空间效率的同时允许插入扩容。

3. **实现简单，且对插删要求不高**：
   - **Bloom Filter** 足够稳定、易用，维护和拓展生态非常成熟。

因此，“最好的”结构并不是一成不变的，而是要看你的**业务场景、更新方式、空间要求、查询负载**以及团队对**实现复杂度**的容忍度。  
在大多数实际系统中，如果你需要**可删、可增**的高效过滤器，Cuckoo Filter 仍然是最常见的首选；如果你要在“只读”数据上追求极致性能和空间利用，那么 XOR Filter 可能是目前最优的选择之一。
