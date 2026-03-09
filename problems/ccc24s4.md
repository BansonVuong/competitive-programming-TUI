[Judge](https://dmoj.ca/problem/ccc24s4)

# ccc24s4 

Time: ~2 hours

Key insight: in any simple cycle, one edge will be grey, all other edges will alternate. This creates a tree.

Second problem: colouring edges is simple if the direction is monotonic: if grey edges only connect an ancestor to a descendant, we can simply colour paths based on depth. Otherwise having to go up and back down ends up with a path that does not alternate - at the turning point, there will be two paths of the same colour.
### How to build this tree?
DFS tree ensures all unused edges connect an ancestor to a descendant. Then, we just colour based on parity during DFS.

Loop through all N to ensure all disconnected components are covered.
