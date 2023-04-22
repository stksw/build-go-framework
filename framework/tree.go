package framework

import (
	"fmt"
	"strings"
)

type TreeNode struct {
	children []*TreeNode
	handler  func(ctx *HttpContext)
	param    string
	parent   *TreeNode
}

func Constructor() *TreeNode {
	return &TreeNode{
		param:    "",
		children: []*TreeNode{},
	}
}

func isGeneral(param string) bool {
	return strings.HasPrefix(param, ":")
}

// nodeに/区切りでpathを登録
func (t *TreeNode) Insert(pathname string, handler func(ctx *HttpContext)) {
	node := t
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
				parent:   node,
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

func (t *TreeNode) findChild(param string) *TreeNode {
	for _, child := range t.children {
		if child.param == param {
			return child
		}
	}
	// 子ノードのparamが一致しなければ見つからないとしてnilを返す
	return nil
}

// 既に登録済のpathなら、そのノードを返す
func (t *TreeNode) Search(pathname string) *TreeNode {
	params := strings.Split(pathname, "/")

	// 未登録のルーティングなら、nil
	// 登録済のルーティングなら、そのノードを返す
	return dfs(t, params)
}

// 深さ優先探索
func dfs(node *TreeNode, params []string) *TreeNode {
	currentParam := params[0]
	isLastParam := len(params) == 1

	// 子ノードを探索
	for _, child := range node.children {
		// 探索するparamが最後のpathか
		if isLastParam {
			// /list/:idが登録された状態で/list/nameのルーティングを探索すると
			// 必ずここを通る
			if isGeneral(child.param) {
				return child
			}

			// 最後のpathが一致したら子ノードを返す
			if child.param == currentParam {
				return child
			}
			// 次の子ノードへ
			continue
		}

		// ノードの中が:idではなく、paramも一致しないなら次の子ノードへ
		if !isGeneral(child.param) && child.param != currentParam {
			continue
		}

		//
		result := dfs(child, params[1:])
		if result != nil {
			return result
		}
	}
	return nil
}

func (t *TreeNode) ParseParams(pathname string) map[string]string {
	node := t
	pathname = strings.TrimSuffix(pathname, "/")
	params := strings.Split(pathname, "/")
	paramDict := make(map[string]string)

	for i := len(params) - 1; i >= 0; i-- {
		if isGeneral(node.param) {
			paramDict[node.param] = params[i]
		}
		node = node.parent
	}

	fmt.Println("list", params)
	return paramDict
}
