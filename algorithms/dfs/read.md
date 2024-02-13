# DFS

깊이 우선 탐색

- 모든 노드를 탐색하기 위하여 사용
- BFS에 비하면 구현이 간단한 편 recursive를 돌거나 for loop를 돌거나
- 단순 검색 속도 자체는 BFS에 비해 느린 감이 있음
- pre-order traversal를 포함한 여러 트리 순회 알고리즘은 모두 DFS의 한 종류
- 한 분기를 모두 처리하고 backtracking을 통해 다음 분기로 넘어가는 방식

## 구현

1. 리스트 (혹은 스택)을 이용한반복 구현
```python
def dfs(graph, start_node):
    ## 기본은 항상 두개의 리스트를 별도로 관리해주는 것
    need_visited, visited = list(), list()
    ## 시작 노드를 시정하기 
    need_visited.append(start_node)
    ## 만약 아직도 방문이 필요한 노드가 있다면,
    while need_visited: 
        ## 그 중에서 가장 마지막 데이터를 추출 (스택 구조의 활용)
        node = need_visited.pop()
        ## 만약 그 노드가 방문한 목록에 없다면
        if node not in visited:
            ## 방문한 목록에 추가하기 
            visited.append(node)
            ## 그 노드에 연결된 노드를 
            need_visited.extend(graph[node])
    return visited
```

2. 재귀 구현

```python
def dfs(graph, start_node, visited = list()):
    visited.append(start_node)
    for node in graph[start_node]:
        if node not in visited:
            ## 언제 조건 처리하냐에 따라 pre-order, in-order, post-order가 결정됨
            dfs(graph, node, visited)
    return visited
```
