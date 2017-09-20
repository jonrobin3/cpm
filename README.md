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
associations between proteins in molecular biology, stocks or
industrial data in economics, groups of people in social networks,
and other associations that can be represented with a graph. A
clique is defined as a complete subgraph of some size k. Below are
examples of complete subgraphs for k=3 and k=4.
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
new community graph. We would then draw edges connecting nodes that
k-1 node in the original graph. This would yield the following
community graph (where the node label's are a concatenation of the
original graph's node labels):

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
# Command line options
# Graph definition file
