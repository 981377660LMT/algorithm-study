interface VirtualDom {
  type: keyof HTMLElementTagNameMap | string
  props: {
    children: Children[]
    [attr: string]: any
  }
}

type Children = VirtualDom | string

function createElement(
  type: keyof HTMLElementTagNameMap,
  props: Record<string, string>,
  ...children: string[]
): VirtualDom
function createElement(
  type: keyof HTMLElementTagNameMap,
  props: Record<string, string>,
  ...children: Children[]
): VirtualDom
/**
 * @param { string } type - valid HTML tag name
 * @param { object } [props] - properties.
 * @param { ...MyNode} [children] - elements as rest arguments
 * @return { MyElement }
 * @description
 * ref、key 等基本功能之外的部分不在本题目考虑范围内
 * Re-render的问题不需要考虑
 */
function createElement(
  type: keyof HTMLElementTagNameMap,
  props: Record<string, string>,
  ...children: Children[]
): VirtualDom {
  // your code here
  return {
    type,
    props: {
      ...props,
      children,
    },
  }
}

/**
 * @param { VirtualDom }
 * @returns { HTMLElement }
 */
function render(json: VirtualDom | string): HTMLElement | Text {
  if (typeof json === 'string') return document.createTextNode(json)

  const {
    type,
    props: { children, ...attrs },
  } = json

  const element = document.createElement(type)
  for (let [key, value] of Object.entries(attrs)) {
    // 当在HTML 文档中的HTML 元素上调用 setAttribute() 方法时，该方法会将其属性名称（attribute name）参数小写化
    // element.setAttribute(key, value)   // 所有名称会转小写
    // @ts-ignore
    element[key] = value
  }

  const childrenArray = Array.isArray(children) ? children : [children]
  for (const child of childrenArray) {
    element.appendChild(render(child))
  }

  return element
}

if (require.main === module) {
  const h = createElement
  render(
    h(
      'div',
      {},
      h('h1', {}, ' this is '),
      h(
        'p',
        { className: 'paragraph' },
        ' a ',
        h('button', {}, ' button '),
        ' from ',
        h('a', { href: 'https://bfe.dev' }, h('b', {}, 'BFE'), '.dev')
      )
    )
  )
}

export {}
