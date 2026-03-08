# codeforces1324F

what is the maximum difference between # of white and # black nodes
if i choose some subtree that contains the vertex v
maximize w-b

honestly this is like 2dfs. if the difference gets bigger as the tree goes up then it will dominate over all of its subtrees

so its monotonic - the answer can only ever increase as we go down the tree

count # of white and black vertices in each subtree, and track the diff in that subtree

got this part in 4 minutes

okay this question is kinda ridiculously worded (15 minutes in)

i dont know if the tree has to be rooted at 1

oh its just **a** connected component. doesn't have to be a subtree. 

so ans[nxt] is 
original best value
plus ans[cur] 
subtract any contribution to ans[cur] from the child so we dont double count

solved in 40 minutes

Claude's read:

**Got on your own:** Two DFS structure, monotonicity insight, bottom-up DP idea, that you track diff not separate counts.

**Missed:** That "subtree" means any connected subgraph (not rooted subtree with all descendants), so you can prune bad branches with `max(0, dp[child])`. Once you misread the problem, your whole approach of counting total white/black per subtree was doomed.

**Needed spoonfed:** The rerooting formula for DFS 2. You had the skeleton but no logic in it.

**TLDR:** You had the right structure but got stuck on a misread for 15 min. The actual DP logic had to be handed to you.
