/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/

// NOTE
// Basically I have a bunch of little scripts that I use from time to time.
// I want to create a cli tool 'launchpad' in go that i can use as either a cli tool with cobra or tui application with charm.
// I could even have outputs go to an editor to then pass back into the launchpad for another script
// I imagine I could create a lot of it with cobra, and then just have the tui run those commands too?
// you could use the editor funciton to write a command first in the editor and then run it, giving you more flexibility in writing long commands

//NOTE
// will most likely need to have a load function, that will generate the available script structure
// or a service running and polling in the background
// do we need to start thinking about databases?
// if we want to allow saving of workflows, maybe it makes sense, that or a flat file, but do we need more than that?
// can still utalise the comments idea

// TODO
// we want a script struct that has a name and args
// NOTE
// maybe add a comment at the top of the scripts that go could read
// then we could build the script structure that would allow for x args to be added
// then it would pass that to the script
// that would be cool

// type script struct {
// 	name string
// 	args []string
// }
//
// var test = script{
// 	name: "felix",
// 	args: []string{"1", "2"},
// }

// func main() {
// TODO
// get the scripts
// also start modularising this
// the getting should be done by a specific module

//NOTE
// https://github.com/spf13/cobra
// you could use this as your entry point maybe?
// so you can still run from the command line, with autosuggest maybe?
// default could bring you to the charm tui
// this could be a great start to a full launchpad
// could even maybe create a ui for it, maybe from a browser?

//NOTE
// go feels like a good choice for this
// rust would be faster, but i do not need speed for this, go is no snail either
// its just a launchpad
// bash could also work, but I feel go would would better with a more complex app
// and its fun to learn

//NOTE
// from initial inspection, this also does not seem to exist yet
// so thats fun

// }

package main

import "launcher/cmd"

// import "launcher/testing"

func main() {
	// testing.Main()

	cmd.Execute()
}
