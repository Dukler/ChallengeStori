package parser

import (
	"log/slog"

	"github.com/Dukler/ChallengeStori/model"
	"github.com/Dukler/ChallengeStori/storage"
	"github.com/google/uuid"
)

type Summarizer struct {
	TxnsByMonth   map[string]int
	Balance       int64
	AvgCredit     uint64
	AvgDebit      uint64
	totalCredit   uint64
	totalDebit    uint64
	txCreditCount uint
	txDebitCount  uint
	inputChannel  chan *validRecord
	executionId   *uuid.UUID
	store         storage.Storage
	bufferSize    int
	isParserDone  bool
	l             *slog.Logger
}

func NewSummarizer(executionId *uuid.UUID, store storage.Storage, logger *slog.Logger) *Summarizer {
	sum := new(Summarizer)
	sum.TxnsByMonth = make(map[string]int)
	sum.inputChannel = make(chan *validRecord)
	sum.store = store
	sum.executionId = executionId
	sum.isParserDone = false
	sum.l = logger
	return sum
}

func (s *Summarizer) Run() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		currentBufferSize := 0
		for {
			select {
			case validRecord, ok := <-s.inputChannel:
				if !ok {
					break
				}
				currentBufferSize++
				s.summaraize(validRecord)
				if currentBufferSize == s.bufferSize && s.isParserDone {
					s.calculateAverages()
					_, err := s.store.CreateSummary(s.summarizerToSummary())
					if err != nil {
						s.l.Info(err.Error())
					}
					return
				}
			}
		}
	}()
	return done
}

func (s *Summarizer) writeChannel(vr *validRecord) {
	s.inputChannel <- vr
}
func (s *Summarizer) closeChannel() {
	close(s.inputChannel)
}

func (s *Summarizer) summaraize(record *validRecord) {
	month := record.Date.Month().String()
	s.TxnsByMonth[month]++
	s.Balance += record.Value
	if record.Value < 0 {
		s.totalDebit += uint64(record.Value * -1)
		s.txDebitCount++
	} else {
		s.totalCredit += uint64(record.Value)
		s.txCreditCount++
	}
}

func (s *Summarizer) calculateAverages() {
	if s.txCreditCount > 0 {
		s.AvgCredit = s.totalCredit / uint64(s.txCreditCount)
	}
	if s.txDebitCount > 0 {
		s.AvgDebit = s.totalDebit / uint64(s.txDebitCount)
	}
}

func (sum *Summarizer) summarizerToSummary() *model.Summary {
	return &model.Summary{
		ExecutionId:   sum.executionId,
		AverageCredit: sum.AvgCredit,
		AverageDebit:  sum.AvgDebit,
		Balance:       sum.Balance,
		TxnsByMonth:   sum.TxnsByMonth,
	}
}

func (sum *Summarizer) setBufferSize(size int) {
	sum.bufferSize = size
	sum.isParserDone = true
}
