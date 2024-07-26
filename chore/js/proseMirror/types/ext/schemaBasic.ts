// 这段代码定义了一个文档模型（schema）用于 ProseMirror，这是一个用于构建富文本编辑器的库。这个模型定义了文档中可以使用的节点（nodes）和标记（marks），以及它们的规范（Specs）。

// ### 节点（Nodes）
// - **doc**: 文档的顶层节点，内容为一个或多个块级节点（block+）。
// - **paragraph**: 段落节点，表示为 `<p>` 元素，内容为零个或多个内联节点（inline*）。
// - **blockquote**: 引用块节点，表示为 `<blockquote>` 元素，包裹一个或多个块级节点。
// - **horizontal_rule**: 水平线节点，表示为 `<hr>` 元素。
// - **heading**: 标题节点，有一个 `level` 属性表示标题级别（1到6），内容为零个或多个内联节点。
// - **code_block**: 代码块节点，表示为 `<pre>` 元素内嵌套 `<code>` 元素，内容为文本节点。
// - **text**: 文本节点。
// - **image**: 图片节点，表示为 `<img>` 元素，有 `src`、`alt` 和 `title` 属性。
// - **hard_break**: 硬换行节点，表示为 `<br>` 元素。

// ### 标记（Marks）
// - **link**: 链接标记，表示为 `<a>` 元素，有 `href` 和 `title` 属性。
// - **em**: 强调标记，表示为 `<em>` 元素。
// - **strong**: 加粗标记，表示为 `<strong>` 元素。
// - **code**: 代码字体标记，表示为 `<code>` 元素。

// ### Schema
// 最后，这段代码使用上述定义的节点和标记创建了一个 `Schema` 实例。这个 schema 大致对应于 CommonMark 规范，但不包括列表元素，这些在 `prosemirror-schema-list` 模块中定义。
// 这个 schema 可以被用来构建一个 ProseMirror 编辑器实例，使其能够处理、显示和编辑上述定义的各种文档内容和格式。通过扩展或从这个 schema 的 `spec.nodes` 和 `spec.marks` 属性读取，可以重用或自定义元素。

import { Schema } from '../model'

// :: Object
// [Specs](#model.NodeSpec) for the nodes defined in this schema.
export const nodes = {
  // :: NodeSpec The top level document node.
  doc: {
    content: 'block+'
  },

  // :: NodeSpec A plain paragraph textblock. Represented in the DOM
  // as a `<p>` element.
  paragraph: {
    content: 'inline*',
    group: 'block',
    parseDOM: [{ tag: 'p' }],
    toDOM() {
      return ['p', 0]
    }
  },

  // :: NodeSpec A blockquote (`<blockquote>`) wrapping one or more blocks.
  blockquote: {
    content: 'block+',
    group: 'block',
    defining: true,
    parseDOM: [{ tag: 'blockquote' }],
    toDOM() {
      return ['blockquote', 0]
    }
  },

  // :: NodeSpec A horizontal rule (`<hr>`).
  horizontal_rule: {
    group: 'block',
    parseDOM: [{ tag: 'hr' }],
    toDOM() {
      return ['hr']
    }
  },

  // :: NodeSpec A heading textblock, with a `level` attribute that
  // should hold the number 1 to 6. Parsed and serialized as `<h1>` to
  // `<h6>` elements.
  heading: {
    attrs: { level: { default: 1 } },
    content: 'inline*',
    group: 'block',
    defining: true,
    parseDOM: [
      { tag: 'h1', attrs: { level: 1 } },
      { tag: 'h2', attrs: { level: 2 } },
      { tag: 'h3', attrs: { level: 3 } },
      { tag: 'h4', attrs: { level: 4 } },
      { tag: 'h5', attrs: { level: 5 } },
      { tag: 'h6', attrs: { level: 6 } }
    ],
    toDOM(node) {
      return ['h' + node.attrs.level, 0]
    }
  },

  // :: NodeSpec A code listing. Disallows marks or non-text inline
  // nodes by default. Represented as a `<pre>` element with a
  // `<code>` element inside of it.
  code_block: {
    content: 'text*',
    marks: '',
    group: 'block',
    code: true,
    defining: true,
    parseDOM: [{ tag: 'pre', preserveWhitespace: 'full' }],
    toDOM() {
      return ['pre', ['code', 0]]
    }
  },

  // :: NodeSpec The text node.
  text: {
    group: 'inline'
  },

  // :: NodeSpec An inline image (`<img>`) node. Supports `src`,
  // `alt`, and `href` attributes. The latter two default to the empty
  // string.
  image: {
    inline: true,
    attrs: {
      src: {},
      alt: { default: null },
      title: { default: null }
    },
    group: 'inline',
    draggable: true,
    parseDOM: [
      {
        tag: 'img[src]',
        getAttrs(dom) {
          return {
            src: dom.getAttribute('src'),
            title: dom.getAttribute('title'),
            alt: dom.getAttribute('alt')
          }
        }
      }
    ],
    toDOM(node) {
      return ['img', node.attrs]
    }
  },

  // :: NodeSpec A hard line break, represented in the DOM as `<br>`.
  hard_break: {
    inline: true,
    group: 'inline',
    selectable: false,
    parseDOM: [{ tag: 'br' }],
    toDOM() {
      return ['br']
    }
  }
}

// :: Object [Specs](#model.MarkSpec) for the marks in the schema.
export const marks = {
  // :: MarkSpec A link. Has `href` and `title` attributes. `title`
  // defaults to the empty string. Rendered and parsed as an `<a>`
  // element.
  link: {
    attrs: {
      href: {},
      title: { default: null }
    },
    inclusive: false,
    parseDOM: [
      {
        tag: 'a[href]',
        getAttrs(dom) {
          return { href: dom.getAttribute('href'), title: dom.getAttribute('title') }
        }
      }
    ],
    toDOM(node) {
      return ['a', node.attrs]
    }
  },

  // :: MarkSpec An emphasis mark. Rendered as an `<em>` element.
  // Has parse rules that also match `<i>` and `font-style: italic`.
  em: {
    parseDOM: [{ tag: 'i' }, { tag: 'em' }, { style: 'font-style=italic' }],
    toDOM() {
      return ['em']
    }
  },

  // :: MarkSpec A strong mark. Rendered as `<strong>`, parse rules
  // also match `<b>` and `font-weight: bold`.
  strong: {
    parseDOM: [
      { tag: 'strong' },
      // This works around a Google Docs misbehavior where
      // pasted content will be inexplicably wrapped in `<b>`
      // tags with a font-weight normal.
      { tag: 'b', getAttrs: node => node.style.fontWeight != 'normal' && null },
      { style: 'font-weight', getAttrs: value => /^(bold(er)?|[5-9]\d{2,})$/.test(value) && null }
    ],
    toDOM() {
      return ['strong']
    }
  },

  // :: MarkSpec Code font mark. Represented as a `<code>` element.
  code: {
    parseDOM: [{ tag: 'code' }],
    toDOM() {
      return ['code']
    }
  }
}

// :: Schema
// This schema rougly corresponds to the document schema used by
// [CommonMark](http://commonmark.org/), minus the list elements,
// which are defined in the [`prosemirror-schema-list`](#schema-list)
// module.
//
// To reuse elements from this schema, extend or read from its
// `spec.nodes` and `spec.marks` [properties](#model.Schema.spec).
export const schema = new Schema({ nodes, marks })
