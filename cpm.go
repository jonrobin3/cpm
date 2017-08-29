//
// BUILD INSTRUCTIONS:
//     go build cpm.go
//
// RUN INSTRUCTIONS:
//     go cpm
//
// PARAMETERS:
// If you want to change the clique size (k), simply change the k
// variable in the main and recompile. I did not have time to make k
// a command line argument.
//
// THEORY OF OPERATION
//    1- first find all cliques of size k in the graph
//    2- then create graph where nodes are cliques of size k
//    3- add edges if two nodes (cliques) share k-1 common nodes
//    4- each connected component is a community
//
// MODEL GRAPH
// Below is the graph that I used for a model while developing the
// clique percolation method (CPM) module. It is sometimes
// referenced in the comments as the Model Graph in order to make
// things more clear. This is the graph that is built up in the main
// function, but obviously the code should work with any graph. 
//
//
//   +----+           +----+
//   | v2 |-----------| v1 |
//   +----++        +-+----+
//         |        |       
//         |        |       
//         +-+----+-+       
//     +-----| v3 |------+  
//     |     +----+      |  
//     |                 |  
//     |                 |  
//     |                 |  
//  +----+            +----+
//  | v4 |------------| v5 |
//  +--+-+     +------+-+--+
//     | |     |        |   
//     | +-----+--------++  
//     |       |        ||  
//     |       |        ||  
//  +--+-+-----+      +-+--+
//  | v6 |------------| v7 |
//  +----+            +----+
//     |                 |  
//     +----+    +-------+  
//          +----+          
//          | v8 |          
//        +-+----++         
//        |       |         
//        |       |         
// +----+-+       +-+----+  
// | v9 |-----------|v10 |  
// +----+           +----+
//


package main

import "fmt"
import "strconv"

type GraphNode struct {
    label string  // any string, but in our model case (v1, v2, ..., v10)
    neighbors []*GraphNode // records edges from this node. 
    associated_clique *Clique // required when building community
                              // graph; not required for starting
                              // graph
}

type CliqueCandidate struct {
    nodes []*GraphNode
    next *CliqueCandidate
}

type Clique struct {
    nodes []*GraphNode
    next *Clique
}

// FUNCTION: NewGraphNode
//
// DESCRIPTION: Creates a new graph node. The assoc_clique is
// required when building the community graph. Each k-clique is
// recorded as a node in the community graph. To determine whether
// or not an edge should be created between nodes in the community
// graph, we must determine if that they have k-1 nodes (vertices)
// in common. This is easier to do when there is a mapping between
// the community node and the clique that caused the node to be
// created. The clique contains a list of vertices.

func NewGraphNode (label string, assoc_clique *Clique) *GraphNode {
    new_node := new(GraphNode)
    new_node.label = label
    new_node.associated_clique = assoc_clique
    return new_node
}

// FUNCTION: AddNeighbor
//
// DESCRIPTION: Adds a neighboring vertex to the graph node's
// neighbor list because there is an edge connecting gn and n.

func (gn *GraphNode) AddNeighbor (n *GraphNode) {
    gn.neighbors = append(gn.neighbors, n)
}

// FUNCTION: PrintGraph
//
// DESCRIPTION: Prints a graph -- vertices and edges.

func PrintGraph(g []*GraphNode) {
    if g == nil {
        fmt.Printf("empty graph\n")
    }
    for _, e := range g {
        fmt.Printf("%s:  ", e.label)
        glen := len(g)
        if glen > 0 {
            for _, n := range e.neighbors {
                fmt.Printf("%s ", n.label)
            }
        }
        fmt.Printf("\n")
    }
}

// FUNCTION: GetCliqueCandidates
//
// PARAMETERS:
// - k int
// - node_list []*GraphNode
//
// DESCRIPTION: A recursive function that generates all the node
// permutations that could form a k-clique for a given
// node. node_list is essentially the neighbor list for some
// node. For example, lets say we want to generate a candidate list
// for v5 with k=3. The original node_list would be {v3, v4, v6, v7}
// (all of v5's neighbors). We remove the first node, v3, and then
// recursive call GetCliqueCandates with node_list = {v4, v6,
// v7}. Again, we remove the first node, v4, and recursively call
// with node_list equal to {v6, v7}. This is our anchor case because
// k=3 and we will add v5 later to this canidate list, thus giving
// us a {v5, v6, v7} candidate clique. When the recursive call
// begins to unwind, we add the removed node in each slot of the
// candidate list. For example, on the first unwind, we start with a
// removed node of v4 and clique_list of {v6, v7}. v4 is substitued
// for each element and this yields a clique_list of {v6, v7}, {v4,
// v7}, and {v6, v4}. And this continues recursively. For v3, we
// would then get a clique_list of:
//
//              {v6, v7}
//              {v4, v7}
//              {v6, v4}
//              {v3, v7}
//              {v6, v3}
//              {v3, v7}  -- duplicate, not added
//              {v4, v3}
//              {v3, v4}  -- duplicate, not added
//              {v6, v3}  -- duplicate, not added
//
// When the function returns, the caller can then make the final candidate
// list using the examination node (v5 in our example). This will yield
// a candidate clique list of:
//
//              {v5, v6, v7}
//              {v5, v4, v7}
//              {v5, v6, v4}
//              {v5, v3, v7}
//              {v5, v6, v3}
//              {v5, v4, v3}
//
//
func GetCliqueCandidates (k int, node_list []*GraphNode) *CliqueCandidate {

    if k < 2 {
        return nil
    }
    if len(node_list) < k - 1 {
        return nil
    }
    if len(node_list) == k - 1 {
        new_candidate := new (CliqueCandidate)
        new_candidate.nodes = node_list
        new_candidate.next = nil
        return new_candidate
    }
    
    node := node_list[0] // the removed node
    clique_list := GetCliqueCandidates(k, node_list[1:])
    var return_clique_list *CliqueCandidate = clique_list
    
    for item := clique_list; item != nil; item = item.next {
        for i, _ := range item.nodes {
            new_candidate := new (CliqueCandidate)
            new_candidate.nodes = make([]*GraphNode, len(item.nodes), len(item.nodes))
            for k, v := range item.nodes {
                new_candidate.nodes[k] = v
            }
            // new_candidate.nodes = item.nodes // copies array
            new_candidate.nodes[i] = node
            
            // only add this candidate list if doesn't already exist
            if new_candidate.IsDuplicate(return_clique_list) == false {
                
                new_candidate.next = return_clique_list
                return_clique_list = new_candidate
            }
        }
    }
    
    return return_clique_list
}

// FUNCTION: MakeCliqueList
//
// PARAMETERS:
//
// - candidate_list *CliqueCandidate
// - examination_node *GraphNode
//
// DESCRIPTION: Determines if nodes on the candidate list are
// completely connected. If the nodes are completely connected, then
// it creates a clique, which includes all candidates and the
// examination node and places them on the clique list to be returned.
//
// For example, the following candidates would be generated for node
// v5 in the Model Graph for a k = 3 clique, and this is what is
// essentially returned from GetCliqueCandidates:
//
//     - v6 v7  
//     - v4 v7  
//     - v6 v4  
//     - v3 v4  
//     - v6 v3  
//     - v3 v7  
//
// These are just candidates to form a clique (k=3) with
// v5. MakeCliqueList determines whether the above nodes do form a
// clique with v5, and if they do then MakeCliqueList places a
// Clique node on the return list. Here's what the determination should be
// for the above candidates:
//
//     - v6 v7  - yes, forms k=3 clique with v5
//     - v4 v7  - yes  "
//     - v6 v4  - yes  "
//     - v3 v4  - yes  "
//     - v6 v3  - no, does not form k=3 clique with v5 
//     - v3 v7  - no, "
//
// So, the return clique list should look like the following in the
// case of examination node v5 and the above candidate list:
//
//     +---------------------+     +--+--+--+
//     |nodes []*GraphNodes--+---->|v5|v6|v7|
//     |next *Clique         |     +--+--+--+
//     |           |         |               
//     +-----------+---------+               
//                 |                         
//                 v                         
//     +---------------------+     +--+--+--+
//     |nodes []*GraphNodes--+---->|v5|v4|v7|
//     |next *Clique         |     +--+--+--+
//     |           |         |               
//     +-----------+---------+               
//                 |                         
//                 v                         
//     +---------------------+     +--+--+--+
//     |nodes []*GraphNodes--+---->|v5|v6|v4|
//     |next *Clique         |     +--+--+--+
//     |           |         |               
//     +-----------+---------+               
//                 |                         
//                 v                         
//     +---------------------+     +--+--+--+
//     |nodes []*GraphNodes--+---->|v5|v4|v3|
//     |next *Clique         |     +--+--+--+
//     |           |         |               
//     +-----------+---------+               
//                 |                         
//                 v                         
//              +----+                      
//              |nil |                      
//              +----+                      
//
// 
func (candidate_list *CliqueCandidate) MakeCliqueList(
    examination_node *GraphNode) *Clique {

    var clique_list *Clique = nil

    for item := candidate_list; item != nil; item = item.next {
        candidate_list_is_clique := true // assumed, not yet determined
        item_nodes_len := len(item.nodes)
        for i := 0; i < item_nodes_len && candidate_list_is_clique; i++ {
            candidate_node := item.nodes[i]
            for j := i + 1; j < item_nodes_len; j++ {
                // if the candidate node is not connected to all other
                // nodes then this candidate does not form clique
                if candidate_node.IsConnected(item.nodes[j]) == false {
                    candidate_list_is_clique = false
                    break
                }
            }
        }
        if (candidate_list_is_clique == true) {
            new_clique := new (Clique)
            new_clique.nodes = make ([]*GraphNode,
                item_nodes_len + 1,
                item_nodes_len + 1)
            copy (new_clique.nodes, item.nodes)
            new_clique.nodes[item_nodes_len] = examination_node
            new_clique.next = clique_list
            clique_list = new_clique
        }
    }
    return clique_list
}

// FUNCTION: IsConnected
//
// DESCRIPTION: Determines if the candidate node (cn) is connected to
// some node (sn)

func (cn *GraphNode) IsConnected (sn *GraphNode) bool {
    is_connected := false
    for _,item := range sn.neighbors {
        if item == cn {
            is_connected = true
            break
        }
    }
    return is_connected
}

// FUNCTION: IsDuplicate
// 
// DESCRIPTION: After examining each node, we will have many
// different duplicate clique candidates (i.e., candidates that all
// have the same vertices). We need to create a candidate list that
// has no duplicates and IsDuplicate determines that.

func (cc *CliqueCandidate) IsDuplicate (clist *CliqueCandidate) bool {
    return_val := false
    
    for item := clist; item != nil; item = item.next {
        match_count := len (cc.nodes)
        for _, ccnode := range cc.nodes {
            for _, list_node := range item.nodes {
                if (list_node == ccnode) {
                    match_count--
                    if match_count == 0 {
                        return true // duplicate list found
                    }
                    break
                }
            }
        }
    }
    
    return return_val
}

// FUNCTION: NotRecorded
//
// DESCRIPTION: Determines whether nor not the clique is already on
// the clique_list. When we merge the candidate lists for each node,
// we will invariably find duplicates, but we only want one unique
// clique recorded -- not multiples.

func (clique *Clique) NotRecorded (clique_list *Clique) bool {
    
    for item := clique_list; item != nil; item = item.next {
        match_count := len(item.nodes)
        if match_count != len(clique.nodes) {
            continue
        }
        for _, node := range item.nodes {
            for _, exam_node := range clique.nodes {
                if exam_node == node {
                    match_count--
                    if match_count == 0 {
                        return false
                    }
                }
            }
        }
    }
    return true
}

// FUNCTION: AppendClique
//
// DESCRIPTION: Appends src_clique_list to dest_clique_list. Returns
// dest_clique_list
//
func AppendClique (dest_clique_list *Clique,
    src_clique_list *Clique) (*Clique) {

        if dest_clique_list == nil {
            return nil
        }

        var last_item *Clique
        for last_item = dest_clique_list;
            last_item.next != nil;
            last_item = last_item.next {
            }

        for clique := src_clique_list; clique != nil; clique = clique.next {
            if clique.NotRecorded(dest_clique_list) == true {
                new_clique := new(Clique)
                new_clique.nodes = clique.nodes
                new_clique.next = nil
                last_item.next = new_clique
                last_item = new_clique
            }
        }
        return dest_clique_list
}

// FUNCTION: CreateLabel
//
// DESCRIPTION: Generates a label for a node in the community
// graph. It does this by concatening the labels of each vertex from
// the origial graph that is in a clique to a single label name.

func CreateLabel (nodes []*GraphNode) string {
    var new_label string
    for i,item := range nodes {
        if i == 0 {
            new_label = item.label
        } else {
            new_label += "," + item.label
        }
    }
    return new_label
}

// FUNCTION: Kminu1CommonNodes
//
// DESCRIPTION: Determines if two nodes on the newly created
// community graph share k-1 vertices from the original graph. If
// so, then there should be an edge from gn and node in the
// community graph.

func (gn *GraphNode) Kminus1CommonNodes (node *GraphNode, k int) bool {
    common_node_count := 0
    return_value := false

    for _, exam_node := range gn.associated_clique.nodes {
        for _, common_node := range node.associated_clique.nodes {
            if exam_node == common_node {
                common_node_count++
                if common_node_count == k - 1 {
                    return true
                }
            }
        }
    }
    return return_value
}

// FUNCTION: AddNeighbors
//
// DESCRIPTION: Determines whether there is an edge between two
// nodes in the generated community graph. If there is, that edge is
// recorded as one of gn's neighbors.

func (gn *GraphNode) AddNeighbors (graph []*GraphNode, k int) {
    for _, node := range graph {
        if gn != node && gn.Kminus1CommonNodes(node, k) {
            gn.AddNeighbor(node)
        }
    }
}

// FUNCTION: CreateCommunityGraph
//
// DESCRIPTION: Transforms every clique created from the original
// graph into a node in the community graph. CreateCommunityGraph
// also determines if there is an edge between vertices and adds the
// appropriate neighbor nodes. For the nodes in the community graph
// to be connected, they must have k-1 vertices in
// common. CreateCommunityGraph returns a valid community graph for
// k.

func CreateCommunityGraph (clique_list *Clique, k int) []*GraphNode {
    var community_graph []*GraphNode
    if clique_list == nil {
        return nil
    }
    for item := clique_list; item != nil; item = item.next {
        label := CreateLabel (item.nodes)
        new_node := NewGraphNode(label, item)
        community_graph = append(community_graph, new_node)
    }

    for _, node := range community_graph {
        node.AddNeighbors(community_graph, k)
    }
    
    return community_graph
}


func main() {
    var graph []*GraphNode
    k := 3

    // Create the vertices
    for i := 1; i <= 10; i++ {
        label := "v" + strconv.Itoa(i)
        v := NewGraphNode(label, nil)
        graph = append(graph, v)
    }

    // Add all the edges
    
    // v1  
    graph[0].neighbors = append(graph[0].neighbors, graph[1])
    graph[0].neighbors = append(graph[0].neighbors, graph[2])

    // v2
    graph[1].neighbors = append(graph[1].neighbors, graph[0])
    graph[1].neighbors = append(graph[1].neighbors, graph[2])
    
    // v3
    graph[2].neighbors = append(graph[2].neighbors, graph[1-1])
    graph[2].neighbors = append(graph[2].neighbors, graph[2-1])
    graph[2].neighbors = append(graph[2].neighbors, graph[4-1])
    graph[2].neighbors = append(graph[2].neighbors, graph[5-1])
    
    // v4
    graph[4-1].neighbors = append(graph[4-1].neighbors, graph[3-1])
    graph[4-1].neighbors = append(graph[4-1].neighbors, graph[5-1])
    graph[4-1].neighbors = append(graph[4-1].neighbors, graph[6-1])
    graph[4-1].neighbors = append(graph[4-1].neighbors, graph[7-1])
    
    // v5
    graph[5-1].neighbors = append(graph[5-1].neighbors, graph[3-1])
    graph[5-1].neighbors = append(graph[5-1].neighbors, graph[4-1])
    graph[5-1].neighbors = append(graph[5-1].neighbors, graph[6-1])
    graph[5-1].neighbors = append(graph[5-1].neighbors, graph[7-1])
    
    // v6
    graph[6-1].neighbors = append(graph[6-1].neighbors, graph[4-1])
    graph[6-1].neighbors = append(graph[6-1].neighbors, graph[5-1])
    graph[6-1].neighbors = append(graph[6-1].neighbors, graph[7-1])
    graph[6-1].neighbors = append(graph[6-1].neighbors, graph[8-1])
    
    // v7
    graph[7-1].neighbors = append(graph[7-1].neighbors, graph[4-1])
    graph[7-1].neighbors = append(graph[7-1].neighbors, graph[5-1])
    graph[7-1].neighbors = append(graph[7-1].neighbors, graph[6-1])
    graph[7-1].neighbors = append(graph[7-1].neighbors, graph[8-1])
    
    // v8
    graph[8-1].neighbors = append(graph[8-1].neighbors, graph[6-1])
    graph[8-1].neighbors = append(graph[8-1].neighbors, graph[7-1])
    graph[8-1].neighbors = append(graph[8-1].neighbors, graph[9-1])
    graph[8-1].neighbors = append(graph[8-1].neighbors, graph[10-1])
    
    // v9
    graph[9-1].neighbors = append(graph[9-1].neighbors, graph[8-1])
    graph[9-1].neighbors = append(graph[9-1].neighbors, graph[10-1])
    
    // v10
    graph[10-1].neighbors = append(graph[10-1].neighbors, graph[8-1])
    graph[10-1].neighbors = append(graph[10-1].neighbors, graph[9-1])

    fmt.Printf("The original graph:\n")
    fmt.Printf("-------------------\n")
    PrintGraph(graph)

    var clique_list *Clique = nil
    for i := 1; i <= 10; i++ {
        candidate_list := GetCliqueCandidates(k, graph[i-1].neighbors)
        if candidate_list != nil {
            temp_clique_list := candidate_list.MakeCliqueList(graph[i-1])
            if clique_list == nil {
                clique_list = temp_clique_list
            } else {
                clique_list = AppendClique(clique_list, temp_clique_list)
            }
        }
    }
 
    community_graph := CreateCommunityGraph(clique_list, k)
    fmt.Printf("Community graph:\n")
    fmt.Printf("----------------\n")
    PrintGraph(community_graph)
}
