from modelscope.msdatasets import MsDataset

# 医疗数据
ds =  MsDataset.load('xiaofengalg/Chinese-medical-dialogue', cache_dir='xxx', subset_name='default', split='train')

# 蒸馏数据
ds =  MsDataset.load('liucong/Chinese-DeepSeek-R1-Distill-data-110k-SFT', cache_dir='xxx', subset_name='default', split='train')