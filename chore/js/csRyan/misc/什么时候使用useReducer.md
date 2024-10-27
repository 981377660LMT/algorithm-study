什么时候使用useReducer

1. 当状态更新逻辑比较复杂的时候，就应该考虑使用useReducer
   - `reducer比setState更加擅长描述“如何更新状态”。`比如，reducer能够读取相关的状态、同时更新多个状态。
   - 组件发 action，reducer 负责更新状态。
