package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
)

type Game struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Thumbnail     string `json:"thumbnail"`
	ShortDesc     string `json:"short_description"`
	Description   string `json:"description"`
	ReleaseDate   string `json:"release_date"`
	Developer     string `json:"developer"`
	Publisher     string `json:"publisher"`
	Genre         string `json:"genre"`
	Platform      string `json:"platform"`
	GameURL       string `json:"game_url"`
	MinimumSystem string `json:"minimum_system_requirements"`
}

func getGames() ([]Game, error) {
	resp, err := http.Get("https://www.freetogame.com/api/games")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var games []Game
	err = json.NewDecoder(resp.Body).Decode(&games)
	if err != nil {
		return nil, err
	}

	return games, nil
}

func filterGamesByName(games []Game, name string, genre string) []Game {
	var filtered []Game
	for _, game := range games {
		if strings.Contains(strings.ToLower(game.Title), strings.ToLower(name)) {
			if genre == "all" || strings.Contains(strings.ToLower(game.Genre), strings.ToLower(genre)) {
				filtered = append(filtered, game)

			}
		}
	}
	return filtered
}

func filterGamesByGenre(games []Game, genre string) []Game {
	var filtered []Game
	for _, game := range games {
		if genre == "all" {
			return games
		} else if genre == "" || strings.Contains(strings.ToLower(game.Genre), strings.ToLower(genre)) {
			filtered = append(filtered, game)
		}
	}
	return filtered
}

func main() {
	games, err := getGames()
	if err != nil {
		panic(err)
	}

	var selectedGenre string

	static := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", static))

	tmpl := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var filtered []Game
		var filteredApplied bool

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			name := r.FormValue("name")
			genre := r.FormValue("genre")

			if genre != "" {
				selectedGenre = genre
			}

			if name != "" {
				filtered = filterGamesByName(games, name, selectedGenre)
				filteredApplied = true
			} else if selectedGenre != "" {
				filtered = filterGamesByGenre(games, selectedGenre)
				filteredApplied = true
			} else {
				filtered = games
				filteredApplied = true
			}

			err = tmpl.Execute(w, filtered)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			err := tmpl.Execute(w, games)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if filteredApplied {
			games = filtered
		} else {
			games, err = getGames()
			if err != nil {
				panic(err)
			}
		}
	})

	http.ListenAndServe(":8080", nil)
}
