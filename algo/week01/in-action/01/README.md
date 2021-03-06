### 88. 合并两个有序数组
- https://leetcode-cn.com/problems/merge-sorted-array/
  
- 给你两个有序整数数组 nums1 和 nums2，请你将 nums2 合并到 nums1中，
  使 nums1 成为一个有序数组。
- 初始化nums1 和 nums2 的元素数量分别为m 和 n 。
- 可以假设nums1 的空间大小等于m + n，这样它就有足够的空间保存来自 nums2 的元素。

```
示例 1：

输入：nums1 = [1,2,3,0,0,0], m = 3, nums2 = [2,5,6], n = 3
输出：[1,2,2,3,5,6]
```

```
示例 2：

输入：nums1 = [1], m = 1, nums2 = [], n = 0
输出：[1]
```
- 思路：两个索引，比较大小，谁小放谁。
- 细节：
  - 处理边界问题。
  - 处理空间问题。
    - 如果直接替换nums1，如果按从前到后，2会替换3，3被替代了。
    - 两种思路
      - 1，新建一个数组，长度=nums1.size(), 都是0，那么不怕替代
      - 2. 还是原位替换nums1，但是从后往前，那么也不怕替换。
           - 此时需要从大往小比较。

