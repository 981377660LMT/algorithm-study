This is the guide to the CodeMirror editor system. It provides a prose description of the system's functionality. For the item-by-item documentation of its interface, see the reference manual.

Architecture Overview
Because CodeMirror is structured quite a bit differently than your classical JavaScript library (including its own previous versions), it is recommended to read at least this section before jumping in, so that you don't waste your time with mismatched expectations.

Modularity
CodeMirror is set up as a collection of separate modules that, together, provide a full-featured text and code editor. On the bright side, this means that you can pick and choose which features you need, and even replace core functionality with a custom implementation if you need to. On the less bright side, this means that setting up an editor requires you to put together a bunch of pieces.

The putting-together part isn't hard, but you will have to install and import the pieces you need. The core packages, without which it'd be hard to set up an editor at all, are:

@codemirror/state, which defines data structures that represent the editor state and changes to that state.

@codemirror/view, a display component that knows how to show the editor state to the user, and translates basic editing actions into state updates.

@codemirror/commands defines a lot of editing commands and some key bindings for them.

This is what a minimal viable editor might look like:

import {EditorState} from "@codemirror/state"
import {EditorView, keymap} from "@codemirror/view"
import {defaultKeymap} from "@codemirror/commands"

let startState = EditorState.create({
doc: "Hello World",
extensions: [keymap.of(defaultKeymap)]
})

let view = new EditorView({
state: startState,
parent: document.body
})
Many things that you'd expect in an editor, such as the line number gutter or undo history, are implemented as extensions to the generic core, and need to be explicitly added to a configuration to be enabled. To make it easy to get started, the codemirror package pulls in most of the things you need for a baseline editor (except a language package).

import {EditorView, basicSetup} from "codemirror"
import {javascript} from "@codemirror/lang-javascript"

let view = new EditorView({
extensions: [basicSetup, javascript()],
parent: document.body
})
The packages are distributed as ES6 modules. This means that it is not currently practical to run the library without some kind of bundler (which packages up a modular program into a single big JavaScript file) or module loader. If you are new to bundling, I recommend looking into rollup or Vite.

Functional Core, Imperative Shell
An attitude that guides the architecture of CodeMirror is that functional (pure) code, which creates new values instead of having side effects, is much easier to work with than imperative code. But the browser DOM is obviously very imperative-minded, as are many of the systems that CodeMirror integrate with.

To resolve this contradiction, the library's state representation is strictly functional—the document and state data structures are immutable, and operations on them are pure functions, whereas the view component and command interface wrap these in an imperative interface.

This means that an old state value stays intact even when the editor moves on to a new state. Having both the old and the new state available is often very useful when dealing with state changes. It also means that directly changing a state value, or writing extensions like additional state fields in an imperative way will not do what you'd hope (and probably just break things).

The TypeScript interface tries to be very clear about this by marking arrays and object properties as readonly. When using plain old JavaScript it can be harder to remember this. But as a general rule, unless explicitly described in the docs, reassignment of properties in objects created by the library is just not supported.

let state = EditorState.create({doc: "123"})
// BAD WRONG NO GOOD CODE:
state.doc = Text.of("abc") // <- DON'T DO THIS
State and Updates
The library handles updates in a way inspired by approaches like Redux or Elm. With a few exceptions (like composition and drag-drop handling), the state of the view is entirely determined by the EditorState value in its state property.

Changes to that state happen in functional code, by creating a transaction that describes the changes to document, selection, or other state fields. Such a transaction can then be dispatched, which tells the view to update its state, at which point it'll synchronize its DOM representation with the new state.

// (Assume view is an EditorView instance holding the document "123".)
let transaction = view.state.update({changes: {from: 0, insert: "0"}})
console.log(transaction.state.doc.toString()) // "0123"
// At this point the view still shows the old state.
view.dispatch(transaction)
// And now it shows the new state.
The view listens for events, such as text input, key presses, or mouse interaction. It translates those into transactions that make the appropriate changes to the editor state, which it dispatches, computing a new state and synchronizing what it is displaying to that new state.

It is of course also possible for transactions to come from other sources. But to affect the editor, they must be dispatched to the view component.

Extension
The core library is rather minimal and generic, and a lot of functionality is implemented in system extensions. Extensions can do all kinds of things, from merely configuring some option, to defining new fields in the state object, to styling the editor, to injecting custom imperative components into the view. The system takes care to allow such extensions to compose without unexpected conflicts.

The set of active extensions is kept in the editor state (and can be changed by a transaction). Extensions are provided as values (usually imported from some package), or arrays of such values. They can be arbitrarily nested (an array containing more arrays is also a valid extension), and are deduplicated during the configuration process. Thus, it is okay for extensions to pull in other extensions—if the same one gets included multiple times, it'll only take effect once.

When relevant, the precedence of extensions is determined first by explicitly set precedence category, and within that, by the position the extension has in the (flattened) collection of extensions passed to the state.

import {keymap} from "@codemirror/view"
import {EditorState, Prec} from "@codemirror/state"

function dummyKeymap(tag) {
return keymap.of([{
key: "Ctrl-Space",
run() { console.log(tag); return true }
}])
}

let state = EditorState.create({extensions: [
dummyKeymap("A"),
dummyKeymap("B"),
Prec.high(dummyKeymap("C"))
]})
In a view using the state from that code, pressing Ctrl-Space will log "C", because, despite being last in the order of extensions, that keymap is tagged with a higher precedence. If that weren't the case, keymap "A" would be the first one to get a chance to handle the key combination, since it occurs before the others.

A later section of the guide goes into more detail on the kinds of extensions that exist, and how to use them. See the configuration example for more example code related to configuration and reconfiguration.

Document Offsets
CodeMirror uses plain numbers to address positions in the document. These represent character counts—more precisely, they count UTF16 code units (so astral characters count as two units). Line breaks always count as a single unit (even when you configure a line separator that is longer than that).

These offsets are used to track the selection, position changes, decorate content, and so on.

It is sometimes necessary to figure out where a position in a start document ends up in a changed document. For this purpose, the library provides a position mapping feature, which, given a transaction (or just a change set) and a start position, can give you the corresponding new position.

import {EditorState} from "@codemirror/state"

let state = EditorState.create({doc: "1234"})
// Delete "23" and insert at "0" at the start.
let tr = state.update({changes: [{from: 1, to: 3}, {from: 0, insert: "0"}]})
// The position at the end of the old document is at 3 in the new document.
console.log(tr.changes.mapPos(4))
The document data structure also indexes by lines, so it is not expensive to look things up by (1-based) line number.

import {Text} from "@codemirror/state"

let doc = Text.of(["line 1", "line 2", "line 3"])
// Get information about line 2
console.log(doc.line(2)) // {from: 7, to: 13, ...}
// Get the line around position 15
console.log(doc.lineAt(15)) // {from: 14, to: 20, ...}
Data Model
CodeMirror, being a text editor, treats the document as a flat string. It stores this string split by line in a tree-shaped data structure to allow cheap updates anywhere in the document (and efficient indexing by line number).

Document Changes
Document changes are themselves values, describing precisely which ranges of the old document are being replaced by which bits of new text. This allows extensions to track precisely what happens to the document, allowing things like an undo history and collaborative editing to be implemented outside the library core.

When creating a change set, all change positions are specified in terms of the original document—they conceptually all happen at once. (If you really need to combine lists of changes where later changes refer to the document created by earlier ones, you can use the change set compose method.)

Selection
Alongside the document, an editor state stores a current selection. Selections may consist of multiple ranges, each of which can be a cursor (empty) or cover a range between its anchor and head. Overlapping ranges are automatically merged, and ranges are sorted, so that a selection's ranges property always holds a sorted, non-overlapping array of ranges.

import {EditorState, EditorSelection} from "@codemirror/state"

let state = EditorState.create({
doc: "hello",
selection: EditorSelection.create([
EditorSelection.range(0, 4),
EditorSelection.cursor(5)
]),
extensions: EditorState.allowMultipleSelections.of(true)
})
console.log(state.selection.ranges.length) // 2

let tr = state.update(state.replaceSelection("!"))
console.log(tr.state.doc.toString()) // "!o!"
One of these ranges is marked as the main one. This is the one that the browser's DOM selection will reflect. The others are drawn and handled entirely by the library.

By default a state will only accept selections with a single range. To get support for multiple selections, you have to include an extension like drawSelection that is able to draw them, and set an option to enable them.

State objects have a convenience method, changeByRange for applying an operation to every selection range separately (which can be a bit awkward to do manually).

import {EditorState, EditorSelection} from "@codemirror/state"

let state = EditorState.create({doc: "abcd", selection: {anchor: 1, head: 3}})
// Upcase the selection
let tr = state.update(state.changeByRange(range => {
let upper = state.sliceDoc(range.from, range.to).toUpperCase()
return {
changes: {from: range.from, to: range.to, insert: upper},
range: EditorSelection.range(range.from, range.from + upper.length)
}
}))
console.log(tr.state.doc.toString()) // "aBCd"
There's also replaceSelection, which creates a transaction that replaces all selection ranges with some piece of text.

Configuration
Each editor state also has a (private) reference to its configuration, which is determined by the extensions that are active for that state. During regular transactions, the configuration stays the same. But it is possible to reconfigure the state using compartments or effects that add to or replace the current configuration.

The direct effects of a state's configuration are the fields it stores and the values associated with facets for that state.

Facets
A facet is an extension point. Different extension values can provide values for the facet. And anyone with access to the state and the facet can read its output value. Depending on the facet, that may just be an array of provided values, or some value computed from those.

The idea behind facets is that most types of extension points allow multiple inputs, but want to compute some coherent combined value from those. How that combining works may differ.

For something like tab size, you need a single output value. So that facet takes the value with the highest precedence and uses that.

When providing event handlers, you want the handlers as an array, sorted by precedence, so that you can try them one at a time until one of them handles the event.

Another common pattern is to compute the logical or of the input values (as in allowMultipleSelections) or reduce them in some other way (say, taking the maximum of the requested undo history depths).

import {EditorState} from "@codemirror/state"

let state = EditorState.create({
extensions: [
EditorState.tabSize.of(16),
EditorState.changeFilter.of(() => true)
]
})
console.log(state.facet(EditorState.tabSize)) // 16
console.log(state.facet(EditorState.changeFilter)) // [() => true]
Facets are defined with Facet.define, which returns a facet value. Such a value can be exported, to allow other code to provide and read it, or it can be kept module-private, in which case only that module can access it. We'll come back to that in the section on writing extensions.

In a given configuration, most facets tend to be static, provided only directly as part of the configuration. But it is also possible to have facet values computed from other aspects of the state.

let info = Facet.define<string>()
let state = EditorState.create({
doc: "abc\ndef",
extensions: [
info.of("hello"),
info.compute(["doc"], state => `lines: ${state.doc.lines}`)
]
})
console.log(state.facet(info))
// ["hello", "lines: 2"]
Such values are automatically recomputed when their declared inputs change.

Facet values are only recomputed when necessary, so you can use an object or array identity test to cheaply check whether a facet changed.

Transactions
Transactions, created with the state's update method, combine a number of effects (all optional):

They can apply document changes.

They can explicitly move the selection. Note that when there are document changes, but no explicit new selection, the selection will be implicitly mapped through these changes.

They can set a flag that instructs the view to scroll the (main) selection head into view.

They can have any number of annotations, which store additional metadata that describes the (entire) transaction. For example, the userEvent annotation can be used to recognize transactions generated for certain common operations like typing or pasting.

They can have effects, which are self-contained additional effects, typically on some extension's state (such as folding code or starting an autocompletion).

They can influence the state's configuration, either by providing a completely new set of extensions, or by replacing specific parts of the configuration.

Transactions are described with specs, which are usually written as object literals, though some methods (such as changeByRange) also return them. Such specs can be passed directly to EditorView.dispatch to create an immediately dispatch a transaction, or to EditorState.update to only create it.

When multiple specs are passed to these methods, they are combined into a single transaction. This can be useful to add some additional fields to a spec created by some helper function.

Changes are described by {from, to, insert} objects (where to and insert are optional), or nested arrays of these. You can also pass in a ChangeSet object, which is the form in which the changes will be represented in the transaction object. The positions given for changes refer to the transaction's start document, even when there are multiple changes.

Positions used in the new selection or in state effects refer to the new document, after the changes.

To completely reset a state—for example to load a new document—it is recommended to create a new state instead of a transaction. That will make sure no unwanted state (such as undo history events) sticks around.

The View
The view tries to be as transparent a layer around the state as possible. Unfortunately, there are some aspects of working with an editor that can't be handled purely with the data in the state.

When dealing with screen coordinates (to figure out where the user clicked, or to find the coordinates of a given position), you need access to the layout, and thus the browser DOM.

The editor takes the text direction from the surrounding document (or its own CSS style, if overridden).

Cursor motion can depend on layout and text direction. Thus, the view provides a number of helper methods for computing different types of motion.

Some state, such as focus and scroll position, isn't stored in the functional state, but left in the DOM.

The library does not expect user code to manipulate the DOM structure it manages. When you do try that, you'll probably just see the library revert your changes right away. See the section on decorations for the proper way to affect the way the content is displayed.

Viewport
One thing to be aware of is that CodeMirror doesn't render the entire document, when that document is big. To keep the editor responsive and resource use low, it will, when updating, detect which part of the content is currently visible (not scrolled out of view), and only render that plus a margin around it. This is called the viewport.

Querying coordinates for positions outside of the current viewport will not work (since they are not rendered, and thus have no layout). The view does track height information (initially estimated, measured accurately when content is drawn) for the entire document, even the parts outside of the viewport.

Long lines (when not wrapped) or chunks of folded code can still make the viewport rather huge. The editor also provides a list of visible ranges, which won't include such invisible content. This can be useful when, for example, highlighting code, where you don't want to do work for text that isn't visible anyway.

Update Cycle
CodeMirror's view makes a serious effort to minimize the amount of DOM reflows it causes. Dispatching a transaction will generally only cause the editor to write to the DOM, without reading layout information. The reading (to check whether the viewport is still valid, whether the cursor needs to be scrolled into view, and so on) is done in a separate measure phase, scheduled using requestAnimationFrame. This phase will, if necessary, follow up with another write phase.

You can schedule your own measure code using the requestMeasure method.

To avoid weird reentrancy situations, the view will raise an error when a new update is initiated while another update is in the process of being (synchronously) applied. Multiple updates applied while a measure phase is still pending are not a problem—those will just cause their measure phases to be combined.

When you are done with a view instance, you must call its destroy method to dispose of it, releasing any resources (global event handlers and mutation observers) that it allocated.

DOM Structure
The editor's DOM structure looks something like this:

<div class="cm-editor [theme scope classes]">
  <div class="cm-scroller">
    <div class="cm-content" contenteditable="true">
      <div class="cm-line">Content goes here</div>
      <div class="cm-line">...</div>
    </div>
  </div>
</div>
The outer (wrap) element is a vertical flexbox. Things like panels and tooltips can be put here by extensions.

Inside of that is the scroller element. If the editor has its own scrollbar, this one should be styled with overflow: auto. But it doesn't have to—the editor also supports growing to accomodate its content, or growing up to a certain max-height and then scrolling.

The scroller is a horizontal flexbox element. When there are gutters, they are added to its start.

Inside that is the content element, which is editable. This has a DOM mutation observer registered on it, and any changes made in there will result in the editor parsing them as document changes and redrawing the affected nodes. This container holds a line element for each line in the viewport, which in turn hold the document text (possibly decorated with styles or widgets).

Styles and Themes
To manage editor-related styles, CodeMirror uses a system to inject styles from JavaScript. Styles can be registered with a facet, which will cause the view to make sure they are available.

Many elements in the editor are assigned classes prefixed with cm-. These can be styled directly in your local CSS. But they can also be targeted by themes. A theme is an extension created with EditorView.theme. It gets its own unique (generated) CSS class (which will be added to the editor when the theme extension is active) and defines styles scoped by that class.

A theme declaration defines any number of CSS rules using style-mod notation. This code creates a crude theme that makes the default text color in the editor orange:

import {EditorView} from "@codemirror/view"

let view = new EditorView({
extensions: EditorView.theme({
".cm-content": {color: "darkorange"},
"&.cm-focused .cm-content": {color: "orange"}
})
})
To allow the automatic class prefixing be done correctly, rules where the first element targets the editor's wrapper element (which is where the theme's unique class will be added), such as the .cm-focused rule in the example, must use an & character to indicate the position of the wrapper element.

Extensions can define base themes to provide default styles for the elements they create. Base themes can use &light (default) and &dark (enabled when there's a dark theme active) placeholders, so that even when they aren't overridden by a theme, they don't look too out of place.

import {EditorView} from "@codemirror/view"

// This again produces an extension value
let myBaseTheme = EditorView.baseTheme({
"&dark .cm-mySelector": { background: "dimgrey" },
"&light .cm-mySelector": { background: "ghostwhite" }
})
When defining editor styling in regular CSS, you must account for the extra prefixed class selector that is added to injected style rules, or your styles will always have a lower precedence. The recommended approach is to include .cm-editor in your rules, which will go in the same place, and produce the same precedence, as the injected styles.

.cm-editor .cm-content { color: purple; }
Commands
Commands are functions with a specific signature: (view: EditorView) => boolean. Their main use is key bindings, but they could also be used for things like menu items or command palettes. A command function represents a user action. It takes a view and returns a boolean, false to indicate it doesn't apply in the current situation, and true to indicate that it successfully executed. The effect of the command is produced imperatively, usually by dispatching a transaction.

When multiple commands are bound to a given key, they are tried one at a time, in order of precedence, until one of them returns true.

Commands that only act on the state, not the entire view, can use the StateCommand type instead, which is a subtype of Command that just expects its argument to have state and dispatch properties. This is mostly useful for being able to test such commands without creating a view.

The @codemirror/commands package exports a lot of different editing commands, along with some keymaps. Keymaps are arrays of KeyBinding objects, that are given to the keymap facet to enable them in an editor.

let myKeyExtension = keymap.of([
{
key: "Alt-c",
run: view => {
view.dispatch(view.state.replaceSelection("?"))
return true
}
}
])
Extending CodeMirror
There are a number of different ways to extend CodeMirror, and picking the proper way for a given use case isn't always obvious. This section goes over the various concepts you'll need to be familiar with to write editor extensions.

State Fields
Extensions often need to store additional information in the state. The undo history needs to store undoable changes, the code folding extension needs to track what has been folded, and so on.

For this purpose, extensions can define additional state fields. State fields, living inside the purely functional state data structure, must store immutable values.

State fields are kept in sync with the rest of the state using something like a reducer. Every time the state updates, a function is called with the field's current value and the transaction, which should return the field's new value.

import {EditorState, StateField} from "@codemirror/state"

let countDocChanges = StateField.define({
create() { return 0 },
update(value, tr) { return tr.docChanged ? value + 1 : value }
})

let state = EditorState.create({extensions: countDocChanges})
state = state.update({changes: {from: 0, insert: "."}}).state
console.log(state.field(countDocChanges)) // 1
You will often want to use annotations or effects to communicate what is happening to your state field.

import {StateField, StateEffect} from "@codemirror/state"

let setFullScreenMode = StateEffect.define<boolean>()

let fullScreenMode = StateField.define({
create() { return false },
update(value, tr) {
for (let e of tr.effects)
if (e.is(setFullScreenMode)) value = e.value
return value
}
})
It can be tempting to try to avoid taking the step of putting state in an actual state field—there's a bit of verbosity involved in declaring one, and firing off a whole transaction for every state change may feel a bit heavyweight. But in almost all cases, it is a really good idea to tie your state into the editor-wide state update cycle, because it makes it a lot easier to keep it in sync with the rest of the editor state.

Affecting the View
View plugins provide a way for extensions to run an imperative component inside the view. This is useful for things like event handlers, adding and managing DOM elements, and doing things that depend on the current viewport.

This simple plugin displays the document size in the editor's corner.

import {ViewPlugin} from "@codemirror/view"

const docSizePlugin = ViewPlugin.fromClass(class {
constructor(view) {
this.dom = view.dom.appendChild(document.createElement("div"))
this.dom.style.cssText =
"position: absolute; inset-block-start: 2px; inset-inline-end: 5px"
this.dom.textContent = view.state.doc.length
}

update(update) {
if (update.docChanged)
this.dom.textContent = update.state.doc.length
}

destroy() { this.dom.remove() }
})
View plugins should generally not hold (non-derived) state. They work best as shallow views over the data kept in the editor state.

When the state is reconfigured, view plugins that aren't part of the new configuration will be destroyed (which is why, if they made changes to the editor, they should define a destroy method that undoes those changes).

When a view plugin crashes, it is automatically disabled to avoid taking down the entire view.

Decorating the Document
When not told otherwise, CodeMirror will draw the document as plain text. Decorations are the mechanism through which extensions can influence what the document looks like. They come in four types:

Mark decorations add style or DOM attributes to the text in a given range.

Widget decorations insert a DOM element at a given position in the document.

Replace decorations hide part of the document or replace it with a given DOM node.

Line decorations can add attributes to a line's wrapping element.

Decorations are provided through a facet. Every time the view is updated, the content of this facet is used to style the visible content.

Decorations are kept in sets, which are again immutable data structures. Such sets can be mapped across changes (adjusting the positions of their content to compensate for the change) or rebuilt on updates, depending on the use case.

There are two ways in which decorations can be provided: directly, by putting the range set value in the facet (often by deriving it from a field), or indirectly, by providing a function from a view to a range set.

Only directly provided decoration sets may influence the vertical block structure of the editor, but only indirectly provided ones can read the editor's viewport (which can be useful if you want to, for example, decorate only the visible content). The reason for this restriction is that the viewport is computed from the block structure, so that must be known before the viewport can be read.

There is a decoration example which implements some common use cases.

Extension Architecture
To create a given piece of editor functionality you often need to combine different kinds of extension: a state field to keep state, a base theme to provide styling, a view plugin to manage in- and output, some commands, maybe a facet for configuration.

A common pattern is to export a function that returns the extension values necessary for your feature to work. Making this a function, even if it doesn't (yet) take parameters is a good idea—it makes it possible to add configuration options later, without breaking backwards compatiblity.

Since extensions can pull in other extensions, it can be useful to consider what happens when your extension is included multiple times. For some kinds of extensions, for example keymaps, it is appropriate to just do the thing it's doing multiple times. But often that would be wasteful or even break something.

It is usually possible to make multiple uses of an extension just do the right thing by relying on the deduplication of identical extension values—if you make sure you create your static extension values (themes, state fields, view plugins, etc) only once, and always return the same instance from your extension constructor function, you'll only get one copy of them in the editor.

But when your extension allows configuration, your other logic will probably need access to that. And what do you do when the different instances of the extension got different configurations?

Sometimes, that's just an error. But often, it is possible to define a strategy for reconciling them. Facets work well for this. You can put the configuration in a module-private facet, and have its combining function either reconcile configurations or thow an error when this is impossible. Then code that needs access to the current configuration can read that facet.

See the zebra stripes example for an illustration of this approach.

Code of Conduct
Report an Issue
