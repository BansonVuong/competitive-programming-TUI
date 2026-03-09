[Judge](https://dmoj.ca/problem/ccc23s4)

# ccc23s4
8:15pm start
~2h
Naive: For every road, eliminate it and then run dijkstra to find out if distance has increased
Pitfalls: how to find out which roads to cut? Cutting a road of cost 9 vs 2 roads of cost 5, etc.
Key insight gained (spoonfed): Kruskals
Works perfectly for subtask 1 where distance does not matter: just run normal kruskals
Next insight (also spoonfed): sort by distance + length
When processing an edge of length L, there are no other shorter edges of that length. That means the path to that point is already optimal.
~~It is impossible to "regret" not adding an edge because there is already a more optimal path to that node given the MST, and if the edge is guaranteed to be useless for any shortest path because if it does not shorten any path between the two nodes it connects, it will not shorten any path at all.~~
Not necessarily true. Therefore, because N and M are so small, we must run dijkstra to make sure this edge will shorten the distance b/t two edges
