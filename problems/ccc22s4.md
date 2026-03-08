# ccc22s4
naive: brute force all 3 combinations (probably probably works for subtask 1 (n<=200))
how to do the geometry to figure out what point is in the middle? i have no clue

key insight: all three arc lengths must be less than half the circumference for (0, 0) to be in the triangle

key insight (spoonfed): instead of counting good triplets, subtract bad triplets

spoonfed: a triplet is bad if the 3 points all fit within some semicircle

sliding window strategy? since max circumference is 1e6, means we can count all in o(n) time, after nlogn sorting and using upper/lower bound

for sliding window, the upper bound may have to be a double (.5) if the circumference is odd

issue: overcounting
for a circle with points 0, 2, 5, 5, and circumference of 10
doubling the array would give 0, 2, 5, 5, 10, 12, 15, 15
sliding window <= C/2 would double count the 5 with the 10s
however we still need to cover the case of triangle with points 0, 2, 5 since this does not 


2 big cases:
case 1: all points fit within window < C/2
bad triplets = all points within the window choose 2

case 2: all points fit within window <=C/2
rest = # of non beginning points 
2a: bad triples = beginning * rest * end
2b: bad triplets = beginning choose 2 * end
2c: bad triples = end choose 2 * beginning
