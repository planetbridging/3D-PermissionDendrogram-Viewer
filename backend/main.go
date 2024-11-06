package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Node represents a single node in the dendrogram.
type Node struct {
	Name        string  `json:"name"`
	Children    []*Node `json:"children,omitempty"`
	OverlapLevel int    `json:"overlapLevel"` // 0 for default, 1 for purple, 2 for yellow, 3 for red
	hasParent   bool    // Internal flag to mark if node has a parent
	hasChildren bool    // Internal flag to mark if node has children
}

// GenerateDendrogram creates a combined dendrogram with marked overlapping and hierarchical nodes
func GenerateDendrogram(users []string) *Node {
	root := &Node{Name: "Permissions", OverlapLevel: 0}
	nodeMap := make(map[string]*Node)

	for _, user := range users {
		userNode := generateUserPermissions(user)
		mergeNodes(root, userNode, nodeMap)
	}
	assignColors(root)
	return root
}

// Generate permissions for each user with a guaranteed hierarchical structure
func generateUserPermissions(user string) *Node {
	root := &Node{Name: user, OverlapLevel: 0}

	// Define unique groups for each user to ensure hierarchical nodes for yellow
	uniqueGroup := &Node{Name: "OU=UniqueGroup-" + user, OverlapLevel: 0, hasParent: true}
	hierarchicalPermission := &Node{Name: "Permission-Unique-" + user, OverlapLevel: 0}

	uniqueGroup.Children = append(uniqueGroup.Children, hierarchicalPermission)
	root.Children = append(root.Children, uniqueGroup)

	// Add overlapping groups to ensure purple and potentially red
	overlapGroup := &Node{Name: "OU=Admins", OverlapLevel: 1} // Shared group to create overlap
	overlapPermission := &Node{Name: "Permission-Admin-Shared", OverlapLevel: 1}
	overlapGroup.Children = append(overlapGroup.Children, overlapPermission)

	root.Children = append(root.Children, overlapGroup)

	// Randomly add more groups to increase variation
	groups := []string{"OU=Finance", "OU=HR", "OU=Engineering", "OU=Product"}
	rand.Shuffle(len(groups), func(i, j int) { groups[i], groups[j] = groups[j], groups[i] })

	for _, group := range groups[:rand.Intn(3)+2] {
		groupNode := &Node{Name: group, OverlapLevel: 0}
		for i := 0; i < rand.Intn(3)+1; i++ {
			permissionNode := &Node{Name: fmt.Sprintf("Permission-%s-%d", group, i), OverlapLevel: 0}
			groupNode.Children = append(groupNode.Children, permissionNode)
		}
		root.Children = append(root.Children, groupNode)
	}
	return root
}

// Merges nodes from a user into the main dendrogram
func mergeNodes(root, userNode *Node, nodeMap map[string]*Node) {
	for _, child := range userNode.Children {
		if existingNode, found := nodeMap[child.Name]; found {
			existingNode.hasParent = true // Set as having a parent for hierarchical coloring
			incrementOverlap(existingNode) // Mark as overlapping
			for _, subChild := range child.Children {
				mergeNodes(existingNode, subChild, nodeMap)
			}
		} else {
			nodeMap[child.Name] = child
			root.Children = append(root.Children, child)
			if len(child.Children) > 0 {
				child.hasChildren = true // Mark as having children
			}
		}
	}
}

// Increment overlap level for overlapping nodes and its children
func incrementOverlap(node *Node) {
	node.OverlapLevel = 1
	for _, child := range node.Children {
		incrementOverlap(child)
	}
}

// Assign colors based on overlap level, parent-child relationship, and combined rules
func assignColors(node *Node) {
	for _, child := range node.Children {
		assignColors(child)

		// Set color based on the rules provided
		if child.OverlapLevel == 1 && child.hasParent && child.hasChildren {
			child.OverlapLevel = 3 // Red for both overlapping and hierarchical
		} else if child.OverlapLevel == 1 {
			child.OverlapLevel = 1 // Purple for overlapping
		} else if child.hasParent && child.hasChildren {
			child.OverlapLevel = 2 // Yellow for hierarchical only
		}
	}
}

var mu sync.Mutex

func main() {
	rand.Seed(time.Now().UnixNano())
	app := fiber.New()

	// WebSocket endpoint to provide the combined dendrogram for both users
	app.Get("/ws/admin", websocket.New(func(c *websocket.Conn) {
		log.Println("Admin WebSocket connection attempt")

		// Generate the combined dendrogram with example users
		mu.Lock()
		dendrogram := GenerateDendrogram([]string{"User Bob", "User Alice"})
		mu.Unlock()

		// Convert the dendrogram to JSON format
		dendrogramJSON, err := json.Marshal(dendrogram)
		if err != nil {
			log.Println("Error marshalling dendrogram:", err)
			return
		}

		// Send the dendrogram JSON to the WebSocket client
		log.Println("Sending dendrogram to client")
		if err = c.WriteMessage(websocket.TextMessage, dendrogramJSON); err != nil {
			log.Println("Error writing message:", err)
			return
		}

		// Keep the WebSocket open and listen for any further client messages if necessary
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				log.Println("WebSocket connection closed by client:", err)
				break
			}
		}
	}))

	log.Fatal(app.Listen(":8432"))
}
