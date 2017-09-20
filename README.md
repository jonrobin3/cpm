# Build instructions

```
go build cpm.go
```

# Run instructions

```
cpm [-k=int] graphDefinitionFile.txt
```

# Description

k-clique percolation method (CPM) is sometimes used to find
associations between stocks, proteins in molecular biology, groups
of people in social networks, and other associations that can be
represented with a graph. A clique is defined as a complete subgraph
of some size k. Below are examples of complete subgraphs for k=3 and
k=4.

```
                       
  Complete graph for k=3  
    +----+      +----+    
    | v0 |------| v1 |    
    +---++      ++---+    
        |        |        
        |        |        
        +-+----+-+        
          | v2 |          
          +----+          
                          
                          
  Complete graph for k=4  
    +----+      +----+    
    | v0 |------| v1 |    
    +--+-\  +---+--+-+    
       |  \ |      |      
       |   -+--\   |      
    +--+-+--+   \--+-+    
    | v2 |------| v3 |    
    +----+      +----+
```

The following four steps describe the CPM:

1. Find all cliques of size k in the graph
2. Create a graph where nodes are cliques of size k
3. Add edges if two nodes (cliques) share k-1 common nodes
4. Each connected component is a community

For example, assume the graph below:

```
   +----+           +----+
   | v2 |-----------| v1 |
   +----++        +-+----+
         |        |       
         |        |       
         +-+----+-+       
     +-----| v3 |------+  
     |     +----+      |  
     |                 |  
     |                 |  
     |                 |  
  +----+            +----+
  | v4 |------------| v5 |
  +--+-+     +------+-+--+
     | |     |        |   
     | +-----+--------++  
     |       |        ||  
     |       |        ||  
  +--+-+-----+      +-+--+
  | v6 |------------| v7 |
  +----+            +----+
     |                 |  
     +----+    +-------+  
          +----+          
          | v8 |          
        +-+----++         
        |       |         
        |       |         
 +----+-+       +-+----+  
 | v9 |-----------|v10 |  
 +----+           +----+
```

If k=3 then we are looking for cliques of size 3, which means we are
looking for three nodes that are fully connected. From our graph, we
would find the following cliques: {v1, v2, v3}, {v3, v4, v5}, {v4,
v5, v6}, {v4, v5, v7}, {v5, v6, v7}, and {v6, v7, v8}. Given those
cliques, we would create a node the corresponds to each clique in a
new community graph. We would then draw edges connecting nodes that,
from the original graph, share k-1 nodes in common. This would yield
the following community graph (where the node labels are a
concatenation of the original graph's node labels):

```
       +-----------+          +-----------+         
       |{v1,v2,v3} |          |{v8,v9,10} |         
       +-----------+          +-----------+         
                                                    
                                                    
                                                    
                   +-----------+                    
               +---|{v3,v4,v5} |---+                
               |   +-----------+   |                
               |                   |                
               |                   |                
               |                   |                
+-----------+--+                   +---+-----------+
|{v4,v5,v6} +--------------------------+{v4,v5,v7} |
+-----+---+-+                          +-+---+-----+
      |   |                              |   |      
      |   | +----------------------------+   |      
      |   +-+-----------------------------+  |      
      |     |                             |  |      
      |     |                             |  |      
+-----+-----+                          +--+--+-----+
|{v4,v6,v7} +--------------------------+{v5,v6,v7} |
+-----------+                       +--+-----------+
            |                       |               
            |                       |               
            |                       |               
            |                       |               
            |      +-----------+    |               
            +----- |{v6,v7,v8} |----+               
                   +-----------+                    
```
If we did the same for k=4 then our community graph would consist of
one node, {v4v5v6v7} beause these are the only four nodes that are
completely connected. 
# Command line options

`-k` is an optional argument that specifies the size of the
clique. If k is not specified, it defaults to k=3.

`graphDefinitionFile` defines the graph to operate on. Vertices
(nodes) are declared on the left hand side (lhs) of the
colon. Vertices on the right hand side (rhs) of the colon define
edges from the definition node to the rhs vertex. For example, from
our previous model graph, v1 is defined as `v1: v2 v3` where `v1`
defines the vertex and `v2` and `v3` define the edges. The entire
graph is defined below:


```
v1: v2 v3
v2: v1 v3
v3: v1 v2 v4 v5
v4: v3 v5 v6 v7
v5: v3 v4 v6 v7
v6: v4 v5 v7 v8
v7: v4 v5 v6 v8
v8: v6 v7 v9 v10
v9: v8 v10
v10: v8 v9
```
