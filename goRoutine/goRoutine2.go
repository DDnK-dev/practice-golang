package goRoutine

import (
	"fmt"
	"sync"
)

var mutex sync.Mutex
var wg2 sync.WaitGroup

type Account struct {
	Balance int
}

func DepositAndWithdraw(account *Account) {
	mutex.Lock()
	defer mutex.Unlock()
	if account.Balance < 0 {
		panic(fmt.Sprintf("Balance should not be negative value: %d", account.Balance))
	}
	account.Balance += 1000
	account.Balance -= 1000
	wg2.Done()
}

func main() {
	account := &Account{0}
	wg2.Add(10)
	for i := 0; i < 10; i++{
		go func(){
			DepositAndWithdraw(account)
		}()
	}
	wg2.Wait()
}
