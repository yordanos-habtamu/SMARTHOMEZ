package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/yordanos-habtamu/realstate/config"
	"github.com/yordanos-habtamu/realstate/db"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Envs

	// Build connection string
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.DB_USER, cfg.DB_PWD, "localhost", cfg.DB_PORT, cfg.DB_NAME)

	database, err := db.NewPostgresStorage(connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("🌱 Starting database seeding...")

	// Seed users (customers and agents)
	if err := seedUsers(database); err != nil {
		log.Fatal("Failed to seed users:", err)
	}

	// Seed houses
	if err := seedHouses(database); err != nil {
		log.Fatal("Failed to seed houses:", err)
	}

	fmt.Println("✅ Database seeding completed successfully!")
}

func seedUsers(db *sql.DB) error {
	fmt.Println("👥 Seeding users...")

	// Hash password for all users
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	users := []struct {
		FirstName string
		LastName  string
		Email     string
		Role      string
		Contact   string
		DoB       string
		Sex       string
	}{
		// Agents
		{"John", "Doe", "john.agent@smarthomez.com", "agent", "+1234567890", "1985-05-15", "male"},
		{"Jane", "Smith", "jane.agent@smarthomez.com", "agent", "+1234567891", "1988-08-20", "female"},
		{"Michael", "Johnson", "michael.agent@smarthomez.com", "agent", "+1234567892", "1982-03-10", "male"},

		// Customers
		{"Alice", "Williams", "alice@example.com", "customer", "+1234567893", "1990-12-01", "female"},
		{"Bob", "Brown", "bob@example.com", "customer", "+1234567894", "1992-07-22", "male"},
		{"Carol", "Davis", "carol@example.com", "customer", "+1234567895", "1995-04-18", "female"},
		{"David", "Miller", "david@example.com", "customer", "+1234567896", "1987-11-30", "male"},
		{"Emma", "Wilson", "emma@example.com", "customer", "+1234567897", "1993-09-14", "female"},
	}

	for _, user := range users {
		_, err := db.Exec(`
			INSERT INTO users (first_name, last_name, email, password, role, contact, dob, sex, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
			ON CONFLICT (email) DO NOTHING
		`, user.FirstName, user.LastName, user.Email, string(hashedPassword), user.Role, user.Contact, user.DoB, user.Sex)

		if err != nil {
			return fmt.Errorf("failed to insert user %s: %w", user.Email, err)
		}
		fmt.Printf("  ✓ Created %s user: %s %s (%s)\n", user.Role, user.FirstName, user.LastName, user.Email)
	}

	return nil
}

func seedHouses(db *sql.DB) error {
	fmt.Println("🏠 Seeding houses...")

	// Get agent IDs
	var agent1ID, agent2ID, agent3ID int
	db.QueryRow("SELECT id FROM users WHERE email = $1", "john.agent@smarthomez.com").Scan(&agent1ID)
	db.QueryRow("SELECT id FROM users WHERE email = $1", "jane.agent@smarthomez.com").Scan(&agent2ID)
	db.QueryRow("SELECT id FROM users WHERE email = $1", "michael.agent@smarthomez.com").Scan(&agent3ID)

	houses := []struct {
		Address      string
		Category     string
		Price        float64
		NumBedrooms  int
		NumBathrooms int
		AreaSqFt     float64
		Description  string
		ImgURL       string
		AgentID      int
		IsSold       bool
	}{
		// Luxury Villas
		{
			"123 Sunset Boulevard, Beverly Hills, CA",
			"villa",
			2500000.00,
			5,
			4,
			4500.00,
			"Stunning luxury villa with panoramic ocean views, infinity pool, and modern amenities. Perfect for those seeking the ultimate in comfort and style.",
			"https://images.unsplash.com/photo-1613490493576-7fde63acd811?w=800",
			agent1ID,
			false,
		},
		{
			"456 Palm Drive, Miami, FL",
			"villa",
			1800000.00,
			4,
			3,
			3800.00,
			"Mediterranean-style villa with lush gardens, private beach access, and state-of-the-art security system.",
			"https://images.unsplash.com/photo-1512917774080-9991f1c4c750?w=800",
			agent2ID,
			false,
		},

		// Apartments
		{
			"789 Park Avenue, New York, NY",
			"apartment",
			850000.00,
			3,
			2,
			1500.00,
			"Modern apartment in prime Manhattan location with stunning city views, 24/7 concierge, and gym facilities.",
			"https://images.unsplash.com/photo-1545324418-cc1a3fa10c00?w=800",
			agent1ID,
			false,
		},
		{
			"321 Downtown Street, Chicago, IL",
			"apartment",
			450000.00,
			2,
			2,
			1200.00,
			"Stylish downtown apartment with hardwood floors, granite countertops, and in-unit washer/dryer.",
			"https://images.unsplash.com/photo-1522708323590-d24dbb6b0267?w=800",
			agent3ID,
			false,
		},
		{
			"555 Lake Shore Drive, Chicago, IL",
			"apartment",
			650000.00,
			3,
			2,
			1600.00,
			"Lakefront apartment with breathtaking views, modern kitchen, and access to building amenities including pool and fitness center.",
			"https://images.unsplash.com/photo-1502672260266-1c1ef2d93688?w=800",
			agent2ID,
			false,
		},

		// Condos
		{
			"888 Beach Road, San Diego, CA",
			"condo",
			720000.00,
			2,
			2,
			1400.00,
			"Beachfront condo with ocean views, updated kitchen, and resort-style amenities including pool and spa.",
			"https://images.unsplash.com/photo-1515263487990-61b07816b324?w=800",
			agent1ID,
			false,
		},
		{
			"999 Mountain View, Denver, CO",
			"condo",
			380000.00,
			2,
			1,
			1000.00,
			"Cozy mountain condo with fireplace, updated bathrooms, and access to hiking trails.",
			"https://images.unsplash.com/photo-1556912998-c57cc6b63cd7?w=800",
			agent3ID,
			false,
		},

		// Houses
		{
			"147 Maple Street, Austin, TX",
			"house",
			550000.00,
			4,
			3,
			2500.00,
			"Charming family home with spacious backyard, modern kitchen, and excellent school district.",
			"https://images.unsplash.com/photo-1568605114967-8130f3a36994?w=800",
			agent2ID,
			false,
		},
		{
			"258 Oak Avenue, Portland, OR",
			"house",
			475000.00,
			3,
			2,
			2000.00,
			"Contemporary house with open floor plan, eco-friendly features, and beautiful garden.",
			"https://images.unsplash.com/photo-1572120360610-d971b9d7767c?w=800",
			agent1ID,
			false,
		},
		{
			"369 Cedar Lane, Seattle, WA",
			"house",
			625000.00,
			4,
			3,
			2800.00,
			"Elegant craftsman home with hardwood floors, updated kitchen, and large deck perfect for entertaining.",
			"https://images.unsplash.com/photo-1600596542815-ffad4c1539a9?w=800",
			agent3ID,
			false,
		},

		// Luxury Apartments
		{
			"777 Penthouse Tower, Los Angeles, CA",
			"luxury",
			3200000.00,
			4,
			4,
			3500.00,
			"Ultra-luxury penthouse with 360-degree city views, private elevator, smart home technology, and rooftop terrace.",
			"https://images.unsplash.com/photo-1600607687939-ce8a6c25118c?w=800",
			agent2ID,
			false,
		},
		{
			"100 Sky Tower, San Francisco, CA",
			"luxury",
			2800000.00,
			3,
			3,
			3000.00,
			"Exclusive luxury residence with floor-to-ceiling windows, Italian marble, and concierge services.",
			"https://images.unsplash.com/photo-1600585154340-be6161a56a0c?w=800",
			agent1ID,
			false,
		},

		// Sold Properties (for testing)
		{
			"SOLD - 500 Main Street, Boston, MA",
			"house",
			890000.00,
			4,
			3,
			2600.00,
			"Recently sold! Historic brownstone with modern updates, located in prime neighborhood.",
			"https://images.unsplash.com/photo-1580587771525-78b9dba3b914?w=800",
			agent3ID,
			true,
		},
	}

	for _, house := range houses {
		_, err := db.Exec(`
			INSERT INTO house (address, category, price, num_bedrooms, num_bathrooms, area_sq_ft, description, img_url, agent_id, is_sold, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		`, house.Address, house.Category, house.Price, house.NumBedrooms, house.NumBathrooms, house.AreaSqFt, house.Description, house.ImgURL, house.AgentID, house.IsSold)

		if err != nil {
			return fmt.Errorf("failed to insert house %s: %w", house.Address, err)
		}

		soldStatus := ""
		if house.IsSold {
			soldStatus = " [SOLD]"
		}
		fmt.Printf("  ✓ Created %s: %s - $%.2f%s\n", house.Category, house.Address, house.Price, soldStatus)
	}

	return nil
}
