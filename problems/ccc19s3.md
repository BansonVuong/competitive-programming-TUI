# ccc19s3

3x3 grid
first look: construction

a   , a+d   , a+2d
a+b , a+b+c , a+b+2c
a+2b, a+2b+e, a+2b+2e

2b+e = 2(b+c-d)
2b+2e = 2(b+2c-2d)

4/15 - at most 3 x symbols - math solve
3/15 - brute force subtask - -10 and 10
4/15 - 7 x symbols - build it yourself
2/15 - even numbers
2/15 - no restrictions

how do i solve this myself?
## sample 1

8 9 10
16 X 20
24 X 30

this one is predetermined. x in the middle is equal to 20-16/2 + 16
= 18

x at the bottom is equal to 30-24/2 + 24
= 27

so if 2/3 elements of a row are pregiven it is guaranteed to be solvable

## sample 2
14 X X
X X 18
X 16 X

much more freedom. 
note: the diffrence between rows can be 0.

possible solution:
14 16 18
14 16 18
14 16 18

## self sample 1

1 4 7
X X X
X X X

could be

1 4 7
1 4 7
1 4 7

## self sample 2
if all Xs

then just 
0 0 0
0 0 0
0 0 0

## self sample 3

if

0 X X
X X X
X X X

then

0 0 0
0 0 0
0 0 0

## self sample 4

if

0 X 2
X X X
X X X

then solve for middle

0 1 2
X X X
X X X

then 

0 1 2
0 1 2
0 1 2

## self sample 5
if
0 X X
X X X
X X 2

then

if you've got one lone number in a column best bet would be to copy it straight down unless you have a constraint from another number down the line

why would we be given a subtask with even numbers? what does that have any effect?

10 minutes left (25 minutes in): grabbed ST1. 4 points.
