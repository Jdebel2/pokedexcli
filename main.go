package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Jdebel2/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	Next     string
	Previous string
	Area     string
	Target   string
	Pokedex  map[string]PokemonResult
}

var commands map[string]cliCommand
var cache pokecache.Cache

func main() {
	config := Config{}
	config.Next = "https://pokeapi.co/api/v2/location-area/"
	config.Previous = ""
	config.Pokedex = map[string]PokemonResult{}
	cache = *pokecache.NewCache(5 * time.Minute)

	commands = map[string]cliCommand{
		"map": {
			name:        "map",
			description: "Displays next 20 map location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 map location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore an area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect caught pokemon",
			callback:    commandInspect,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()
		cleanedInput := cleanInput(userInput)
		commandString := cleanedInput[0]
		command, ok := commands[commandString]
		if ok {
			if commandString == "explore" {
				config.Area = cleanedInput[1]
			}
			if commandString == "catch" {
				config.Target = cleanedInput[1]
			}
			if commandString == "inspect" {
				config.Target = cleanedInput[1]
			}
			command.callback(&config)
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func cleanInput(text string) []string {
	lowerString := strings.ToLower(text)
	return strings.Fields(lowerString)
}

func commandMap(config *Config) error {
	val, ok := cache.Get(config.Next)
	if ok {
		var mapResults []MapResult
		if err := json.Unmarshal(val, &mapResults); err != nil {
			log.Fatal(err)
		}
		for _, v := range mapResults {
			fmt.Println(v.Name)
		}
		return nil
	}
	cacheKey := config.Next
	content, err := json.Marshal(getMapResults(config))
	if err != nil {
		log.Fatal(err)
	}
	cache.Add(cacheKey, content)
	return nil
}

func commandMapb(config *Config) error {
	val, ok := cache.Get(config.Previous)
	if ok {
		var mapResults []MapResult
		if err := json.Unmarshal(val, &mapResults); err != nil {
			log.Fatal(err)
		}
		for _, v := range mapResults {
			fmt.Println(v.Name)
		}
		return nil
	}
	cacheKey := config.Previous
	content, err := json.Marshal(getMapbResults(config))
	if err != nil {
		log.Fatal(err)
	}
	cache.Add(cacheKey, content)
	return nil
}

func commandExplore(config *Config) error {
	val, ok := cache.Get(config.Area)
	if ok {
		var pokemonResult []PokemonEncounterResult
		if err := json.Unmarshal(val, &pokemonResult); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Exploring %s...\n", config.Area)
		fmt.Println("Found Pokemon:")
		for _, v := range pokemonResult {
			fmt.Printf("- %s\n", v.Pokemon.Name)
		}
		return nil
	}
	cacheKey := config.Area
	content, err := json.Marshal(getExploreResults(config))
	if err != nil {
		log.Fatal(err)
	}
	cache.Add(cacheKey, content)
	return nil
}

func commandCatch(config *Config) error {
	var p PokemonResult
	val, ok := cache.Get(config.Target)
	if ok {
		var pokemonResult PokemonResult
		if err := json.Unmarshal(val, &pokemonResult); err != nil {
			log.Fatal(err)
		}
		p = pokemonResult
		catchPokemon(config, p)
		return nil
	}
	cacheKey := config.Target
	p = getPokemonByName(config)
	content, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	cache.Add(cacheKey, content)
	catchPokemon(config, p)
	return nil
}

func commandInspect(config *Config) error {
	inspectPokemon(config)
	return nil
}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, value := range commands {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	return nil
}
