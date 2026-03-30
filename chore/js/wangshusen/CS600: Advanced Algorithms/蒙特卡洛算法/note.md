https://www.youtube.com/watch?v=XRGquU0ZJok&list=PLvOO0btloRns2Wnn2MPQ-Z8viHoBfxGGJ

Monte Carlo (蒙特卡洛) 是一类随机算法，通过随机抽样来近似解决计算问题。蒙特卡洛算法在统计学、机器学习、金融、物理学等领域有广泛应用。
这节课通过 5 个具体的例子来介绍蒙特卡洛。
0:12 一、均匀抽样估算 Pi 值
5:27 二、布封投针 (Buffon's Needle Problem) 估算 Pi 值
9:56 三、求阴影部分面积
13:40 四、蒙特卡洛积分 (近似求定积分)
20:18 五、近似求期望

---

这段文本是一堂关于**蒙特卡洛方法 (Monte Carlo Methods)** 的深入浅出课程的逐字稿。蒙特卡洛方法是计算科学、物理学、统计学以及深度学习中极其重要的基石算法。

以下是对这节课内容的深入分析与讲解，结合了数学原理与代码实现的视角。

---

# 蒙特卡洛方法 (Monte Carlo Methods) 深度解析

## 1. 核心概念与本质

**定义**：蒙特卡洛不是某一种具体的算法，而是一大类**随机算法**的总称。
**本质**：通过**重复的随机采样**来估算难以通过解析解获得的数值结果。
**数学保证**：**大数定律 (Law of Large Numbers)**。当样本数量 $N$ 趋于无穷大时，经验均值（随机模拟的结果）依概率收敛于理论期望值（真实值）。

与之相对的还有：

- **拉斯维加斯算法 (Las Vegas)**：输出总是正确的，但运行时间是随机的（例如随机快速排序）。
- **蒙特卡洛算法**：运行时间固定，输出是近似的（带有随机误差），但误差随样本量增加而减小。

---

## 2. 经典案例解析

### 2.1 估算 $\pi$ 值 (几何法)

这是演示蒙特卡洛最直观的例子。

- **场景**：
  - 一个边长为 2 的正方形（$x \in [-1, 1], y \in [-1, 1]$），面积 $A_{square} = 4$。
  - 内切一个单位圆（半径 $r=1$），面积 $A_{circle} = \pi r^2 = \pi$。
- **方法**：
  1.  在正方形内均匀随机生成 $N$ 个点 $(X, Y)$。
  2.  统计落在圆内的点数 $M$（判断标准：$X^2 + Y^2 \le 1$）。
- **推导**：
  $$ \frac{M}{N} \approx \frac{A*{circle}}{A*{square}} = \frac{\pi}{4} \implies \pi \approx \frac{4M}{N} $$
- **误差分析**：
  - 精度收敛速度慢，误差上界与 $\frac{1}{\sqrt{N}}$ 成正比。
  - 意味着：想要精度提高 10 倍，计算量（样本量）需要增加 100 倍。

### 2.2 布丰投针 (Buffon's Needle)

这是历史上早期的物理蒙特卡洛实验。

- **场景**：一组间距为 $D$ 的平行线，投掷长度为 $L$ 的针（$L \le D$）。
- **概率公式**（微积分推导结论）：
  $$ P(\text{相交}) = \frac{2L}{\pi D} $$
- **估算 $\pi$**：
  通过实验投掷 $N$ 次，观察到相交 $M$ 次，则 $P \approx \frac{M}{N}$。
  $$ \pi \approx \frac{2LN}{DM} $$

### 2.3 估算不规则图形面积

这是数值积分的前置概念。

- **方法**：拒绝采样 (Rejection Sampling) 的变体。
  1.  找到一个包含目标不规则图形的规则区域（如正方形），其面积 $A_{total}$ 已知。
  2.  均匀采样，计算落在目标区域内的比例 $\frac{M}{N}$。
  3.  目标面积 $\approx A_{total} \times \frac{M}{N}$。
- **价值**：当目标区域由复杂的方程组（如圆与扇形的交集）定义，无法直接计算面积时，这种方法非常有效。

---

## 3. 蒙特卡洛的核心应用

### 3.1 近似定积分 (Numerical Integration)

在工程和科学中，很多函数的原函数找不到（无法解析积分），只能求数值解。

**一元函数定积分**：
计算 $I = \int_a^b f(x) dx$。

1.  在 $[a, b]$ 上均匀采样 $N$ 个点 $x_i$。
2.  计算函数平均高度 $\bar{y} = \frac{1}{N} \sum_{i=1}^N f(x_i)$。
3.  积分值 $\approx \text{区间宽度} \times \text{平均高度} = (b-a) \times \frac{1}{N} \sum_{i=1}^N f(x_i)$。

**多元函数定积分**：
计算 $I = \int_{\Omega} f(\mathbf{x}) d\mathbf{x}$。

1.  定义积分区域 $\Omega$ 的体积为 $V$（前提是 $\Omega$ 形状简单，体积易求，如长方体、球）。
2.  在 $\Omega$ 内均匀采样 $N$ 个向量 $\mathbf{x}_i$。
3.  $I \approx V \times \frac{1}{N} \sum_{i=1}^N f(\mathbf{x}_i)$。

### 3.2 估算数学期望 (Estimation of Expectation)

**这在机器学习中尤为关键**。

- **定义**：随机变量 $X$ 服从概率分布 $p(x)$，求函数 $f(X)$ 的期望：
  $$ \mathbb{E}[f(X)] = \int f(x) p(x) dx $$
- **蒙特卡洛方法**：
  1.  **关键区别**：不再是均匀采样，而是**根据概率密度函数 $p(x)$ 进行采样**得到 $x_1, ..., x_N$。
  2.  直接求函数值的算术平均：
      $$ \mathbb{E}[f(X)] \approx \frac{1}{N} \sum\_{i=1}^N f(x_i) $$

**深度学习中的应用**：

- **SGD (随机梯度下降)**：本质上就是蒙特卡洛方法。损失函数是所有样本梯度的平均（期望），SGD 通过随机抽取一小部分样本（Batch）来估算整体梯度。虽然有噪声（随机误差），但计算量大幅降低，且方向大体正确。

---

## 4. 总结与代码示例

课程强调了蒙特卡洛虽然结果“通常是错的”（是近似值），在统计学和机器学习中“足够好”且“计算可行”才是最重要的。

以下是基于课程内容编写的 Python 代码，对应 **例 1 (求 Pi)** 和 **应用 1 (求定积分)**：

```python
import numpy as np

def monte_carlo_pi(n_samples):
    """
    通过蒙特卡洛方法估算 Pi 值
    """
    # 1. 在 -1 到 1 之间均匀生成随机点
    x = np.random.uniform(-1, 1, n_samples)
    y = np.random.uniform(-1, 1, n_samples)

    # 2. 计算点到原点的距离平方
    dist_sq = x**2 + y**2

    # 3. 统计落在圆内的点数 (距离平方 <= 1)
    m_inside = np.sum(dist_sq <= 1)

    # 4. 计算 Pi ≈ 4 * (M / N)
    pi_estimate = 4 * m_inside / n_samples
    return pi_estimate

def monte_carlo_integration(func, a, b, n_samples):
    """
    通过蒙特卡洛方法估算一元函数定积分
    """
    # 1. 在积分区间 [a, b] 上均匀采样
    x_samples = np.random.uniform(a, b, n_samples)

    # 2. 计算函数值的平均值
    y_mean = np.mean(func(x_samples))

    # 3. 积分值 = 区间长度 * 平均高度
    integral = (b - a) * y_mean
    return integral

# --- 测试 ---
if __name__ == "__main__":
    N = 10_000_000

    # 测试 1: 求 Pi
    pi_val = monte_carlo_pi(N)
    print(f"Sample N={N}, Pi estimate: {pi_val:.6f}")

    # 测试 2: 求积分 ∫(1/(1+sin(x)*ln(x)^2)) dx 在 [0.8, 3]
    # (这是课程中提到的那个复杂函数)
    def complex_func(x):
        return 1 / (1 + np.sin(x) * (np.log(x)**2))

    integral_val = monte_carlo_integration(complex_func, 0.8, 3.0, N)
    print(f"Integral estimate: {integral_val:.6f}")
```

### 关键记忆点

1.  **大数定律**是理论保证。
2.  **精度与 $\sqrt{N}$ 成反比**，为了提高精度需要付出平方级的计算成本。
3.  **应用核心**：求积分、求期望（机器学习基础）。
4.  **采样策略**：普通积分用均匀采样，求期望需按分布 $p(x)$ 采样。
