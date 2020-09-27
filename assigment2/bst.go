package main

import "fmt"

//BinaryNode to store patient node
type BinaryNode struct {
	patient *Patient
	left    *BinaryNode
	right   *BinaryNode
}

//BST to store patient in alphabetical order
type BST struct {
	root *BinaryNode
}

func (bst *BST) insertNode(t **BinaryNode, patient *Patient) error {

	if *t == nil {
		newNode := &BinaryNode{
			patient: patient,
			left:    nil,
			right:   nil,
		}
		*t = newNode
		return nil
	}

	if patient.name < (*t).patient.name {
		bst.insertNode(&((*t).left), patient)
	} else {
		bst.insertNode(&((*t).right), patient)
	}

	return nil
}
func (bst *BST) insert(patient *Patient) {
	bst.insertNode(&bst.root, patient)

}

func (bst *BST) inOrderTraverse(t *BinaryNode) {
	if t != nil {
		bst.inOrderTraverse(t.left)

		if t.patient.name != "admin" {
			fmt.Printf("%-10v %v ", t.patient.name, t.patient.hasBooking)
			if t.patient.hasBooking == true {
				fmt.Printf("%v slot%v %v\n", t.patient.timeslot.date, t.patient.timeslot.slot, t.patient.timeslot.doctor)
			} else {
				fmt.Println("")
			}
		}

		bst.inOrderTraverse(t.right)
	}
}

func (bst *BST) inOrder() {
	fmt.Println("\nusername   hasBooking date slot   doctor")
	bst.inOrderTraverse(bst.root)
}

func (bst *BST) searchNode(t *BinaryNode, name string) *Patient {
	if t == nil {
		return nil
	}
	if t.patient.name == name {
		return t.patient
	}
	if name < t.patient.name {
		return bst.searchNode(t.left, name)
	}
	return bst.searchNode(t.right, name)

}

func (bst *BST) search(name string) *Patient {
	return bst.searchNode(bst.root, name)
}
