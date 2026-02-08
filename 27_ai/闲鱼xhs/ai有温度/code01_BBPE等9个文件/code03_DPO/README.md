> 小红书 @AI有温度

此dpo消耗显存约17G

### 模型下载

此项目为了让大家都可以把模型跑起来，所以使用qwen2.5-0.5B：

- [https://www.modelscope.cn/models/Qwen/Qwen2.5-0.5B-Instruct](https://www.modelscope.cn/models/Qwen/Qwen2.5-0.5B-Instruct)

### 训练数据

data/train_data.json 此数据仅用于学习使用，这里可以清楚的看到DPO需要的数据格式，大家也可以寻找开源数据集进行训练。

### config.py配置文件

使用前需要修改模型地址与模型储存地址

### 运行顺序

1. main.py进行训练
2. inference.py进行推理
3. dpo.py实现了训练的过程与步骤

### 实际应用

1. 可以自己按格式构建数据集
2. 可以更换更大更好的模型
3. config中`is_reload_trained_params`可以选择是否根据上次的结果是否继续训练