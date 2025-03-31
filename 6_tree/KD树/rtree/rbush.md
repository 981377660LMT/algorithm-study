## 1. 概述

这份代码实现了 [RBush](https://github.com/mourner/rbush) —— 一个高性能的 JavaScript 库，用于 2D 空间索引。核心数据结构为 **R 树**，一种多叉树，常用于快速查询二维空间中的矩形或点。它支持快速的区域搜索、碰撞检测以及批量加载数据，适合用于地理信息系统（GIS）、地图渲染、游戏中碰撞检测以及其他需要快速空间查询的场景。

---

## 2. 数据结构与基本原理

### R 树简介

- **R 树**：一种层次化的数据结构，每个节点都有一个边界框（bounding box），包含其所有子节点的空间范围。查询时通过比较查询区域与边界框的相交情况，可以快速排除不相关的区域。
- **叶子节点**：存储具体数据项（通常也是矩形或点）。
- **非叶子节点**：存储子节点的信息，其边界框为所有子节点边界框的并集。

---

## 3. RBush 类结构

### 构造函数与属性

```js
constructor(maxEntries = 9) {
  // 设置每个节点的最大子节点数（默认9），同时计算最小填充数（大约40%）
  this._maxEntries = Math.max(4, maxEntries)
  this._minEntries = Math.max(2, Math.ceil(this._maxEntries * 0.4))
  this.clear()
}
```

- **maxEntries**：控制每个节点允许存放的最大子节点数。默认值为 9，数值越大，树的高度可能越小，但每次遍历的节点数会增加。
- **minEntries**：节点最小应有子节点数（约为最大数的 40%），用来保证树的平衡性。
- **clear()**：初始化树的数据结构。

---

## 4. 主要方法讲解

### 4.1 遍历与查询

- **all() 方法**

  ```js
  all() {
    return this._all(this.data, [])
  }
  ```

  该方法递归遍历整棵树，返回所有存储的数据项。内部调用了辅助方法 `_all`。

- **search(bbox) 方法**

  ```js
  search(bbox) {
    // ...检查当前节点是否与查询框有交集
    // 遍历树的节点，若节点为叶子且与查询框相交，则将该数据项加入结果数组
    // 若子节点完全被查询框包含，则直接把子节点下的所有数据项全部加入结果
    // 否则将符合条件的子节点加入待搜索队列
    return result
  }
  ```

  该方法根据传入的边界框（bbox）返回所有与之相交的数据项。利用了边界框相交和包含的判断，优化搜索过程。

- **collides(bbox) 方法**

  用于检测树中是否存在与给定边界框发生碰撞的数据项，适用于碰撞检测场景。

### 4.2 数据加载与插入

- **load(data) 方法**

  ```js
  load(data) {
    // 对于数据量较小，逐个调用 insert() 插入
    // 否则使用批量加载算法 _build() 从零构建树结构
    return this
  }
  ```

  支持批量加载数据。当数据量较大时，利用 OMT（Optimal Minimal Tree）算法构建一棵平衡的树，比逐个插入更高效。

- **insert(item) 方法**

  插入单个数据项，调用内部的 `_insert` 方法，将数据插入到合适的叶子节点，并根据需要调整父节点的边界框和节点分裂。

### 4.3 删除与清空

- **remove(item, equalsFn) 方法**

  采用深度优先遍历的方式查找目标数据项，并在找到后删除。删除后，会通过 `_condense` 方法向上调整，更新边界框，清除空节点。

- **clear() 方法**

  重置树结构，将数据清空。

### 4.4 辅助方法与内部逻辑

- **toBBox(item) 方法**

  默认实现为直接返回数据项本身，假设数据项已经包含 `minX`、`minY`、`maxX` 和 `maxY` 属性。如果数据格式不同，可以覆盖此方法以转换数据为边界框。

- **\_build(items, left, right, height) 方法**

  用于批量构建树。该方法根据数据量和目标高度递归地将数据分块，构建出一棵平衡树。使用了 **multiSelect** 算法来对数据进行部分排序。

- **\_chooseSubtree(bbox, node, level, path) 方法**

  在插入时，选择最佳的子树。根据子节点边界框与待插入边界框的重叠面积和扩展面积来决定最佳位置。

- **\_insert(item, level, isNode) 方法**

  插入时的内部实现，确定插入路径、更新父节点边界框以及处理节点溢出（调用 \_split）。

- **\_split(insertPath, level) 方法**

  当某个节点的子节点数超过最大允许数时进行分裂。选择分裂轴（\_chooseSplitAxis）和分裂位置（\_chooseSplitIndex）来最小化重叠面积和总面积。

- **\_adjustParentBBoxes(bbox, path, level) 方法**

  插入数据后，需要向上调整经过的所有父节点的边界框，使其能正确包含新增的数据。

- **\_condense(path) 方法**

  删除数据后，用于修正树的结构：移除空节点、更新边界框，从而保证树的平衡性。

### 4.5 辅助函数

- **几何计算**

  - `calcBBox`：根据子节点计算父节点的边界框。
  - `extend`：扩展一个边界框以包含另一个边界框。
  - `bboxArea`、`bboxMargin`：计算边界框的面积和周长，用于选择分裂方式。
  - `enlargedArea`、`intersectionArea`：计算两个边界框扩展后增加的面积或相交面积，辅助决定最佳插入位置。

- **排序与选择**

  - `multiSelect` 与 `quickselect`：结合选择算法与分治策略，对数据进行部分排序，便于批量加载时分块构建树。
  - 比较函数 `compareMinX`、`compareMinY` 用于对数据按照 X 或 Y 坐标排序。

- **几何关系判断**

  - `contains` 与 `intersects`：分别判断一个边界框是否包含或与另一个边界框相交。

---

## 5. 使用场景

该库适合以下场景：

- **地图与地理信息系统（GIS）**  
  在地图应用中，常需要根据经纬度范围快速查询和渲染数据。例如，在缩放或平移地图时，查询当前视图范围内的所有地理对象。

- **碰撞检测**  
  在游戏或物理引擎中，利用空间索引快速判断对象之间是否发生碰撞，从而提高性能。

- **数据可视化**  
  在海量数据可视化中，根据视图窗口只加载或显示当前区域的数据。

- **空间数据库索引**  
  在存储和查询空间数据时，利用 R 树可以大大加快范围查询和最近邻搜索。

---

## 6. 使用方法

### 6.1 初始化与数据格式

假设你的数据项都是对象，并且每个对象都具有 `minX`、`minY`、`maxX`、`maxY` 属性。你可以这样使用 RBush：

```js
// 创建 RBush 实例，maxEntries 参数可选，默认 9
const tree = new RBush()

// 示例数据：矩形对象
const data = [
  { minX: 10, minY: 10, maxX: 20, maxY: 20 },
  { minX: 15, minY: 15, maxX: 25, maxY: 25 }
  // ...其他数据
]

// 批量加载数据
tree.load(data)

// 或逐个插入数据
data.forEach(item => tree.insert(item))
```

### 6.2 查询数据

```js
// 定义查询边界框
const searchBBox = { minX: 12, minY: 12, maxX: 18, maxY: 18 }

// 搜索所有与查询边界框相交的对象
const results = tree.search(searchBBox)
console.log(results)
```

### 6.3 碰撞检测

```js
const collision = tree.collides({ minX: 14, minY: 14, maxX: 16, maxY: 16 })
console.log(collision) // true 表示存在碰撞
```

### 6.4 删除数据

```js
// 使用 remove 方法删除数据项，可以传入自定义比较函数
tree.remove(data[0], (a, b) => a.minX === b.minX && a.minY === b.minY)
```

### 6.5 序列化与反序列化

可以通过 `toJSON()` 和 `fromJSON()` 方法保存和恢复树的结构，这对于持久化和传输数据非常有用。

```js
// 序列化
const jsonData = tree.toJSON()

// 反序列化到新的 RBush 实例
const newTree = new RBush()
newTree.fromJSON(jsonData)
```

---

## 7. 总结

这份代码提供了一个完整的 R 树实现，用于高效管理二维空间数据。通过方法如 `load`、`insert`、`search` 和 `remove`，你可以构建一个能够快速进行空间查询、碰撞检测和范围查询的数据结构。常见的使用场景包括地图应用、GIS 系统、游戏开发及其它需要空间数据索引的领域。

使用时只需要确保数据项符合边界框格式，或者通过重写 `toBBox` 方法来转换数据格式。通过批量加载数据（load）以及动态插入和删除数据（insert/remove），可以保证树结构始终保持高效查询性能。
