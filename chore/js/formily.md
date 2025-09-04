https://fishedee.com/2021/07/21/Formily%E7%9A%84React%E5%BA%93%E7%BB%8F%E9%AA%8C%E6%B1%87%E6%80%BB/#json-schema-1

好的，下面将对 Formily 进行深度、详细、具体的讲解，涵盖其核心原理、常用 API、典型用法、与传统表单的区别、进阶特性（如联动、动态表单、校验、扩展性等），并配合代码示例，帮助你从零到一理解和掌握 Formily。

---

# 一、Formily 的核心原理

Formily 是一个响应式、声明式的表单框架。它的核心思想是**数据驱动视图**，即你只需要声明表单的数据结构、校验规则和联动逻辑，Formily 会自动帮你管理表单的状态、渲染和交互。

## 1.1 表单状态管理

Formily 用一个“表单树”来管理所有字段的状态，包括字段值、校验状态、显示/隐藏、只读/禁用等。每个字段都是一个响应式对象，任何变化都会自动驱动 UI 更新。

## 1.2 响应式机制

Formily 内部用 [observable](https://mobx.js.org/observable-state.html) 响应式系统（类似 MobX），保证数据和视图的双向绑定。

---

# 二、Formily 的基本用法

## 2.1 安装

```bash
npm install @formily/core @formily/react
# 如果用 Ant Design 组件库
npm install @formily/antd antd
```

## 2.2 创建表单

```jsx
import { createForm } from '@formily/core'

const form = createForm()
```

## 2.3 表单上下文

```jsx
import { FormProvider } from '@formily/react'
;<FormProvider form={form}>{/* 这里放表单字段 */}</FormProvider>
```

## 2.4 声明字段

```jsx
import { Field } from '@formily/react'
;<Field name="username" title="用户名" required>
  {({ value, onChange }) => <input value={value} onChange={e => onChange(e.target.value)} />}
</Field>
```

- `name`：字段名
- `title`：标签
- `required`：是否必填
- `value`、`onChange`：字段值和变更事件

## 2.5 获取表单数据

```jsx
import { FormConsumer } from '@formily/react'
;<FormConsumer>{() => <pre>{JSON.stringify(form.values, null, 2)}</pre>}</FormConsumer>
```

---

# 三、与传统表单的区别

| 传统表单       | Formily              |
| -------------- | -------------------- |
| 手动管理 state | 自动响应式管理       |
| 校验逻辑分散   | 校验声明式集中       |
| 联动复杂难维护 | 联动声明式、集中管理 |
| 组件复用性差   | 组件高度解耦、可复用 |

---

# 四、进阶特性

## 4.1 配合 UI 组件库

Formily 支持主流 UI 库（如 Ant Design、Element Plus），通过 `x-component` 属性自动渲染对应组件。

```jsx
import { Form, Field } from '@formily/antd'
;<Form form={form}>
  <Field
    name="email"
    title="邮箱"
    required
    x-component="Input"
    x-component-props={{ placeholder: '请输入邮箱' }}
  />
</Form>
```

## 4.2 JSON Schema 声明表单

Formily 支持用 JSON Schema 声明整个表单结构，适合动态表单场景。

```jsx
import { createSchemaField } from '@formily/react'
import { Input } from 'antd'

const SchemaField = createSchemaField({
  components: { Input }
})

const schema = {
  type: 'object',
  properties: {
    username: {
      type: 'string',
      title: '用户名',
      required: true,
      'x-component': 'Input'
    }
  }
}

<SchemaField schema={schema} />
```

## 4.3 校验

支持多种校验方式：

- 内置校验（required、pattern、min/max 等）
- 自定义校验函数

```jsx
<Field
  name="age"
  title="年龄"
  required
  validator={[
    { required: true, message: '必填' },
    { min: 18, message: '必须大于18岁' },
    value => (value % 2 === 0 ? '' : '必须是偶数')
  ]}
  x-component="InputNumber"
/>
```

## 4.4 字段联动

通过 `effects` 实现字段间的联动：

```jsx
import { onFieldValueChange } from '@formily/core'

form.addEffects('联动', () => {
  onFieldValueChange('province', field => {
    form.setFieldState('city', state => {
      state.visible = !!field.value
    })
  })
})
```

## 4.5 动态表单

可以动态增减字段、数组项等：

```jsx
<Field name="users" title="用户列表" type="array" x-component="ArrayTable">
  <Field name="name" title="姓名" x-component="Input" />
  <Field name="age" title="年龄" x-component="InputNumber" />
</Field>
```

---

# 五、常用 API 总结

- `createForm()`：创建表单实例
- `FormProvider`：表单上下文
- `Field`：声明字段
- `FormConsumer`：消费表单状态
- `createSchemaField`：基于 JSON Schema 声明表单
- `form.values`：获取表单数据
- `form.submit()`：提交表单
- `form.setFieldState()`：设置字段状态
- `form.addEffects()`：添加联动逻辑

---

# 六、最佳实践

1. **优先用 Schema 声明表单**，便于动态扩展和维护。
2. **表单逻辑集中管理**，如校验、联动都放在 effects 里。
3. **组件解耦**，字段组件只关心自身 UI，逻辑交给 Formily 管理。
4. **善用响应式**，避免手动管理 state。

---

# 七、学习资源

- [Formily 官方文档](https://formilyjs.org/zh-CN/docs/guide)
- [Formily 示例代码](https://codesandbox.io/s/formily-demo-6k6wz)
- [Formily Github](https://github.com/alibaba/formily)

---

在 Formily 中，**decorator** 和 **component** 是两个非常核心的概念。它们分别负责表单字段的“结构包装”和“具体渲染”。理解这两个概念，有助于你灵活定制表单的外观和行为。

---

## 1. component（组件）

**component** 指的是字段实际渲染的输入控件，比如 Input、Select、Checkbox 等。  
在 Schema 或 Field 里通过 `x-component` 指定。

**作用**：负责字段的输入、展示和交互。

**示例：**

```jsx
<Field
  name="username"
  title="用户名"
  x-component="Input" // 这里的 Input 就是 component
/>
```

如果用 JSON Schema：

```js
{
  "type": "string",
  "title": "用户名",
  "x-component": "Input"
}
```

---

## 2. decorator（装饰器）

**decorator** 是用来包裹字段组件的结构性组件，通常用于布局、标签、校验提示等。  
比如常见的 Form.Item（Antd）、FormField（Element Plus）等。

**作用**：负责字段的结构、布局、标签、错误提示等。

**示例：**

```jsx
<Field
  name="username"
  title="用户名"
  x-decorator="FormItem" // 这里的 FormItem 就是 decorator
  x-component="Input"
/>
```

JSON Schema 写法：

```js
{
  "type": "string",
  "title": "用户名",
  "x-decorator": "FormItem",
  "x-component": "Input"
}
```

---

## 3. 关系与区别

- **component** 只负责“内容”，比如输入框本身。
- **decorator** 负责“外壳”，比如表单项的标签、校验提示、布局等。
- 通常一个字段会有一个 decorator 和一个 component。

---

## 4. 直观结构

比如下面的结构：

```jsx
<FormItem label="用户名">
  <Input />
</FormItem>
```

在 Formily 里就是：

```jsx
<Field name="username" title="用户名" x-decorator="FormItem" x-component="Input" />
```

---

## 5. 自定义 decorator 和 component

你可以自定义 decorator 或 component，然后在 SchemaField 注册：

```jsx
const SchemaField = createSchemaField({
  components: { Input, MyDecorator }
})

<SchemaField
  schema={{
    type: 'object',
    properties: {
      username: {
        type: 'string',
        title: '用户名',
        'x-decorator': 'MyDecorator',
        'x-component': 'Input'
      }
    }
  }}
/>
```

---

## 总结

- **component**：输入控件，负责内容和交互。
- **decorator**：结构包装，负责布局、标签、提示等。
- 两者配合，让表单既有良好的结构，又有丰富的交互。

---

# Formily 深度解析

## 一、整体架构与核心理念

Formily 是一套完整的表单解决方案，由多个核心包组成，形成了一个层次分明的架构。

### 1.1 整体架构

```
@formily/core          - 表单核心库，管理状态、校验、联动
@formily/react         - React适配层，提供React组件和Hooks
@formily/reactive      - 响应式库，类似MobX，提供响应式能力
@formily/path          - 路径解析库，处理表单字段路径
@formily/json-schema   - JSON Schema协议解析库
@formily/antd          - Ant Design适配层
@formily/element       - Element UI适配层
```

### 1.2 核心理念

1. **响应式表单**：用可观察数据模型驱动 UI 渲染
2. **声明式表单**：用 Schema 描述表单结构和行为
3. **效率至上**：局部更新+按需渲染，确保性能
4. **分层架构**：关注点分离，提高复用性和可维护性

## 二、核心包详解

### 2.1 @formily/core

作为最核心的包，管理整个表单的状态、校验和联动。

```jsx
import { createForm } from '@formily/core'

// 创建表单实例
const form = createForm({
  // 表单初始值
  initialValues: {
    username: 'admin',
    password: ''
  },
  // 表单校验器
  validateFirst: true, // 是否在第一个校验失败后停止校验
  // 表单效果
  effects() {
    // 监听字段变化并做出响应
  }
})
```

**重要 API：**

- `form.createField()` - 创建字段
- `form.setValues()` - 设置表单值
- `form.submit()` - 提交表单
- `form.validate()` - 触发校验

### 2.2 @formily/reactive

Formily 自研的响应式库，类似 MobX，但更轻量和性能更好。

```jsx
import { observable, autorun } from '@formily/reactive'

// 创建响应式对象
const obs = observable({
  name: 'formily'
})

// 自动响应变化
autorun(() => {
  console.log(obs.name) // 当obs.name变化时，自动执行
})

// 修改会触发autorun
obs.name = 'formily2'
```

### 2.3 @formily/react

将表单核心与 React 连接，提供上下文、Hook 和组件。

```jsx
import { FormProvider, Field, useField } from '@formily/react'

// FormProvider提供上下文
;<FormProvider form={form}>
  {/* Field组件处理字段渲染 */}
  <Field name="username" component={[Input]} />
</FormProvider>

// useField获取字段状态
const field = useField()
console.log(field.value, field.errors)
```

## 三、表单状态管理深度解析

### 3.1 表单状态树

Formily 维护一个表单状态树，每个字段是树的一个节点。

```jsx
// 表单状态树示例
{
  values: { username: 'admin', address: { city: 'Beijing' } },
  initialValues: { username: '', address: { city: '' } },
  valid: true,
  invalid: false,
  errors: [],
  warnings: [],
  fields: {
    username: { /* 字段状态 */ },
    'address.city': { /* 字段状态 */ }
  }
}
```

### 3.2 字段状态

字段是表单的核心，每个字段包含丰富的状态属性：

```jsx
// 字段状态示例
{
  name: 'username',
  path: 'username',
  value: 'admin',
  initialValue: '',
  required: true,
  visible: true,
  display: 'visible', // 'none' | 'hidden' | 'visible'
  disabled: false,
  readOnly: false,
  valid: true,
  invalid: false,
  validating: false,
  errors: [],
  warnings: []
}
```

### 3.3 状态变更机制

Formily 使用批处理和路径系统优化状态变更：

```jsx
// 批量更新
form.setFieldState('username', state => {
  state.value = 'new value'
  state.required = true
})

// 使用路径批量更新
form.setFieldState('address.*', state => {
  state.visible = false
})
```

## 四、Schema 驱动表单

### 4.1 JSON Schema 基础

使用 JSON 描述表单结构和行为：

```jsx
const schema = {
  type: 'object',
  properties: {
    username: {
      type: 'string',
      title: '用户名',
      required: true,
      'x-component': 'Input',
      'x-decorator': 'FormItem'
    },
    age: {
      type: 'number',
      title: '年龄',
      'x-component': 'InputNumber',
      'x-decorator': 'FormItem'
    }
  }
}
```

### 4.2 创建 Schema 组件

```jsx
import { createSchemaField } from '@formily/react'
import { Form, FormItem, Input, NumberPicker } from '@formily/antd'

// 创建Schema渲染器
const SchemaField = createSchemaField({
  components: {
    FormItem,
    Input,
    InputNumber: NumberPicker
  }
})

// 使用Schema
<Form form={form}>
  <SchemaField schema={schema} />
</Form>
```

### 4.3 Schema 表达式能力

Schema 支持表达式动态计算属性：

```jsx
{
  "type": "string",
  "title": "邮箱",
  "required": "{{$deps[0] === 'admin'}}",
  "x-component": "Input",
  "x-reactions": [
    {
      "dependencies": ["username"],
      "fulfill": {
        "state": {
          "visible": "{{$deps[0] !== ''}}"
        }
      }
    }
  ]
}
```

## 五、校验机制详解

### 5.1 校验方式

Formily 支持多种校验方式：

```jsx
// 1. Schema验证规则
{
  "type": "string",
  "required": true,
  "format": "email",
  "minLength": 5,
  "maxLength": 20
}

// 2. 自定义校验器
{
  "type": "string",
  "validator": (value) => {
    if (!/^1\d{10}$/.test(value)) {
      return '手机号格式不正确'
    }
    return ''
  }
}

// 3. 异步校验
{
  "type": "string",
  "validator": async (value) => {
    const res = await checkUsernameExist(value)
    if (res.exists) {
      return '用户名已存在'
    }
    return ''
  }
}
```

### 5.2 校验策略

```jsx
createForm({
  validateFirst: true, // 遇到第一个错误就停止校验
  // 或者
  validateFirst: 'parallel' // 并行校验所有规则，收集所有错误
})
```

## 六、联动机制详解

### 6.1 路径系统

Formily 的路径系统非常强大，支持通配符、相对路径等：

```
username           - 精确匹配
address.*          - 匹配address下所有字段
*.name             - 匹配任何层级下的name字段
address[*].street  - 匹配address数组内所有项的street
```

### 6.2 字段联动

#### 6.2.1 响应器 (Reactions)

```jsx
// Schema中定义联动
{
  "x-reactions": [
    {
      "dependencies": ["username"],
      "fulfill": {
        "state": {
          "visible": "{{$deps[0] !== ''}}"
        }
      }
    }
  ]
}
```

#### 6.2.2 Effects API

```jsx
import { createForm, onFieldValueChange } from '@formily/core'

const form = createForm({
  effects() {
    // 监听字段变化
    onFieldValueChange('username', field => {
      // 根据username字段值控制password字段
      form.setFieldState('password', state => {
        state.visible = field.value !== ''
      })
    })
  }
})
```

## 七、高级特性

### 7.1 表单分片

大型表单可以分片渲染提升性能：

```jsx
// 只有激活的分片会被渲染和校验
<FormTab>
  <FormTab.TabPane name="tab1" title="基本信息">
    <SchemaField schema={basicSchema} />
  </FormTab.TabPane>
  <FormTab.TabPane name="tab2" title="高级设置">
    <SchemaField schema={advancedSchema} />
  </FormTab.TabPane>
</FormTab>
```

### 7.2 字段递归

支持无限嵌套表单：

```jsx
// 递归组件示例
const RecursiveField = observer(({ field }) => {
  if (field.dataType === 'object') {
    return (
      <div>
        {field.fields.map(childField => (
          <RecursiveField key={childField.name} field={childField} />
        ))}
      </div>
    )
  }
  return <Input {...field} />
})
```

### 7.3 性能优化

Formily 通过多种方式优化表单性能：

1. **表单分片**：只渲染可见部分
2. **异步校验防抖**：避免频繁校验
3. **路径系统**：精准更新字段
4. **局部渲染**：只更新发生变化的组件

```jsx
// 使用observer包装组件，实现精准更新
import { observer } from '@formily/reactive-react'

const MyField = observer(({ field }) => {
  // 只有field相关属性变化才会重新渲染
  return <Input value={field.value} onChange={field.onChange} />
})
```

## 八、实际案例分析

### 8.1 动态数组表单

```jsx
const schema = {
  type: 'object',
  properties: {
    users: {
      type: 'array',
      'x-component': 'ArrayTable',
      items: {
        type: 'object',
        properties: {
          name: {
            type: 'string',
            'x-decorator': 'FormItem',
            'x-component': 'Input'
          },
          age: {
            type: 'number',
            'x-decorator': 'FormItem',
            'x-component': 'InputNumber'
          },
          address: {
            type: 'string',
            'x-decorator': 'FormItem',
            'x-component': 'Input'
          }
        }
      }
    }
  }
}

// 自定义数组操作
<ArrayTable>
  <ArrayTable.Column title="姓名">
    <Field name="name" component={[Input]} />
  </ArrayTable.Column>
  <ArrayTable.Column title="年龄">
    <Field name="age" component={[InputNumber]} />
  </ArrayTable.Column>
  <ArrayTable.Column title="操作">
    <Space>
      <ArrayTable.Remove />
      <ArrayTable.MoveDown />
      <ArrayTable.MoveUp />
    </Space>
  </ArrayTable.Column>
  <ArrayTable.Addition title="添加用户" />
</ArrayTable>
```

### 8.2 复杂联动表单

```jsx
const schema = {
  type: 'object',
  properties: {
    delivery: {
      type: 'string',
      title: '配送方式',
      enum: [
        { label: '快递', value: 'express' },
        { label: '自提', value: 'self' }
      ],
      'x-decorator': 'FormItem',
      'x-component': 'Select'
    },
    address: {
      type: 'object',
      title: '收货地址',
      'x-decorator': 'FormItem',
      'x-component': 'FormLayout',
      'x-reactions': {
        dependencies: ['delivery'],
        fulfill: {
          state: {
            visible: '{{$deps[0] === "express"}}'
          }
        }
      },
      properties: {
        province: {
          type: 'string',
          title: '省份',
          'x-decorator': 'FormItem',
          'x-component': 'Select',
          enum: [
            { label: '北京', value: 'bj' },
            { label: '上海', value: 'sh' }
          ]
        },
        city: {
          type: 'string',
          title: '城市',
          'x-decorator': 'FormItem',
          'x-component': 'Select',
          'x-reactions': {
            dependencies: ['.province'],
            fulfill: {
              state: {
                visible: '{{!!$deps[0]}}',
                dataSource: '{{getCities($deps[0])}}'
              }
            }
          }
        }
      }
    }
  }
}
```

## 九、Formily 开发最佳实践

1. **拆分大型表单**：使用分片和布局组件
2. **优先使用 Schema**：维护性和可扩展性更好
3. **集中管理校验和联动**：在 effects 中统一处理复杂逻辑
4. **善用表达式和响应器**：减少命令式代码
5. **组件分层设计**：decorator 负责结构，component 负责内容
6. **按需渲染优化**：使用 observer 和合理设计状态依赖
7. **抽象复用逻辑**：将常用逻辑封装为 Hook 或组件

## 十、总结

Formily 是一个功能全面的表单解决方案，其核心优势在于：

1. **响应式状态管理**：精准更新，高性能
2. **声明式表单设计**：低代码、易维护
3. **强大的联动机制**：满足复杂业务场景
4. **丰富的校验能力**：支持同步异步、多种规则
5. **扩展性极强**：可适配任何 UI 库和业务场景

掌握 Formily 需要理解其核心概念：表单状态、字段、Schema、联动和校验机制，以及如何合理使用这些功能构建高性能、易维护的表单系统。

---

好的，我们来深入讲解一下 Formily 以及其核心概念 `decorator` 和 `component`。

### Formily 是什么？

Formily 是一个由阿里巴巴开源的、专注于解决复杂表单场景的通用表单解决方案。它的核心思想是**将表单的数据、逻辑与 UI 完全分离**。

传统的表单库（如 antd Form）通常将数据状态、校验逻辑和 UI 渲染强绑定在一起，当表单变得非常复杂（例如，字段之间有复杂的联动、动态增删、不同状态下显隐/禁用等）时，代码会变得难以维护。

Formily 通过引入**面向对象的领域模型 (Formily Core)** 和 **MVVM (Model-View-ViewModel) 设计模式**来解决这个问题：

1.  **Model (模型层):** Formily Core 负责管理整个表单的数据状态、校验状态、联动逻辑等。它是一个纯粹的、与 UI 无关的数据逻辑核心。你可以把它想象成一个看不见的、在后台运行的表单“大脑”。
2.  **View (视图层):** 这是用户能看到的 UI 部分，比如输入框、选择框、标签、错误提示等。在你的代码中，这部分由 `@ecom/lander-ui-setting` 里的 `Input`、`Select` 等组件构成。
3.  **ViewModel (视图模型层):** 这是连接 Model 和 View 的桥梁。在 Formily 中，`Field`、`ArrayField`、`ObjectField` 等组件就扮演了这个角色。它们从 Model 层获取状态（如 `value`, `error`, `disabled`），然后将这些状态传递给具体的 UI 组件（View），同时监听 UI 组件的用户操作（如 `onChange`），再将这些操作反馈给 Model 层。

这种分离带来了巨大的优势：

- **高可维护性:** UI 和逻辑分离，修改联动逻辑不需要动 UI 代码，反之亦然。
- **高性能:** 只在必要时精确更新受影响的 UI 组件，避免不必要的重复渲染。
- **高扩展性:** 可以轻松地对接任何 UI 库（Ant Design, Fusion Design, Element UI 等），因为核心逻辑是独立的。

---

### `component` 和 `decorator` 是什么？

在 Formily 的 `Field` 组件中，`component` 和 `decorator` 是两个核心属性，它们共同定义了一个表单字段的完整 UI。这正是 Formily UI 分离思想的精髓体现。

我们可以用一个“**给照片装裱**”的例子来理解：

- **`component` (组件):** 相当于**照片本身**。它是表单字段的核心交互元素，是用户直接操作的部分。比如 `Input` 输入框、`Select` 选择器、`Checkbox` 复选框等。它的职责是接收 `value` 并通过 `onChange` 事件将用户的输入值传出去。它不关心自己是否有标题、是否有错误提示、布局是怎样的。

- **`decorator` (装饰器):** 相当于**相框、玻璃和衬纸**。它负责包裹 `component`，并为其提供“装饰性”或“容器性”的功能。比如 `FormItem`，它的职责包括：
  - 显示字段的标题（`label`）。
  - 显示必填星号（`*`）。
  - 显示校验错误信息。
  - 控制标题和输入框的布局（上下布局、左右布局）。

**总结一下：`component` 是“内容”，`decorator` 是“容器”或“包装”。**

#### 代码中的体现

让我们来看你代码中的例子：

```tsx
<Field
  name="key"
  initialValue={item.key}
  decorator={[FormItem, { style: { marginBottom: 0, flex: '1 1 0%', overflow: 'hidden' } }]}
  component={[
    Input,
    { placeholder: 'key', size: 'small', onChange: (e: any) => handleKeyChange(e, index) }
  ]}
  validator={keyValidator}
/>
```

这里 `Field` 组件的工作流程是这样的：

1.  **连接模型:** `Field` 根据 `name="key"` 在 Formily 的数据模型中找到对应的字段状态对象。
2.  **状态分发:**
    - 它从模型中取出 `value` 和 `onChange` 等核心交互属性，传递给 `component` 数组中指定的 `Input` 组件。
    - 它从模型中取出 `title` (label)、`required`、`errors` (校验信息) 等外围属性，传递给 `decorator` 数组中指定的 `FormItem` 组件。
3.  **UI 渲染:**
    - `FormItem` (`decorator`) 先渲染出自己的布局结构。
    - 然后，`FormItem` 在自己的结构内部，将 `Input` (`component`) 渲染出来。
    - 如果 `validator` 校验失败，模型会更新错误信息，`Field` 将错误信息传递给 `FormItem`，`FormItem` 就会在 `Input` 下方渲染出红色的错误提示文本。

#### 为什么这样设计？

这种设计极为灵活。想象一下，如果你想把一个普通的输入框换成一个带搜索功能的输入框，你只需要：

- **更换 `component`:** 把 `Input` 换成 `SearchInput`。
- **保持 `decorator` 不变:** `FormItem` 依然负责显示标题和错误信息，你完全不需要改动它。

反之，如果你想改变所有表单项的布局，比如从上下布局改成左右对齐，你只需要修改 `decorator` (`FormItem`) 的样式，而完全不需要关心里面包裹的是 `Input` 还是 `Select`。

这就是 Formily 通过 `component` 和 `decorator` 实现 UI 与逻辑、内容与容器分离的强大之处。
