https://github.com/wangshusen/DRL?tab=readme-ov-file#deep-reinforcement-learning
https://github.com/wangshusen/deep-rl/tree/master
https://www.bilibili.com/video/BV12o4y197US

# 深度强化学习 (Deep Reinforcement Learning)

## 1. 概要 (Overview)

- [强化学习基础 (Reinforcement Learning)](#11-强化学习基础-reinforcement-learning)
- [基于价值的学习 (Value-Based Learning)](#12-基于价值的学习-value-based-learning)
- [基于策略的学习 (Policy-Based Learning)](#13-基于策略的学习-policy-based-learning)
- [Actor-Critic 方法 (Actor-Critic Methods)](#14-actor-critic-方法-actor-critic-methods)
- [AlphaGo](#15-alphago)

## 2. 时序差分学习 (TD Learning)

- [Sarsa](#21-sarsa)
- [Q-learning](#22-q-learning)
- [多步时序差分目标 (Multi-Step TD Target)](#23-多步时序差分目标-multi-step-td-target)

## 3. 基于价值学习的高级话题 (Advanced Topics on Value-Based Learning)

- [经验回放与优先经验回放 (Experience Replay & Prioritized ER)](#31-经验回放与优先经验回放-experience-replay--prioritized-er)
- [高估问题、目标网络与 Double DQN (Overestimation, Target Network, & Double DQN)](#32-高估问题目标网络与-double-dqn-overestimation-target-network--double-dqn)
- [Dueling Networks](#33-dueling-networks)

## 4. 带基准的策略梯度 (Policy Gradient with Baseline)

- [带基准的策略梯度 (Policy Gradient with Baseline)](#41-带基准的策略梯度-policy-gradient-with-baseline)
- [带基准的 REINFORCE (REINFORCE with Baseline)](#42-带基准的-reinforce-reinforce-with-baseline)
- [Advantage Actor-Critic (A2C)](#43-advantage-actor-critic-a2c)
- [REINFORCE 与 A2C 对比 (REINFORCE versus A2C)](#44-reinforce-与-a2c-对比-reinforce-versus-a2c)

## 5. 基于策略学习的高级话题 (Advanced Topics on Policy-Based Learning)

- [置信区域策略优化 (Trust-Region Policy Optimization - TRPO)](#51-置信区域策略优化-trpo)
- [部分观测与 RNN (Partial Observation and RNNs)](#52-部分观测与-rnn-partial-observation-and-rnns)

## 6. 连续动作空间处理 (Dealing with Continuous Action Space)

- [离散控制与连续控制 (Discrete versus Continuous Control)](#61-离散控制与连续控制-discrete-versus-continuous-control)
- [连续控制的确定性策略梯度 (DPG for Continuous Control)](#62-连续控制的确定性策略梯度-dpg-for-continuous-control)
- [连续控制的随机性策略梯度 (Stochastic Policy Gradient for Continuous Control)](#63-连续控制的随机性策略梯度-stochastic-policy-gradient-for-continuous-control)

## 7. 多智能体强化学习 (Multi-Agent Reinforcement Learning)

- [基础与挑战 (Basics and Challenges)](#71-基础与挑战-basics-and-challenges)
- [中心化与去中心化 (Centralized VS Decentralized)](#72-中心化与去中心化-centralized-vs-decentralized)

## 8. 模仿学习 (Imitation Learning)

- [逆强化学习 (Inverse Reinforcement Learning)](#81-逆强化学习-inverse-reinforcement-learning)
- [生成对抗模仿学习 (Generative Adversarial Imitation Learning - GAIL)](#82-生成对抗模仿学习-gail)
