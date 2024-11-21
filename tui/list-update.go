package tui

import (
	backend "github.com/Doctor-Scott/launcher/backend"
	C "github.com/Doctor-Scott/launcher/globalConstants"
	"strconv"

	"github.com/spf13/viper"
)

func (m model) setSelectedScriptsInView() model {
	if len(m.chain) == 0 {
		return deselectAllItems(m)
	}
	for i, listItem := range m.lists.scripts.Items() {
		if item, ok := listItem.(item); ok {
			for _, chainScript := range m.chain {
				if item.script.Name == chainScript.Name {
					m.lists.scripts.SetItem(i, selectItem(m, item))
					break
				} else {
					m.lists.scripts.SetItem(i, deselectItem(item))
				}
			}
		}
	}

	return m
}

func generateFailedItemView(m model) model {
	for i, listItem := range m.lists.scripts.Items() {
		if item, ok := listItem.(item); ok {
			if item.script.Name == m.lastFaildScriptName {
				item.failed = true
				m.lists.scripts.SetItem(i, item)
			}
		}
	}

	return m
}

func selectItem(m model, item item) item {
	item.selected = true
	indexes := findScriptIndexes(m.chain, item.script)
	desc := generatePositionString(indexes, len(m.chain))
	item.desc = desc
	return item
}

func deselectItem(item item) item {
	item.selected = false
	if item.title != "Input" {
		item.desc = ""
	} else {
		item.desc = C.INPUT_COMMAND_DESC
	}
	return item

}

func deselectAllItems(m model) model {
	for i, listItem := range m.lists.scripts.Items() {
		if item, ok := listItem.(item); ok {
			m.lists.scripts.SetItem(i, deselectItem(item))
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
	separator := viper.GetString(C.SelectedScriptDescriptionConfig.ChainSeparator.Name)
	desc := viper.GetString(C.SelectedScriptDescriptionConfig.Prefix.Name) + strconv.Itoa(indexes[0]+1)
	if len(indexes) != 1 {

		for i, index := range indexes {
			if i == 0 {
				continue
			}
			if viper.GetBool(C.SelectedScriptDescriptionConfig.UseAnd.Name) {
				if i != len(indexes)-1 {
					desc += separator
				} else {
					desc += " and "
				}
			} else {
				desc += separator
			}

			desc += strconv.Itoa(index + 1)
		}
	}
	desc += viper.GetString(C.SelectedScriptDescriptionConfig.ChainTotalSeparator.Name) + strconv.Itoa(chainLength)
	return desc

}

func maybeSetLastFailedScript(m model, scriptResult backend.ScriptResult) model {
	if !scriptResult.Success {
		m.lastFaildScriptName = scriptResult.Script.Name
		m = generateFailedItemView(m)

		return m
	}
	m.lastFaildScriptName = ""
	return m
}
