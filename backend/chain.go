package backend

import (
	C "github.com/Doctor-Scott/launcher/globalConstants"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Chain []Script

type ChainResult []ScriptResult

type ChainItem struct {
	Name  string
	Chain Chain
}

func GetChainStructure(path string) []ChainItem {
	if path == "" {
		path = viper.GetString(C.PathConfig.LauncherDir.Name) + "/custom/"
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

func runChain(stdin []byte, chain Chain, scriptResults []ScriptResult) ChainResult {
	if !(len(scriptResults) == 0) && !scriptResults[len(scriptResults)-1].Success {
		// we stop after the first failure
		return scriptResults
	}

	if len(chain) == 0 {
		// finished chain
		if viper.GetBool(C.ClearChainAfterRun.Name) {
			// we save the empty chain here, so when it is read again it will be cleared
			MaybeAutoSaveChain(chain)
		}
		return scriptResults
	}
	scriptResult := RunScript(chain[0], stdin)
	return runChain(scriptResult.Stdout, chain[1:], append(scriptResults, scriptResult))
}

func RunChain(stdin []byte, chain Chain) ChainResult {
	return runChain(stdin, chain, []ScriptResult{})
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
		shouldRemoveInput := chain[i].Name == C.INPUT_COMMAND_NAME && scriptToRemove.Name == C.INPUT_COMMAND_NAME
		if shouldRemoveScript || shouldRemoveInput {
			//pop the item
			return MaybeAutoSaveChain(append(chain[0:i], chain[i+1:]...))
		}
	}
	// item not found in chain, so just return the chain
	return MaybeAutoSaveChain(chain)
}

func DeleteChainConfig(name string) {
	path := viper.GetString(C.PathConfig.LauncherDir.Name) + "/custom/" + name + ".json"
	err := DeleteFile(path)
	if err != nil {
		log.Fatal(err)
	}

}

func MaybeAutoSaveChain(chain Chain) Chain {
	if viper.GetBool(C.Autosave.Name) {
		err := Save(viper.GetString(C.PathConfig.LauncherDir.Name)+"/"+C.CHAIN_AUTOSAVE_FILE_NAME, chain)
		if err != nil {
			log.Fatal(err)
		}
	}
	return chain
}

func ClearAutoSave() Chain {
	chain := Chain{}

	err := Save(viper.GetString(C.PathConfig.LauncherDir.Name)+"/"+C.CHAIN_AUTOSAVE_FILE_NAME, chain)
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
	Load(viper.GetString(C.PathConfig.LauncherDir.Name)+"/"+C.CHAIN_AUTOSAVE_FILE_NAME, &chainConfig)
	return chainConfig
}
