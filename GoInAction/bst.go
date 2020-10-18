package main

//BinaryNode to store user node
type BinaryNode struct {
	user  *User
	left  *BinaryNode
	right *BinaryNode
}

//BST to store user in alphabetical order
type BST struct {
	root *BinaryNode
}

func (bst *BST) insertNode(t **BinaryNode, user *User) error {

	if *t == nil {
		newNode := &BinaryNode{
			user:  user,
			left:  nil,
			right: nil,
		}
		*t = newNode
		return nil
	}

	if user.Name < (*t).user.Name {
		bst.insertNode(&((*t).left), user)
	} else {
		bst.insertNode(&((*t).right), user)
	}

	return nil
}
func (bst *BST) insert(user *User) {
	bst.insertNode(&bst.root, user)

}

func (bst *BST) inOrderTraverse(t *BinaryNode, userMap *map[string]Timeslot) {
	if t != nil {
		bst.inOrderTraverse(t.left, userMap)

		if t.user.Name != "admin" {
			if t.user.HasBooking == true {
				(*userMap)[t.user.Name] = *(t.user.Timeslot)
			} else {
				(*userMap)[t.user.Name] = Timeslot{}
			}
		}

		bst.inOrderTraverse(t.right, userMap)
	}
}

func (bst *BST) inOrder() map[string]Timeslot {
	userMap := &(map[string]Timeslot{})
	bst.inOrderTraverse(bst.root, userMap)
	return *userMap
}

func (bst *BST) searchNode(t *BinaryNode, name string) *User {
	if t == nil {
		return nil
	}
	if t.user.Name == name {
		return t.user
	}
	if name < t.user.Name {
		return bst.searchNode(t.left, name)
	}
	return bst.searchNode(t.right, name)

}

func (bst *BST) search(name string) *User {
	return bst.searchNode(bst.root, name)
}
