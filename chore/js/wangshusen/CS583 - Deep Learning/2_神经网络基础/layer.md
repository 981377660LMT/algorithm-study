在设计神经网络时，选择合适的层主要取决于**数据的类型**和**任务的目标**。

以下是常见层的使用场景和最佳实践指南：

### 1. 核心层 (Core Layers)

- **全连接层 (Dense / Fully Connected Layer)**

  - **场景**：处理非结构化数据（如表格数据），或作为网络的最后几层用于分类/回归。
  - **作用**：整合前面提取的特征，映射到最终输出。
  - **最佳实践**：通常放在卷积层或循环层之后。参数量巨大，容易过拟合，建议配合 Dropout 使用。

- **激活层 (Activation Layer)**

  - **ReLU**：隐层的首选，计算快，缓解梯度消失。
  - **Sigmoid / Softmax**：仅用于输出层。二分类用 Sigmoid，多分类用 Softmax。
  - **Tanh**：有时用于 RNN 中，因为其输出均值为 0。

- **Dropout 层**
  - **场景**：几乎所有容易过拟合的深度网络。
  - **作用**：随机丢弃神经元，防止过拟合。
  - **最佳实践**：通常跟在 Dense 层之后，比率设为 0.2 到 0.5。

### 2. 卷积层 (Convolutional Layers)

- **Conv2D**

  - **场景**：图像处理、视觉任务（CV）。
  - **作用**：提取空间特征（边缘、纹理、形状）。
  - **最佳实践**：通常使用 $3 \times 3$ 或 $5 \times 5$ 的小卷积核。随着层数加深，增加 Filter 数量（例如 32 -> 64 -> 128）。

- **Conv1D**
  - **场景**：序列数据（文本、时间序列、音频）。
  - **作用**：提取局部序列特征。
  - **最佳实践**：处理文本分类时有时比 RNN 更快更高效。

### 3. 池化层 (Pooling Layers)

- **MaxPooling**

  - **场景**：紧跟卷积层之后。
  - **作用**：降低数据维度，减少计算量，提取最显著特征（平移不变性）。

- **GlobalAveragePooling**
  - **场景**：替代全连接层之前的 Flatten 操作。
  - **作用**：极大减少参数量，防止过拟合，常用于现代 CNN（如 ResNet）。

### 4. 循环层 (Recurrent Layers)

- **LSTM / GRU**
  - **场景**：时间序列预测、NLP、语音识别。
  - **作用**：处理变长序列，捕捉长距离依赖。
  - **最佳实践**：GRU 比 LSTM 计算更少，效果往往相当。现在很多 NLP 任务逐渐被 Transformer (Attention) 结构取代。

### 5. 归一化层 (Normalization Layers)

- **Batch Normalization (BN)**

  - **场景**：CNN 和深层全连接网络。
  - **作用**：加速收敛，允许更大的学习率。
  - **最佳实践**：通常放在卷积层/全连接层之后，激活函数之前（关于放在激活前还是后有争议，但放在激活前是原始论文建议）。

- **Layer Normalization (LN)**
  - **场景**：RNN 和 Transformer。
  - **作用**：不仅依赖 Batch 大小，适合序列模型。

### 6. 常见任务架构速查表

| 任务类型         | 推荐输入处理    | 核心架构                                | 输出层                  | 损失函数                        |
| :--------------- | :-------------- | :-------------------------------------- | :---------------------- | :------------------------------ |
| **图像分类**     | Conv2D          | CNN (ResNet, VGG) + BatchNorm + Pooling | Dense + Softmax         | Categorical Crossentropy        |
| **文本分类**     | Embedding       | LSTM/GRU 或 Conv1D 或 Transformer       | Dense + Softmax         | Categorical Crossentropy        |
| **时间序列回归** | 如果是滑动窗口  | LSTM/GRU                                | Dense (线性激活)        | MSE / MAE                       |
| **表格数据分类** | 归一化/独热编码 | 多层 Dense + Dropout                    | Dense + Sigmoid/Softmax | Binary/Categorical Crossentropy |

### 示例代码 (Keras/TensorFlow)

这是一个典型的图像分类模型结构，展示了层的组合逻辑：

```python
import tensorflow as tensorflow
from tensorflow.keras import layers, models

model = models.Sequential()

# 1. 特征提取部分 (卷积基)
# Conv2D 用于提取特征
# BatchNormalization 加速收敛
# Activation 引入非线性
# MaxPooling 降维
model.add(layers.Conv2D(32, (3, 3), input_shape=(64, 64, 3), use_bias=False))
model.add(layers.BatchNormalization())
model.add(layers.Activation('relu'))
model.add(layers.MaxPooling2D((2, 2)))

model.add(layers.Conv2D(64, (3, 3), use_bias=False))
model.add(layers.BatchNormalization())
model.add(layers.Activation('relu'))
model.add(layers.MaxPooling2D((2, 2)))

# 2. 分类头部分
# GlobalAveragePooling2D 将三维特征图压缩为一维向量，替代 Flatten
model.add(layers.GlobalAveragePooling2D())

# Dropout 防止过拟合
model.add(layers.Dropout(0.5))

# Dense 输出最终分类概率 (假设有 10 类)
model.add(layers.Dense(10, activation='softmax'))

model.summary()
```
