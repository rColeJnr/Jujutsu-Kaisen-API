package handler

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	lg "github.com/rcolejnr/jjk-api/internal/logger"
	"github.com/rcolejnr/jjk-api/internal/repository"
	"github.com/rcolejnr/jjk-api/internal/service"
	"github.com/rs/cors"
)

// Mux creates a multiplexer with al the endpoints for our services.
func Mux(db *sql.DB, logger *slog.Logger, port string) http.Handler {
	cRepo := repository.NewCharacterRepository(db)
	ctRepo := repository.NewCursedTechniqueRepository(db)
	dRepo := repository.NewDomainExpansionRepository(db)
	eRepo := repository.NewEpisodeRepository(db)
	tRepo := repository.NewTerminologyRepository(db)

	cSvc := service.NewCharacterService(cRepo)
	ctSvc := service.NewCursedTechniqueService(cRepo, ctRepo)
	dSvc := service.NewDomainExpansionService(cRepo, dRepo)
	eSvc := service.NewEpiesodeService(eRepo)
	tSvc := service.NewTerminologyService(tRepo)

	cHandler := NewCharacterHandler(cSvc)
	ctHandler := NewCursedTechniqueHandler(ctSvc)
	dHandler := NewDomainExpansioneHandler(dSvc)
	eHandler := NewEpisodeHandler(eSvc)
	tHandler := NewTerminologyHandler(tSvc)

	mux := http.NewServeMux()
	// <-- Character -->
	mux.HandleFunc(http.MethodGet+" "+CharacterRoute, cHandler.FindAll)
	mux.HandleFunc(http.MethodGet+" "+CharacterByID, cHandler.FindByID)
	mux.HandleFunc(http.MethodPost+" "+CharacterRoute, cHandler.Insert)
	mux.HandleFunc(http.MethodPatch+" "+CharacterByID, cHandler.Update)
	mux.HandleFunc(http.MethodDelete+" "+CharacterByID, cHandler.Delete)

	// <-- Cursed Technique -->
	mux.HandleFunc(http.MethodGet+" "+CursedTechniqueRoute, ctHandler.FindAll)
	mux.HandleFunc(http.MethodGet+" "+CursedTechniqueByID, ctHandler.FindByID)
	mux.HandleFunc(http.MethodPost+" "+CursedTechniqueRoute, ctHandler.Insert)
	mux.HandleFunc(http.MethodPatch+" "+CursedTechniqueByID, ctHandler.Update)

	// <-- Domain Expansion -->
	mux.HandleFunc(http.MethodGet+" "+DomainExpansionRoute, dHandler.FindAll)
	mux.HandleFunc(http.MethodGet+" "+DomainExpansionByID, dHandler.FindByID)
	mux.HandleFunc(http.MethodPost+" "+DomainExpansionRoute, dHandler.Insert)
	mux.HandleFunc(http.MethodPatch+" "+DomainExpansionByID, dHandler.Update)

	// <-- Episode -->
	mux.HandleFunc(http.MethodGet+" "+EpisodeRoute, eHandler.FindAll)
	mux.HandleFunc(http.MethodGet+" "+EpisodeByID, eHandler.FindByID)
	mux.HandleFunc(http.MethodPost+" "+EpisodeRoute, eHandler.Insert)
	mux.HandleFunc(http.MethodPatch+" "+EpisodeByID, eHandler.Update)
	mux.HandleFunc(http.MethodDelete+" "+EpisodeByID, eHandler.Delete)

	// <-- Terminology -->
	mux.HandleFunc(http.MethodGet+" "+TerminologyRoute, tHandler.FindAll)
	mux.HandleFunc(http.MethodGet+" "+TerminologyByID, tHandler.FindByID)
	mux.HandleFunc(http.MethodPost+" "+TerminologyRoute, tHandler.Insert)
	mux.HandleFunc(http.MethodPatch+" "+TerminologyByID, tHandler.Update)

	l := lg.Logger(logger)(mux)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{fmt.Sprintf("http://localhost:%d", port)},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return corsMiddleware.Handler(l)
}

const (
	// <-- CHARACTER -->
	// CharacterID is the character'a identifier
	CharacterID = "id"

	// CharacterRoute is When http.MethodPut the path to create a new character.
	// when http.MethodGet the path to get a list of characters
	CharacterRoute = "/api/character"

	// CharacterByID is the path to interact with a specific character.
	// Support http.Mehotds Update, Delete & Get
	CharacterByID = "/api/character/{" + CharacterID + "}"

	// <-- CURSED TECHNIQUE -->
	// CursedTechniqueID is the cursed technique identifier
	CursedTechniqueID = "id"

	// CursedTechniqueRoute is When http.MethodPut the path to create a new cursed technique.
	// when http.MethodGet the path to get a list of cursed techniques.
	CursedTechniqueRoute = "/api/cursedtechnique"

	// CursedTechniqueByID is the path to interact with a specific CursedTechnique.
	// Support http.Mehotds Update & Get
	CursedTechniqueByID = "/api/cursedtechnique/{" + CursedTechniqueID + "}"

	// <-- DOMAIN EXPANSION -->
	// DomainExpansionID is the domain expansion identifier
	DomainExpansionID = "id"

	// DomainExpansionRoute is When http.MethodPut the path to create a new DomainExpansion.
	// when http.MethodGet the path to get a list of DomainExpansions.
	DomainExpansionRoute = "/api/domainexpansion"

	// DomainExpansionByID is the path to interact with a specific DomainExpansion.
	// Support http.Mehotds Update & Get
	DomainExpansionByID = "/api/domainexpansion/{" + DomainExpansionID + "}"

	// <-- EPISODE -->
	// EpisodeID is the episode identifier
	EpisodeID = "id"

	// EpisodeRoute is When http.MethodPut the path to create a new episode.
	// when http.MethodGet the path to get a list of episodes.
	EpisodeRoute = "/api/episode"

	// EpisodeByID is the path to interact with a specific episode.
	// Supports http.Mehotds Update & Get
	EpisodeByID = "/api/episode/{" + EpisodeID + "}"

	// <-- TERMINOLOGY -->
	// TerminologyID is the terminology identifier
	TerminologyID = "id"

	// TerminologyRoute is When http.MethodPut the path to create a new Terminology.
	// when http.MethodGet the path to get a list of Terminology.
	TerminologyRoute = "/api/terminology"

	// TerminologyByID is the path to interact with a specific Terminology.
	// Supports http.Mehotds Update & Get
	TerminologyByID = "/api/terminology/{" + TerminologyID + "}"
)
