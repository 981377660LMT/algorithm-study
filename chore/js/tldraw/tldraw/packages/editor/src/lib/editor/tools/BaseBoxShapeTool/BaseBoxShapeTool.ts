import { TLShape } from '@tldraw/tlschema'
import { StateNode, TLStateNodeConstructor } from '../StateNode'
import { Idle } from './children/Idle'
import { Pointing } from './children/Pointing'

/** @public */
export abstract class BaseBoxShapeTool extends StateNode {
	static override id = 'box'
	static override initial = 'idle'
	static override children(): TLStateNodeConstructor[] {
		return [Idle, Pointing]
	}

	abstract override shapeType: string

	onCreate?(_shape: TLShape | null): void | null
}
