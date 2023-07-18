## Testlib 术语

相关文章：[Testlib 简介](https://oi-wiki.org/tools/testlib/)

- Generator `数据生成器`，用来生成输入数据(Input)
- Validator `数据校验器`，判断生成的输入数据(Input)是否符合题目要求，如数据范围、格式等。
- Interactor `交互器`，用于特殊评测的交互题。
- Checker`检查器`，用于特殊评测的普通题。

--- 

### 检查器输出结果

| 结果 |	Testlib  别名	| 含义 |
|--- | --- | --- |
|Ok|_ok | 答案正确。|
|Wrong Answer|_wa |答案错误。|
|Presentation Error|_pe / _dirt|答案格式错误。注意包括 Codeforces 在内的许多 OJ 并不区分 PE 和 WA。|
|Partially Correct| _pc(score) |答案部分正确。仅限于有部分分的测试点，其中 score 为一个正整数，从  （没分）到  （可能的最大分数）。|
|Partially Correct| _pc(score) |答案部分正确。仅限于有部分分的测试点，其中 score 为一个正整数，从  （没分）到  （可能的最大分数）。|
|Fail|fail|validator中表示输入不合法，不通过校验。checker 中表示程序内部错误、标准输出有误或选手输出比标准输出更优，需要裁判/出题人关注。（也就是题目锅了）|
| Unexpected EOF | _unexpected_eof | 不可预料的文件末尾，视作PE |

\* 阅读源码可以发现存在` _points`别名，此为PCMS 2软件的一种评测结果，故不作支持。