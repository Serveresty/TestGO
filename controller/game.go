package controller

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type DiceData struct {
	BetAmount        string
	Range            string
	Multiply         string
	WinChance        string
	Profit           string
	Number1          int
	Number2          int
	Number3          int
	Number4          int
	NotEnoughMoney   bool
	OutOfRange       bool
	OutOfMultiply    bool
	UnknownWinChance bool
	Balance          string
}

type SendDiceData struct {
	BetAmount        string
	Range            string
	Multiply         string
	WinChance        string
	Profit           string
	Number1          int
	Number2          int
	Number3          int
	Number4          int
	NotEnoughMoney   bool
	OutOfRange       bool
	OutOfMultiply    bool
	UnknownWinChance bool
	Balance          string
}

func (s *DataBase) GetDiceData(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	conditionsMap := map[string]any{}
	var data DiceData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var send_data SendDiceData
	user, status := s.GetUser(rw, r)
	_ = status
	conditionsMap["balance"] = user.balance
	usern := user.username
	uID := user.id

	bal := conditionsMap["balance"].(string)
	n, _ := strconv.ParseFloat(bal, 64)

	//Проверка наличия суммы ставки на балансе
	if k, _ := strconv.ParseFloat(data.BetAmount, 64); n >= k {
		send_data.BetAmount = data.BetAmount
	} else {
		send_data.NotEnoughMoney = true
	}

	//Проверка диапазона умножения
	if k, _ := strconv.ParseFloat(data.Multiply, 64); k <= 9500 && k >= 1.0106 {
		send_data.Multiply = fmt.Sprintf("%.2f", k)
	} else {
		send_data.OutOfMultiply = true
	}

	//Проверка диапазона ставки
	if k, _ := strconv.ParseFloat(data.Range, 64); k <= 94.0 && k >= 0.01 {
		send_data.Range = fmt.Sprintf("%.2f", k)
	} else {
		send_data.OutOfRange = true
	}

	//Проверка шанса на победу
	if k, _ := strconv.ParseFloat(data.WinChance, 64); k <= 94.0 && k >= 0.01 {
		send_data.WinChance = fmt.Sprintf("%.2f", k)
	} else {
		send_data.UnknownWinChance = true
	}

	//Проверка на наличие ошибок при ставке
	if send_data.NotEnoughMoney == false && send_data.OutOfRange == false && send_data.OutOfMultiply == false && send_data.UnknownWinChance == false {
		bet_am, _ := strconv.ParseFloat(send_data.BetAmount, 64)
		mult, _ := strconv.ParseFloat(data.Multiply, 64)

		send_data.Profit = fmt.Sprintf("%.2f", bet_am*mult-bet_am)
		send_data.Number1 = rand.Intn(9 + 1)
		send_data.Number2 = rand.Intn(9 + 1)
		send_data.Number3 = rand.Intn(9 + 1)
		send_data.Number4 = rand.Intn(9 + 1)
		num := fmt.Sprintf("%d%d.%d%d", send_data.Number1, send_data.Number2, send_data.Number3, send_data.Number4)
		balance := s.transactionDice(num, send_data.Range, send_data.Profit, send_data.BetAmount, n, uID, usern)
		send_data.Balance = fmt.Sprintf("%.2f", balance)
	} else {
		send_data.Number1 = 0
		send_data.Number2 = 0
		send_data.Number3 = 0
		send_data.Number4 = 0
		send_data.BetAmount = data.BetAmount
		send_data.Profit = data.Profit
		send_data.Balance = user.balance
	}
	rw.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(rw).Encode(send_data)
	return
}

func (s *DataBase) transactionDice(n string, rng string, pt string, bt string, bal float64, idd interface{}, usrnm string) float64 {
	conditionsMap := map[string]any{}
	numb, _ := strconv.ParseFloat(n, 64)
	rang, _ := strconv.ParseFloat(rng, 64)
	prof, _ := strconv.ParseFloat(pt, 64)
	bet, _ := strconv.ParseFloat(bt, 64)
	var status string
	var summ string
	if rang > numb {
		bal += prof
		status = "win"
		summ = "+" + fmt.Sprint(prof)
	} else {
		bal -= bet
		status = "lost"
		summ = "-" + fmt.Sprint(bet)
	}

	perem := `
			INSERT INTO transactions (user_id, type, bet_amount, stat, summ) VALUES(?, ?, ?, ?, ?)
		`

	insert, errdb := s.Data.Query(perem, idd, "DiceBet", bet, status, summ)
	defer func() {
		if insert != nil {
		}
	}()
	if errdb != nil {
	}

	conditionsMap["balance"] = fmt.Sprint(bal)

	perem2 := `
			UPDATE users_account SET balance = ? WHERE username = ?;
		`

	insert2, errdb2 := s.Data.Query(perem2, conditionsMap["balance"], usrnm)
	defer func() {
		if insert2 != nil {
		}
	}()
	if errdb2 != nil {
	}
	return bal
}
