下面这段代码来自 [**FastFilter/xorfilter**](https://github.com/FastFilter/xorfilter) 项目，展示了两种常见的**基于 XOR 的近似集合**数据结构：

1. **Xor8**：也就是本文最先介绍的“XOR Filter”，只存储 8-bit 指纹（fingerprint），假阳性率（false-positive rate）约为 **0.3%**；
2. **BinaryFuse**：更复杂的一种 XOR-based 结构，通常能获得更好（更低）的假阳性率或更快构造速度等特性。

它们都用 **三重哈希 (3-wise) + XOR** 的构造方式，对给定的静态集合建立一个紧凑的过滤器，在查询时能以极小的概率发生假阳性，但不会出现假阴性。下面会依照代码注释进行**详细剖析**。

---

# 目录

- [目录](#目录)
  - [核心数据结构](#核心数据结构)
  - [Xor8 的实现](#xor8-的实现)
    - [数据结构](#数据结构)
    - [主要函数简介](#主要函数简介)
      - [主要辅助数据结构](#主要辅助数据结构)
    - [Populate 过程](#populate-过程)
      - [核心思想](#核心思想)
    - [Contains 方法](#contains-方法)
  - [BinaryFuse 的实现](#binaryfuse-的实现)
    - [数据结构 (BinaryFuse)](#数据结构-binaryfuse)
    - [PopulateBinaryFuse8](#populatebinaryfuse8)
  - [总体原理与流程总结](#总体原理与流程总结)
    - [小结](#小结)
  - [1. 基础示例：使用 `Xor8`](#1-基础示例使用-xor8)
    - [1.1 构建过滤器（Populate）](#11-构建过滤器populate)
      - [解释](#解释)
  - [2. 进阶示例：使用 `BinaryFuse8`](#2-进阶示例使用-binaryfuse8)
  - [3. 其他注意事项](#3-其他注意事项)
  - [4. 序列化 / 存储](#4-序列化--存储)
  - [5. 总结](#5-总结)

---

## 核心数据结构

在 XOR-based 过滤器里，一般都会有以下要素：

- **Seed**：随机种子，用来混入哈希计算（减少哈希冲突带来的构造失败概率）；
- **BlockLength 或 SegmentLength**：数组的大小，也可以理解为构造时要分多少个桶或者段；
- **Fingerprints**：存储每个槽位（bucket）对应的指纹（指纹可能占 8-bit、16-bit 或更多，看具体需求）。
- **Populate**：一次性把所有要存的 key（元素）填充进过滤器，构造成功后就不再支持动态插入或删除（是**静态**结构）。
- **Contains(key) bool**：查询某个 key 是否**可能**在集合里，若返回 `false` 则一定不在，否则以极小概率假阳性。

在本代码中，Xor8 只保存 8-bit 指纹，而 BinaryFuse 可以针对 8-bit / 16-bit / 32-bit 等等泛型 `T Unsigned` 做更通用的处理。

---

## Xor8 的实现

### 数据结构

```go
type Xor8 struct {
    Seed         uint64   // 随机种子
    BlockLength  uint32   // 块的大小(= capacity/3)
    Fingerprints []uint8  // 指纹数组长度=capacity, 其中 capacity ~ 1.23*len(keys), 向下取到3的倍数
}
```

- 由于采用三重哈希（3-wise hashing），把整个指纹数组分成 3 段，每段大小都是 `BlockLength`。
  - 段 0 : `[0, BlockLength-1]`
  - 段 1 : `[BlockLength, 2*BlockLength - 1]`
  - 段 2 : `[2*BlockLength, 3*BlockLength - 1]`
- 当构造完毕后，**Contains** 方法会把 key 映射到这 3 个位置，然后做 XOR 验证。

### 主要函数简介

```go
func (filter *Xor8) Contains(key uint64) bool
```

- 给定一个 key，先用 `mixsplit(key, filter.Seed)` 做哈希；
- 然后计算三重哈希：`h0, h1, h2`，映射到三个段中对应的位置；
- 取出 `filter.Fingerprints[h0] ^ filter.Fingerprints[h1] ^ filter.Fingerprints[h2]`，对比与 key 的“fingerprint”是否相同。若相同，则可能包含；若不同，则一定不包含。

```go
func Populate(keys []uint64) (*Xor8, error)
```

- 构造一个 Xor8 过滤器，内部做若干次迭代，若出现构造冲突，就换个 seed 重试；
- 其主要过程包括：
  1. 分配一个 `capacity ~ 1.23 * size`（并对齐到 3 的倍数）；
  2. 把 3 段(`sets0, sets1, sets2`)中每个位置都统计 “xormask”和 “count”；
  3. 反向构造 stack，一一确定指纹的最终值；
  4. 回填到 `filter.Fingerprints`。

#### 主要辅助数据结构

- `xorset`：带有 `(xormask uint64, count uint32)`，用来暂存每个槽位的 XOR 结果和计数。
- `keyindex`： `(hash uint64, index uint32)`，在填充时用来记录某个被找到只有单一 count 的位置。
- `Q0, Q1, Q2`：临时队列，维护那些计数为 1 的槽位（段0、1、2 分别对应的队列）。

### Populate 过程

1. **capacity** = \( \lceil 1.23 \* size \rceil \)；再对齐到能被 3 整除。
2. **BlockLength** = capacity / 3。数组 `Fingerprints` 长度就是 capacity；
3. **sets0, sets1, sets2** = size=BlockLength 的 `[]xorset`，先都设为 `{0,0}`；
4. 先把所有 key 的三重哈希算出，累积到对应的 xorset 里：
   ```go
   for i := 0; i < size; i++ {
       key := keys[i]
       hs := filter.geth0h1h2(key)
       sets0[hs.h0].xormask ^= hs.h
       sets0[hs.h0].count++
       sets1[hs.h1].xormask ^= hs.h
       sets1[hs.h1].count++
       sets2[hs.h2].xormask ^= hs.h
       sets2[hs.h2].count++
   }
   ```
5. 在 sets0/sets1/sets2 里，所有 **count==1** 的槽位，依次加入到 `Q0/Q1/Q2`：
   ```go
   Q0, Q0size := scanCount(Q0, sets0) // count==1
   Q1, Q1size := scanCount(Q1, sets1)
   Q2, Q2size := scanCount(Q2, sets2)
   ```
6. 反复“弹出”这几个队列中的位置（count==1），将对应 key（其实就是 `xormask`）拆分出去，并减少另两个段中的 count。
   - 如果对另两个段操作后又出现了 `count==1`，就再加回队列里，这形成一个“链式反应”，直到没有新的 count==1 为止。
   - 这些“先被弹出”的 keyindex 就压到 `stack` 里，以备最后反向填充指纹值。
7. 如果最终弹出的元素数量 == size，说明成功把所有 key“分配”完毕（没冲突）。若不等，就重置后换 seed 重来。
8. 最后一步：根据 stack 里的记录，**反向计算**每个位置的 fingerprint，并填入 `filter.Fingerprints` 对应槽位里。

#### 核心思想

- 通过找出某个位置 count==1 的槽位，可知“这个槽位只对应了唯一一个 key 的 XOR”，因此能直接确定这个 key 的指纹；
- 再把这个 key 的信息“去除”掉（相当于把它 XOR 掉）另两个槽位；如此，链式下去，就能最终确定所有指纹。
- 若构造过程碰上冲突（如出现多个 key 纠缠在一起导致无法 count==1 的槽位），就随机化 seed 并重试。

### Contains 方法

```go
func (filter *Xor8) Contains(key uint64) bool {
    hash := mixsplit(key, filter.Seed)
    f := uint8(fingerprint(hash))
    r0 := uint32(hash)
    r1 := uint32(rotl64(hash, 21))
    r2 := uint32(rotl64(hash, 42))
    h0 := reduce(r0, filter.BlockLength)
    h1 := reduce(r1, filter.BlockLength) + filter.BlockLength
    h2 := reduce(r2, filter.BlockLength) + 2*filter.BlockLength
    return f == (filter.Fingerprints[h0] ^ filter.Fingerprints[h1] ^ filter.Fingerprints[h2])
}
```

- 算出 key 的 64-bit 哈希；
- 计算**指纹** `f = fingerprint(hash)` (取 `hash ^ (hash>>32)` 的低 8 bits)；
- 计算 3 个桶位置 `h0, h1, h2`；
- 比较 `f` 和 `Fingerprints[h0] ^ Fingerprints[h1] ^ Fingerprints[h2]` 是否相等；
  - 若不相等 => 一定不在集合；
  - 若相等 => “可能在集合”（假阳性概率约 0.3%）。

---

## BinaryFuse 的实现

在代码后半部分，我们可以看到 **BinaryFuse** 结构以及相应的 `PopulateBinaryFuse8()` 等函数。它是另一个 XOR-based Filter 的**更先进**实现，思路上类似但结构更复杂。

### 数据结构 (BinaryFuse)

```go
type BinaryFuse[T Unsigned] struct {
    Seed               uint64
    SegmentLength      uint32
    SegmentLengthMask  uint32
    SegmentCount       uint32
    SegmentCountLength uint32
    Fingerprints       []T
}
```

- **T**：可以是 `uint8`, `uint16`, or `uint32`，代表指纹占多少位。
- 和 Xor8 类似，也会进行 3 次哈希，但是把整个数组分成 `SegmentCount` 个段（每段大小 `SegmentLength`），让构造和查找更具并行友好和 cache 友好性。

### PopulateBinaryFuse8

```go
func PopulateBinaryFuse8(keys []uint64) (*BinaryFuse8, error) {
    filter, err := NewBinaryFuse[uint8](keys)
    if err != nil {
        return nil, err
    }
    return (*BinaryFuse8)(filter), nil
}
```

- 相当于对 `BinaryFuse[uint8]` 的包装，最终构造一个 8-bit 指纹的 BinaryFuse filter。
- 原理和 Xor8 类似，也是通过一次性对所有 key 做哈希、分配、反向构造等过程，不过 BinaryFuse 使用了段式布局和更细致的逻辑来提升构建效率或查询效率。

---

## 总体原理与流程总结

1. **XOR Filter 原理**

   - 对每个 key 做哈希 -> 得到 3 个槽位 -> 只需在指纹数组中存储足够信息，使得 `fingerprint(key)` 可以用 XOR 来重现。
   - 构造过程：把所有 key 的 XOR 信息累积起来，并通过计数为 1 的“单一”槽位来逐步确定这些 key 的最终指纹。
   - 在查询时，对 key 做同样的哈希映射并 XOR 3 个槽位的指纹，若与 `fingerprint(key)` 相等则返回 `true`。
   - 因为只存 8-bit 指纹，所以一定概率会出现不同 key 的指纹 XOR 相同 => **假阳性**。但不会出现假阴性。

2. **Xor8 vs BinaryFuse**

   - Xor8：较为直接、简单，每个 key 只计算三次哈希，对应在 `[0, BlockLength-1]`, `[BlockLength, 2*BlockLength-1]`, `[2*BlockLength, 3*BlockLength - 1]` 3 个位置；
   - BinaryFuse：更优化的结构，分段式设计，让构造及查询可借助更好的缓存局部性；对大规模数据时表现更好。

3. **主要函数**

   - **`Populate(keys []uint64)`** / **`PopulateBinaryFuse8(keys []uint64)`**：把一个静态集合构建成对应的过滤器；
   - **`Contains(key uint64) bool`**：查询是否可能包含 `key`。
   - 若需要更小的假阳性率，可以用更大的指纹长度（如 `BinaryFuse16`, `BinaryFuse32`），但占用空间也更多。

4. **使用方式**
   - 与常见的 Bloom Filter 用法类似：
     1. 一次性收集全部 key；
     2. 调用 `Populate` 构造；
     3. 用 `Contains(key)` 做 membership 测试。
   - 与 Bloom/Cuckoo Filter 不同的是，XOR Filter 适用于**静态**或**极少更新**的场景，若要插入新元素，需要重新构造或用分层策略。

---

### 小结

- **Xor8**（和一般 XOR Filter）对仅插入一次后不再修改的集合（静态集合）非常合适，拥有**极快查询**和**较小空间**（比同等假阳率的 Bloom Filter 通常更优）。
- **BinaryFuse** 进一步优化了构造与查询性能，适合更大规模的集合或对吞吐量要求更高的环境。
- 两者内部都采用**三重哈希 + XOR** 原理，以及**反向处理**（先找 count=1 槽位）的方式构造指纹数组，确保在查询时只要做几次 XOR 和指纹比较即可。
- 假阳性率取决于指纹位数（8-bit 大约 0.3%～0.4%），如果需要更低，可以用更多位，但会增大内存占用。

本项目给出了可完整应用于生产的参考实现，若要使用的话，最常见流程就是：

1. **`Populate()`** 传入一个去重后的 key 列表，得到 `filter`;
2. **`filter.Contains(k)`** 做 membership 测试；
3. 如果需要迁移或存档，也可直接序列化 `filter.Seed`、`filter.BlockLength`、以及 `filter.Fingerprints[]`，在另一个进程中恢复。

这就完成了一个高性能、低误报率的**静态近似集合**（Approximate Set Membership）功能。

---

下面给出一个简要示例，演示如何使用该项目（[**FastFilter/xorfilter**](https://github.com/FastFilter/xorfilter)）提供的 **Xor8** 以及 **BinaryFuse8** 来构建并查询一个**静态近似集合**。假设你的项目使用 Go 模块，并已在 `go.mod` 中通过 `go get github.com/FastFilter/xorfilter` 引入该库。

---

## 1. 基础示例：使用 `Xor8`

**Xor8** 过滤器适合在只读（静态）场景下，对一个给定的 `uint64` 集合构建近似查询。

### 1.1 构建过滤器（Populate）

```go
package main

import (
    "fmt"
    "log"

    "github.com/FastFilter/xorfilter"
)

func main() {
    // 1.准备你的静态数据：这里举例 20 个随机 uint64 值
    keys := []uint64{1, 42, 123, 999, 1000, 1234, 5555, 9999, 7777, 55555,
        66666, 88888, 101010, 202020, 303030, 404040, 505050, 606060, 707070, 808080}

    // 2. 构建 Xor8 过滤器
    filter, err := xorfilter.Populate(keys)
    if err != nil {
        log.Fatal("Populate failed:", err)
    }

    // 3. 使用过滤器进行查询
    testKeys := []uint64{42, 999, 99999}
    for _, k := range testKeys {
        mightContain := filter.Contains(k)
        fmt.Printf("Key %d => mightContain=%v\n", k, mightContain)
    }
}
```

#### 解释

1. **`keys`**：实际场景中，最好先去重，因为输入集合越大、越有重复，会影响构造性能或需要额外去重逻辑。
2. **`xorfilter.Populate(keys)`**：用全部 key 一次性构建一个 `*Xor8`。
3. **`filter.Contains(k)`**：若返回 `false` 则一定不在集合；若 `true` 则表示“极大概率存在”（仅有约 0.3% 的假阳性几率）。

因为 `Xor8` 的指纹是 8-bit，对于大规模数据集，这种假阳性率在实际中通常是能接受的。如果需要更低的误报，可考虑 `BinaryFuse` 并使用更大的位宽（如 16 位、32 位指纹）。

---

## 2. 进阶示例：使用 `BinaryFuse8`

**BinaryFuse** 结构号称在构造和查询性能上都有更优表现，特别适用于大规模数据（数百万级或以上）。与 `Xor8` 的使用方式基本一致，差异在于调用 `PopulateBinaryFuse8()` 而非 `Populate()`。

```go
package main

import (
    "fmt"
    "log"

    "github.com/FastFilter/xorfilter"
)

func main() {
    // 1. 准备你的数据
    keys := []uint64{1, 42, 123, 999, 1000, 1234, /* ... */ 808080}

    // 2. 使用 BinaryFuse8 来构建过滤器
    fuseFilter, err := xorfilter.PopulateBinaryFuse8(keys)
    if err != nil {
        log.Fatal("PopulateBinaryFuse8 failed:", err)
    }

    // 3. 测试查询
    testKeys := []uint64{42, 999, 99999}
    for _, k := range testKeys {
        mightContain := fuseFilter.Contains(k)
        fmt.Printf("Key %d => mightContain=%v\n", k, mightContain)
    }
}
```

- **`xorfilter.PopulateBinaryFuse8()`**：生成一个指纹为 8 位 (`uint8`) 的 BinaryFuse 过滤器；对大量数据（上百万）时，也能保持较高的构建成功率和查询速度。
- 查询方式同 `Xor8`：`fuseFilter.Contains(k)`。

---

## 3. 其他注意事项

1. **静态集合**：
   - XOR-based 过滤器天生只适用于**不再变动的集合**。若要插入或删除元素，需要整体重建或考虑多级过滤器方案。
2. **去重**：
   - 建议在调用 `Populate` 前对 `keys` 去重，如果重复元素很多，会占用额外内存并提高冲突概率。
   - 库内部在构造失败后，可能会尝试 `pruneDuplicates()`，但那也会增加构建成本。
3. **空间和误报率**：
   - `Xor8` 约 0.3% 的假阳性率；
   - `BinaryFuse8` 约 0.4% 左右；
   - 如果你需要更低误报率，可以用更大的指纹类型（如 `BinaryFuse16`, `BinaryFuse32` 等，需手动写泛型实例化或参考该库其他函数），但占用空间会相应增大。
4. **性能测试**：
   - 如果你有数百万元素，可以事先对比 `Xor8` 和 `BinaryFuse8` 的构造用时、查询用时、空间大小、实际误报率，然后选择最合适的一种。

---

## 4. 序列化 / 存储

如果需要把已构造好的过滤器保存到文件或传输给其他服务，需要将下列字段序列化：

- 对于 `Xor8`：`Seed`, `BlockLength`, `Fingerprints`（长度=capacity）。
- 对于 `BinaryFuse`：`Seed`, `SegmentLength`, `SegmentCount`, 以及 `Fingerprints` 数组。

在生产环境中，你可以把它们存成二进制文件或放到数据库里，然后在应用启动时读出并重建过滤器对象即可。该项目暂未提供内置的“自动序列化”方法，但你可以自己实现（如 `encoding/binary` 或 `gob` 方式）：

```go
// 伪码示例
// 序列化时
data := make([]byte, 0)
// 先写 seed, blockLength...
// 再写 fingerprints 数组

// 反序列化时
// 读 seed, blockLength...
// 分配一个 Xor8，然后把 fingerprints 读进来。
filter := &xorfilter.Xor8{
    Seed:        storedSeed,
    BlockLength: storedBlockLength,
    Fingerprints: storedFingerprints,
}
```

---

## 5. 总结

**使用核心流程**：

1. **收集所有 key**（类型为 uint64），并去重；
2. 调用 `xorfilter.Populate()` 或 `xorfilter.PopulateBinaryFuse8()` 一次性构建过滤器；
3. 使用 `filter.Contains(key)` 检查任何 key 是否在集合中：
   - `false` => 一定不在；
   - `true` => “可能在”，有极低假阳性概率。

这样就完成了对**静态集合**的高效 membership 查询，兼具**极快查询**和**较小内存占用**的特点。与 Bloom / Cuckoo Filter 相比，XOR Filter 在只读场景下通常**空间效率更佳**、查询也更快，非常适合各种离线索引、日志/请求指纹集合等需要低误报的场景。
