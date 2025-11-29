package db

import (
	"context"
	"log"

	"github.com/ecetinerdem/forseerv2/internal/store"
)

func Seed(store *store.Storage) {

	ctx := context.Background()

	users := generateUsers(100)

	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user: ", err)
			return
		}
	}

	//portfolios := generatePortfolios(200, users)

	return
}

func generateUsers(numUsers int) []*store.User {

	users := make([]*store.User, numUsers)

	for i := 0; i < numUsers; i++ {
		users[i] = &store.User{
			FirstName: SeedUsers[i].FirstName,
			LastName:  SeedUsers[i].LastName,
			Username:  SeedUsers[i].Username,
			Email:     SeedUsers[i].Email,
			Password:  SeedUsers[i].Password,
		}

	}
	return users

}

//func generatePortfolios(numPortfolios int, users []*store.User) []*store.Portfolio {

//}

//func generateStocks(numStocks int) []*store.Stock {

//}
