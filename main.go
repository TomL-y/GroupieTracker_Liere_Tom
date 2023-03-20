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

func getPCGames() ([]Game, error) {
	resp, err := http.Get("https://www.freetogame.com/api/games?platform=pc")
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

func getWebGames() ([]Game, error) {
	resp, err := http.Get("https://www.freetogame.com/api/games?platform=browser")
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

func filterGamesByPlatform(games []Game, platform string) []Game {
	var filtered []Game
	for _, game := range games {
		if strings.Contains(strings.ToLower(game.Platform), strings.ToLower(platform)) {
			filtered = append(filtered, game)
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
	allGames, err := getGames()
	if err != nil {
		panic(err)
	}

	var selectedGenre string

	static := http.FileServer(http.Dir("src"))
	http.Handle("/src/", http.StripPrefix("/src/", static))

	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl2 := template.Must(template.ParseFiles("src/template/main.html"))
	tmpl3 := template.Must(template.ParseFiles("src/template/jeuxpc.html"))
	tmpl4 := template.Must(template.ParseFiles("src/template/jeuxweb.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, allGames)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		games := allGames
		var filtered []Game

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
			} else if selectedGenre != "" {
				filtered = filterGamesByGenre(games, selectedGenre)
			} else {
				filtered = games
			}

		} else {
			filtered = games
		}

		err := tmpl2.Execute(w, filtered)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/jeuxpc", func(w http.ResponseWriter, r *http.Request) {
		games := allGames
		var filtered []Game

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			name := r.FormValue("name")
			genre := r.FormValue("genre")

			if name != "" || genre != "" {
				filtered = filterGamesByName(games, name, genre)
				filtered = filterGamesByPlatform(filtered, "PC")
			} else {
				filtered = filterGamesByPlatform(games, "PC")
			}

		} else {
			filtered = filterGamesByPlatform(games, "PC")
		}

		err := tmpl3.Execute(w, filtered)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/jeuxweb", func(w http.ResponseWriter, r *http.Request) {
		games := allGames
		var filtered []Game

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			name := r.FormValue("name")
			genre := r.FormValue("genre")

			if name != "" || genre != "" {
				filtered = filterGamesByName(games, name, genre)
				filtered = filterGamesByPlatform(filtered, "Web")
			} else {
				filtered = filterGamesByPlatform(games, "Web")
			}

		} else {
			filtered = filterGamesByPlatform(games, "Web")
		}

		err := tmpl4.Execute(w, filtered)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.ListenAndServe(":8080", nil)
}
