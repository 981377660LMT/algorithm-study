下面给出这段 **LSH（局部敏感哈希）** 代码的**详细解读**。它以面向 **L2 距离**（欧几里得距离）的 LSH 为例，实现了多种变体：**Basic LSH**、**Multiprobe LSH**、以及 **LSH Forest**。核心目标是：给定高维向量数据，通过 LSH 将相似（欧几里得距离小）的点映射到同样或相似的哈希值，以便在近似最近邻搜索中快速筛选候选点。

---

## 一、背景与概念

1. **LSH（局部敏感哈希，Locality-Sensitive Hashing）**

   - 一种在高维数据上实现**近似最近邻**搜索的常用方法。
   - 思想：将高维空间中的相似点（距离小）映射到**同一个或相似的哈希桶**，减少在 NN 查询时需要检查的点数。
   - 对 \(\ell_2\) 距离，常采用 Datar 等人提出的“随机超平面投影 + 分桶宽度 \( w \)”等方式。

2. **主要参数**

   - **dim**：数据点维度。
   - **l**：哈希表（Hash Table）的数量。
   - **m**：每个哈希表由多少个哈希函数组合。将 \(m\) 个哈希函数的结果拼成一个哈希键，记为 hashTableKey。
   - **w**：单个哈希函数的宽度，用于将投影值量化到桶中。

3. **各变体**
   - **Basic LSH**：单纯地把哈希值拼接后存入多张表，每张表存不同随机投影。查询时只取各表对应的桶。
   - **Multiprobe LSH**：多探针，除了 base 哈希值外，还会对哈希值做若干“扰动(perturbation)”以探测相邻桶，以获得更多候选点。
   - **LSH Forest**：用**前缀树**存储哈希序列，实现一种层次化索引和逐层扩展搜索。

---

## 二、通用部分：`lshParams`

### 2.1 `lshParams` 结构

```go
type lshParams struct {
	dim int     // 数据维度
	l   int     // 哈希表数
	m   int     // 每个哈希表拼接哈希函数数
	w   float64 // 量化宽度

	a [][]Point   // 随机向量a, 大小 [l][m], 表示 l*m 个投影向量
	b [][]float64 // 随机偏移 b, 大小 [l][m]
}
```

- 每张哈希表有 `m` 个哈希函数：
  - 一个哈希函数 \(h\) 常形如：  
    \[
    h(point) = \left\lfloor \frac{(point \cdot a) + b}{w} \right\rfloor
    \]
  - 这里 `a[i][j]` 是维度与 `point` 同维的随机向量；
  - `b[i][j]` 是在 `[0, w)` 范围内的随机偏移；
  - `w` 用来决定桶宽度。

### 2.2 `newLshParams` 构造

```go
func newLshParams(dim, l, m int, w float64) *lshParams {
	a := make([][]Point, l)
	b := make([][]float64, l)
	random := rand.New(rand.NewSource(rand_seed))

	for i := range a {
		a[i] = make([]Point, m)
		b[i] = make([]float64, m)
		for j := 0; j < m; j++ {
			// a[i][j] 是随机高斯向量, 用 NormFloat64()
			a[i][j] = make(Point, dim)
			for d := 0; d < dim; d++ {
				a[i][j][d] = random.NormFloat64()
			}
			// b[i][j] 是 [0,w) 之间均匀分布
			b[i][j] = random.Float64() * w
		}
	}

	return &lshParams{dim: dim, l: l, m: m, w: w, a: a, b: b}
}
```

- `Point` 实际定义是 `[]float64`，表示一个向量。

### 2.3 `hash` 函数

```go
func (lsh *lshParams) hash(point Point) []hashTableKey {
	hvs := make([]hashTableKey, lsh.l)
	for i := range hvs {
		s := make(hashTableKey, lsh.m)
		for j := 0; j < lsh.m; j++ {
			hv := (point.Dot(lsh.a[i][j]) + lsh.b[i][j]) / lsh.w
			s[j] = int(math.Floor(hv))
		}
		hvs[i] = s
	}
	return hvs
}
```

- 对输入 `point`，在第 `i` 张表里，对 `m` 个哈希函数逐个计算：  
  \[
  hv = \left\lfloor \frac{p \cdot a[i][j] + b[i][j]}{w}\right\rfloor
  \]
- 拼成一个长度为 `m` 的 `hashTableKey`（其实就是 `[]int`），然后收集所有 `l` 张表的结果。

---

## 三、Basic LSH

### 3.1 结构定义

```go
type BasicLsh struct {
	*lshParams             // 嵌入
	tables []hashTable     // l 个哈希表
}

type hashTable map[basicHashTableKey]hashTableBucket

type basicHashTableKey string
type hashTableBucket []string
```

- `BasicLsh` 有 `tables`，其中每个 table 是一个 `map[basicHashTableKey]hashTableBucket`：
  - key: 转换成字符串的哈希键；
  - value: 一组数据点 id。

### 3.2 `NewBasicLsh`

```go
func NewBasicLsh(dim, l, m int, w float64) *BasicLsh {
	tables := make([]hashTable, l)
	for i := range tables {
		tables[i] = make(hashTable)
	}
	return &BasicLsh{
		lshParams: newLshParams(dim, l, m, w),
		tables:    tables,
	}
}
```

- 创建 LSH 参数 + `l` 个空哈希表。

### 3.3 `Insert`

```go
func (index *BasicLsh) Insert(point Point, id string) {
	hvs := index.toBasicHashTableKeys(index.hash(point))
	var wg sync.WaitGroup
	wg.Add(len(index.tables))

	for i := range index.tables {
		hv := hvs[i]
		table := index.tables[i]
		go func(table hashTable, hv basicHashTableKey) {
			if _, exist := table[hv]; !exist {
				table[hv] = make(hashTableBucket, 0)
			}
			table[hv] = append(table[hv], id)
			wg.Done()
		}(table, hv)
	}
	wg.Wait()
}
```

1. 对 `point` 调用 `hash(...)` 得到 `l` 个 `hashTableKey`；
2. 用 `toBasicHashTableKeys` 把每个 `hashTableKey`(其实是 `[]int`) 转成一个字符串 `basicHashTableKey`；
3. 把 `id` 插入到各个 `table[hv]` 里。

### 3.4 `Query`

```go
func (index *BasicLsh) Query(q Point) []string {
	hvs := index.toBasicHashTableKeys(index.hash(q))
	seen := make(map[string]bool)
	for i, table := range index.tables {
		if candidates, exist := table[hvs[i]]; exist {
			for _, id := range candidates {
				seen[id] = true
			}
		}
	}
	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids
}
```

1. 计算 `hash(q)`；
2. 在每张哈希表中取到对应桶里的所有 `id`；
3. 去重后返回。

### 3.5 `Delete`

```go
func (index *BasicLsh) Delete(id string) {
	var wg sync.WaitGroup
	wg.Add(len(index.tables))
	for _, table := range index.tables {
		go func(table hashTable) {
			for tableIndex, bucket := range table {
				for idx, identifier := range bucket {
					if identifier == id {
						table[tableIndex] = remove(bucket, idx)
						if len(table[tableIndex]) == 0 {
							delete(table, tableIndex)
						}
					}
				}
			}
			wg.Done()
		}(table)
	}
	wg.Wait()
}
```

- 不关心 point，只在每个哈希表内遍历所有桶、所有 id 找到后移除。

---

## 四、Multiprobe LSH

Qin Lv 等人提出的**多探针 LSH**：在查询时，除了哈希到的 base 桶，还会对哈希键做**小的扰动(perturbation)** 以探测相邻桶，从而获取更多候选点。

### 4.1 结构

```go
type MultiprobeLsh struct {
	*BasicLsh
	t         int         // 探针个数
	scores    []float64
	perturbSets []perturbSet
	perturbVecs [][][]int
}
```

- 继承了 `BasicLsh`，因此包含 `tables` 等；
- `t`: 需要多少个扰动组合；
- `perturbSets`: 一系列“扰动集合”；
- `perturbVecs`: 把这组扰动转换成对 hashTableKey 的**加减 1**操作等。

### 4.2 生成扰动序列

```go
func (index *MultiprobeLsh) initProbeSequence() {
	m := index.m
	index.scores = make([]float64, 2*m)
	// ...
	index.genPerturbSets()
	index.genPerturbVecs()
}
```

1. `scores`：给可能的扰动(1..2m)赋一些分值(heuristic)，越小表示更优先尝试；
2. `genPerturbSets()`：生成 `t` 个 `perturbSet`；
3. `genPerturbVecs()`：对每个 `perturbSet`，在不同表 `l` 上有一个“加/减 1”的映射到 hashTableKey。

### 4.3 `perturbSet`

```go
type perturbSet map[int]bool
```

- 例如 `{1: true, 3: true}` 表示“在位置 1 做 +1 or -1(视情况)”，“在位置 3 做 +1 or -1”。
- 代码里 `shift()`、`expand()`、`isValid()` 等操作是用来生成并校验下一个可能扰动集合。

### 4.4 `Query`

```go
func (index *MultiprobeLsh) Query(q Point) []string {
	baseKey := index.hash(q)
	results := make(chan string)
	go func() {
		defer close(results)
		// 先查 baseKey
		for i := 0; i < len(index.perturbVecs)+1; i++ {
			perturbedTableKeys := baseKey
			if i != 0 {
				perturbedTableKeys = index.perturb(baseKey, index.perturbVecs[i-1])
			}
			index.queryHelper(perturbedTableKeys, results)
		}
	}()
	// 收集结果
	seen := map[string]bool{}
	for id := range results {
		seen[id] = true
	}
	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids
}
```

- 多探针的关键：**在 baseKey 基础上**，应用若干 `perturbVecs[i-1]` 得到新的哈希键，然后都在 `tables` 中查询合并。
- 因此可得到更多候选。

---

## 五、LSH Forest

Mayank Bawa 等人提出的 **LSH Forest**。不同于前两种把哈希值当作 map key，这里把哈希序列(类似 `[hashVal_1, hashVal_2, ..., hashVal_m]`)视为一条路径插入到前缀树（prefixTree）中。

### 5.1 结构

```go
type LshForest struct {
	*lshParams
	trees []prefixTree
}
```

- `trees` 数量 = `l`，每棵树对应一组哈希函数组合。

### 5.2 `prefixTree` & `treeNode`

```go
type prefixTree struct {
	count int
	root  *treeNode
}

type treeNode struct {
	hashKey int
	ids     []string
	children map[int]*treeNode
}
```

- 类似 Trie/PrefixTree，但是这里节点是 “hashKey = int”，表示某一层的哈希值(单个整型)。
- 如果到达叶子，存储若干 `ids`。

### 5.3 `Insert`

```go
func (index *LshForest) Insert(point Point, id string) {
	hvs := index.hash(point)
	var wg sync.WaitGroup
	wg.Add(len(index.trees))
	for i := range index.trees {
		hv := hvs[i] // 这是一个[]int of length m
		tree := &(index.trees[i])
		go func(tree *prefixTree, hv hashTableKey) {
			tree.insertIntoTree(id, hv)
			wg.Done()
		}(tree, hv)
	}
	wg.Wait()
}
```

- `hv` 是一个 `hashTableKey` (length m).
- 在 `prefixTree.insertIntoTree` 中依层递归插入到 `children` 里，如果节点不存在就新建，直到最终 leaf。

### 5.4 `Query(q, k)`

```go
func (index *LshForest) Query(q Point, k int) []string {
	hvs := index.hash(q) // [l][m]
	results := make(chan string)
	done := make(chan struct{})

	go func() {
		// 先从 m 开始, 逐渐减少 maxLevel
		for maxLevels := index.m; maxLevels >= 0; maxLevels-- {
			select {
			case <-done:
				return
			default:
				index.queryHelper(maxLevels, hvs, done, results)
			}
		}
		close(results)
	}()
	// 收集结果
	seen := map[string]bool{}
	for id := range results {
		if len(seen) >= k {
			break
		}
		seen[id] = true
	}
	close(done)

	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids
}
```

- **核心**：从 `maxLevels = m` 开始，一层层减少**最大前缀长度**(top-down approach)。
- 每次 `queryHelper` 会在每棵树的前缀树里**匹配前 `maxLevels` 层**，然后把子树下的所有 `id` 都送到 `results`。
- 当搜到 `k` 个候选就停止。

---

## 六、点与距离

```go
type Point []float64

func (p Point) Dot(q Point) float64 { ... }
func (p Point) L2(q Point) float64  { ... }
```

- `Point` 是 `[]float64`，提供 `Dot`(内积) 与 `L2`(欧几里得距离) 用来辅助 LSH 生成和真实距离计算等。

---

## 七、运行流程示例

虽然本代码 `main()` 函数是空的，实际用法可能如下：

1. 用户选择某个 LSH 实现，例如 `NewBasicLsh(dim=128, l=10, m=4, w=1.5)`。
2. 对数据集中每个点 `point`，调 `Insert(point, "id_xxx")`。
3. 查询阶段，给定查询点 `q`, 调 `Query(q)` 拿到一批**候选 id**；
4. 再对这批候选点做精确距离计算找最近邻(可选)。

对应 Multiprobe LSH / LSH Forest 也是类似。

---

## 八、总结

- 该库在 **L2 空间** 下实现了多种 LSH 变体：

  1. **Basic LSH**：最简单，多表并存，查询只找各表桶。
  2. **Multiprobe LSH**：对桶做相邻扰动，探测更多桶。
  3. **LSH Forest**：利用树状结构存储哈希前缀，更灵活地控制搜索深度。

- **核心参数**：`dim`, `l`, `m`, `w`；
- **核心原理**：
  - 将点用 `l` 组随机向量(`a[i][j]`) 和随机偏移(`b[i][j]`)投影并分桶（`floor((p·a + b)/w)`）；
  - 对查询点也做同样投影，进入相同桶的点被视为候选近邻；
  - Variants 主要在**如何**组织、合并、扩展桶做改进。

这段代码展示了**LSH**在 Go 中的一个示例性实现，适用于近似最近邻搜索场景。
