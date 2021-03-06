### 26. 删除有序数组中的重复项
- https://leetcode-cn.com/problems/remove-duplicates-from-sorted-array/
- 一个有序数组 nums ，原地删除重复出现的元素，使每个元素只出现一次，
  返回删除后数组的新长度。
- 不要使用额外的数组空间，你必须在**原地**修改输入数组
  并在使用 O(1) 额外空间的条件下完成。

```
示例 1：

输入：nums = [1,1,2]
输出：2, nums = [1,2]
解释：函数应该返回新的长度 2 ，并且原数组 nums 的前两个元素被修改为 1, 2 。不需要考虑数组中超出新长度后面的元素。
```

```
示例 2：

输入：nums = [0,0,1,1,1,2,2,3,3,4]
输出：5, nums = [0,1,2,3,4]
解释：函数应该返回新的长度 5 ， 并且原数组 nums 的前五个元素被修改为 0, 1, 2, 3, 4 。不需要考虑数组中超出新长度后面的元素。
```
- 思路：
  - 和上题思路一样，这回更简单，一个指针表示空间索引，一个指针表示元素索引。
  - 对每一个元素，看空间上是否已经有该元素。
  - 从小到大即可。
  - 返回的n就是空间索引的最后值。