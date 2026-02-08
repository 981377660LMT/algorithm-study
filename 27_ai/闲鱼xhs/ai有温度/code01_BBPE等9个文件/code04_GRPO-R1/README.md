> 小红书 @AI有温度

### 视频讲解

[http://xhslink.com/n/6nxt7jc9IUN](http://xhslink.com/n/6nxt7jc9IUN)



### 数据下载

也可以直接解压

```Python
from modelscope.msdatasets import MsDataset
ds =  MsDataset.load('testUser/GSM8K_zh', subset_name='default', split='train')
```



### 使用Transformers的trl库实现

大家最好可以自己debug一下，会更清楚