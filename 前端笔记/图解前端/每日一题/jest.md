```js
// https://juejin.cn/post/6844903888512876558
// https://xin-tan.com/2019-05-04-jest-base/

// 安装
npm install -D jest ts-jest @types/jest

// 初始化配置文件
npx ts-jest config:init

// stmts是语句覆盖率（statement coverage）：每个语句是否都执行了
// Branch分支覆盖率（branch coverage）：条件语句是否都执行了复制代码
// Funcs函数覆盖率（function coverage）：函数是否全都调用了复制代码
//  Lines行覆盖率（line coverage）：未执行的代码行数


```

```JSON
  "scripts": {
    "test": "jest --config ./jest.config.js",
    "coverage": "npm run test -- --coverage"  // --表示给jest传参
  },
```

找不到名称“jest”:ts 的选项 types 需要加入["jest"]
**最好的做法是不去管 types 选项** jest 会自动被识别

集成测试
持续继承测试我们借助 https://travis-ci.org/ 这个平台，它的工作流程非常简单：

1. 在它平台上授权 github 仓库的权限，github 仓库下配置 .travis.yml 文件
2. 每次 commit 推上新代码的时候，travis-ci 平台都会接收到通知
3. 读取 .travis.yml 文件，然后创建一个虚拟环境，来跑配置好的脚本（比如启动测试脚本）

vue-cli 的 jest 配置
https://github.com/vuejs/vue-cli/tree/dev/packages/%40vue/cli-plugin-unit-jest

常见的 ci
travis.ci
github action (.github/workflows/test-coverage.yml)
