package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type LocationResult struct {
	Next     string      `json:"next"`
	Previous string      `json:"previous"`
	Results  []MapResult `json:"results"`
	Count    int         `json:"count"`
}

type ExploreResult struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []PokemonEncounterResult `json:"pokemon_encounters"`
}

type MapResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonEncounterResult struct {
	Pokemon struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"pokemon"`
	VersionDetails []struct {
		EncounterDetails []struct {
			Chance          int   `json:"chance"`
			ConditionValues []any `json:"condition_values"`
			MaxLevel        int   `json:"max_level"`
			Method          struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"method"`
			MinLevel int `json:"min_level"`
		} `json:"encounter_details"`
		MaxChance int `json:"max_chance"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"version_details"`
}

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
