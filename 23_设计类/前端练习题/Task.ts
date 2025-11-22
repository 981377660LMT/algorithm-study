// 可调度、可取消、可监控、支持依赖关系的任务运行时系统
// C# 的 Task、Rust 的 Future 以及现代微内核架构中的任务调度思想
// Task、Schedular、Cancellation
//
// - 状态机管理：Pending -> Running -> Completed/Failed/Canceled。
// - 任务调度器 (Scheduler)：控制并发度，支持优先级。
// - 取消机制 (Cancellation)：支持级联取消（父任务取消，子任务自动取消）。
// - 依赖管理：支持 continueWith (链式调用) 和 Task.whenAll。
