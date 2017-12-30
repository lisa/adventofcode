package main

/* Day 7 part A:

Determine a tree structure given an ascii representation:

pbga (66)
xhth (57)
ebii (61)
havc (66)
ktlj (57)
fwft (72) -> ktlj, cntj, xhth
qoyq (66)
padx (45) -> pbga, havc, qoyq
tknk (41) -> ugml, padx, fwft
jptl (61)
ugml (68) -> gyxo, ebii, jptl
gyxo (61)
cntj (57)

tknk is at the bottom with children ugml, padx, fwft. ugml has children gyxo, ebii, jptl. padx has children pbga, havc, qoyq. fwft has children ktlj, cntj, xhth. The outmost children have no children.
Weights of each node are in parenthesis following its declaration. At this time they appear to be unused.

With the given input programmatically determine what node is the base holding up everything else.

Part B:

For any node with children, each of that node's children forms a sub-tree. Each of those sub-trees are supposed to be the same weight, or the node itself isn't balanced. The weight of a tower is the sum of the weights of the nodes in that tower.

In the example above, this means that for ugml's disc to be balanced, gyxo, ebii, and jptl must all have the same weight, and they do: 61.

However, for tknk to be balanced, each of its child nodes and all of its grandchildren must each match. This means that the following sums must all be the same:

ugml + (gyxo + ebii + jptl) = 68 + (61 + 61 + 61) = 251
padx + (pbga + havc + qoyq) = 45 + (66 + 66 + 66) = 243
fwft + (ktlj + cntj + xhth) = 72 + (57 + 57 + 57) = 243

ugml is unbalancing this which means that ugml itself has the incorrect weight. To correct the weight, subtract the difference (8) from ugml's weight.

Exactly one node has the incorrect weight: Identify the unbalanced node and determine what its weight should be to restore balance to the tower?

*/

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var inputFile = flag.String("inputFile", "./inputs/day07-example.txt", "Input file")
var partB = flag.Bool("partB", false, "Perform part B solution?")

type Node struct {
	Name       string
	Children   []*Node
	Weight     int
	TreeWeight int
	Parent     *Node
}

func NewNode(name string, weight int) *Node {
	return &Node{
		Name:       name,
		Weight:     weight,
		TreeWeight: weight,
	}
}

/* Return both the parent and child because both need updating at the caller */
func (n *Node) AddChild(c *Node) (*Node, *Node) {
	n.Children = append(n.Children, c)
	c.Parent = n
	return n, c
}

// Compute the weights of all trees under n
func (n *Node) ComputeTreeWeight() (*Node, int) {
	var sum, subtreeSum int
	sum = n.Weight // Start with the sum as the current node's weight
	for _, node := range n.Children {
		node, subtreeSum = node.ComputeTreeWeight()
		sum += subtreeSum
	}
	n.TreeWeight = sum
	return n, n.TreeWeight
}

// Return a node's sibling
func (n *Node) GetSibling() *Node {
	var ret *Node
	for _, childNode := range n.Parent.Children {
		if childNode != n {
			ret = childNode
		}
	}
	return ret
}

// Return true if the node is balanced
// Return false if it is not along with the offending node causing unbalance and its offset relative to the correct weights
func (n *Node) IsTreeBalanced() (bool, *Node, int) {
	sum := 0
	for _, childNode := range n.Children {
		sum += childNode.TreeWeight
	}
	// Histogram will be weight => [Nodes with this weight]
	weightHistogram := make(map[int][]*Node)
	for _, child := range n.Children {
		weightHistogram[child.TreeWeight] = append(weightHistogram[child.TreeWeight], child)
	}
	for _, weight := range weightHistogram {
		if len(weight) == 1 {
			return false, weight[0], weight[0].TreeWeight - weight[0].GetSibling().TreeWeight
		}
	}
	return true, nil, 0
}
func main() {
	flag.Parse()
	input, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Couldn't read file: %s\n", err)
		os.Exit(1)
	}
	defer input.Close()

	// Create a map of name -> Node for quicker access, especially when building
	tower := make(map[string]*Node)

	// Map of parents => children. Keys are names of nodes whose children are the values
	children := make(map[string][]string)

	lineReader := bufio.NewScanner(input)
	for lineReader.Scan() {
		// loop over tokens separated by spaces
		line := lineReader.Text()
		var nodeName string
		var nodeWeight int

		for n, token := range strings.Split(line, " ") {
			switch {
			case n == 0:
				// name
				nodeName = token
			case n == 1:
				// weight
				fmt.Sscanf(token, "(%d)", &nodeWeight)
			case n == 2:
				// ->
				continue
			case n > 2:
				// list of children whose name may end in ,
				childName := strings.TrimSuffix(token, ",")
				children[nodeName] = append(children[nodeName], childName)
			} // end switch, which means we have all the fields we need to make a new node
			tower[nodeName] = NewNode(nodeName, nodeWeight)
		} //we've read every line
	} // EOF
	// Go through and turn references to names into pointers
	for parentName, childNames := range children {
		// name => list of name's children
		for _, child := range childNames {
			tower[parentName], tower[child] = tower[parentName].AddChild(tower[child])
		}
	}
	// at this point the tower has parentage, so, the only Node without a parent is the base.
	var rootNode *Node
	for _, node := range tower {
		if node.Parent == nil {
			rootNode = node
		}
	}

	rootNode, _ = rootNode.ComputeTreeWeight()

	if *partB {
		balanced, offender, offset := rootNode.IsTreeBalanced()

		var lastOffender *Node
		var lastOffset int
		workingNode := offender
		for !balanced {
			lastOffset = offset
			lastOffender = offender
			balanced, offender, offset = workingNode.IsTreeBalanced()
			if !balanced {
				workingNode = offender
			}
		}
		if lastOffender != nil {

			fmt.Printf("%s is not balanced! Adjust its weight by %d (should be %d)\n", lastOffender.Name, -1*lastOffset, lastOffender.Weight+(-1*lastOffset))
		}

	} else {

		fmt.Printf("%s is the base\n", rootNode.Name)
	}

}
