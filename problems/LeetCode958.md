# LeetCode 958

Approach: DFS greedy - leaves shouldn't have a camera because then a camera would not be optimal. If a leaf has a parent then obviously putting the camera on the parent would cover the grandparent and its siblings too.

Therefore, use one DFS. Compute on the way up.
