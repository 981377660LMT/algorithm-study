# ahooks analysis

https://gpingfeng.github.io/ahooks-analysis/guide/blog

- React hooks utils 库主要解决的两个问题如下：
  公共逻辑的抽象。
  解决 React hooks 存在的弊端，比如闭包等。
- 使用 useUrlState 这个 hook，需要独立安装 @ahooksjs/use-url-state，其源码在 packages/use-url-state 中。我理解官方的用意应该是这个功能依赖于 react-router，可能有一些项目不需要用到，把它提出来有助于减少包的大小。
- peerDependencies 的目的是`提示宿主环境去安装满足插件 peerDependencies 所指定依赖的包`，然后在插件 import 或者 require 所依赖的包的时候，永远都是引用宿主环境统一安装的 npm 包，最终解决插件与所依赖包不一致的问题。`这里的宿主环境一般指的就是我们自己的项目本身了。`

- DOM 规范
  https://ahooks.js.org/zh-CN/guide/dom/

  - target 支持三种类型 React.MutableRefObject（通过 useRef 保存的 DOM）、HTMLElement、() => HTMLElement（一般运用于 SSR 场景）
    https://github.com/GpingFeng/hooks/blob/guangping/read-code/packages/hooks/src/utils/domTarget.ts
  - DOM 类 Hooks 的 target 是支持动态变化的

- 一个优秀的工具库应该有自己的**一套输入输出规范**，一来能够支持更多的场景，二来可以更好的在内部进行封装处理，三来使用者能够更加快速熟悉和使用相应的功能，能做到举一反三。
