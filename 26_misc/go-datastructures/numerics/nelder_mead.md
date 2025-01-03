下面给出对这段 **Nelder-Mead** 代码的详细分析，并附带一个示例用法，以便读者理解该算法的实现及如何在实际项目中使用。

---

## 一、Nelder-Mead 方法简介

- **Nelder-Mead** 是一种常用的无约束优化方法，也称作单纯形法（Simplex method），尤其适合在多维空间中寻找局部最优解。
- 原理：
  1. 以初始猜测（guess）构造一个包含 \(n+1\) 个顶点的简单形（对于 n 维问题）。
  2. 计算每个顶点处的函数值，按质量好坏排序（这里的“好”可能是距离某个目标值最近，或满足最大/最小等需求）。
  3. 然后做一系列几何操作（反射、扩展、收缩、缩放等）来迭代更新单纯形，期望逐步逼近一个最优解或满足收敛条件。

本代码在此基础上还进行了多重“重启”（Restarts）及一系列自定义逻辑，以在一定程度上尝试搜索全局极值。

---

## 二、核心结构与重要概念

### 1. `NelderMeadConfiguration`

```go
type NelderMeadConfiguration struct {
    Target float64
    Fn     func([]float64) (float64, bool)
    Vars   []float64
}
```

- **Target**：优化目标值，可以是：
  - 具体数值（如 0.0），表示我们想逼近这个值；
  - `math.Inf(1)`：表示要寻找最大值；
  - `math.Inf(-1)`：表示要寻找最小值。
- **Fn**：需要被优化的函数，输入一个 `[]float64`（作为自变量），返回 `(float64, bool)`：
  1. `float64` 为函数的值；
  2. `bool` 表示这个点是否“有效”（满足约束），如不满足则可能要继续寻找别的点。
- **Vars**：初始猜测，一个浮点切片，其长度决定了问题的维度。

### 2. 顶点（`nmVertex`）与顶点列表（`vertices`）

```go
type nmVertex struct {
    vars             []float64 // 当前点的自变量向量
    distance, result float64   // distance 与 result 用于与目标的距离、函数值
    good             bool      // 是否满足约束
}
```

- **vars**：在 n 维问题中，表示一个具体的候选解。
- **result**：调用 `Fn(vars)` 得到的函数值。
- **distance**：如果 `Target` 不是 \(\pm\infty\)，就表示 \(|Fn(vars) - Target|\)；否则用 `result` 做比较。

`vertices` 是 `[]*nmVertex` 的别名，用来存储多个顶点，并可对其排序或做几何操作（加减乘等）。

### 3. `NelderMead` 主函数

```go
func NelderMead(config NelderMeadConfiguration) []float64 {
    nm := newNelderMead(config)
    nm.evaluate()
    return nm.results.vertices[0].vars
}
```

- 外部调用者只需给出一个 `NelderMeadConfiguration`，函数返回优化后的最优解（即 `nm.results.vertices[0].vars`）。
- 内部 `newNelderMead(config)` 会创建一个 `nelderMead` 对象并初始化相关数据结构；再由 `nm.evaluate()` 执行完整的迭代流程。

### 4. 关键常量与几何参数

```go
const (
    alpha = 1     // Reflection coefficient
    beta  = 2     // Expansion coefficient, must be > 1
    gamma = 0.5   // Contraction coefficient, 0 < gamma < 1
    sigma = 0.5   // Shrink coefficient, 0 < sigma < 1

    delta         = 0.0001 // 收敛阈值
    maxRuns       = 130
    maxIterations = 5
)
```

- `alpha, beta, gamma, sigma`：Nelder-Mead 算法中的标准参数，分别对应反射、扩展、收缩、整体收缩等操作。
- `delta`：判断收敛的阈值，如果多次迭代后变化量 < `delta`，就认为收敛或停止。
- `maxRuns`：单次搜索的最大迭代次数。
- `maxIterations`：Nelder-Mead 的“重启”次数上限，用于尝试在多处收敛点中获取更好解。

### 5. `evaluate()` 主循环

```go
func (nm *nelderMead) evaluate() {
    vertices := nm.results.grab(len(nm.config.Vars) + 1)
    // ... 若初始顶点不满足约束，则直接返回 ...

    // 外循环 maxIterations 次（支持多重“重启”）
    for i := 0; i < maxIterations; i++ {
        // 内循环 maxRuns 次（单次收敛迭代）
        for j := 0; j < maxRuns; j++ {
            vertices.evaluate(nm.config) // 调用 Fn, 按距离或结果排序
            best := vertices[0]

            // 如果已满足收敛条件，则退出
            if !nm.checkIteration(vertices) {
                break
            }

            midpoint := findMidpoint(vertices[:len(vertices)-1]...)
            reflection := nm.reflect(vertices, midpoint)
            // ... 进一步判断 reflection 的好坏，选择 expand / contract / shrink ...
            // ... 若满意则结束当前循环 ...

        }
        // “重启”后的结果记入 nm.results，再从里面获取下一批顶点
        nm.results.reSort(vertices[0])
        vertices = nm.results.grab(len(nm.config.Vars) + 1)
    }
}
```

- **多次重启**：每轮迭代在收敛后，会将最优点保存，并从结果集 `nm.results` 中获取新的随机点继续迭代，以此尝试找到更好的全局解。
- **核心几何操作**：
  1. **Reflect** 反射：如果某个顶点过差，则把它反射到单纯形另一边；
  2. **Expand** 扩展：如果反射得到的点更好，继续沿相同方向走得更远；
  3. **Contract** 收缩：若反射不如理想，则在中点和反射点之间收缩；
  4. **Shrink** 整体缩放：如果无论如何都没得到更好的点，就把所有顶点往最好点收缩。

### 6. `results` 与多点采样（`generateRandomVerticesFromGuess`）

```go
type results struct {
    vertices vertices
    config   NelderMeadConfiguration
    pbs      pbs // pbs = []*vertexProbabilityBundle
}
```

- **results**：维护了多个历史结果/随机猜测点，并且通过一些概率分布 (`calculateVVP`) 做候选点的再排序。
- `generateRandomVerticesFromGuess`：从初始 guess 点随机生成一批顶点（默认 1000 个），以形成候选点池。其逻辑是给每个维度一个随机数，并在某种程度上去重。
- `grab(num int)`：实际从这个池子里取若干顶点，构成新的搜索单纯形；搜索后再重排，将最好结果插回 `results`，依此反复。

**特别注意**：这部分代码实现了一种额外的随机探测策略，相比于标准 Nelder-Mead，可能更适合于多模态函数的搜索。

---

## 三、主要方法解读

### 1. `NelderMead(config) []float64`

- **入口函数**：用户只需传入 `NelderMeadConfiguration`，得到优化后的自变量向量。
- 内部做了：
  1. 创建 `nelderMead` 结构（封装初始顶点、随机顶点等）；
  2. 调用 `evaluate()` 进行最多 `maxIterations * maxRuns` 次迭代或直到收敛；
  3. 返回最优顶点的 `vars`。

### 2. `evaluateWithConstraints`

```go
func (nm *nelderMead) evaluateWithConstraints(vertices vertices, vertex *nmVertex) *nmVertex {
    vertex.evaluate(nm.config)
    return vertex
    if vertex.good {
        return vertex
    }
    // ...
}
```

- 这段似乎还留有一定“如果点不满足约束就回退”的逻辑，但在当前版本中仅保留了 `vertex.evaluate(...)` 这一行，然后直接 `return vertex`。
- 如果将来需要在外部函数里检查更多约束，可以在这里扩展。

### 3. 几何操作示例：`reflect, expand, outsideContract, insideContract, shrink`

- **Reflect**

  ```go
  func (nm *nelderMead) reflect(vertices vertices, midpoint *nmVertex) *nmVertex {
      toScalar := midpoint.subtract(nm.lastVertex(vertices))
      toScalar = toScalar.multiply(alpha) // alpha
      toScalar = midpoint.add(toScalar)
      return nm.evaluateWithConstraints(vertices, toScalar)
  }
  ```

  - 取最差点（`nm.lastVertex(vertices)`) 与中点做“反射”，生成新的候选点，然后评估其好坏。

- **Expand**

  ```go
  func (nm *nelderMead) expand(vertices vertices, midpoint, reflection *nmVertex) *nmVertex {
      toScalar := reflection.subtract(midpoint)
      toScalar = toScalar.multiply(beta) // beta
      toScalar = midpoint.add(toScalar)
      return nm.evaluateWithConstraints(vertices, toScalar)
  }
  ```

  - 如果反射点很好，则再往同方向走 `beta` 倍，尝试拿到更好的极值。

- **Shrink**
  ```go
  func (nm *nelderMead) shrink(vertices vertices) {
      one := vertices[0]
      for i := 1; i < len(vertices); i++ {
          toScalar := vertices[i].subtract(one)
          toScalar = toScalar.multiply(sigma) // sigma
          vertices[i] = one.add(toScalar)
      }
  }
  ```
  - 当无法找到好点，就把所有顶点往最优点（`vertices[0]`）收缩。

### 4. 收敛判定：`checkIteration`

```go
func (nm *nelderMead) checkIteration(vertices vertices) bool {
    // 1. 如果 best 与 Target 的误差 < delta，则收敛
    // 2. 如果不是无穷大/小目标，则判断所有点的距离是否都收敛在 delta 内
    // 3. 判断所有点与最好点在欧几里得距离上是否也小于 delta
    ...
}
```

- 返回 `false` 表示不再继续迭代（已满足收敛或误差足够小）。

---

## 四、使用示例

假设我们要最小化一个函数 \( f(x, y) = (x-3)^2 + (y+2)^2 \)，其最小值是 \((3, -2)\) 时为 0。可以使用下面的示例代码：

```go
package main

import (
    "fmt"
    "math"
    "github.com/your_repo/optimization" // 假设把代码放在这个 import
)

func main() {
    // 目标函数
    fn := func(vars []float64) (float64, bool) {
        x, y := vars[0], vars[1]
        val := (x-3)*(x-3) + (y+2)*(y+2)
        // 不带额外约束，这里直接返回 true
        return val, true
    }

    // 配置：因为我们想找最小值，所以把 Target = math.Inf(-1)
    config := optimization.NelderMeadConfiguration{
        Target: math.Inf(-1),
        Fn:     fn,
        Vars:   []float64{0, 0}, // 初始猜测 (0,0)
    }

    // 执行 Nelder-Mead
    result := optimization.NelderMead(config)
    fmt.Println("Optimized result:", result)

    // 一般我们期望它接近 [3, -2]
}
```

- 在 `fn` 中，我们直接返回函数值和 `true`。若有约束可以返回 `false` 时 Nelder-Mead 会尝试其它点。
- `Target = math.Inf(-1)` 表示最小化问题，如果想要最大化，设为 `math.Inf(1)`；若想逼近某个固定值（如 0），则把 Target 设成该值。
- 执行后算法会多次迭代、随机重启，最终结果一般会在 \([3, -2]\) 附近（存在一定随机性）。

---

## 五、总结

1. **核心流程**：

   - 基于初始猜测及随机点构造初始单纯形；
   - 执行最多 `maxRuns` 次迭代，通过反射（Reflect）、扩展（Expand）、收缩（Contract）、整体收缩（Shrink）来改进单纯形；
   - 若收敛或达迭代上限，则保存最佳解并重启新的随机猜测，再次迭代；
   - 最终在多轮“重启”后输出全局最优可能性更高的解。

2. **拓展性**：该实现增加了随机点池（`generateRandomVerticesFromGuess`）和多轮重启策略，帮助跳出局部极值，适合多模态函数。
3. **使用场景**：适用于维数不太高、函数较平滑的情形（或可容忍相对较慢的收敛），在数值优化、机器学习、工程调参等领域都能应用。
4. **限制**：
   - Nelder-Mead 不保证全局最优，它常被用来寻找局部最优；多重重启能一定程度提高概率，但仍可能陷入局部极值；
   - 在高维度（几十或上百维）问题可能效率较低；
   - 若有复杂的非线性约束，需要在 `evaluateWithConstraints` 中自定义处理。

通过以上分析可见，本代码是 **Nelder-Mead** 算法的一个相对完整的实现，并且带有一定的扩展策略。使用时只需编写好目标函数、设置目标值和初始猜测，即可得到一个逼近解。
