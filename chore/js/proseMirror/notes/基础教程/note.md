## 背景

![alt text](image.png)

- 马金-哈弗贝克
- CodeMirror、ProseMirror、Lezer、Acorn、Eloquent JavaScript
- 模块化、定制化

## 四大核心模块

![alt text](image-1.png)

## 文档(doc)

1. 与 dom 区别
   - 保存行内元素的方式不同
   - 可持久化
     ![alt text](image-2.png)
2. 三种创建方式
   - 直接创建
   - 从 JSON 解析
   - 从 DOM 解析
3. 索引、位置、偏移量

## Schema

1. 创建
   ![alt text](image-3.png)
2. 内容表达式
   ![alt text](image-4.png)
   - 节点顺序很重要
     ![alt text](image-5.png)
3. 序列化与解析
   ![alt text](image-6.png)

## 文档变换

1. 为什么要使用文档变换
   ![因为collab](image-7.png)
2. 步骤：Step
   ![alt text](image-8.png)
3. 变换：Transform
   ![alt text](image-9.png)
4. 位置映射：StepMap/Mapping
   ![alt text](image-10.png)
   ![alt text](image-11.png)
5. 理解 Rebase
   ![alt text](image-12.png)
   ![alt text](image-13.png)

## State

1. 编辑器状态中包含什么
   ![alt text](image-14.png)
   ![alt text](image-15.png)
2. 选区
   ![alt text](image-16.png)
3. 事务
   ![alt text](image-17.png)
   ![alt text](image-18.png)
4. 插件
   ![alt text](image-19.png)
   ![alt text](image-20.png)
   ![alt text](image-21.png)

## View

三种

1. 视图组件概述
   ![保持DOM选区与编辑器状态同步](image-22.png)
2. 数据流与高效更新
   ![MVC](image-23.png)
3. 视图属性
   ![props](image-24.png)
4. 视图装饰
   ![alt text](image-25.png)
   ![alt text](image-26.png)
5. 节点视图
   ![NodeView](image-27.png)
   ![alt text](image-28.png)
   ![alt text](image-29.png)

## Command

![alt text](image-30.png)
![alt text](image-31.png)

## Collab

1. 协同算法简介
   ![alt text](image-32.png)
2. 算法步骤
   rebase
3. 中心服务
   ![alt text](image-33.png)
4. collab 模块
   ![alt text](image-34.png)
