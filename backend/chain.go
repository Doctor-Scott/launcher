package backend

import (
	C "launcher/globalConstants"
	"log"
)

type Chain []Script

func RunChain(stdin []byte, chain Chain) []byte {
	if len(chain) == 0 {
		if C.CLEAR_CHAIN_AFTER_RUN && C.AUTO_SAVE {
			SaveChain(chain)
		}
		return stdin
	}
	stdout := RunScript(chain[0], stdin)
	return RunChain(stdout, chain[1:])

}

func AddScriptToChain(scriptToAdd Script, chain Chain) Chain {
	if C.AUTO_SAVE {
		return SaveChain(append(chain, scriptToAdd))
	}
	return append(chain, scriptToAdd)
}

func RemoveScriptFromChain(scriptToRemove Script, chain Chain) Chain {
	for i := len(chain) - 1; i >= 0; i-- {
		shouldRemoveScript := chain[i].Name == scriptToRemove.Name && chain[i].Path == scriptToRemove.Path
		shouldRemoveInput := chain[i].Name == C.INPUT_SCRIPT_NAME && scriptToRemove.Name == C.INPUT_SCRIPT_NAME
		if shouldRemoveScript || shouldRemoveInput {
			//pop the item
			if C.AUTO_SAVE {
				return SaveChain(append(chain[0:i], chain[i+1:]...))
			}
			return append(chain[0:i], chain[i+1:]...)
		}
	}
	// item not found in chain, so just return the chain
	if C.AUTO_SAVE {
		return SaveChain(chain)
	}
	return chain
}

func SaveChain(chain Chain) Chain {
	err := Save(ResolvePath("~")+"/"+C.CHAIN_SAVE_FILE, chain)
	if err != nil {
		log.Fatal(err)
	}
	return chain
}

func ReadChainConfig() []Script {
	// Safely handle chain configuration
	var chainConfig []Script
	Load(ResolvePath("~")+"/"+C.CHAIN_SAVE_FILE, &chainConfig)
	return chainConfig
}
