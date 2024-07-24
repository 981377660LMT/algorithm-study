import * as prosemirror_state from '../state'
import { Command, Plugin } from '../state'
import { MenuItem, Dropdown, MenuElement } from './menu'
import { Schema } from '../model'

declare type MenuItemResult = {
  /**
    A menu item to toggle the [strong mark](https://prosemirror.net/docs/ref/#schema-basic.StrongMark).
    */
  toggleStrong?: MenuItem
  /**
    A menu item to toggle the [emphasis mark](https://prosemirror.net/docs/ref/#schema-basic.EmMark).
    */
  toggleEm?: MenuItem
  /**
    A menu item tose toggle the [code font mark](https://prosemirror.net/docs/ref/#schema-basic.CodeMark).
    */
  toggleCode?: MenuItem
  /**
    A menu item to toggle the [link mark](https://prosemirror.net/docs/ref/#schema-basic.LinkMark).
    */
  toggleLink?: MenuItem
  /**
    A menu item to insert an [image](https://prosemirror.net/docs/ref/#schema-basic.Image).
    */
  insertImage?: MenuItem
  /**
    A menu item to wrap the selection in a [bullet list](https://prosemirror.net/docs/ref/#schema-list.BulletList).
    */
  wrapBulletList?: MenuItem
  /**
    A menu item to wrap the selection in an [ordered list](https://prosemirror.net/docs/ref/#schema-list.OrderedList).
    */
  wrapOrderedList?: MenuItem
  /**
    A menu item to wrap the selection in a [block quote](https://prosemirror.net/docs/ref/#schema-basic.BlockQuote).
    */
  wrapBlockQuote?: MenuItem
  /**
    A menu item to set the current textblock to be a normal
    [paragraph](https://prosemirror.net/docs/ref/#schema-basic.Paragraph).
    */
  makeParagraph?: MenuItem
  /**
    A menu item to set the current textblock to be a
    [code block](https://prosemirror.net/docs/ref/#schema-basic.CodeBlock).
    */
  makeCodeBlock?: MenuItem
  /**
    Menu items to set the current textblock to be a
    [heading](https://prosemirror.net/docs/ref/#schema-basic.Heading) of level _N_.
    */
  makeHead1?: MenuItem
  makeHead2?: MenuItem
  makeHead3?: MenuItem
  makeHead4?: MenuItem
  makeHead5?: MenuItem
  makeHead6?: MenuItem
  /**
    A menu item to insert a horizontal rule.
    */
  insertHorizontalRule?: MenuItem
  /**
    A dropdown containing the `insertImage` and
    `insertHorizontalRule` items.
    */
  insertMenu: Dropdown
  /**
    A dropdown containing the items for making the current
    textblock a paragraph, code block, or heading.
    */
  typeMenu: Dropdown
  /**
    Array of block-related menu items.
    */
  blockMenu: MenuElement[][]
  /**
    Inline-markup related menu items.
    */
  inlineMenu: MenuElement[][]
  /**
    An array of arrays of menu elements for use as the full menu
    for, for example the [menu
    bar](https://github.com/prosemirror/prosemirror-menu#user-content-menubar).
    */
  fullMenu: MenuElement[][]
}
/**
Given a schema, look for default mark and node types in it and
return an object with relevant menu items relating to those marks.
*/
declare function buildMenuItems(schema: Schema): MenuItemResult

/**
Inspect the given schema looking for marks and nodes from the
basic schema, and if found, add key bindings related to them.
This will add:

* **Mod-b** for toggling [strong](https://prosemirror.net/docs/ref/#schema-basic.StrongMark)
* **Mod-i** for toggling [emphasis](https://prosemirror.net/docs/ref/#schema-basic.EmMark)
* **Mod-`** for toggling [code font](https://prosemirror.net/docs/ref/#schema-basic.CodeMark)
* **Ctrl-Shift-0** for making the current textblock a paragraph
* **Ctrl-Shift-1** to **Ctrl-Shift-Digit6** for making the current
  textblock a heading of the corresponding level
* **Ctrl-Shift-Backslash** to make the current textblock a code block
* **Ctrl-Shift-8** to wrap the selection in an ordered list
* **Ctrl-Shift-9** to wrap the selection in a bullet list
* **Ctrl->** to wrap the selection in a block quote
* **Enter** to split a non-empty textblock in a list item while at
  the same time splitting the list item
* **Mod-Enter** to insert a hard break
* **Mod-_** to insert a horizontal rule
* **Backspace** to undo an input rule
* **Alt-ArrowUp** to `joinUp`
* **Alt-ArrowDown** to `joinDown`
* **Mod-BracketLeft** to `lift`
* **Escape** to `selectParentNode`

You can suppress or map these bindings by passing a `mapKeys`
argument, which maps key names (say `"Mod-B"` to either `false`, to
remove the binding, or a new key name string.
*/
declare function buildKeymap(
  schema: Schema,
  mapKeys?: {
    [key: string]: false | string
  }
): {
  [key: string]: Command
}

/**
A set of input rules for creating the basic block quotes, lists,
code blocks, and heading.
*/
declare function buildInputRules(schema: Schema): prosemirror_state.Plugin<{
  transform: prosemirror_state.Transaction
  from: number
  to: number
  text: string
} | null>

/**
Create an array of plugins pre-configured for the given schema.
The resulting array will include the following plugins:

 * Input rules for smart quotes and creating the block types in the
   schema using markdown conventions (say `"> "` to create a
   blockquote)

 * A keymap that defines keys to create and manipulate the nodes in the
   schema

 * A keymap binding the default keys provided by the
   prosemirror-commands module

 * The undo history plugin

 * The drop cursor plugin

 * The gap cursor plugin

 * A custom plugin that adds a `menuContent` prop for the
   prosemirror-menu wrapper, and a CSS class that enables the
   additional styling defined in `style/style.css` in this package

Probably only useful for quickly setting up a passable
editor—you'll need more control over your settings in most
real-world situations.
*/
declare function exampleSetup(options: {
  /**
    The schema to generate key bindings and menu items for.
    */
  schema: Schema
  /**
    Can be used to [adjust](https://prosemirror.net/docs/ref/#example-setup.buildKeymap) the key bindings created.
    */
  mapKeys?: {
    [key: string]: string | false
  }
  /**
    Set to false to disable the menu bar.
    */
  menuBar?: boolean
  /**
    Set to false to disable the history plugin.
    */
  history?: boolean
  /**
    Set to false to make the menu bar non-floating.
    */
  floatingMenu?: boolean
  /**
    Can be used to override the menu content.
    */
  menuContent?: MenuItem[][]
}): (
  | Plugin<any>
  | Plugin<{
      transform: prosemirror_state.Transaction
      from: number
      to: number
      text: string
    } | null>
)[]

export { buildInputRules, buildKeymap, buildMenuItems, exampleSetup }
