https://github.com/wangshusen/DeepLearning
https://www.bilibili.com/video/BV1YK4y1G7jw

# 深度学习 (Deep Learning)

## 1. 机器学习基础 (Machine Learning Basics)

- [机器学习概览 (ML basics)](#11-机器学习概览-ml-basics)
- [回归 (Regression)](#12-回归-regression)
- [分类 (Classification)](#13-分类-classification)
- [正则化 (Regularizations)](#14-正则化-regularizations)
- [聚类 (Clustering)](#15-聚类-clustering)
- [降维 (Dimensionality reduction)](#16-降维-dimensionality-reduction)
- [科学计算库 (Scientific computing libraries)](#17-科学计算库-scientific-computing-libraries)

## 2. 神经网络基础 (Neural Network Basics)

- [多层感知机与反向传播 (Multilayer perceptron and backpropagation)](#21-多层感知机与反向传播-multilayer-perceptron-and-backpropagation)
- [Keras 框架](#22-keras-框架)
- [激活函数、参数初始化与优化算法 (Further reading)](#23-激活函数参数初始化与优化算法)

## 3. 卷积神经网络 (Convolutional Neural Networks - CNNs)

- [CNN 基础 (CNN basics)](#31-cnn-基础-cnn-basics)
- [提升测试准确率的技巧 (Tricks for improving test accuracy)](#32-提升测试准确率的技巧-tricks-for-improving-test-accuracy)
- [特征缩放与批归一化 (Feature scaling and batch normalization)](#33-特征缩放与批归一化-feature-scaling-and-batch-normalization)
- [CNN 高级话题 (Advanced topics on CNNs)](#34-cnn-高级话题-advanced-topics-on-cnns)
- [主流 CNN 架构 (Popular CNN architectures)](#35-主流-cnn-架构-popular-cnn-architectures)

## 4. 循环神经网络 (Recurrent Neural Networks - RNNs)

- [类别特征处理 (Categorical feature processing)](#41-类别特征处理-categorical-feature-processing)
- [文本处理与词向量 (Text processing and word embedding)](#42-文本处理与词向量-text-processing-and-word-embedding)
- [RNN 基础 (RNN basics)](#43-rnn-基础-rnn-basics)
- [长短期记忆网络 (LSTM)](#44-长短期记忆网络-lstm)
- [提升 RNN 效果 (Making RNNs more effective)](#45-提升-rnn-效果-making-rnns-more-effective)
- [文本生成 (Text generation)](#46-文本生成-text-generation)
- [机器翻译 (Machine translation)](#47-机器翻译-machine-translation)
- [注意力机制 (Attention)](#48-注意力机制-attention)
- [自注意力机制 (Self-attention)](#49-自注意力机制-self-attention)
- [图像描述生成 (Image caption generation)](#410-图像描述生成-image-caption-generation)

## 5. Transformer 模型 (Transformer Models)

- [Transformer (1/2)：不含 RNN 的注意力机制](#51-transformer-12不含-rnn-的注意力机制)
- [Transformer (2/2)：从浅层到深层](#52-transformer-22从浅层到深层)
- [BERT：预训练 Transformer](#53-bert预训练-transformer)
- [Vision Transformer (ViT)](#54-vision-transformer-vit)

## 6. 自编码器 (Autoencoders)

- [用于降维的自编码器 (Autoencoder for dimensionality reduction)](#61-用于降维的自编码器-autoencoder-for-dimensionality-reduction)
- [用于图像生成的变分自编码器 (VAEs)](#62-用于图像生成的变分自编码器-vaes)

## 7. 生成对抗网络 (Generative Adversarial Networks - GANs)

- [DC-GAN](#71-dc-gan)

## 8. 深度强化学习 (Deep Reinforcement Learning)

- [强化学习基础 (Reinforcement learning basics)](#81-强化学习基础-reinforcement-learning-basics)
- [基于价值的学习 (Value-based learning)](#82-基于价值的学习-value-based-learning)
- [基于策略的学习 (Policy-based learning)](#83-基于策略的学习-policy-based-learning)
- [Actor-Critic 方法](#84-actor-critic-方法)
- [AlphaGo 与蒙特卡洛树搜索 (AlphaGo and MCTS)](#85-alphago-与蒙特卡洛树搜索-alphago-and-mcts)

## 9. 并行计算 (Parallel Computing)

- [基础与 MapReduce (Basics and MapReduce)](#91-基础与-mapreduce-basics-and-mapreduce)
- [参数服务器与去中心化网络 (Parameter server and decentralized network)](#92-参数服务器与去中心化网络-parameter-server-and-decentralized-network)
- [TensorFlow 镜像策略与 Ring All-reduce](#93-tensorflow-镜像策略与-ring-all-reduce)
- [联邦学习 (Federated learning)](#94-联邦学习-federated-learning)

## 10. 对抗鲁棒性 (Adversarial Robustness)

- [数据规避攻击与防御 (Data evasion attack and defense)](#101-数据规避攻击与防御-data-evasion-attack-and-defense)
- [数据投毒攻击 (Data poisoning attack)](#102-数据投毒攻击-data-poisoning-attack)

## 11. 元学习 (Meta Learning)

- [小样本学习：基本概念 (Few-shot learning: basic concepts)](#111-小样本学习基本概念-few-shot-learning-basic-concepts)
- [孪生网络 (Siamese network)](#112-孪生网络-siamese-network)
- [预训练 + 微调 (Pretraining + fine tuning)](#113-预训练--微调-pretraining--fine-tuning)

## 12. 神经网络架构搜索 (Neural Architecture Search - NAS)

- [NAS 基础 (Basics)](#121-nas-基础-basics)
- [RNN + 强化学习 (RNN + Reinforcement Learning)](#122-rnn--强化学习-rnn--reinforcement-learning)
- [可微架构搜索 (Differentiable NAS)](#123-可微架构搜索-differentiable-nas)
