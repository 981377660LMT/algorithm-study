> 小红书 @AI有温度

## LoRA微调

### 数据集

**这个数据集较大，请自行下载**

使用医疗领域垂域数据集。

地址：[https://www.modelscope.cn/datasets/xiaofengalg/Chinese-medical-dialogue](https://www.modelscope.cn/datasets/xiaofengalg/Chinese-medical-dialogue)

```Python
from modelscope.msdatasets import MsDataset
ds =  MsDataset.load('xiaofengalg/Chinese-medical-dialogue', cache_dir='xxx', subset_name='default', split='train')
```

### 更换自己的模型地址

```Python
model_name = "/data/xxx/LLMs/Qwen/Qwen2.5-0.5B-Instruct"
output_dir="/data/xxx/lora_output"

```



## 数据蒸馏

> 所需显存:14G

使用DeepSeek数据蒸馏Qwen模型，大模型目前所说的蒸馏就是全参SFT。

### 数据集

**这个数据集较大，请自行下载**

地址：[https://www.modelscope.cn/datasets/liucong/Chinese-DeepSeek-R1-Distill-data-110k-SFT](https://www.modelscope.cn/datasets/liucong/Chinese-DeepSeek-R1-Distill-data-110k-SFT)

```Python
from modelscope.msdatasets import MsDataset
ds =  MsDataset.load('liucong/Chinese-DeepSeek-R1-Distill-data-110k-SFT', cache_dir='xxx', subset_name='default', split='train')

```