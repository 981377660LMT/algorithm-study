/**
 * @type {import('eslint').Linter.Config}
 */
module.exports = {
  parser: '@typescript-eslint/parser',
  extends: ['airbnb-base', 'plugin:@typescript-eslint/recommended'],
  plugins: ['@typescript-eslint', 'import'],
  env: {
    browser: true,
    node: true,
    es2021: true
  },
  parserOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module'
  },
  rules: {
    '@typescript-eslint/no-extra-semi': 0, // 不警告行首分号
    '@typescript-eslint/no-non-null-assertion': 1, // 警告非空断言

    semi: 0, // 行尾不需要分号
    'comma-dangle': 0, // 不检查尾随逗号
    'linebreak-style': 0, // 不检查换行符
    'import/prefer-default-export': 0, // 导出时可以不必使用默认导出
    'space-before-function-paren': 0,
    'function-paren-newline': 0,
    'no-use-before-define': 0, // 允许interface、type不同的定义顺序
    'no-underscore-dangle': 0, // 允许使用下划线开头的变量名
    'object-curly-newline': 0, // 允许对象的属性换行
    'import/no-unresolved': 0, // 模块识别
    'import/extensions': 0, // 导入可以不加后缀名
    'arrow-parens': 0, // 一个变量的箭头函数不需要括号
    'no-plusplus': 0, // 允许使用++
    'no-unused-expressions': 0, // 允许使用未使用的表达式 例如使用++
    'no-sequences': 0, // 允许使用逗号操作符
    'no-continue': 0, // 允许使用continue
    'no-void': 0, // 允许使用void
    indent: 0, // 不检查缩进(编辑器经常会报错)
    'no-bitwise': 0, // 允许使用位运算,
    'operator-linebreak': 0, // 允许操作符结尾
    'lines-between-class-members': 0, // 允许类中间有空行
    'max-classes-per-file': 0, // 允许一个文件中定义多个类
    'prefer-const': 0, // 关闭自动转换为const
    'no-restricted-syntax': 0, // 关闭使用迭代器遍历数组的警告(这样很慢/heavyweight)

    'no-shadow': 1, // 警告声明变量名与已声明变量名重名
    eqeqeq: 1, // 警告使用 == (在判空的时候兼容undefined和null).
    'no-param-reassign': 1, // 警告不允许修改函数参数,
    'class-methods-use-this': 1 // 警告类中的方法不使用this
  }
}
