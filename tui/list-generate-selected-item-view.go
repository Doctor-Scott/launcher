package tui

import (
	backend "launcher/backend"
	C "launcher/globalConstants"
	"strconv"
)

func generateSelectedItemView(m model) model {
	if len(m.chain) == 0 {
		return deselectAllItems(m)
	}
	for i, listItem := range m.list.Items() {
		if item, ok := listItem.(item); ok {
			for _, chainScript := range m.chain {
				if item.script.Name == chainScript.Name {
					m.list.SetItem(i, selectItem(m, item))
					break
				} else {
					m.list.SetItem(i, deselectItem(item))
				}
			}
		}
	}

	return m
}

func selectItem(m model, item item) item {
	item.script.Selected = true
	indexes := findScriptIndexes(m.chain, item.script)
	desc := generatePositionString(indexes, len(m.chain))
	item.desc = desc
	return item
}

func deselectItem(item item) item {
	item.script.Selected = false
	if item.title != "Input" {
		item.desc = ""
	} else {
		item.desc = C.INPUT_SCRIPT_DESC
	}
	return item

}

func deselectAllItems(m model) model {
	for i, listItem := range m.list.Items() {
		if item, ok := listItem.(item); ok {
			m.list.SetItem(i, deselectItem(item))
		}
	}
	return m
}

func findScriptIndexes(chain backend.Chain, script backend.Script) []int {
	indexes := []int{}
	for i, chainScript := range chain {
		if chainScript.Name == script.Name {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func generatePositionString(indexes []int, chainLength int) string {
	desc := "Position: " + strconv.Itoa(indexes[0]+1)
	if len(indexes) != 1 {

		for i, index := range indexes {
			if i == 0 {
				continue
			}
			if C.USE_AND_IN_DESC {
				if i != len(indexes)-1 {
					desc += ", "
				} else {
					desc += " and "
				}
			} else {
				desc += ", "
			}

			desc += strconv.Itoa(index + 1)
		}
	}
	desc += " of " + strconv.Itoa(chainLength)
	return desc

}
