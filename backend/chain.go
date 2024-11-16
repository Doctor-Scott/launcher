package backend

import (
	C "launcher/globalConstants"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Chain []Script

type ChainItem struct {
	Name  string
	Chain Chain
}

func GetChainStructure(path string) []ChainItem {
	if path == "" {
		path = viper.GetString(C.LauncherDir.Name) + "/custom/"
	}
	files := getFiles(path)
	chainItems := []ChainItem{}
	for _, file := range files {
		name := strings.TrimSuffix(file, ".json")
		chain := LoadCustomChain(path, name)
		chainItems = append(chainItems, ChainItem{Name: name, Chain: chain})
	}
	return chainItems
}

func RunChain(stdin []byte, chain Chain) []byte {
	if len(chain) == 0 {
		if viper.GetBool("clearChainAfterRun") {
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
		shouldRemoveScript := chain[i].Name == scriptToRemove.Name && chain[i].Command == scriptToRemove.Command
		shouldRemoveInput := chain[i].Name == C.INPUT_SCRIPT_NAME && scriptToRemove.Name == C.INPUT_SCRIPT_NAME
		if shouldRemoveScript || shouldRemoveInput {
			//pop the item
			return MaybeAutoSaveChain(append(chain[0:i], chain[i+1:]...))
		}
	}
	// item not found in chain, so just return the chain
	return MaybeAutoSaveChain(chain)
}

func DeleteChainConfig(name string) {
	path := viper.GetString(C.LauncherDir.Name) + "/custom/" + name + ".json"
	err := DeleteFile(path)
	if err != nil {
		log.Fatal(err)
	}

}

func MaybeAutoSaveChain(chain Chain) Chain {
	if viper.GetBool("autosave") {
		err := Save(viper.GetString(C.LauncherDir.Name)+"/"+C.CHAIN_AUTOSAVE_FILE_NAME, chain)
		if err != nil {
			log.Fatal(err)
		}
	}
	return chain
}

func ClearAutoSave() Chain {
	chain := Chain{}

	err := Save(viper.GetString(C.LauncherDir.Name)+"/"+C.CHAIN_AUTOSAVE_FILE_NAME, chain)
	if err != nil {
		log.Fatal(err)
	}
	return chain
}

func SaveCustomChain(chain Chain, path string, name string) Chain {
	err := Save(path+name+".json", chain)
	if err != nil {
		log.Fatal(err)
	}
	return chain

}

func LoadCustomChain(path string, name string) Chain {
	var chain Chain
	Load(path+name+".json", &chain)
	return chain
}

func ReadChainConfig() Chain {
	// Safely handle chain configuration
	var chainConfig Chain
	Load(viper.GetString(C.LauncherDir.Name)+"/"+C.CHAIN_AUTOSAVE_FILE_NAME, &chainConfig)
	return chainConfig
}
