package weight

import (
	"math"
)

type WeightedMatchLong struct {
	V           int   // the number of vertices in the graph
	E           int   // the number of edges    in the graph
	dummyVertex int   // artifical vertex for boundary conditions
	dummyEdge   int   // artifical edge   for boundary conditions
	a           []int // adjacency list
	end         []int
	mate        []int
	weight      []int64

	base       []int
	lastEdge   [3]int
	lastVertex []int
	link       []int
	nextDelta  []int64
	nextEdge   []int
	nextPair   []int
	nextVertex []int
	y          []int64

	delta, lastDelta                       int64
	newBase, nextBase, stopScan, pairPoint int
	neighbor, newLast, nextPoint           int
	oldFirst, secondMate                   int
	f, nxtEdge, nextE, nextU               int
	e, v, i                                int
}

func NewWeightedMatchLong() *WeightedMatchLong {
	return &WeightedMatchLong{}
}

const (
	MINIMIZE  = true
	MAXIMIZE  = false
	UNMATCHED = 0
)

func (w *WeightedMatchLong) WeightedMatchLong(costs [][]int64, minimizeWeight bool) []int {
	w.input(costs)

	// W1. Initialize.
	w.initialize(costs, minimizeWeight)
	for {
		// W2. Start a new search.
		w.delta = 0
		for w.v = 1; w.v <= w.V; w.v++ {
			if w.mate[w.v] == w.dummyEdge {
				// Link all exposed vertices.
				w.pointer(w.dummyVertex, w.v, w.dummyEdge)
			}
		}

		// W3. Get next edge.
		for {
			w.i = 1
			for j := 2; j <= w.V; j++ {
				/* !!! Dissertation, p. 213, it is nextDelta[i] < nextDelta[j]
				 * When I make it <, the routine seems to do nothing.
				 */
				if w.nextDelta[w.i] > w.nextDelta[j] {
					w.i = j
				}
			}

			// delta is the minimum slack in the next edge.
			w.delta = w.nextDelta[w.i]

			if w.delta == w.lastDelta {
				// W8. Undo blossoms.
				w.setBounds()
				w.unpairAll()
				for w.i = 1; w.i <= w.V; w.i++ {
					w.mate[w.i] = w.end[w.mate[w.i]]
					if w.mate[w.i] == w.dummyVertex {
						w.mate[w.i] = UNMATCHED
					}
				}

				// W9.
				return w.mate
			}

			// W4. Assign pair links.
			w.v = w.base[w.i]
			if w.link[w.v] >= 0 {
				if w.pair() {
					break
				}
			} else {
				// W5. Assign pointer link.
				ww := w.bmate(w.v) // blossom w is matched with blossom v.
				if w.link[ww] < 0 {
					// w is unlinked.
					w.pointer(w.v, ww, w.oppEdge(w.nextEdge[w.i]))
				} else {
					// W6. Undo a pair link.
					w.unpair(w.v, ww)
				}
			}
		}

		// W7. Enlarge the matching.
		w.lastDelta -= w.delta
		w.setBounds()

		g := w.oppEdge(w.e)
		w.rematch(w.bend(w.e), g)
		w.rematch(w.bend(g), w.e)
	}
}

// Begin 5 simple functions
//
func (w *WeightedMatchLong) bend(e int) int  { return w.base[w.end[e]] }
func (w *WeightedMatchLong) blink(v int) int { return w.base[w.end[w.link[v]]] }
func (w *WeightedMatchLong) bmate(v int) int { return w.base[w.end[w.mate[v]]] }
func (w *WeightedMatchLong) oppEdge(e int) int {
	return func() int {
		if (e-w.V)%2 == 0 {
			return e - 1
		} else {
			return e + 1
		}
	}()
}
func (w *WeightedMatchLong) slack(e int) int64 {
	return w.y[w.end[e]] + w.y[w.end[w.oppEdge(e)]] - w.weight[e]
}

func (w *WeightedMatchLong) initialize(costs [][]int64, minimizeWeight bool) {
	// initialize basic data structures
	w.setUp(costs)
	w.dummyVertex = w.V + 1
	w.dummyEdge = w.V + 2*w.E + 1
	w.end[w.dummyEdge] = w.dummyVertex

	minWeight, maxWeight := int64(math.MaxInt64), int64(math.MinInt64)
	for i := 0; i < w.V; i++ {
		for j := i + 1; j < w.V; j++ {
			var cost = 2 * costs[i][j]
			if cost > maxWeight {
				maxWeight = cost
			}
			if cost < minWeight {
				minWeight = cost
			}
		}
	}

	// If minimize costs, invert weights
	if minimizeWeight {
		if w.V%2 != 0 {
			panic("111111!!")
		}
		maxWeight += 2 // Don't want all 0 weight
		for i := w.V + 1; i <= w.V+2*w.E; i++ {
			w.weight[i] = maxWeight - w.weight[i]
		}
		maxWeight = maxWeight - minWeight
	}

	w.lastDelta = maxWeight / 2

	allocationSize := w.V + 2
	w.mate = make([]int, allocationSize)

	w.link = make([]int, allocationSize)
	w.base = make([]int, allocationSize)
	w.nextVertex = make([]int, allocationSize)
	w.lastVertex = make([]int, allocationSize)
	w.y = make([]int64, allocationSize)
	w.nextDelta = make([]int64, allocationSize)
	w.nextEdge = make([]int, allocationSize)

	allocationSize = w.V + 2*w.E + 2
	w.nextPair = make([]int, allocationSize)

	for w.i = 1; w.i <= w.V+1; w.i++ {
		w.mate[w.i] = w.dummyEdge
		w.nextEdge[w.i] = w.dummyEdge
		w.nextVertex[w.i] = 0
		w.link[w.i] = -w.dummyEdge
		w.base[w.i] = w.i
		w.lastVertex[w.i] = w.i
		w.y[w.i] = w.lastDelta
		w.nextDelta[w.i] = w.lastDelta
	}
}

func (w *WeightedMatchLong) input(costs [][]int64) {
	w.V = len(costs)
	w.E = w.V * (w.V - 1) / 2

	allocationSize := w.V + 2*w.E + 2
	w.a = make([]int, allocationSize)
	w.end = make([]int, allocationSize)
	w.weight = make([]int64, allocationSize)
	for i := 0; i < allocationSize; i++ {
		w.a[i] = 0
		w.end[i] = 0
		w.weight[i] = 0
		//System.out.println("input: i: " + i + ", a: " + a[i] + " " + end[i] + " " + weight[i] );
	}
}

func (w *WeightedMatchLong) insertPair() {
	var deltaE int64 // !! check declaration.

	// IP1. Prepare to insert.
	deltaE = w.slack(w.e) / 2

	w.nextPoint = w.nextPair[w.pairPoint]

	// IP2. Fint insertion point.
	for ; w.end[w.nextPoint] < w.neighbor; w.nextPoint = w.nextPair[w.nextPoint] {
		w.pairPoint = w.nextPoint
	}

	if w.end[w.nextPoint] == w.neighbor {
		// IP3. Choose the edge.
		if deltaE >= w.slack(w.nextPoint)/2 { // !!! p. 220. reversed in diss.
			return
		}
		w.nextPoint = w.nextPair[w.nextPoint]
	}

	// IP4.
	w.nextPair[w.pairPoint] = w.e
	w.pairPoint = w.e
	w.nextPair[w.e] = w.nextPoint

	// IP5. Update best linking edge.
	if w.nextDelta[w.newBase] > deltaE {
		w.nextDelta[w.newBase] = deltaE
	}
}

func (w *WeightedMatchLong) linkPath(e int) {
	var u int // !! declaration?

	// L1. Done?
	for /* L1. */ w.v = w.bend(e); w.v != w.newBase; w.v = w.bend(e) {
		// L2. Link next vertex.
		u = w.bmate(w.v)
		w.link[u] = w.oppEdge(e)

		// L3. Add vertices to blossom list.
		w.nextVertex[w.newLast] = w.v
		w.nextVertex[w.lastVertex[w.v]] = u
		w.newLast = w.lastVertex[u]
		w.i = w.v

		// L4. Update base.
		f := func() {
			w.base[w.i] = w.newBase
			w.i = w.nextVertex[w.i]
		}
		f()
		for w.i != w.dummyVertex {
			f()
		}

		// L5. Get next edge.
		e = w.link[w.v]
	}
}

func (w *WeightedMatchLong) mergePairs(v int) {
	// MP1. Prepare to merge.
	w.nextDelta[v] = w.lastDelta

	w.pairPoint = w.dummyEdge
	for w.f = w.nextEdge[v]; w.f != w.dummyEdge; {
		// MP2. Prepare to insert.
		w.e = w.f
		w.neighbor = w.end[w.e]
		w.f = w.nextPair[w.f]

		// MP3. Insert edge.
		if w.base[w.neighbor] != w.newBase {
			w.insertPair()
		}
	}
}

func (ww *WeightedMatchLong) pair() bool {

	var u, w, temp int

	// PA1. Prepare to find edge.
	ww.e = ww.nextEdge[ww.v]

	// PA2. Find edge.
	for ww.slack(ww.e) != 2*ww.delta {
		ww.e = ww.nextPair[ww.e]
	}

	// PA3. Begin flagging vertices.
	w = ww.bend(ww.e)
	ww.link[ww.bmate(w)] = -ww.e // Flag bmate(w)

	u = ww.bmate(ww.v)

	// PA4. Flag vertices.
	for ww.link[u] != -ww.e { // u is NOT FLAGGED
		ww.link[u] = -ww.e

		if ww.mate[w] != ww.dummyEdge {
			temp = ww.v
			ww.v = w
			w = temp
		}
		ww.v = ww.blink(ww.v)
		u = ww.bmate(ww.v)
	}

	// PA5. Augmenting path?
	if u == ww.dummyVertex && ww.v != w {
		return true // augmenting path found
	}

	// PA6. Prepare to link vertices.
	ww.newLast = ww.v
	ww.newBase = ww.v
	ww.oldFirst = ww.nextVertex[ww.v]

	// PA7. Link vertices.
	ww.linkPath(ww.e)
	ww.linkPath(ww.oppEdge(ww.e))

	// PA8. Finish linking.
	ww.nextVertex[ww.newLast] = ww.oldFirst
	if ww.lastVertex[ww.newBase] == ww.newBase {
		ww.lastVertex[ww.newBase] = ww.newLast
	}

	// PA9. Start new pair list.
	ww.nextPair[ww.dummyEdge] = ww.dummyEdge
	ww.mergePairs(ww.newBase)
	ww.i = ww.nextVertex[ww.newBase]

	// PA10. Merge subblossom's pair list
	f := func() {
		ww.mergePairs(ww.i)
		ww.i = ww.nextVertex[ww.lastVertex[ww.i]]

		// PA11. Scan subblossom.
		ww.scan(ww.i, 2*ww.delta-ww.slack(ww.mate[ww.i]))
		ww.i = ww.nextVertex[ww.lastVertex[ww.i]]
	}
	f()
	// PA12. More blossoms?
	for ww.i != ww.oldFirst {
		f()
	}

	// PA14.
	return false
}

func (w *WeightedMatchLong) pointer(u int, v int, e int) {
	var i int
	var del int64
	w.link[u] = -w.dummyEdge
	w.nextVertex[w.lastVertex[u]] = w.dummyVertex
	w.nextVertex[w.lastVertex[v]] = w.dummyVertex
	if w.lastVertex[u] != u {
		// u's blossom contains other vertices
		i = w.mate[w.nextVertex[u]]
		del = -w.slack(i) / 2
	} else {
		del = w.lastDelta
	}
	i = u

	// PT3.
	for ; i != w.dummyVertex; i = w.nextVertex[i] {
		w.y[i] += del
		w.nextDelta[i] += del
	}

	// PT4. Link v & scan.

	if w.link[v] < 0 {
		// v is unlinked.
		w.link[v] = e

		w.nextPair[w.dummyEdge] = w.dummyEdge
		w.scan(v, w.delta)
	} else {
		/* Yes, it looks like this statement can be factored out, and put
		 * after if condition, eliminating the else.
		 * However, link is a global variable used in scan:
		 *
		 * I'm not fooling with it!
		 */
		w.link[v] = e
	}
}

func (w *WeightedMatchLong) rematch(firstMate int, e int) {
	// R1. Start rematching.
	w.mate[firstMate] = e
	w.nextE = -w.link[firstMate]

	// R2. Done?
	for w.nextE != w.dummyEdge {
		// R3. Get next edge.
		e = w.nextE
		w.f = w.oppEdge(e)
		firstMate = w.bend(e)
		w.secondMate = w.bend(w.f)
		w.nextE = -w.link[firstMate]

		// R4. Relink and rematch.
		w.link[firstMate] = -w.mate[w.secondMate]
		w.link[w.secondMate] = -w.mate[firstMate]

		w.mate[firstMate] = w.f
		w.mate[w.secondMate] = e
	}
}

func (w *WeightedMatchLong) scan(x int, del int64) {
	var u int
	var delE int64

	// SC1. Initialize.
	w.newBase = w.base[x]
	w.stopScan = w.nextVertex[w.lastVertex[x]]
	for ; x != w.stopScan; x = w.nextVertex[x] /* SC7. */ {
		// SC2. Set bounds & initialize for x.
		w.y[x] += del
		w.nextDelta[x] = w.lastDelta

		w.pairPoint = w.dummyEdge
		w.e = w.a[x] // !!! in dissertation: if there are no edges, go to SC7.
		for ; w.e != 0; w.e = w.a[w.e] /* SC6. */ {
			// SC3. Find a neighbor.
			w.neighbor = w.end[w.e]
			u = w.base[w.neighbor]

			// SC4. Pair link edge.
			if w.link[u] < 0 {

				if w.link[w.bmate(u)] < 0 || w.lastVertex[u] != u {
					delE = w.slack(w.e)
					if w.nextDelta[w.neighbor] > delE {
						w.nextDelta[w.neighbor] = delE
						w.nextEdge[w.neighbor] = w.e

					}
				}
			} else {
				// SC5.
				if u != w.newBase {
					w.insertPair()
				}
			}
		}
	}

	// SC8.
	w.nextEdge[w.newBase] = w.nextPair[w.dummyEdge]
}

func (w *WeightedMatchLong) setBounds() {
	var del int64

	// SB1. Examine each vertex.
	for w.v = 1; w.v <= w.V; w.v++ {
		// SB2. Is vertex a linked base?
		if w.link[w.v] < 0 || w.base[w.v] != w.v {
			// SB8. Update nextDelta.
			w.nextDelta[w.v] = w.lastDelta

			continue
		}

		// SB3. Begin processing linked blossom.
		w.link[w.v] = -w.link[w.v]

		w.i = w.v

		// SB4. Update y in linked blossom.
		// !! discrepancy: dissertation (do-while); Rothberg (while)
		for w.i != w.dummyVertex {
			w.y[w.i] -= w.delta
			w.i = w.nextVertex[w.i]
		}

		// SB5. Is linked blossom matched?
		w.f = w.mate[w.v]
		if w.f != w.dummyEdge {
			// SB6. Begin processing unlinked blossom.
			w.i = w.bend(w.f)
			del = w.slack(w.f)

			// SB7. Update y in unlinked blossom.
			// !! discrepancy: dissertation (do-while); Rothberg (while)
			for w.i != w.dummyVertex {
				w.y[w.i] -= del
				w.i = w.nextVertex[w.i]
			}
		}
		w.nextDelta[w.v] = w.lastDelta
	}
}

func (w *WeightedMatchLong) unlink(oldBase int) {
	// UL1. Prepare to unlink paths.
	w.i = w.nextVertex[oldBase]
	w.newBase = w.i
	w.nextBase = w.nextVertex[w.lastVertex[w.newBase]]
	w.e = w.link[w.nextBase]

	// Loop is executed twice, for the 2 paths containing the subblossom.
	for j := 1; j <= 2; j++ {
		// UL2. Get next path edge.
		w.nxtEdge = w.oppEdge(w.link[w.newBase])

		for k := 1; k <= 2; k++ {
			// UL3. Unlink blossom base.
			w.link[w.newBase] = -w.link[w.newBase]

			// UL4. Update base array.
			f := func() {
				w.base[w.i] = w.newBase
				w.i = w.nextVertex[w.i]
			}
			f()
			for w.i != w.nextBase {
				f()
			}

			// UL5. Get next vertex.
			w.newBase = w.nextBase
			w.nextBase = w.nextVertex[w.lastVertex[w.newBase]]
		}
		// UL6. More vertices?
		for w.link[w.nextBase] == w.nxtEdge {
			// UL2. Get next path edge.
			w.nxtEdge = w.oppEdge(w.link[w.newBase])

			for k := 1; k <= 2; k++ {
				// UL3. Unlink blossom base.
				w.link[w.newBase] = -w.link[w.newBase]

				// UL4. Update base array.

				f := func() {
					w.base[w.i] = w.newBase
					w.i = w.nextVertex[w.i]
				}
				f()
				for w.i != w.nextBase {
					w.base[w.i] = w.newBase
					w.i = w.nextVertex[w.i]
				}

				// UL5. Get next vertex.
				w.newBase = w.nextBase
				w.nextBase = w.nextVertex[w.lastVertex[w.newBase]]
			}
		}

		// UL7. End of path.
		if j == 1 {
			w.lastEdge[1] = w.nxtEdge
			w.nxtEdge = w.oppEdge(w.e)
			if w.link[w.nextBase] == w.nxtEdge {
				continue // check the control flow logic.
			}
		}
		break
	}
	w.lastEdge[2] = w.nxtEdge

	// UL8. Update blossom list.
	if w.base[w.lastVertex[oldBase]] == oldBase {
		w.nextVertex[oldBase] = w.newBase
	} else {
		w.nextVertex[oldBase] = w.dummyVertex
		w.lastVertex[oldBase] = oldBase
	}
}

func (w *WeightedMatchLong) unpair(oldBase int, oldMate int) {
	var e, newbase, u int // !! Are these global (i.e., static)?

	// UP1. Unlink vertices.
	w.unlink(oldBase)

	// UP2. Rematch a path.
	newbase = w.bmate(oldMate)
	if newbase != oldBase {
		w.link[oldBase] = -w.dummyEdge
		w.rematch(newbase, w.mate[oldBase])
		w.link[w.secondMate] = func() int {
			if w.f == w.lastEdge[1] {
				return -w.lastEdge[2]
			} else {
				return -w.lastEdge[1]
			}
		}()
	}

	// UP3. Examine the linking edge.
	e = w.link[oldMate]
	u = w.bend(w.oppEdge(e))
	if u == newbase {
		// UP7. Relink oldmate.
		w.pointer(newbase, oldMate, e)
		return
	}
	w.link[w.bmate(u)] = -e
	// UP4. missing from dissertation.
	f := func() {
		e = -w.link[u]
		w.v = w.bmate(u)
		w.pointer(u, w.v, -w.link[w.v])

		// UP6. Get next blossom.
		u = w.bend(e)
	}
	f()
	for u != newbase {
		f()
	}
	e = w.oppEdge(e)

	// UP7. Relink oldmate
	w.pointer(newbase, oldMate, e)
}

func (w *WeightedMatchLong) setUp(costs [][]int64) {
	currentEdge := w.V + 2
	//System.out.println("setUp: initial currentEdge: " + currentEdge);
	for i := w.V; i >= 1; i-- {
		for j := i - 1; j >= 1; j-- {
			cost := 2 * costs[i-1][j-1]
			w.weight[currentEdge-1] = cost
			w.weight[currentEdge] = cost
			w.end[currentEdge-1] = i
			w.end[currentEdge] = j
			w.a[currentEdge] = w.a[i]
			w.a[i] = currentEdge
			w.a[currentEdge-1] = w.a[j]
			w.a[j] = currentEdge - 1

			currentEdge += 2
		}
	}
}

func (w *WeightedMatchLong) unpairAll() {
	var u int

	// UA1. Unpair each blossom.
	for w.v = 1; w.v <= w.V; w.v++ {
		if w.base[w.v] != w.v || w.lastVertex[w.v] == w.v {
			continue
		}

		// UA2. Prepare to unpair.
		w.nextU = w.v
		w.nextVertex[w.lastVertex[w.nextU]] = w.dummyVertex

		for {
			// UA3. Get next blossom to unpair.
			u = w.nextU
			w.nextU = w.nextVertex[w.nextU]

			// UA4. Unlink a blossom.
			w.unlink(u)
			if w.lastVertex[u] != u {
				// UA5. List subblossoms to unpair.
				w.f = func() int {
					if w.lastEdge[2] == w.oppEdge(w.e) {
						return w.lastEdge[1]
					} else {
						return w.lastEdge[2]
					}
				}()
				w.nextVertex[w.lastVertex[w.bend(w.f)]] = u
			}

			// UA6. Rematch blossom.
			w.newBase = w.bmate(w.bmate(u))
			if w.newBase != w.dummyVertex && w.newBase != u {
				w.link[u] = -w.dummyEdge
				w.rematch(w.newBase, w.mate[u])
			}

			// UA7. Find next blossom to unpair.
			for w.lastVertex[w.nextU] == w.nextU && w.nextU != w.dummyVertex {
				w.nextU = w.nextVertex[w.nextU]
			}
			if w.lastVertex[w.nextU] == w.nextU && w.nextU == w.dummyVertex {
				break
			}
		}
	}
}
