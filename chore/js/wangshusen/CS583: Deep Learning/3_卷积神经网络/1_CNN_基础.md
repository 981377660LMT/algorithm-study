卷积神经网络（CNN）。本部分重点介绍 CNN 及其在计算机视觉问题中的应用。

# CNN\_基础

这是一份关于**卷积神经网络 (CNN)** 基础知识及其在 **Keras** 中实现的教学课件（Slides 7_CNN_1.pdf）。

通过对文档内容的分析与解构，我们可以将其分为**理论基础**、**实战演练 (MNIST)** 和 **实验结果分析** 三个主要部分。该课件通过对比简单的 Softmax 分类器与 CNN 模型，展示了 CNN 在图像处理任务上的强大能力。

以下是详细的深度解构：

### 1. 理论基础解构 (Core Concepts)

这部分主要定义了 CNN 的核心组件和运作机制。

- **卷积操作 (Convolution)**

  - **定义**: 介绍了矩阵卷积和张量卷积的概念。
  - **核心元素**:
    - **Filter / Kernel (滤波器/卷积核)**: 用于提取特征的小矩阵。
    - **Patch (图块)**: 输入图像中与滤波器进行运算的局部区域。
    - **Feature Map (特征图)**: 卷积操作的输出结果。
  - **控制参数**:
    - **Stride (步幅)**: 滤波器滑动的步长，默认通常为 1。
    - **Zero-padding (零填充)**: 在图像边缘填充 0，用于控制输出尺寸（默认为无填充）。

- **CNN 架构组件**
  - **卷积层 (Convolutional Layer)**: 包含滤波器数量、形状、步幅和填充方式等超参数。
  - **激活函数 (Activation Functions)**: 引入非线性（如 ReLU）。
  - **池化层 (Pooling Layer)**: 用于降维和提取主要特征。
    - **类型**: MaxPool (最大池化), AveragePool (平均池化)。
    - **参数**: 池化尺寸 (Pool size) 和 步幅 (Pool stride)。

### 2. 实战演练解构 (Implementation via Keras)

课件使用经典的 **MNIST 手写数字数据集** 进行了演示，并通过对比 **Softmax 线性分类器** 和 **CNN** 的性能差异来突出 CNN 的优势。

#### A. 数据准备 (Data Preparation)

- **加载数据**: 使用 `keras.datasets.mnist` 加载数据。
- **向量化 (Vectorization)**:
  - 对于 Softmax 模型，将 28x28 的图像展平 (Reshape) 为 784 维的向量 (`60000, 784`)。
  - _注：CNN 输入通常保持 (28, 28, 1) 的空间结构，但在 Softmax 部分代码中明确展示了展平操作。_
- **One-hot 编码**: 将标签 (Label) 转换为 10 维的向量（例如数字 3 变为 `[0,0,0,1,0,...]`）。
- **数据集划分**: 训练集 (50,000)、验证集 (10,000)、测试集 (10,000)。

#### B. 基准模型：Softmax 分类器 (The Baseline)

- **模型构建**:
  - 使用 `Sequential` 模型。
  - 仅包含一个全连接层 (`Dense`)，输入 784 维，输出 10 维。
  - 激活函数：`softmax`。
- **参数量**: 784 (权重) \* 10 + 10 (偏置) = **7,850** 个参数。
- **训练配置**:
  - 优化器: RMSprop (学习率 0.0001)。
  - 损失函数: `categorical_crossentropy`。
  - Epochs: 50, Batch size: 128。
- **性能表现**:
  - 训练过程显示 accuracy 逐渐上升，但最终稳定在 **92%** 左右 (Training ~92.3%, Validation ~90.8%)。
  - **局限性**: 也就意味着丢失了图像的空间结构信息。

#### C. 进阶模型：卷积神经网络 (The CNN Solution)

_(虽然提取的文本中 Softmax 的代码更完整，但根据头部的高精度日志和 Summary，可以重构出 CNN 部分的关键信息)_

- **模型构建**:
  - 构建包含 `Conv`, `Pool`, `FC` (全连接) 层的网络。
- **性能表现**:
  - 日志显示训练过程准确率极高。
  - **Training accuracy**: **99.3%**
  - **Validation accuracy**: **98.8%**
  - **Test accuracy**: **98.8%**
- **结论**: 相比 Softmax 基准模型 (~91%)，CNN 将错误率降低了绝大部分，证明了其提取局部空间特征的有效性。

### 3. Keras 工具链解构

课件展示了使用 Keras 进行深度学习开发的标准流程：

1.  **Import**: 引入 `models`, `layers`, `optimizers`。
2.  **Dataset**: `batch_size`, `input_shape`, `to_one_hot` 处理。
3.  **Build**: `model.add(...)` 堆叠层。
4.  **Inspect**: `model.summary()` 查看参数量和 output shape。
5.  **Compile**: 指定 Optimizer, Loss, Metrics。
6.  **Fit**: 传入训练数据和验证数据 (`validation_data`) 进行训练，保存 `history`。
7.  **Evaluate**: 使用 `model.evaluate` 在测试集上计算 Loss 和 Accuracy。
8.  **Visualize**: 使用 `matplotlib` 绘制 Loss/Accuracy 随 Epoch 变化的曲线，以及观察预测错误的样本。
9.  **Internals**: 使用 `K.get_session().run(model.trainable_weights)` 获取并可视化权重（虽然对于 CNN 滤波器可视化更有意义，此处展示了 Dense 层的权重矩阵）。

### 总结

这份课件是一个非常标准的**CNN 入门教程**。它没有直接抛出复杂的数学公式，而是采取了**“原理概念 -> 代码实现 -> 效果对比”**的教学路径。通过 Softmax (全连接) 和 CNN 在同一任务上的巨大性能鸿沟 (91% vs 99%)，直观地让学生理解卷积层在处理图像数据时的必要性。
