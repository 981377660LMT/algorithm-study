# Markdown 中 LaTex 数学公式命令
https://juejin.cn/post/6844903733248131080#heading-1
https://www.zybuluo.com/codeep/note/163962#5%E6%B7%BB%E5%8A%A0%E5%88%A0%E9%99%A4%E7%BA%BF

## Markdown 中使用 LaTeX 基础语法

LaTeX 公式有两种，一种是用在正文中的，一种是单独显示的。

- 行内公式：用 $formula$ 表示
  例如: $\sum_{i=0}^{n}i^2$
- 独立公式：用 $$formula$$ 
  例如: $$ \sum_{i=0}^{n}i^2 $$

- 带编号的公式
  若需要`手动编号`，可在公式后使用 `\tag{编号}` 语句。
  $$ \sum_{i=0}^{n}i^2 \tag{1-1} $$
  自动编号后的公式可在全文任意处使用 \eqref{eq:公式名} 语句引用。

## 常用数学表达命令

### 上下标表示
如果上下标的内容多于一个字符，需要用 {} 将这些内容括成一个整体
- 上标: 用 ^ 后的内容表示上标
  例如: $x^2$
- 下标: 用 _ 后的内容表示下标
  例如: $x_2$， $x_{2}$
- 上下标混用，例如: $x_1^2$ , $x^{y^{z} }$ , $x^{y_z}$

### 分数表示
- 分数: 用 \frac{分子}{分母} 表示分数
  例如: $\frac{1}{2}$

### 根式表示
- 根式: 用 \sqrt{内容} 表示根式
  例如: $\sqrt{2}$
- n 次根式: 用 \sqrt[n]{内容} 表示 n 次根式
  例如: $\sqrt[3]{2}$

### 向量表示
- 向量: 用 \vec{内容} 表示向量
  例如: $\vec{a}$

### 空白间距 - 占位宽度
- 空白间距: 用 \quad 表示空白间距
  例如: $a\quad b$
- 1/3 个空白间距
  $a\ b$
- 两个空格
   $a\qquad b$

### 省略号
- 省略号: 用 \cdots \dots 表示省略号
  例如: $a\cdots b$， $a\dots b$

### 分支公式
- 分支公式: 用 \begin{cases} \end{cases} 表示分支公式
  例如: $\begin{cases} x=1 & \text{if } y=1 \\ x=2 & \text{if } y=2 \end{cases}$

### 矩阵
- 矩阵: 用 \begin{matrix} \end{matrix} 表示矩阵
  生成矩阵的命令中每一行以 \\\ 结束，矩阵的元素之间用 & 来分隔开
  例如: $\begin{matrix} 1 & 2 \\ 3 & 4 \end{matrix}$

### 积分
- 积分: 用 \int_{下限}^{上限} 表示积分
  例如: $\int_{0}^{1} x^2 dx$

### 极限
- 极限: 用 \lim_{x \to 0} 表示极限
  例如: $\lim_{x \to 0} x^2$
  $$ \lim_{n \to \infty} \frac{1}{n(n+1)} \quad and \quad \lim_{x\leftarrow{示例}} \frac{1}{n(n+1)} $$