package utils

import "github.com/wawafc/go-utils/money"

func ZeroMoney() money.Money {
	return money.NewMoneyFromFloat(0)
}

func OneMoney() money.Money {
	return money.NewMoneyFromFloat(1)
}

func MonthsInYearMoney() money.Money {
	return money.NewMoneyFromFloat(12)
}

func DaysInYearMoney() money.Money {
	return money.NewMoneyFromFloat(365)
}

func PercentageMoney() money.Money {
	return money.NewMoneyFromFloat(100)
}
