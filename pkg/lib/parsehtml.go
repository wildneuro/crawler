package lib

import (
	"errors"

	"golang.org/x/net/html"
)

var (
	// ErrNoElementNodes ...
	ErrNoElementNodes = errors.New("Element node not found")
)

// GetElementNode ...
func GetElementNode(node *html.Node, tag string) (find html.Node, err error) {
	if node.Type == html.ElementNode && node.Data == tag {
		return *node, nil
	}
	for subNode := node.FirstChild; subNode != nil; subNode = subNode.NextSibling {
		find, err = GetElementNode(subNode, tag)
		if err == ErrNoElementNodes {
			continue
		}
		return
	}
	return find, ErrNoElementNodes
}

// GetElementsNode ...
func GetElementsNode(node *html.Node, tag string) (find []html.Node, err error) {
	if node.Type == html.ElementNode && node.Data == tag {
		find = append(find, *node)
	}
	for subNode := node.FirstChild; subNode != nil; subNode = subNode.NextSibling {
		findOther, errOther := GetElementsNode(subNode, tag)
		if errOther == ErrNoElementNodes {
			continue
		}
		find = append(find, findOther...)
	}
	if len(find) == 0 {
		return find, ErrNoElementNodes
	}
	return find, nil
}

// GetElementNodeAttrValueString ...
func GetElementNodeAttrValueString(node *html.Node, attr string) (flt string, err error) {
	for _, val := range node.Attr {
		if val.Key == attr {
			return val.Val, nil
		}
	}
	return
}
