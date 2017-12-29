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
*/

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var inputFile = flag.String("inputFile", "./inputs/day07-example.txt", "Input file")

// Just refer by name to next nodes and not pointers. It's crude, but effective.
type Node struct {
	Children []string
	Parent   string
	Weight   int
	Name     string
}

func (n Node) SetParent(parent string) Node {
	n.Parent = parent
	return n
}
func (n Node) AddChildNode(child string) Node {
	n.Children = append(n.Children, child)
	return n
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
	tower := make(map[string]Node)

	lineReader := bufio.NewScanner(input)
	for lineReader.Scan() {
		// loop over tokens separated by spaces
		line := lineReader.Text()
		node := Node{}
		for n, token := range strings.Split(line, " ") {
			switch {
			case n == 0:
				// name
				node.Name = token
			case n == 1:
				// weight
				var weight int
				fmt.Sscanf(token, "(%d)", &weight)
				node.Weight = weight
			case n == 2:
				// ->
				continue
			case n > 2:
				// list of children whose name may end in ,
				name := strings.TrimSuffix(token, ",")
				node.Children = append(node.Children, name)
			} // end switch
			tower[node.Name] = node
		} //we've read every line
	} // EOF

	// Go through and turn references to names into pointers
	for name, node := range tower {
		for _, childName := range node.Children {
			tower[childName] = tower[childName].SetParent(name)
		}
	}
	// at this point the tower has parentage, so, the only Node without a parent is the base.
	for _, node := range tower {
		if node.Parent == "" {
			fmt.Printf("%s is the base\n", node.Name)
		}
	}

}
