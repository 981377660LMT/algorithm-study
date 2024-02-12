/* eslint-disable @typescript-eslint/ban-types */
// 虚拟dom

interface VDom {
  type: keyof HTMLElementTagNameMap | (string & {})
  props: {
    children?: Children[] | Children
    [attr: string]: any
  }
}

type Children = VDom | string

function virtualize(element: HTMLElement | Text): VDom | string {
  if (element instanceof Text) {
    return element.data
  }

  const res: VDom = { type: element.tagName.toLowerCase(), props: {} }

  // props里的非children属性
  for (const attr of Array.from(element.attributes)) {
    const name = attr.name === 'class' ? 'className' : attr.name
    res.props[name] = attr.value
  }

  // props里children属性
  const children: (string | VDom)[] = []
  for (const node of Array.from(element.childNodes)) {
    // text node  (1:元素节点，2：属性节点，废弃，3：文本节点)
    if (node.nodeType === 3) children.push(node.textContent!)
    else children.push(virtualize(node as HTMLElement))
  }

  res.props.children = children.length === 1 ? children[0] : children

  return res
}

function render(json: VDom | string): HTMLElement | Text {
  // 1.文本节点
  if (typeof json === 'string') return document.createTextNode(json)

  const {
    type,
    props: { children, ...attrs }
  } = json

  // 2.根节点的非children属性
  const element = document.createElement(type)
  for (const [attr, value] of Object.entries(attrs)) {
    element.setAttribute(attr, value)
  }

  // 3. 根节点的children
  const childrenArray = Array.isArray(children) ? children : [children]
  for (let child of childrenArray) {
    if (child === void 0) continue
    element.append(render(child))
  }

  return element
}

const json: VDom = {
  type: 'div',
  props: {
    children: [
      {
        type: 'h1',
        props: {
          children: ' this is '
        }
      },
      {
        type: 'p',
        props: {
          className: 'paragraph',
          children: [
            ' a ',
            {
              type: 'button',
              props: {
                children: ' button '
              }
            },
            ' from',
            {
              type: 'a',
              props: {
                href: 'https://bfe.dev',
                children: [
                  {
                    type: 'b',
                    props: {
                      children: 'BFE'
                    }
                  },
                  '.dev'
                ]
              }
            }
          ]
        }
      }
    ]
  }
}

console.log(render(json))
console.log(virtualize(render(json)))

export {}
// @ts-ignore
// console.log(Array.from(undefined))
