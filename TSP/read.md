> 本文内容由[旅行推销员问题](https://zh.wikipedia.org/wiki/旅行推销员问题)转换而来。

**[旅行商问题](../Page/旅行推销员问题.md "wikilink")**（英語：Travelling salesman problem, TSP）是[组合优化](../Page/组合优化.md "wikilink")中的一个[NP 困难](../Page/NP困难.md "wikilink")问题，在[运筹学](../Page/運籌學.md "wikilink")和[理论计算机科学](../Page/理論計算機科學.md "wikilink")中非常重要。问题内容为“给定一系列城市和每对城市之间的距离，求解访-{}-问每一座城市一次并回到起始城市的最短回路。”

[GLPK_solution_of_a_travelling_salesman_problem.svg](https://zh.wikipedia.org/wiki/File:GLPK_solution_of_a_travelling_salesman_problem.svg "GLPK_solution_of_a_travelling_salesman_problem.svg") TSP 是与[车辆路径问题 ⓦ](https://zh.wikipedia.org/wiki/车辆路径问题 "wikilink")的一种特殊情况。

作为[计算复杂性理论](../Page/計算複雜性理論.md "wikilink")中的一个典型的判定性问题，TSP 的一个版本是给定一个图和长度 _L_，要求回答图中是否存在比 _L_ 短的回路（英语：circuit 或 tour）。该问题被划分为[NP 完全](../Page/NP完全.md "wikilink")问题。已知 TSP 算法[最坏情况下的 ⓦ](https://zh.wikipedia.org/wiki/計算時間 "wikilink")[时间复杂度随着城市数量的增多而成超多项式](../Page/时间复杂度.md "wikilink")（可能是）级别增长。

问题在 1930 年首次被形式化，并且是在最优化中研究最深入的问题之一。许多优化方法都用它作为一个基准。尽管问题在计算上很困难，但已经有了大量的[启发式](../Page/启发法.md "wikilink")和精确方法，因此可以完全求解城市数量上万的实例，并且甚至能在误差 1%范围内估计上百万个城市的问题。[^1]

甚至纯粹形式的 TSP 都有若干应用，如[企划 ⓦ](https://zh.wikipedia.org/wiki/企划 "wikilink")、[物流](../Page/物流.md "wikilink")、[芯片](../Page/集成电路.md "wikilink")制造。稍作修改，就是[DNA 测序](../Page/DNA測序.md "wikilink")等许多领域的一个子问题。在这些应用中，“城市”的概念用来表示客户、焊接点或 DNA 片段，而“距离”的概念表示旅行时间或成本或 DNA 片段之间的相似性度量。TSP 还用在天文学中，观察很多光源的天文学家會希望减少在不同光源之间转动望远镜的时间。在许多应用場景中（如资源或时间窗口有限等等），可能会需要加入额外的约束條件。

## 描述

### 作为图论问题

[Weighted_K4.svg](https://zh.wikipedia.org/wiki/File:Weighted_K4.svg "Weighted_K4.svg") 可以用[无向加权图来对 TSP 建模](<../Page/图_(数学).md> "wikilink")，则城市是图的[顶点](<../Page/顶点_(图论).md> "wikilink")，道路是图的[边](../Page/图论术语.md "wikilink")，道路的距离就是该边的长度。它是起点和终点都在一个特定[顶点](<../Page/顶点_(图论).md> "wikilink")，访-{}-问每个顶点恰好一次的最小化问题。通常，该模型是一个[完全圖](../Page/完全圖.md "wikilink")（即每对顶点由一条边连接）。如果两个城市之间不存在路径，則增加一条非常长的边就可以完成图，而不影响计算最优回路。

### 非对称和对称

在*对称 TSP 问题*中，两座城市之间来回的距离是相等的，形成一个[无向图](<../Page/图_(数学).md> "wikilink")。这种对称性将解的数量减少了一半。在*非对称 TSP 问题*中，可能不是双向的路径都存在，或是来回的距离不同，形成了[有向图 ⓦ](https://zh.wikipedia.org/wiki/有向图 "wikilink")。[交通事故 ⓦ](https://zh.wikipedia.org/wiki/車禍 "wikilink")、[单行道](../Page/单行道.md "wikilink")和出发与到达某些城市机票价格不同等都是打破这种对称性的例子。

### 相关问题

- [图论](../Page/图论.md "wikilink")中的一个等价形式是：给定一个[加权完全图](../Page/图论术语.md "wikilink")（顶点表示城市，边表示道路，权重就会是道路的成本或距离）, 求一权值最小的[哈密尔顿回路](../Page/哈密顿图.md "wikilink")。

<!-- -->

- 返回到起始城市的要求不会改变问题的[计算复杂度](../Page/計算複雜性理論.md "wikilink")，见[哈密頓路徑問題](../Page/哈密顿路径问题.md "wikilink")。

<!-- -->

- 另一个相关问题是（bottleneck TSP）：求[加权图中权重最大的](../Page/图论术语.md "wikilink")[边最小的哈密尔顿回路](<../Page/图_(数学).md> "wikilink")。问题在运输和物流之外都有相当广泛的实际意义。一个典型的例子是在[印刷电路板](../Page/印刷电路板.md "wikilink")制造中：规划[打孔机在 PCB 版上钻孔的路线](../Page/电钻.md "wikilink")。在机械加工或钻孔应用中，“城市”是需要加工的部分或需要钻的（不同大小）的孔，而“旅行成本”包括更换机具所用的时间（单机作业排序问题）。

<!-- -->

- [旅行推銷員问题 ⓦ](https://zh.wikipedia.org/wiki/Set_TSP_problem "wikilink")，处理“国家”中有（一个或多个）“城市”，而旅行商需要在每个“国家”访-{}-问恰好一座“城市”。其中一种应用是在求解时，想要最小化刀具改变次数中。另一种应用与[半导体](../Page/半导体.md "wikilink")制造业中的打孔有关，见 US patent:7054798。令人惊喜的是，Behzad 与 Modarres 证明了广义旅行商问题可以转化为相同城市数量的标准旅行商问题 ，只是要改变[距离矩阵 ⓦ](https://zh.wikipedia.org/wiki/距离矩阵 "wikilink")。[^2]

<!-- -->

- 优先顺序旅行推销员问题处理城市之间存在访-{}-问次序的问题。

<!-- -->

- 旅行购买者问题涉及购买一系列产品的购买者。他可以在若干城市购买这些产品，但价格会有不同，也不是所有城市都有售相同的商品。目标是在城市的子集中间找到一条路径，使得总成本（旅行成本 + 购买成本）最小，并且能够买到所有需求的商品。

## 整数线性规划形式

单钻头的运动可以看成是典型的 TSP 问题。TSP 可以用[整数线性规划来形式化](../Page/整数规划.md "wikilink")。[^3][^4][^5] 用数字 0, ..., _n_ 标记这些城市（打孔位置），并定义：

$$x_{ij} = \begin{cases} 1 & \text{the path goes from city } i \text{ to city } j \\ 0 & \text{otherwise} \end{cases}$$

对于 _i_ = 0, ..., _n_，令 $u_i$ 为一人工变量，最后把 $c_{ij}$ 作为从城市 _i_ 到 _j_ 的距离。那么 TSP 可以写成下面的整数线性规划问题：

$$
\begin{align}
\min &\sum_{i=0}^n \sum_{j\ne i,j=0}^nc_{ij}x_{ij} &&  \\
     & 0 \le x_{ij} \le 1  && i,j=0, \cdots, n  \\
     & u_{i} \in \mathbf{Z} && i=0, \cdots, n \\
     & \sum_{i=0,i\ne j}^n x_{ij} = 1 && j=0, \cdots, n \\
     & \sum_{j=0,j\ne i}^n x_{ij} = 1 && i=0, \cdots, n \\
&u_i-u_j +nx_{ij} \le n-1 && 1 \le i \ne j \le n
\end{align}
$$

第一组等式要求每个城市都能另一个城市前来，而第二组等式要求每个城市都能出发。最后的约束迫使覆盖所有城市的路径只有一条，而不是两条或者多条分散的路径在一起覆盖的。要证明这一点，下面会去证 (1)每个可行解包含只有一条封闭城市序列，以及(2)对于每条覆盖所有城市的单独路径，[虚拟变量](../Page/虚拟变量.md "wikilink") $u_i$ 有值可以满足限制式。

证明可行解中的每个子回路经过 0 号城市（注意到等式保证了只有一条这样的路径），就能证明所有可行解只包含一个封闭城市序列。对于若我们对所有 $x_{ij}=1$ 对应的不等式求和的话，对 _k_ 步不经过 0 号城市的任何子回路，我们得到：

$$nk \leq (n-1)k,$$

这构成矛盾。

必须证明对每个覆盖所有城市的单独回路，虚拟变量 $u_i$ 有值可以满足约束。

为了不失一般性，定义起始点为 0 号城市。如果在第 _t_ 步访-{}-问城市 _i_ 后 (_i_, _t_ = 1, 2, ..., n) 选取 $u_{i}=t$。则

$$u_i-u_j\le n-1,$$

由于 $u_i$ 不大于 _n_ 而 $u_j$ 不小于 1；因此，每当 $x_{ij}=0$ 时满足约束。对于 $x_{ij}=1$，我们有：

$$u_{i} - u_{j} + nx_{ij} = (t) - (t+1) + n = n-1,$$

满足约束。

## 参见

-

- [邮递员问题](../Page/邮递员问题.md "wikilink")

- [柯尼斯堡七桥问题](../Page/柯尼斯堡七桥问题.md "wikilink")

- [车辆路径问题 ⓦ](https://zh.wikipedia.org/wiki/车辆路径问题 "wikilink")

## 注释

## 参考文献

## 延伸阅读

## 外部链接

- [Traveling Salesman Problem](http://www.math.uwaterloo.ca/tsp/index.html) at [University of Waterloo](../Page/滑鐵盧大學.md "wikilink")
- [TSPLIB](https://web.archive.org/web/20170325123309/http://www.iwr.uni-heidelberg.de/groups/comopt/software/TSPLIB95/) at the [University of Heidelberg](../Page/海德堡大学.md "wikilink")
- _[Traveling Salesman Problem](http://demonstrations.wolfram.com/TravelingSalesmanProblem/) _ by Jon McLoone at the Wolfram Demonstrations Project
- [Traveling Salesman movie (on IMDB)](http://www.imdb.com/title/tt1801123/)
- [MAOS-TSP](https://github.com/xfxie/MAOS-TSP) JAVA TSP Solver

[Category:NP 完全问题](https://zh.wikipedia.org/wiki/Category:NP完全问题 "wikilink") [Category:運籌學](https://zh.wikipedia.org/wiki/Category:運籌學 "wikilink") [Category:图算法](https://zh.wikipedia.org/wiki/Category:图算法 "wikilink")

[^1]: 参见已解出精确解 0.05%范围内的 TSP 世界巡游问题。[1](http://www.math.uwaterloo.ca/tsp/world/)
[^2]:
[^3]: , pp.308-309.
[^4]: Tucker, A. W. (1960), "On Directed Graphs and Integer Programs", IBM Mathematical research Project (Princeton University)
[^5]: Dantzig, George B. (1963), _Linear Programming and Extensions_, Princeton, NJ: PrincetonUP, pp. 545–7, ISBN 0-691-08000-3, sixth printing, 1974.
