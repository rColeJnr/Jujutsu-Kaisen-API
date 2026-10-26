package handler

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rcolejnr/jjk-api/internal/config"
	"github.com/rcolejnr/jjk-api/internal/jwt"
	lg "github.com/rcolejnr/jjk-api/internal/logger"
	"github.com/rcolejnr/jjk-api/internal/middleware"
	"github.com/rcolejnr/jjk-api/internal/repository"
	"github.com/rcolejnr/jjk-api/internal/service"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

const (
	// <-- CHARACTER -->
	// CharacterID is the character identifier
	CharacterID = "id"

	// CharacterRoute is When http.MethodPOST the path to create a new character.
	// when http.MethodGet the path to get a list of characters
	CharacterRoute = "/api/character"

	// CharacterByID is the path to interact with a specific character.
	// Support http.Mehotds Update, Delete & Get
	CharacterByID = "/api/character/{" + CharacterID + "}"

	// <-- CURSED TECHNIQUE -->
	// CursedTechniqueID is the cursed technique identifier
	CursedTechniqueID = "id"

	// CursedTechniqueRoute is When http.MethodPOST the path to create a new cursed technique.
	// when http.MethodGet the path to get a list of cursed techniques.
	CursedTechniqueRoute = "/api/cursedtechnique"

	// CursedTechniqueByID is the path to interact with a specific CursedTechnique.
	// Support http.Mehotds Update & Get
	CursedTechniqueByID = "/api/cursedtechnique/{" + CursedTechniqueID + "}"

	// <-- DOMAIN EXPANSION -->
	// DomainExpansionID is the domain expansion identifier
	DomainExpansionID = "id"

	// DomainExpansionRoute is When http.MethodPOST the path to create a new DomainExpansion.
	// when http.MethodGet the path to get a list of DomainExpansions.
	DomainExpansionRoute = "/api/domainexpansion"

	// DomainExpansionByID is the path to interact with a specific DomainExpansion.
	// Support http.Mehotds Update & Get
	DomainExpansionByID = "/api/domainexpansion/{" + DomainExpansionID + "}"

	// <-- EPISODE -->
	// EpisodeID is the episode identifier
	EpisodeID = "id"

	// EpisodeRoute is When http.MethodPOST the path to create a new episode.
	// when http.MethodGet the path to get a list of episodes.
	EpisodeRoute = "/api/episode"

	// EpisodeByID is the path to interact with a specific episode.
	// Supports http.Mehotds Update & Get
	EpisodeByID = "/api/episode/{" + EpisodeID + "}"

	// <-- TERMINOLOGY -->
	// TerminologyID is the terminology identifier
	TerminologyID = "id"

	// TerminologyRoute is When http.MethodPOST the path to create a new Terminology.
	// when http.MethodGet the path to get a list of Terminology.
	TerminologyRoute = "/api/terminology"

	// TerminologyByID is the path to interact with a specific Terminology.
	// Supports http.Mehotds Update & Get
	TerminologyByID = "/api/terminology/{" + TerminologyID + "}"
)

// Mux creates a multiplexer with al the endpoints for our services.
func Mux(db *sql.DB, logger *slog.Logger, cfg *config.Config) http.Handler {
	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.ExpireHours)

	apiLimiter := middleware.NewRateLimiter(rate.Limit(cfg.RateLimit.RPS), cfg.RateLimit.Burst, cfg.RateLimit.TrustHeader)
	authLimiter := middleware.NewRateLimiter(rate.Limit(cfg.RateLimit.AuthPerMin/60), cfg.RateLimit.AuthBurst, cfg.RateLimit.TrustHeader)

	aRepo := repository.NewAuthRepository(db)
	cRepo := repository.NewCharacterRepository(db)
	ctRepo := repository.NewCursedTechniqueRepository(db)
	dRepo := repository.NewDomainExpansionRepository(db)
	eRepo := repository.NewEpisodeRepository(db)
	tRepo := repository.NewTerminologyRepository(db)

	aSvc := service.NewAuthService(aRepo, jwtManager)
	cSvc := service.NewCharacterService(cRepo)
	ctSvc := service.NewCursedTechniqueService(cRepo, ctRepo)
	dSvc := service.NewDomainExpansionService(cRepo, dRepo)
	eSvc := service.NewEpiesodeService(eRepo)
	tSvc := service.NewTerminologyService(tRepo)

	aHandler := NewAuthHandler(aSvc)
	cHandler := NewCharacterHandler(cSvc)
	ctHandler := NewCursedTechniqueHandler(ctSvc)
	dHandler := NewDomainExpansioneHandler(dSvc)
	eHandler := NewEpisodeHandler(eSvc)
	tHandler := NewTerminologyHandler(tSvc)

	mux := http.NewServeMux()

	// <-- Auth -->
	mux.HandleFunc(http.MethodPost+" "+"/api/auth/register", authLimiter.Middleware(aHandler.Register))
	mux.HandleFunc(http.MethodPost+" "+"/api/auth/login", authLimiter.Middleware(aHandler.Login))
	mux.HandleFunc(http.MethodPatch+" "+"/api/auth/update", authLimiter.Middleware(jwtManager.Middleware(aHandler.Update)))

	// <-- Character -->
	mux.HandleFunc(http.MethodGet+" "+CharacterRoute, cHandler.FindAll)
	mux.HandleFunc(http.MethodGet+" "+CharacterByID, cHandler.FindByID)
	mux.HandleFunc(http.MethodPost+" "+CharacterRoute, jwtManager.Middleware(cHandler.Insert))
	mux.HandleFunc(http.MethodPatch+" "+CharacterByID, jwtManager.Middleware(cHandler.Update))
	mux.HandleFunc(http.MethodDelete+" "+CharacterByID, jwtManager.Middleware(cHandler.Delete))

	// <-- Cursed Technique -->
	mux.HandleFunc(http.MethodGet+" "+CursedTechniqueRoute, ctHandler.FindAll)
	mux.HandleFunc(http.MethodGet+" "+CursedTechniqueByID, ctHandler.FindByID)
	mux.HandleFunc(http.MethodPost+" "+CursedTechniqueRoute, jwtManager.Middleware(ctHandler.Insert))
	mux.HandleFunc(http.MethodPatch+" "+CursedTechniqueByID, jwtManager.Middleware(ctHandler.Update))

	// <-- Domain Expansion -->
	mux.HandleFunc(http.MethodGet+" "+DomainExpansionRoute, dHandler.FindAll)
	mux.HandleFunc(http.MethodGet+" "+DomainExpansionByID, dHandler.FindByID)
	mux.HandleFunc(http.MethodPost+" "+DomainExpansionRoute, jwtManager.Middleware(dHandler.Insert))
	mux.HandleFunc(http.MethodPatch+" "+DomainExpansionByID, jwtManager.Middleware(dHandler.Update))

	// <-- Episode -->
	mux.HandleFunc(http.MethodGet+" "+EpisodeRoute, eHandler.FindAll)
	mux.HandleFunc(http.MethodGet+" "+EpisodeByID, eHandler.FindByID)
	mux.HandleFunc(http.MethodPost+" "+EpisodeRoute, jwtManager.Middleware(eHandler.Insert))
	mux.HandleFunc(http.MethodPatch+" "+EpisodeByID, jwtManager.Middleware(eHandler.Update))
	mux.HandleFunc(http.MethodDelete+" "+EpisodeByID, jwtManager.Middleware(eHandler.Delete))

	// <-- Terminology -->
	mux.HandleFunc(http.MethodGet+" "+TerminologyRoute, tHandler.FindAll)
	mux.HandleFunc(http.MethodGet+" "+TerminologyByID, tHandler.FindByID)
	mux.HandleFunc(http.MethodPost+" "+TerminologyRoute, jwtManager.Middleware(tHandler.Insert))
	mux.HandleFunc(http.MethodPatch+" "+TerminologyByID, jwtManager.Middleware(tHandler.Update))

	l := apiLimiter.Middleware(middleware.RequestID(lg.Logger(logger)(mux)))

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{fmt.Sprintf("http://localhost:%s", cfg.Server.Port)},
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
			middleware.HeaderRequestID,
		},
		ExposedHeaders:   []string{middleware.HeaderRequestID, "Retry-After"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return corsMiddleware.Handler(l)
}
