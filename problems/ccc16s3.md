# ccc16s3

first read: graph theory

n-1 edges: tree

we can ignore all edges/nodes that do not form part of the path to/from a pho restaurant. do a dfs from a pho restaurant, if we dont hit a pho restaurant, cut that edge/node out on the way back. then do a tree diameter and find distance of each branch. (5 minutes in)

learned: minimum distance to visit all nodes in a tree is 2*edges-diameter

Claude's read:
**What you got on your own:** The core observation — prune non-pho leaves, work on the minimal subtree. That's the hard part conceptually, and you saw it in 5 minutes.

**What you didn't get:** The `2E - D` formula. You had the diameter and the pruned tree but couldn't connect them into an answer. You went down a rabbit hole trying to compute per-branch distances to the "trunk" (diameter path), which is overcomplicating it.

**What you needed spoonfed:** Just the formula and the intuition behind it (DFS = 2E, save one diameter path by not returning). Plus two small bugs (missing `vis[start]=1`, not resetting `far`).

**Bottleneck:** Not the observation, not the pruning, not the implementation — it's that you didn't have `2E - D` as a known pattern. Now you do. This is a classic "visit all nodes in a subtree" result that shows up repeatedly. Bank it.
