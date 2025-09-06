export {}

// 为什么只需要一个 PGroup.empty 值，而不是每次都使用一个创建新空映射的函数？
// class PGroup {
//   // Your code here
// }
//
// let a = PGroup.empty.add("a");
// let ab = a.add("b");
// let b = ab.delete("a");

// console.log(b.has("b"));
// // → true
// console.log(a.has("b"));
// // → false
// console.log(b.has("a"));
// // → false
//
// 要将empty属性添加到构造函数中，可以将其声明为静态属性。
//
// !你只需要一个empty实例，因为所有空组都是一样的，而且类的实例不会改变。你可以从这个单一的空组中创建许多不同的组，而不会影响它。

const roads = [
  "Alice's House-Bob's House",
  "Alice's House-Cabin",
  "Alice's House-Post Office",
  "Bob's House-Town Hall",
  "Daria's House-Ernie's House",
  "Daria's House-Town Hall",
  "Ernie's House-Grete's House",
  "Grete's House-Farm",
  "Grete's House-Shop",
  'Marketplace-Farm',
  'Marketplace-Post Office',
  'Marketplace-Shop',
  'Marketplace-Town Hall',
  'Shop-Town Hall'
]

function buildGraph(edges: string[]) {
  function addEdge(from: string, to: string) {
    if (from in graph) {
      graph[from].push(to)
    } else {
      graph[from] = [to]
    }
  }

  const graph: Record<string, string[]> = Object.create(null)
  for (const [from, to] of edges.map(r => r.split('-'))) {
    addEdge(from, to)
    addEdge(to, from)
  }
  return graph
}

const roadGraph = buildGraph(roads)

class VillageState {
  static random(parcelCount = 5): VillageState {
    const parcels: { place: string; address: string }[] = []
    for (let i = 0; i < parcelCount; i++) {
      const address = randomPick(Object.keys(roadGraph))
      let place: string
      do {
        place = randomPick(Object.keys(roadGraph))
      } while (place === address)
      parcels.push({ place, address })
    }
    return new VillageState('Post Office', parcels)
  }

  readonly place: string
  readonly parcels: { place: string; address: string }[]

  constructor(place: string, parcels: { place: string; address: string }[]) {
    this.place = place
    this.parcels = parcels
  }

  move(destination: string) {
    if (!roadGraph[this.place].includes(destination)) {
      return this
    } else {
      const parcels = this.parcels
        .map(p => {
          if (p.place != this.place) return p
          return { place: destination, address: p.address }
        })
        .filter(p => p.place != p.address)
      return new VillageState(destination, parcels)
    }
  }
}

function runRobot(
  state: VillageState,
  robot: (state: VillageState, memory: string[]) => { direction: string; memory: string[] },
  memory: string[]
): void {
  for (let turn = 0; ; turn++) {
    if (state.parcels.length == 0) {
      console.log(`Done in ${turn} turns`)
      break
    }
    const action = robot(state, memory)
    state = state.move(action.direction)
    memory = action.memory
    console.log(`Moved to ${action.direction}`)
  }
}

function randomRobot(
  state: VillageState,
  memory: string[]
): { direction: string; memory: string[] } {
  return { direction: randomPick(roadGraph[state.place]), memory }
}

function randomPick<T>(array: T[]): T {
  const choice = Math.floor(Math.random() * array.length)
  return array[choice]
}

const mailRoute = [
  "Alice's House",
  'Cabin',
  "Alice's House",
  "Bob's House",
  'Town Hall',
  "Daria's House",
  "Ernie's House",
  "Grete's House",
  'Shop',
  "Grete's House",
  'Farm',
  'Marketplace',
  'Post Office'
]

function routeRobot(
  state: VillageState,
  memory: string[]
): { direction: string; memory: string[] } {
  if (memory.length == 0) {
    memory = mailRoute
  }
  return { direction: memory[0], memory: memory.slice(1) }
}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  runRobot(VillageState.random(), randomRobot, [])
  runRobot(VillageState.random(), routeRobot, [])
}
