export {}

// 像素画编辑器
//
// 该应用程序的界面在顶部显示一个大的 <canvas> 元素，下方是一些表单字段。
// 用户通过从一个 <select> 字段中选择工具，然后在画布上点击、触摸或拖动来绘制图像。
// 有用于绘制单个像素或矩形的工具，用于填充区域的工具，以及用于从图像中拾取颜色的工具。
// 我们将把编辑器界面构建为多个组件，这些组件负责 DOM 的一部分，并且可能在内部包含其他组件。
// !应用程序的状态由当前图像、选定的工具和选定的颜色组成。
// 应用程序状态将是一个带有 picture、tool 和 color 属性的对象。
// 我们将进行设置，使状态存在于单个值中，并且界面组件始终根据当前状态来调整其外观。
//
// !我们将对数据流保持严格。存在一个状态，界面根据该状态进行绘制。界面组件可能会通过更新状态来响应用户操作，此时组件有机会与这个新状态同步。
// !状态的更新以对象的形式表示，我们将它们称为操作。组件可以创建此类操作并分发它们——将它们传递给一个中心状态管理函数。该函数计算下一个状态，之后界面组件更新自己以适应这个新状态。
// DOM 事件改变状态的唯一方式是将操作分发到状态。
// !状态更改应该通过一个明确定义的通道进行，而不是到处发生.
// #region Type Definitions

type Color = string

interface Position {
  x: number
  y: number
}

interface Pixel {
  x: number
  y: number
  color: Color
}

// The main state of the application
interface State {
  picture: Picture
  tool: string
  color: Color
  done: Picture[]
  doneAt: number
}

// Actions that can be dispatched to update the state
interface UpdateAction {
  picture?: Picture
  tool?: string
  color?: Color
  undo?: never // Ensure undo is not present
}

interface UndoAction {
  undo: true
}

type Action = UpdateAction | UndoAction

// The dispatch function
type Dispatch = (action: Action) => void

// A tool is a function that handles pointer events
type Tool = (
  pos: Position,
  state: State,
  dispatch: Dispatch
) => ((pos: Position, state: State) => void) | void | undefined

// A map of tool names to tool functions
interface Tools {
  [name: string]: Tool
}

// The main editor configuration
interface EditorConfig {
  tools: Tools
  controls: Control[]
  dispatch: Dispatch
}

// A control is a class that creates a UI component
interface Control {
  new (state: State, config: EditorConfig): {
    dom: HTMLElement
    syncState(state: State): void
  }
}

// #endregion

class Picture {
  width: number
  height: number
  pixels: Color[]

  constructor(width: number, height: number, pixels: Color[]) {
    this.width = width
    this.height = height
    this.pixels = pixels
  }

  static empty(width: number, height: number, color: Color): Picture {
    let pixels = new Array(width * height).fill(color)
    return new Picture(width, height, pixels)
  }

  pixel(x: number, y: number): Color {
    return this.pixels[x + y * this.width]
  }

  draw(pixels: Pixel[]): Picture {
    let copy = this.pixels.slice()
    for (let { x, y, color } of pixels) {
      copy[x + y * this.width] = color
    }
    return new Picture(this.width, this.height, copy)
  }
}

function elt(type: string, props: object | null, ...children: (Node | string)[]): HTMLElement {
  let dom = document.createElement(type)
  if (props) Object.assign(dom, props)
  for (let child of children) {
    if (typeof child != 'string') dom.appendChild(child)
    else dom.appendChild(document.createTextNode(child))
  }
  return dom
}

const scale = 10

class PictureCanvas {
  dom: HTMLCanvasElement
  picture: Picture | undefined

  constructor(picture: Picture, pointerDown: (pos: Position) => ((pos: Position) => void) | void) {
    this.dom = elt('canvas', {
      onmousedown: (event: MouseEvent) => this.mouse(event, pointerDown),
      ontouchstart: (event: TouchEvent) => this.touch(event, pointerDown)
    }) as HTMLCanvasElement
    this.syncState(picture)
  }

  syncState(picture: Picture) {
    if (this.picture == picture) return

    // Optimization: only redraw changed pixels
    if (
      this.picture &&
      this.picture.width == picture.width &&
      this.picture.height == picture.height
    ) {
      let cx = this.dom.getContext('2d')!
      for (let i = 0; i < picture.pixels.length; i++) {
        if (this.picture.pixels[i] !== picture.pixels[i]) {
          let x = i % picture.width
          let y = Math.floor(i / picture.width)
          cx.fillStyle = picture.pixels[i]
          cx.fillRect(x * scale, y * scale, scale, scale)
        }
      }
    } else {
      // Full redraw if size changes or no previous picture
      drawPicture(picture, this.dom, scale)
    }
    this.picture = picture
  }

  mouse(downEvent: MouseEvent, onDown: (pos: Position) => ((pos: Position) => void) | void) {
    if (downEvent.button != 0) return
    let pos = pointerPosition(downEvent, this.dom)
    let onMove = onDown(pos)
    if (!onMove) return
    let move = (moveEvent: MouseEvent) => {
      if (moveEvent.buttons == 0) {
        this.dom.removeEventListener('mousemove', move)
      } else {
        let newPos = pointerPosition(moveEvent, this.dom)
        if (newPos.x == pos.x && newPos.y == pos.y) return
        pos = newPos
        onMove(newPos)
      }
    }
    this.dom.addEventListener('mousemove', move)
  }

  touch(startEvent: TouchEvent, onDown: (pos: Position) => ((pos: Position) => void) | void) {
    let pos = pointerPosition(startEvent.touches[0], this.dom)
    let onMove = onDown(pos)
    startEvent.preventDefault()
    if (!onMove) return
    let move = (moveEvent: TouchEvent) => {
      let newPos = pointerPosition(moveEvent.touches[0], this.dom)
      if (newPos.x == pos.x && newPos.y == pos.y) return
      pos = newPos
      onMove(newPos)
    }
    let end = () => {
      this.dom.removeEventListener('touchmove', move)
      this.dom.removeEventListener('touchend', end)
    }
    this.dom.addEventListener('touchmove', move)
    this.dom.addEventListener('touchend', end)
  }
}

function pointerPosition(
  pos: { clientX: number; clientY: number },
  domNode: HTMLElement
): Position {
  let rect = domNode.getBoundingClientRect()
  return {
    x: Math.floor((pos.clientX - rect.left) / scale),
    y: Math.floor((pos.clientY - rect.top) / scale)
  }
}

function drawPicture(picture: Picture, canvas: HTMLCanvasElement, scale: number) {
  canvas.width = picture.width * scale
  canvas.height = picture.height * scale
  let cx = canvas.getContext('2d')!

  for (let y = 0; y < picture.height; y++) {
    for (let x = 0; x < picture.width; x++) {
      cx.fillStyle = picture.pixel(x, y)
      cx.fillRect(x * scale, y * scale, scale, scale)
    }
  }
}

class PixelEditor {
  state: State
  canvas: PictureCanvas
  controls: { dom: HTMLElement; syncState(state: State): void }[]
  dom: HTMLElement

  constructor(state: State, config: EditorConfig) {
    let { tools, controls, dispatch } = config
    this.state = state

    this.canvas = new PictureCanvas(state.picture, pos => {
      let tool = tools[this.state.tool]
      let onMove = tool(pos, this.state, dispatch)
      if (onMove) return (pos: Position) => onMove(pos, this.state)
    })
    this.controls = controls.map(Control => new Control(state, config))
    this.dom = elt(
      'div',
      {},
      this.canvas.dom,
      elt('br'),
      ...this.controls.reduce((a, c) => a.concat(' ', c.dom), [] as (string | HTMLElement)[])
    )
  }
  syncState(state: State) {
    this.state = state
    this.canvas.syncState(state.picture)
    for (let ctrl of this.controls) ctrl.syncState(state)
  }
}

// #region Controls
class ToolSelect {
  select: HTMLSelectElement
  dom: HTMLElement

  constructor(state: State, { tools, dispatch }: EditorConfig) {
    this.select = elt(
      'select',
      {
        onchange: () => dispatch({ tool: this.select.value })
      },
      ...Object.keys(tools).map(name =>
        elt(
          'option',
          {
            selected: name == state.tool
          },
          name
        )
      )
    ) as HTMLSelectElement
    this.dom = elt('label', null, '🖌 Tool: ', this.select)
  }
  syncState(state: State) {
    this.select.value = state.tool
  }
}

class ColorSelect {
  input: HTMLInputElement
  dom: HTMLElement

  constructor(state: State, { dispatch }: { dispatch: Dispatch }) {
    this.input = elt('input', {
      type: 'color',
      value: state.color,
      onchange: () => dispatch({ color: this.input.value })
    }) as HTMLInputElement
    this.dom = elt('label', null, '🎨 Color: ', this.input)
  }
  syncState(state: State) {
    this.input.value = state.color
  }
}

class SaveButton {
  picture: Picture
  dom: HTMLButtonElement

  constructor(state: State) {
    this.picture = state.picture
    this.dom = elt(
      'button',
      {
        onclick: () => this.save()
      },
      '💾 Save'
    ) as HTMLButtonElement
  }
  save() {
    let canvas = elt('canvas', null) as HTMLCanvasElement
    drawPicture(this.picture, canvas, 1)
    let link = elt('a', {
      href: canvas.toDataURL(),
      download: 'pixelart.png'
    })
    document.body.appendChild(link)
    link.click()
    link.remove()
  }
  syncState(state: State) {
    this.picture = state.picture
  }
}

class LoadButton {
  dom: HTMLButtonElement

  constructor(_: State, { dispatch }: { dispatch: Dispatch }) {
    this.dom = elt(
      'button',
      {
        onclick: () => startLoad(dispatch)
      },
      '📁 Load'
    ) as HTMLButtonElement
  }
  syncState() {}
}

function startLoad(dispatch: Dispatch) {
  let input = elt('input', {
    type: 'file',
    onchange: () => finishLoad(input.files && input.files[0], dispatch)
  }) as HTMLInputElement
  document.body.appendChild(input)
  input.click()
  input.remove()
}

function finishLoad(file: File | null, dispatch: Dispatch) {
  if (file == null) return
  let reader = new FileReader()
  reader.addEventListener('load', () => {
    let image = elt('img', {
      onload: () =>
        dispatch({
          picture: pictureFromImage(image)
        }),
      src: reader.result as string
    }) as HTMLImageElement
  })
  reader.readAsDataURL(file)
}

function pictureFromImage(image: HTMLImageElement): Picture {
  let width = Math.min(100, image.width)
  let height = Math.min(100, image.height)
  let canvas = elt('canvas', { width, height }) as HTMLCanvasElement
  let cx = canvas.getContext('2d')!
  cx.drawImage(image, 0, 0)
  let pixels = []
  let { data } = cx.getImageData(0, 0, width, height)

  function hex(n: number) {
    return n.toString(16).padStart(2, '0')
  }
  for (let i = 0; i < data.length; i += 4) {
    let [r, g, b] = data.slice(i, i + 3)
    pixels.push('#' + hex(r) + hex(g) + hex(b))
  }
  return new Picture(width, height, pixels)
}

class UndoButton {
  dom: HTMLButtonElement

  constructor(state: State, { dispatch }: { dispatch: Dispatch }) {
    this.dom = elt(
      'button',
      {
        onclick: () => dispatch({ undo: true }),
        disabled: state.done.length == 0
      },
      '⮪ Undo'
    ) as HTMLButtonElement
  }
  syncState(state: State) {
    this.dom.disabled = state.done.length == 0
  }
}

function historyUpdateState(state: State, action: Action): State {
  if (action.undo == true) {
    if (state.done.length == 0) return state
    return {
      ...state,
      picture: state.done[0],
      done: state.done.slice(1),
      doneAt: 0
    }
  } else if (action.picture && state.doneAt < Date.now() - 1000) {
    return {
      ...state,
      ...action,
      done: [state.picture, ...state.done],
      doneAt: Date.now()
    }
  } else {
    return { ...state, ...action }
  }
}
// #endregion

// #region Tools
function draw(pos: Position, state: State, dispatch: Dispatch) {
  function drawPixel({ x, y }: Position, state: State) {
    let drawn = { x, y, color: state.color }
    dispatch({ picture: state.picture.draw([drawn]) })
  }
  drawPixel(pos, state)
  return drawPixel
}

function rectangle(start: Position, state: State, dispatch: Dispatch) {
  function drawRectangle(pos: Position) {
    let xStart = Math.min(start.x, pos.x)
    let yStart = Math.min(start.y, pos.y)
    let xEnd = Math.max(start.x, pos.x)
    let yEnd = Math.max(start.y, pos.y)
    let drawn: Pixel[] = []
    for (let y = yStart; y <= yEnd; y++) {
      for (let x = xStart; x <= xEnd; x++) {
        drawn.push({ x, y, color: state.color })
      }
    }
    dispatch({ picture: state.picture.draw(drawn) })
  }
  drawRectangle(start)
  return drawRectangle
}
const around = [
  { dx: -1, dy: 0 },
  { dx: 1, dy: 0 },
  { dx: 0, dy: -1 },
  { dx: 0, dy: 1 }
]

function fill({ x, y }: Position, state: State, dispatch: Dispatch) {
  let targetColor = state.picture.pixel(x, y)
  let drawn = [{ x, y, color: state.color }]
  let visited = new Set<string>()
  for (let done = 0; done < drawn.length; done++) {
    for (let { dx, dy } of around) {
      let currentX = drawn[done].x + dx,
        currentY = drawn[done].y + dy
      if (
        currentX >= 0 &&
        currentX < state.picture.width &&
        currentY >= 0 &&
        currentY < state.picture.height &&
        !visited.has(currentX + ',' + currentY) &&
        state.picture.pixel(currentX, currentY) == targetColor
      ) {
        drawn.push({ x: currentX, y: currentY, color: state.color })
        visited.add(currentX + ',' + currentY)
      }
    }
  }
  dispatch({ picture: state.picture.draw(drawn) })
}

function pick(pos: Position, state: State, dispatch: Dispatch) {
  dispatch({ color: state.picture.pixel(pos.x, pos.y) })
}
// #endregion

const startState: State = {
  tool: 'draw',
  color: '#000000',
  picture: Picture.empty(60, 30, '#f0f0f0'),
  done: [],
  doneAt: 0
}

const baseTools: Tools = { draw, fill, rectangle, pick }

const baseControls: Control[] = [ToolSelect, ColorSelect, SaveButton, LoadButton, UndoButton]

function startPixelEditor({
  state = startState,
  tools = baseTools,
  controls = baseControls
}: Partial<{
  state: State
  tools: Tools
  controls: Control[]
}>): HTMLElement {
  let app = new PixelEditor(state, {
    tools,
    controls,
    dispatch(action: Action) {
      state = historyUpdateState(state, action)
      app.syncState(state)
    }
  })
  return app.dom
}

document.querySelector('div')?.appendChild(startPixelEditor({}))
