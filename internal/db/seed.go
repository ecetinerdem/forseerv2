package db

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/ecetinerdem/forseerv2/internal/store"
)

func Seed(store *store.Storage) {
	// Initialize random generator with new seed method
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	ctx := context.Background()

	// Generate users
	users := generateUsers(100)

	// Create users - we need to handle transactions properly
	for _, user := range users {
		// Use CreateAndInvite since we don't have direct transaction access
		// Using a dummy token and expiration for seeding
		token := "seed-token-" + strconv.Itoa(int(time.Now().UnixNano()))
		if err := store.Users.CreateAndInvite(ctx, user, token, time.Hour); err != nil {
			log.Println("Error creating user: ", err)
			return
		}
	}

	log.Printf("Created %d users\n", len(users))

	// Generate portfolios for each user (5 per user)
	portfolios := generatePortfolios(users, r)

	// Create portfolios with stocks
	for _, portfolio := range portfolios {
		if err := store.Portfolio.CreatePortfolioWithStocks(ctx, portfolio); err != nil {
			log.Println("Error creating portfolio: ", err)
			return
		}
	}

	log.Printf("Created %d portfolios with stocks\n", len(portfolios))
}

func generateUsers(numUsers int) []*store.User {
	users := make([]*store.User, numUsers)

	for i := 0; i < numUsers; i++ {
		baseUser := SeedUsers[i%len(SeedUsers)]

		// Create unique usernames and emails by appending numbers
		uniqueSuffix := strconv.Itoa(i)
		if i < len(SeedUsers) {
			// For the first batch, use original names
			users[i] = &store.User{
				FirstName: baseUser.FirstName,
				LastName:  baseUser.LastName,
				Username:  baseUser.Username,
				Email:     baseUser.Email,
			}
		} else {
			// For additional users, make unique
			users[i] = &store.User{
				FirstName: baseUser.FirstName,
				LastName:  baseUser.LastName,
				Username:  baseUser.Username + uniqueSuffix,
				Email:     "user" + uniqueSuffix + "@example.com",
			}
		}

		// Set password for all users
		users[i].Password.Set(baseUser.Password)
	}
	return users
}

func generatePortfolios(users []*store.User, r *rand.Rand) []*store.Portfolio {
	var portfolios []*store.Portfolio
	portfolioNames := []string{
		"Growth Portfolio", "Dividend Portfolio", "Tech Investments",
		"Long Term Holdings", "Trading Account", "Retirement Fund",
		"Blue Chip Stocks", "Emerging Markets", "Value Picks", "Income Portfolio",
	}

	for _, user := range users {
		for j := 0; j < 5; j++ {
			portfolio := &store.Portfolio{
				UserID: user.ID,
				Name:   portfolioNames[j%len(portfolioNames)] + " " + strconv.Itoa(j+1),
				Stocks: generateStocks(5, r),
			}
			portfolios = append(portfolios, portfolio)
		}
	}
	return portfolios
}

func generateStocks(numStocks int, r *rand.Rand) []store.Stock {
	// Common stock symbols with their typical price ranges
	stockData := []struct {
		Symbol   string
		MinPrice float64
		MaxPrice float64
	}{
		{"AAPL", 150.0, 200.0},
		{"GOOGL", 120.0, 150.0},
		{"MSFT", 280.0, 350.0},
		{"AMZN", 120.0, 180.0},
		{"TSLA", 180.0, 300.0},
		{"META", 300.0, 400.0},
		{"NVDA", 400.0, 500.0},
		{"JPM", 150.0, 180.0},
		{"JNJ", 150.0, 170.0},
		{"V", 220.0, 250.0},
		{"PG", 140.0, 160.0},
		{"UNH", 450.0, 550.0},
		{"HD", 300.0, 350.0},
		{"DIS", 80.0, 120.0},
		{"NFLX", 350.0, 450.0},
		{"ADBE", 500.0, 600.0},
		{"PYPL", 55.0, 70.0},
		{"CSCO", 45.0, 55.0},
		{"INTC", 30.0, 40.0},
		{"BA", 180.0, 220.0},
	}

	stocks := make([]store.Stock, numStocks)

	// Shuffle the stock data to get random selection
	shuffledStocks := make([]struct {
		Symbol   string
		MinPrice float64
		MaxPrice float64
	}, len(stockData))
	copy(shuffledStocks, stockData)

	// Use the provided random generator instead of global rand
	r.Shuffle(len(shuffledStocks), func(i, j int) {
		shuffledStocks[i], shuffledStocks[j] = shuffledStocks[j], shuffledStocks[i]
	})

	for i := 0; i < numStocks && i < len(shuffledStocks); i++ {
		stock := shuffledStocks[i]

		// Generate random shares between 10 and 1000 as float64
		shares := float64(r.Intn(991)+10) + r.Float64()

		// Generate random price within the stock's range
		priceRange := stock.MaxPrice - stock.MinPrice
		randomPrice := stock.MinPrice + r.Float64()*priceRange

		// Round to 2 decimal places for currency
		randomPrice = float64(int(randomPrice*100)) / 100
		shares = float64(int(shares*100)) / 100

		stocks[i] = store.Stock{
			Symbol:       stock.Symbol,
			Shares:       shares,
			AveragePrice: randomPrice,
		}
	}

	return stocks
}
