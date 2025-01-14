package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func getMapResults(config *Config) []MapResult {
	var results LocationResult
	res, err := http.Get(config.Next)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&results); err != nil {
		log.Fatal(err)
	}
	for _, v := range results.Results {
		fmt.Println(v.Name)
	}
	config.Next = results.Next
	config.Previous = results.Previous
	return results.Results
}

func getMapbResults(config *Config) []MapResult {
	if config.Previous == "" {
		fmt.Println("you're on the first page")
		return []MapResult{}
	}
	var results LocationResult
	res, err := http.Get(config.Previous)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&results); err != nil {
		log.Fatal(err)
	}
	for _, v := range results.Results {
		fmt.Println(v.Name)
	}
	config.Next = results.Next
	config.Previous = results.Previous
	return results.Results
}

func getExploreResults(config *Config) []PokemonEncounterResult {
	fullUrl := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", config.Area)
	res, err := http.Get(fullUrl)
	if err != nil {
		log.Fatal(err)
	}
	var exploreResult ExploreResult
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&exploreResult); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Exploring %s...\n", config.Area)
	fmt.Println("Found Pokemon:")
	for _, v := range exploreResult.PokemonEncounters {
		fmt.Printf("- %s\n", v.Pokemon.Name)
	}
	return exploreResult.PokemonEncounters
}

func getPokemonByName(config *Config) PokemonResult {
	fullUrl := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", config.Target)
	res, err := http.Get(fullUrl)
	if err != nil {
		log.Fatal(err)
	}
	var pokemon PokemonResult
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&pokemon); err != nil {
		log.Fatal(err)
	}
	return pokemon
}

func catchPokemon(config *Config, pokemon PokemonResult) {
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	catchIndex := rand.Intn(pokemon.BaseExperience)
	if catchIndex <= 50 {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		config.Pokedex[pokemon.Name] = pokemon
		return
	}
	fmt.Printf("%s escaped!\n", pokemon.Name)
}

func inspectPokemon(config *Config) {
	pokemon, ok := config.Pokedex[config.Target]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, v := range pokemon.Stats {
		fmt.Printf("-%s: %v\n", v.Stat.Name, v.BaseStat)
	}
	fmt.Println("Types:")
	for _, v := range pokemon.Types {
		fmt.Printf("- %s\n", v.Type.Name)
	}
}
