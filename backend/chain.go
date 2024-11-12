package backend

import (
	C "launcher/globalConstants"
	"log"
	"strings"
)

type Chain []Script

type ChainItem struct {
	Name  string
	Chain Chain
}

func GetChainStructure() []ChainItem {
	path := ResolvePath("~") + ".launcher/custom/"
	files := getFiles(path)
	chainItems := []ChainItem{}
	for _, file := range files {
		name := strings.TrimSuffix(file, ".json")
		chain := LoadCustomChain(name)
		chainItems = append(chainItems, ChainItem{Name: name, Chain: chain})
	}
	return chainItems
}

func RunChain(stdin []byte, chain Chain) []byte {
	if len(chain) == 0 {
		if C.CLEAR_CHAIN_AFTER_RUN {
			MaybeAutoSaveChain(chain)
		}
		return stdin
	}
	stdout := RunScript(chain[0], stdin)
	return RunChain(stdout, chain[1:])

}

func AddScriptToChain(scriptToAdd Script, chain Chain) Chain {
	return MaybeAutoSaveChain(append(chain, scriptToAdd))
}

func AddChainToChain(chainToAdd Chain, chain Chain) Chain {
	return MaybeAutoSaveChain(append(chain, chainToAdd...))
}

func RemoveScriptFromChain(scriptToRemove Script, chain Chain) Chain {
	for i := len(chain) - 1; i >= 0; i-- {
		shouldRemoveScript := chain[i].Name == scriptToRemove.Name && chain[i].Path == scriptToRemove.Path
		shouldRemoveInput := chain[i].Name == C.INPUT_SCRIPT_NAME && scriptToRemove.Name == C.INPUT_SCRIPT_NAME
		if shouldRemoveScript || shouldRemoveInput {
			//pop the item
			return MaybeAutoSaveChain(append(chain[0:i], chain[i+1:]...))
		}
	}
	// item not found in chain, so just return the chain
	return MaybeAutoSaveChain(chain)
}

func MaybeAutoSaveChain(chain Chain) Chain {
	if C.AUTO_SAVE {
		err := Save(ResolvePath("~")+"/"+C.CHAIN_SAVE_FILE, chain)
		if err != nil {
			log.Fatal(err)
		}
	}
	return chain
}

func ClearAutoSave() Chain {
	chain := Chain{}

	err := Save(ResolvePath("~")+"/"+C.CHAIN_SAVE_FILE, chain)
	if err != nil {
		log.Fatal(err)
	}
	return chain
}

func SaveCustomChain(chain Chain, name string) Chain {
	err := Save(ResolvePath("~")+".launcher/custom/"+name+".json", chain)
	if err != nil {
		log.Fatal(err)
	}
	return chain

}

func LoadCustomChain(name string) Chain {
	var chain Chain
	Load(ResolvePath("~")+".launcher/custom/"+name+".json", &chain)
	return chain
}

func ReadChainConfig() Chain {
	// Safely handle chain configuration
	var chainConfig Chain
	Load(ResolvePath("~")+"/"+C.CHAIN_SAVE_FILE, &chainConfig)
	return chainConfig
}
