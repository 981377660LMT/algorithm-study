# 教育中的开发容器：讲师指南

链接：https://code.visualstudio.com/blogs/2020/08/24/dev-containers-education

## 深入分析

继续聚焦教育场景，这次针对讲师（老师和培训师）。

### 教育中的痛点
1. **环境差异** - 学生的电脑操作系统各异，讲师需要支持Windows、Mac、Linux
2. **安装地狱** - 安装各种开发工具耗时且容易出错
3. **版本不一致** - 同学之间的库版本不同，导致代码在某些电脑上跑不通

### DevContainer的解决方案
- 讲师只需提供一个Dockerfile，学生所有人的环境都完全相同
- 学生不需要理解Dockerfile的细节，只需运行一条命令启动容器

### 具体用例
1. **编程入门课** - Python、JavaScript基础
2. **Web开发课** - Node.js、Django、Rails等
3. **数据科学课** - Jupyter + 数据分析库

### 长期影响
- 这一推动，使得DevContainer逐渐成为教学的标准工具
- 许多大学的计算机科学系开始在课程中采用这一方法
