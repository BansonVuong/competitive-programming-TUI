# ccc22s5

ridiculous amount of time spent on this. almost a week.

first impression: graph theory

~~kruskal + dp?~~

~~balance: highest degree compared to lowest cost~~

spoonfed: tree dp (n-1 edges)
if S is the set of paid students, every student that starts as N must be a neighbor to an S student
If we keep only paid nodes and edges b/t paid nodes, every component must have a Y within it.

dp:
base case: the child is already Y (nothing needed to be done)
if the child is N:
    parent needs to be paid, or the grandchild needs to be paid
    compare cost of getting paid from parent or child

for a given node:
A: is u paid?
B: does u already become Y without needing the parent? (either u was Y, or one of its child is paid)
C: if u is paid: does the paid component containing u already contain an initially Y node?

A no, C yes: not possible

#1 A yes, B yes, C yes: u is Y: valid | u is N: valid (paid, already has starter)

#2 A no, B yes, C no: u is Y: valid | u is N: valid (not paid, either influenced from below or is itself Y)

#3 A no, B no, C no u is Y: invalid | u is N:valid (not paid, not influenced from below (needs parent))

#5 A yes, B no, C no u is Y: invalid | u is N: valid (paid, needs parent)


Leaf base cases:
State 1: possible when u is Y, cost: C_u
State 2: possible when u is Y, cost: 0
State 3: possible when u is N, cost: 0
State 5: possible when u is N, cost: C_u

Throwing up the white flag: Need to learn: Tree DP

Learned: iterative DFS

Learned: solve integer overflow by setting INF to lower number (1e12)

Ended up just focusing on implementing the dp transition and being done.

Honestly I didn't really get this problem, but I have a general sense of what to do. I've spent too much time stuck on this problem.
