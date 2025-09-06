export {}
// 我们希望有一种易于人类阅读和编辑的方式来指定关卡
// 句点表示空的空间，井号 (#) 字符表示墙壁，加号表示熔岩。玩家的起始位置是 at 符号 (@)。
// 每个 O 字符都是一枚硬币，顶部的等号 (=) 是一块左右移动的熔岩块。
// 我们将支持两种额外的移动熔岩类型：竖线 (|) 创建垂直移动的熔岩块，而 v 表示滴落熔岩——垂直移动的熔岩，它不会来回弹跳，只会向下移动，当它撞击地面时跳回到其起始位置。
// 整个游戏由玩家必须完成的多个关卡组成。当所有硬币都被收集时，关卡就完成了。如果玩家碰到熔岩，当前关卡将恢复到其起始位置，玩家可以再次尝试。

const simpleLevelPlan = `
......................
..#................#..
..#..............=.#..
..#.........o.o....#..
..#.@......#####...#..
..#####............#..
......#++++++++++++#..
......##############..
......................`

type Actor = Player | Lava | Coin

function overlap(actor1: Actor, actor2: Actor) {
  return (
    actor1.pos.x + actor1.size.x > actor2.pos.x &&
    actor1.pos.x < actor2.pos.x + actor2.size.x &&
    actor1.pos.y + actor1.size.y > actor2.pos.y &&
    actor1.pos.y < actor2.pos.y + actor2.size.y
  )
}

class Vec {
  readonly x: number
  readonly y: number

  constructor(x: number, y: number) {
    this.x = x
    this.y = y
  }

  plus(other: Vec) {
    return new Vec(this.x + other.x, this.y + other.y)
  }

  times(factor: number) {
    return new Vec(this.x * factor, this.y * factor)
  }
}

class State {
  static start(level: Level) {
    return new State(level, level.startActors, 'playing')
  }

  readonly level: Level
  readonly actors: Actor[]
  readonly status: string

  constructor(level: Level, actors: Actor[], status: string) {
    this.level = level
    this.actors = actors
    this.status = status
  }

  update(time: number, keys: { [key: string]: boolean }): State {
    let actors = this.actors.map(actor => actor.update(time, this, keys))
    let newState = new State(this.level, actors, this.status)

    if (newState.status !== 'playing') return newState

    let player = newState.player
    if (this.level.touches(player.pos, player.size, 'lava')) {
      return new State(this.level, actors, 'lost')
    }

    for (let actor of actors) {
      if (actor !== player && overlap(actor, player)) {
        newState = (actor as Lava | Coin).collide(newState)
      }
    }
    return newState
  }

  get player() {
    return this.actors.find(a => a.type == 'player')!
  }
}

class Level {
  readonly height: number
  readonly width: number
  readonly startActors: Actor[]
  readonly rows: string[][]

  constructor(plan: string) {
    let rows = plan
      .trim()
      .split('\n')
      .map(l => [...l])
    this.height = rows.length
    this.width = rows[0].length
    this.startActors = []

    this.rows = rows.map((row, y) => {
      return row.map((ch, x) => {
        let type = levelChars[ch as keyof typeof levelChars]
        if (typeof type != 'string') {
          let pos = new Vec(x, y)
          // @ts-ignore
          this.startActors.push(type.create(pos, ch))
          type = 'empty'
        }
        return type
      })
    })
  }

  /**
   * 一个矩形（由位置和大小指定）是否接触到给定类型的网格元素.
   */
  touches(pos: Vec, size: Vec, type: string) {
    let xStart = Math.floor(pos.x)
    let xEnd = Math.ceil(pos.x + size.x)
    let yStart = Math.floor(pos.y)
    let yEnd = Math.ceil(pos.y + size.y)

    for (let y = yStart; y < yEnd; y++) {
      for (let x = xStart; x < xEnd; x++) {
        let isOutside = x < 0 || x >= this.width || y < 0 || y >= this.height
        let here = isOutside ? 'wall' : this.rows[y][x]
        if (here == type) return true
      }
    }
    return false
  }
}

const playerXSpeed = 7
const gravity = 30
const jumpSpeed = 17
class Player {
  static create(pos: Vec) {
    return new Player(pos.plus(new Vec(0, -0.5)), new Vec(0, 0))
  }
  private static readonly _sharedSize = new Vec(0.8, 1.5)

  readonly pos: Vec
  readonly speed: Vec

  constructor(pos = new Vec(0, 0), speed = new Vec(0, 0)) {
    this.pos = pos
    this.speed = speed
  }

  update(time: number, state: State, keys: { [key: string]: boolean }): Player {
    let xSpeed = 0
    if (keys['ArrowLeft']) xSpeed -= playerXSpeed
    if (keys['ArrowRight']) xSpeed += playerXSpeed
    let pos = this.pos
    let movedX = pos.plus(new Vec(xSpeed * time, 0))
    if (!state.level.touches(movedX, this.size, 'wall')) {
      pos = movedX
    }

    let ySpeed = this.speed.y + time * gravity
    let movedY = pos.plus(new Vec(0, ySpeed * time))
    if (!state.level.touches(movedY, this.size, 'wall')) {
      pos = movedY
    } else if (keys['ArrowUp'] && ySpeed > 0) {
      ySpeed = -jumpSpeed
    } else {
      ySpeed = 0
    }
    return new Player(pos, new Vec(xSpeed, ySpeed))
  }

  get type() {
    return 'player'
  }

  get size(): Vec {
    return Player._sharedSize
  }
}
// @ts-ignore
// Player.prototype.size = new Vec(0.8, 1.5)

class Lava {
  static create(pos = new Vec(0, 0), ch = '=') {
    if (ch == '=') {
      return new Lava(pos, new Vec(2, 0))
    } else if (ch == '|') {
      return new Lava(pos, new Vec(0, 2))
    } else if (ch == 'v') {
      return new Lava(pos, new Vec(0, 3), pos)
    }
  }

  private static readonly _sharedSize = new Vec(1, 1)

  readonly pos: Vec
  readonly speed: Vec
  readonly reset: Vec | null

  constructor(pos: Vec, speed: Vec, reset: Vec | null = null) {
    this.pos = pos
    this.speed = speed
    this.reset = reset
  }

  update(time: number, state: State): Lava {
    let newPos = this.pos.plus(this.speed.times(time))
    if (!state.level.touches(newPos, this.size, 'wall')) {
      return new Lava(newPos, this.speed, this.reset)
    } else if (this.reset) {
      return new Lava(this.reset, this.speed, this.reset)
    } else {
      return new Lava(this.pos, this.speed.times(-1))
    }
  }

  collide(state: State): State {
    return new State(state.level, state.actors, 'lost')
  }

  get type() {
    return 'lava'
  }

  get size(): Vec {
    return Lava._sharedSize
  }
}

const wobbleSpeed = 8,
  wobbleDist = 0.07
class Coin {
  static create(pos = new Vec(0, 0)) {
    let basePos = pos.plus(new Vec(0.2, 0.1))
    return new Coin(basePos, basePos, Math.random() * Math.PI * 2)
  }

  private static readonly _sharedSize = new Vec(0.6, 0.6)

  readonly pos: Vec
  readonly basePos: Vec
  readonly wobble: number

  constructor(pos: Vec, basePos: Vec, wobble: number) {
    this.pos = pos
    this.basePos = basePos
    this.wobble = wobble
  }

  update(time: number, state: State): Coin {
    let wobble = this.wobble + time * wobbleSpeed
    let wobblePos = Math.sin(wobble) * wobbleDist
    return new Coin(this.basePos.plus(new Vec(0, wobblePos)), this.basePos, wobble)
  }

  collide(state: State): State {
    let filtered = state.actors.filter(a => a !== this)
    let status = state.status
    if (!filtered.some(a => a.type == 'coin')) status = 'won'
    return new State(state.level, filtered, status)
  }

  get type() {
    return 'coin'
  }

  get size(): Vec {
    return Coin._sharedSize
  }
}

const levelChars = {
  '.': 'empty',
  '#': 'wall',
  '+': 'lava',
  '@': Player,
  o: Coin,
  '=': Lava,
  '|': Lava,
  v: Lava
}

let simpleLevel = new Level(simpleLevelPlan)
console.log(`${simpleLevel.width} by ${simpleLevel.height}`)
// → 22 by 9

function elt(name: string, attrs: { [key: string]: string }, ...children: HTMLElement[]) {
  let dom = document.createElement(name)
  for (let attr of Object.keys(attrs)) {
    dom.setAttribute(attr, attrs[attr])
  }
  for (let child of children) {
    dom.appendChild(child)
  }
  return dom
}

class DOMDisplay {
  readonly dom: HTMLElement
  actorLayer: HTMLElement | null

  constructor(parent: HTMLElement, level: Level) {
    this.dom = elt('div', { class: 'game' }, drawGrid(level))
    this.actorLayer = null
    parent.appendChild(this.dom)
  }

  clear() {
    this.dom.remove()
  }

  syncState(state: State): void {
    if (this.actorLayer) {
      this.actorLayer.remove()
    }
    this.actorLayer = drawActors(state.actors)
    this.dom.appendChild(this.actorLayer)
    this.dom.className = `game ${state.status}`
    /** 它确保如果关卡突出到视窗之外，我们将滚动视窗以确保玩家位于其中心附近. */
    this.scrollPlayerIntoView(state)
  }

  scrollPlayerIntoView(state: State): void {
    let width = this.dom.clientWidth
    let height = this.dom.clientHeight
    let margin = width / 3

    // The viewport
    let left = this.dom.scrollLeft,
      right = left + width
    let top = this.dom.scrollTop,
      bottom = top + height

    let player = state.player
    let center = player.pos.plus(player.size.times(0.5)).times(scale)

    if (center.x < left + margin) {
      this.dom.scrollLeft = center.x - margin
    } else if (center.x > right - margin) {
      this.dom.scrollLeft = center.x + margin - width
    }
    if (center.y < top + margin) {
      this.dom.scrollTop = center.y - margin
    } else if (center.y > bottom - margin) {
      this.dom.scrollTop = center.y + margin - height
    }
  }
}

/** 单个单位在屏幕上占据的像素数量. */
const scale = 20

function drawGrid(level: Level) {
  return elt(
    'table',
    {
      class: 'background',
      style: `width: ${level.width * scale}px`
    },
    ...level.rows.map(row =>
      elt('tr', { style: `height: ${scale}px` }, ...row.map(type => elt('td', { class: type })))
    )
  )
}

// .background    { background: rgb(52, 166, 251);
//   table-layout: fixed;
//   border-spacing: 0;              }
// .background td { padding: 0;                     }
// .lava          { background: rgb(255, 100, 100); }
// .wall          { background: white;              }
// .actor  { position: absolute;            }
// .coin   { background: rgb(241, 229, 89); }
// .player { background: rgb(64, 64, 64);   }
// .lost .player {
//   background: rgb(160, 64, 64);
// }
// .won .player {
//   box-shadow: -4px -7px 8px white, 4px -7px 8px white;
// }
// .game {
//   overflow: hidden;
//   max-width: 600px;
//   max-height: 450px;
//   position: relative;
// }
function drawActors(actors: Actor[]) {
  return elt(
    'div',
    {},
    ...actors.map(actor => {
      let rect = elt('div', { class: `actor ${actor.type}` })
      rect.style.width = `${actor.size.x * scale}px`
      rect.style.height = `${actor.size.y * scale}px`
      rect.style.left = `${actor.pos.x * scale}px`
      rect.style.top = `${actor.pos.y * scale}px`
      return rect
    })
  )
}

function trackKeys(keys: string[]): { [key: string]: boolean } {
  let down = Object.create(null)
  function track(event: KeyboardEvent) {
    if (keys.includes(event.key)) {
      down[event.key] = event.type == 'keydown'
      event.preventDefault()
    }
  }
  window.addEventListener('keydown', track)
  window.addEventListener('keyup', track)
  return down
}

const arrowKeys = trackKeys(['ArrowLeft', 'ArrowRight', 'ArrowUp'])

function runAnimation(frameFunc: (time: number) => boolean) {
  let lastTime: number | null = null
  function frame(time: number) {
    if (lastTime != null) {
      let timeStep = Math.min(time - lastTime, 100) / 1000
      if (frameFunc(timeStep) === false) return
    }
    lastTime = time
    requestAnimationFrame(frame)
  }
  requestAnimationFrame(frame)
}

function runLevel(level: Level, Display: typeof DOMDisplay) {
  let display = new Display(document.body, level)
  let state = State.start(level)
  let ending = 1
  return new Promise(resolve => {
    runAnimation(time => {
      state = state.update(time, arrowKeys)
      display.syncState(state)
      if (state.status == 'playing') {
        return true
      } else if (ending > 0) {
        ending -= time
        return true
      } else {
        display.clear()
        resolve(state.status)
        return false
      }
    })
  })
}

async function runGame(plans: string[], Display: typeof DOMDisplay) {
  for (let level = 0; level < plans.length; ) {
    let status = await runLevel(new Level(plans[level]), Display)
    if (status == 'won') level++
  }
  console.log("You've won!")
}

{
  const simpleLevel = new Level(simpleLevelPlan)
  const display = new DOMDisplay(document.body, simpleLevel)
  display.syncState(State.start(simpleLevel))
}
