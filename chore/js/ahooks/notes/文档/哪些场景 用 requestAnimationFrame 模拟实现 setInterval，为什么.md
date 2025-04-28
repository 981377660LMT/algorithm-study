用 `requestAnimationFrame` 模拟实现 `setInterval` 的常见场景有：

---

### 1. **动画相关的定时任务**

- 比如页面元素平滑移动、canvas 绘制、进度条动画等。
- 需要和浏览器刷新率同步，保证动画流畅不卡顿。

### 2. **高性能、低功耗需求**

- `requestAnimationFrame` 只有页面可见时才执行，页面切到后台会自动暂停，节省资源。
- `setInterval` 即使页面不可见也会继续执行，浪费性能。

### 3. **避免掉帧和卡顿**

- `requestAnimationFrame` 会自动和浏览器的渲染节奏对齐，减少掉帧。
- `setInterval` 可能和渲染不同步，导致动画不流畅。

---

#### 为什么用 `requestAnimationFrame` 替代 `setInterval`？

- **动画更流畅**：`requestAnimationFrame` 保证每一帧都在浏览器渲染前执行，动画不会撕裂或卡顿。
- **节能省电**：页面不可见时自动暂停，省资源。
- **时间间隔更精准**：不会因为主线程阻塞而积压回调。

---

**一句话总结：**  
**动画、性能敏感、需要和浏览器刷新同步的场景，优先用 `requestAnimationFrame` 替代 `setInterval`。**

---

请注意，如下两种情况下很可能是不适用的，优先考虑普通的 setInterval ：

- 时间间隔小于 16ms
- 希望页面不渲染的情况下依然执行定时器
