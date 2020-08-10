package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/mroth/weightedrand"
	"io/ioutil"
	"math/rand"
	"time"
	log "github.com/sirupsen/logrus"
)

// Bot Constants
const (
	debugMessage = "Papa tell me your secret"
	fullDebugMessage = "Mama tell me your secret"
)

// query is the data structure for the Chat query
type query string

// response is the data structure for the Chat response. It includes the possible responses and the weight of the responses
type response map[string]int

// dataset is the data structure that unifies the learned queries and responses
type dataset map[query]response

// Bot is the data structure of the ChatBot.
type Bot struct {
	ds       dataset
	strategy string
}

// NewBot creates New Bot with empty learned phrases and chosen Strategy
func NewBot (strategy string, path string) (*Bot, error) {
	initQuery := "Hi There!"
	initResponse := make(map[string]int)
	initDataset := make(map[query]response)
	initResponse["Hello! How are you?"] = 1
	initDataset[query(initQuery)] = initResponse

	bot := &Bot {
		ds:       initDataset,
		strategy: strategy,
	}

	ok := bot.ValidStrategy()
	if !ok {
		return bot, invalidStrategyError
	}

	if path == "" {
		return bot, nil
	} else {
		err := bot.Load(path)
		return bot, err
	}
}

// getMaximum chooses the response with the highest weight, given multiple responses
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

// getWeighted chooses the response randomly according to its weights, given multiple responses
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

// getRandom chooses one response randomly, given multiple responses
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


// contains check whether the query exists within the dataset
func (bot *Bot) contains(q string) bool {
	_, ok := bot.ds[query(q)]
	return ok
}

// Debug displays the learned queries and possible responses in a nice format
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

// FullDebug displays the learned queries, possible responses and its weights in a nice format
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

// LenQueries get the total number of learned queries by the ChatBot
func (bot *Bot) LenQueries() int {
	return len(bot.ds)
}

// LenResponses get the total number of learned responses for a query by the ChatBot
func (bot *Bot) LenResponses(q string) (int, error) {
	respDict, ok := bot.ds[query(q)]
	if !ok {
		return 0, invalidQueryError
	}
	return len(respDict), nil
}

// ShowQueries shows all of the learned queries by the ChatBot
func (bot *Bot) ShowQueries() []string {
	queries := make([]string, 0, len(bot.ds))
	for k := range bot.ds {
		queries = append(queries, string(k))
	}
	return queries
}

// ShowResponses shows all of the learned responses by the ChatBot
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

// Add appends new phrase (query and response) into the ChatBot database
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

// Adds appends multiple phrases (queries, responses) into the ChatBot Database
func (bot *Bot) Adds(qs []string, rs []string) error {
	tuples, err := Zip(qs, rs)
	if err != nil {
		return err
	}
	for _, tuple := range tuples {
		bot.Add(tuple.Q, tuple.R)
	}
	return nil
}

// RemoveQuery delete a specific query from the ChatBot Database
func (bot *Bot) RemoveQuery(q string) error {
	_, ok := bot.ds[query(q)]
	if !ok {
		return invalidQueryError
	}
	delete(bot.ds, query(q))
	return nil
}

// RemoveQueries delete multiple queries from the ChatBot Database
func (bot *Bot) RemoveQueries(qs []string) error {
	for _, q := range qs {
		err := bot.RemoveQuery(q)
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveResponse delete a specific response for a specific query from the DataBase
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

// RemoveResponses delete responses from the respective queries from the DataBase
func (bot *Bot) RemoveResponses(qs []string, rs []string) error {
	tuples, err := Zip(qs, rs)
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

// Get allows the ChatBot to reply a query, according to its intrinsic Strategy
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
		case "RANDOM":
			return bot.getRandom(response)
		case "WEIGHTED":
			return bot.getWeighted(response)
		case "MAXIMUM":
			return bot.getMaximum(response)
		default:
			return "", invalidStrategyError
		}
	} else {
		return q, nil
	}
}

func (bot *Bot) Save(path string) error {
	message, err := json.MarshalIndent(bot.ds, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, message, 0644)
	log.Info(fmt.Sprintf("Saving file to %v is successful", path))
	if err != nil {
		return err
	}

	return nil
}

func (bot *Bot) Load(path string) error {
	file, err := ioutil.ReadFile(path)
	log.Info(fmt.Sprintf("Loading file from %v is successful", path))

	if err != nil {
		return err
	}

	loadedDs := dataset{}

	err = json.Unmarshal([]byte(file), &loadedDs)
	if err != nil {
		return err
	}

	bot.ds = loadedDs
	return nil
}

func (bot *Bot) ValidStrategy() bool {
	_, ok := Find(supportedStrategy, bot.strategy)
	return ok
}