package main

import (
	"log/slog"

	"github.com/Dukler/ChallengeStori/email"
	"github.com/Dukler/ChallengeStori/parser"
	"github.com/Dukler/ChallengeStori/storage"
	"github.com/google/uuid"
)

func main() {
	logger := slog.Default().With(
		slog.String("app", "StoriChallenge"),
		slog.String("version", "1.0"),
	)

	es := email.NewEmailService()
	store, err := storage.NewPostgresStore(10, 1, logger)
	if err != nil {
		panic(err)
	}
	err = store.Init()
	if err != nil {
		panic(err)
	}

	executionId := uuid.New()
	summarizer := parser.NewSummarizer(&executionId, store, logger)
	summarizerDone := summarizer.Run()

	p := parser.New(&executionId, store, summarizer, logger)

	parserDone, err := p.Run()
	if err != nil {
		panic(err)
	}

	<-parserDone
	<-summarizerDone
	es.SendSummary(summarizer)
}
