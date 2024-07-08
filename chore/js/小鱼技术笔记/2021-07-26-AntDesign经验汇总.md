2021-07-26-AntDesign 经验汇总

Ant Design 是较为完整的 React 组件库。我们平时开发 UI 都是 div+flex+css 使用到底的，这其实是一种低层次的抽象。组件库的目的的是，`让开发者的抽象提高上来`，接触的不是 div，而是 List，Card，Tab 这样的组件，尽量避免使用 css 和 flex，而是要考虑如何用现有的组件去实现它。在学习的过程中，我们要注意，不同组件之间的抽象的共性，例如，Ant Design 的组件都有 title，extra，footer，icon，shape，children 的这些共同的抽象，尝试理解这些共性，才能更好地使用它。

一个后台管理系统的组件库，它主要包括以下几个模块：

- 页面布局组件，ProCard，PageContainer，ProLayout
- 基础展示组件，Button，Badge，Tag，Card，Statistic
- 数组展示组件，TreeSelect，List，Table
- 输入组件，Input，InputNumber，CheckBox，Select，DatePicker 等等。
