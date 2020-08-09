package pkg

import (
	"fmt"
	"github.com/adityasidharta/gombal/internal"
	"github.com/mroth/weightedrand"
	"math/rand"
	"time"
)

const (
	debugMessage = "Papa tell me your secret"
	fullDebugMessage = "Mama tell me your secret"
	RANDOM = 0
	WEIGHTED = 1
	MAXIMUM = 2
)

type query string
type response map[string]int
type dataset map[query]response

type Bot struct {
	ds       dataset
	strategy int
}

func NewBot (strategy int) *Bot {
	initQuery := "Hi There!"
	initResponse := make(map[string]int)
	initDataset := make(map[query]response)

	initResponse["Hello! How are you?"] = 1
	initDataset[query(initQuery)] = initResponse

	return &Bot {
		ds:       initDataset,
		strategy: strategy,
	}
}

func (bot *Bot) getMaximum(response map[string]int) (string, error) {
	if len(response) == 0 {
		return "", emptyResponseError
	}

	var maxVal string
	var maxWeight int

	for k, v := range response {
		if v > maxWeight {
			maxVal = k
			maxWeight = v
		}
	}
	return maxVal, nil
}

func (bot *Bot) getWeighted(response map[string]int) (string, error) {
	rand.Seed(time.Now().UTC().UnixNano())

	if len(response) == 0 {
		return "", emptyResponseError
	}

	choices := make([]weightedrand.Choice, 0, len(response))

	for k, v := range response {
		choices = append(choices, weightedrand.Choice{Item: k, Weight: uint(v)})
	}

	c := weightedrand.NewChooser(
        choices...
    )

    result := c.Pick().(string)
    return result, nil
}


func (bot *Bot) getRandom(response map[string]int) (string, error) {
	rand.Seed(time.Now().UTC().UnixNano())

	if len(response) == 0 {
		return "", emptyResponseError
	}

	choices := make([]weightedrand.Choice, 0, len(response))

	for k := range response {
		choices = append(choices, weightedrand.Choice{Item: k, Weight: 1})
	}

	c := weightedrand.NewChooser(
        choices...
    )

    result := c.Pick().(string)
    return result, nil
}


func (bot *Bot) contains(q string) bool {
	_, ok := bot.ds[query(q)]
	return ok
}

func (bot *Bot) Debug() string {
	result := ""
	for q, rDict := range bot.ds {
		result = result + fmt.Sprintf("%v :", q) + "\n"
		for k := range rDict {
			result = result + fmt.Sprintf("  - %v", k) + "\n"
		}
	}
	return result
}

func (bot *Bot) FullDebug() string {
	result := ""
	for q, rDict := range bot.ds {
		result = result + fmt.Sprintf("%v :", q) + "\n"
		for k, v := range rDict {
			result = result + fmt.Sprintf("  - %v : %v", k, v) + "\n"
		}
	}
	return result
}


func (bot *Bot) LenQueries() int {
	return len(bot.ds)
}

func (bot *Bot) LenResponses(q string) (int, error) {
	respDict, ok := bot.ds[query(q)]
	if !ok {
		return 0, invalidQueryError
	}
	return len(respDict), nil
}

func (bot *Bot) ShowQueries() []string {
	queries := make([]string, 0, len(bot.ds))
	for k := range bot.ds {
		queries = append(queries, string(k))
	}
	return queries
}

func (bot *Bot) ShowResponses(q string) ([]string, error) {
	respDict, ok := bot.ds[query(q)]
	if !ok {
		return make([]string, 0, 0), invalidQueryError
	}
	responses := make([]string, 0, len(respDict))
	for k := range respDict {
		responses = append(responses, k)
	}
	return responses, nil
}

func (bot *Bot) Add(q string, r string) {
	if q == fullDebugMessage {
		return
	}
	if q == debugMessage {
		return
	}

	respDict, ok := bot.ds[query(q)]
	if !ok {
		newRespDict := make(response)
		newRespDict[r] = 1
		bot.ds[query(q)] = newRespDict
	} else {
		_, ok = respDict[r]
		if !ok {
			bot.ds[query(q)][r] = 1
		} else {
			bot.ds[query(q)][r] = bot.ds[query(q)][r] + 1
		}
	}
}

func (bot *Bot) Adds(qs []string, rs []string) error {
	tuples, err := internal.Zip(qs, rs)
	if err != nil {
		return err
	}
	for _, tuple := range tuples {
		bot.Add(tuple.Q, tuple.R)
	}
	return nil
}

func (bot *Bot) RemoveQuery(q string) error {
	_, ok := bot.ds[query(q)]
	if !ok {
		return invalidQueryError
	}
	delete(bot.ds, query(q))
	return nil
}

func (bot *Bot) RemoveQueries(qs []string) error {
	for _, q := range qs {
		err := bot.RemoveQuery(q)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bot *Bot) RemoveResponse(q string, r string) error {
	respDict, ok := bot.ds[query(q)]
	if !ok {
		return invalidQueryError
	}
	_, ok = respDict[r]
	if !ok {
		return invalidResponseError
	}
	delete(bot.ds[query(q)], r)
	return nil
}

func (bot *Bot) RemoveResponses(qs []string, rs []string) error {
	tuples, err := internal.Zip(qs, rs)
	if err != nil {
		return err
	}
	for _, tuple := range tuples {
		err := bot.RemoveResponse(tuple.Q, tuple.R)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bot *Bot) Get(q string) (string, error) {
	if q == debugMessage {
		return bot.Debug(), nil
	}
	if q == fullDebugMessage {
		return bot.FullDebug(), nil
	}

	if bot.contains(q) {
		response, ok := bot.ds[query(q)]
		if !ok {
			return "", invalidQueryError
		}
		if len(response) == 0 {
			return "", emptyResponseError
		}
		switch strategy := bot.strategy; strategy {
		case RANDOM:
			return bot.getRandom(response)
		case WEIGHTED:
			return bot.getWeighted(response)
		case MAXIMUM:
			return bot.getMaximum(response)
		default:
			return "", invalidStrategyError
		}
	} else {
		return q, nil
	}
}
