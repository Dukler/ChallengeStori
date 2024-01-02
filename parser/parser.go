package parser

import (
	"encoding/csv"
	"io"
	"log/slog"
	"sync"

	"github.com/Dukler/ChallengeStori/model"
	"github.com/Dukler/ChallengeStori/resources"
	"github.com/Dukler/ChallengeStori/storage"
	"github.com/google/uuid"
)

type Parser struct {
	executionId *uuid.UUID
	store       storage.Storage
	sum         *Summarizer
	l           *slog.Logger
}

func New(executionId *uuid.UUID, store storage.Storage, summarizer *Summarizer, logger *slog.Logger) *Parser {
	return &Parser{
		executionId: executionId,
		store:       store,
		sum:         summarizer,
		l:           logger,
	}
}

func (p *Parser) Run() (chan struct{}, error) {
	file := resources.Get("txns.csv")
	defer file.Close()
	reader := csv.NewReader(file)
	// Read and discard the header
	_, err := reader.Read()
	if err != nil {
		panic(err)
	}
	records := 0
	var wg sync.WaitGroup
	done := make(chan struct{})
	// Iterate through the records
	for {
		record, err := reader.Read()
		if err == io.EOF {
			p.sum.setBufferSize(records)
			break
		}
		if err != nil {
			p.l.Error(err.Error())
		}
		vr := new(validRecord)
		// Validating every value of the record before further processing.
		err = vr.validateRecord(record)
		if err != nil {
			p.l.Error(err.Error())
			//skip the record as it's not valid
			continue
		}
		records++
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := p.store.CreateTransaction(p.validRecordToTx(vr))
			if err != nil {
				p.l.Info(err.Error())
			}
		}()
		go p.sum.writeChannel(vr)
	}
	go func() {
		wg.Wait()
		close(done)
	}()

	return done, nil
}

func (p *Parser) validRecordToTx(vr *validRecord) *model.Transaction {
	return &model.Transaction{
		ExecutionId: p.executionId,
		ExternalId:  vr.Id,
		Value:       vr.Value,
		Date:        vr.Date,
	}
}
