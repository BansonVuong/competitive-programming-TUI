[Judge](https://dmoj.ca/problem/ccc22s3)

# ccc22s3

first instinct:

max pitch M
good sample - a subsequence that has all unique numbers

sample output:

n=3
max = 2
samples = 5

1 2 1
(1)
(1, 2)
(1, 2, 1)
(2, 1)
(1)

n=1
max = 1
samples = 1
(1)

n=2
max = 1
samples = 2
(1)
(1)

n=2
max = 2;
samples = 3
(1)
(1, 2)
(2)

n=3
max=1;
samples = 3

## if max = n, samples = n.

n=3
max=2

1 1 1 (3 samples)
1 2 1 (5 samples)
1 2 2
(1), (2), (2)
(1, 2)
-> 4 samples

n=3
max=3
1 1 1 (3)
1 1 2 (4, same as 1 2 2)
1 2 1 (5)
1 2 3 (6)

n=4
max=4
1 1 1 1 (4)
1 1 1 2 (5)
1 1 2 2 (5) (1 2 2)+1
1 1 2 1 (6) (1 2 1)+1
1 1 2 3 (7)
1 3 2 1 (9)

im looking at this and it seems like dp to me? (10 minutes in)

when is it impossible?
k < n: impossible
m=2: k > 

ill try subtask first

subtask 1: m=2

n=1
1
n=2
1 1 (2)
1 2 (3)
n=3
1 1 1 (3)
1 1 2 (4)
1 2 1 (5)
n=4
1 1 1 1 (4)
1 1 1 2 (5)
1 1 2 1 (6)
1 2 1 2 (7)
n=5
1 1 1 1 1 (5)
2 1 1 1 1
1 2 1 1 1
2 1 2 1 1
1 2 1 2 1
max = n*2-1

grabbed subtask 1 in 20 minutes

okay for each what is the max?

m=1 max = n
m=2 max = 2n-1
m=3
n=3
3 of length 1
2 of length 2
1 of length 3
n=10
10 of length 1
9 of length 2
8 of length 3

~~max of m is mn-m~~

so how do i remove from an optimal one?

n=6
1 2 3 4 5 6 (21) (6*1, 5*2, 4*3, 3*4, 2*5, 1*6)
1 2 3 4 5 1 (20)
1 2 3 4 1 2 (18)
1 2 3 1 2 3 (15)

we can make k of form
n(n+1)/2 - (n-k)(n-k+1)/2

its greedy. essentially, at each position, you have a certain budget remaining. try to maximize the budget.

n=5, m=5, k=8
1 used: 1 budget: 4
1 2 used: 3 budget: 5
1 2 3 used: 6 budget: 6
1 2 3 3 used: 7 budget: 7
1 2 3 3 3 used: 8 budget: 8

n=5, m=5, k=9
1 u1b5
1 2 u3b6
1 2 3 u6b7
~~1 2 3 4 u10b8~~
~~1 2 3 1 u9b8~~
1 2 3 3 u7n8
1 2 3 3 1 u9b9
12, 23, 31
123

n=5 m=5 k=9
b5u1 1
b6u3 1 2
b7u6 1 2 3
b8u8 1 2 3 3
b9u9 1 2 3 3 1

building:
1 2 3 1 1

DONE:
Claude's read on it:
**Got on your own:** The feasibility bounds (when to print -1), the range of valid K, the greedy intuition of "maximize budget at each step, leave room for the rest."

**Needed a push:** The L_i framing — that total good samples = sum of longest-distinct-suffix lengths. You were thinking about it combinatorially (counting samples directly) instead of decomposing per-position.

**Kept tripping on during implementation:** Separating L values from pitches. You mixed them in the same array three times. The conceptual gap was clear (you understood the greedy) but translating "L_i says repeat the pitch L_i positions back" into clean code took too long.

**Takeaway for Tuesday:** When you see "count subarrays with property X," immediately think per-endpoint contribution (how many valid subarrays *end* at each position). That's the observation pattern you were missing. The implementation lesson: if your algorithm has two distinct concepts (L values vs pitches), use two distinct variables from the start.
