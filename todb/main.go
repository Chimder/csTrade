package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func initDb() (*pgxpool.Pool, error) {
	ctx := context.Background()
	url := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("Failed to parse config DBPOOL: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database pool created successfully")
	return pool, nil

}

func main() {
	ctx := context.Background()

	pg, err := initDb()
	if err != nil {
		panic(err)
	}
	defer pg.Close()

	resp, err := http.Get("https://raw.githubusercontent.com/ByMykel/CSGO-API/main/public/api/en/skins.json")
	if err != nil {
		log.Printf("Error getting  %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		log.Printf("HTTP error %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		panic(err)
	}

	var skins []Skins
	err = json.Unmarshal(body, &skins)
	if err != nil {
		panic(err)
	}

	query := `
INSERT INTO skins (
  id, name, description, category_id, category_name, pattern_id, pattern_name, min_float, max_float, rarity_id, rarity_name, rarity_color, stattrak, souvenir, paint_index, legacy_model, image, phase, team_id, team_name
  ) VALUES (
    @id, @name, @description, @category_id, @category_name, @pattern_id, @pattern_name, @min_float, @max_float, @rarity_id, @rarity_name, @rarity_color, @stattrak, @souvenir, @paint_index, @legacy_model, @image, @phase, @team_id, @team_name
  )
	`

	skinsBatchInsert := &pgx.Batch{}

	for _, v := range skins {
		skinsBatchInsert.Queue(query, pgx.NamedArgs{
			"id":            v.ID,
			"name":          v.Name,
			"description":   v.Description,
			"category_id":   v.Category.ID,
			"category_name": v.Category.Name,
			"pattern_id":    v.Pattern.ID,
			"pattern_name":  v.Pattern.Name,
			"min_float":     v.MinFloat,
			"max_float":     v.MaxFloat,
			"rarity_id":     v.Rarity.ID,
			"rarity_name":   v.Rarity.Name,
			"rarity_color":  v.Rarity.Color,
			"stattrak":      v.Stattrak,
			"souvenir":      v.Souvenir,
			"paint_index":   v.PaintIndex,
			"legacy_model":  v.LegacyModel,
			"image":         v.Image,
			"phase":         v.Phase,
			"team_id":       v.Team.ID,
			"team_name":     v.Team.Name,
		})
	}

	res := pg.SendBatch(ctx, skinsBatchInsert)
	defer res.Close()

	for i := 0; i < skinsBatchInsert.Len(); i++ {
		_, err := res.Exec()
		if err != nil {
			log.Printf("Insert err %d: %v", i, err)
		}
	}

	log.Printf("Inserted %d skins", len(skins))
	log.Print("all done")
}

type Skins struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	Pattern struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"pattern"`
	MinFloat float64 `json:"min_float"`
	MaxFloat float64 `json:"max_float"`
	Rarity   struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Color string `json:"color"`
	} `json:"rarity"`
	Stattrak   bool   `json:"stattrak"`
	Souvenir   bool   `json:"souvenir,omitempty"`
	PaintIndex string `json:"paint_index"`
	Wears      []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"wears,omitempty"`
	Collections []interface{} `json:"collections,omitempty"`
	Crates      []struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Image string `json:"image"`
	} `json:"crates"`
	Team struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"team"`
	LegacyModel  bool   `json:"legacy_model"`
	Image        string `json:"image"`
	Phase        string `json:"phase,omitempty"`
	SpecialNotes []struct {
		Source string `json:"source"`
		Text   string `json:"text"`
	} `json:"special_notes,omitempty"`
}
