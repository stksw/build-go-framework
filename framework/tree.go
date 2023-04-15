package framework

import (
	"net/http"
	"strings"
)

type TreeNode struct {
	children []*TreeNode
	handler  func(w http.ResponseWriter, r *http.Request)
	param    string
}

func Constructor() TreeNode {
	return TreeNode{
		param:    "",
		children: []*TreeNode{},
	}
}

func (this *TreeNode) Insert(pathname string, handler func(w http.ResponseWriter, r *http.Request)) {
	node := this
	params := strings.Split(pathname, "/")

	// 該当のノードに到達するか、新規ノードを作るまでループ
	for _, param := range params {
		// 子ノードが合致するか確認
		child := node.findChild(param)

		if child == nil {
			// なければ、新規に子ノードをつくる
			child = &TreeNode{
				param:    param,
				children: []*TreeNode{},
			}
			node.children = append(node.children, child)
		}
		// 子ノードを現在のノードに更新
		node = child
	}

	// 作成したノードにhandlerを代入
	// もし同じノードがあった場合、上書き
	node.handler = handler
}

func (this *TreeNode) findChild(param string) *TreeNode {
	for _, child := range this.children {
		if child.param == param {
			return child
		}
	}
	// 子ノードのparamが一致しなければ見つからないとしてnilを返す
	return nil
}

// 一致するpathのhandlerを返す
func (this *TreeNode) Search(pathname string) func(w http.ResponseWriter, r *http.Request) {
	params := strings.Split(pathname, "/")
	node := this
	for _, param := range params {
		child := node.findChild(param)
		if child == nil {
			return nil
		}
		node = child
	}

	// 到達したノードのhandlerを返す
	return node.handler
}
