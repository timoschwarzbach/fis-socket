package sequences

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SequenceService struct {
	db      *gorm.DB
	current *Sequences
	next    *Sequences
}

func CreateSequenceService() *SequenceService {
	db := connectToSQLite()
	c := &SequenceService{
		db:      db,
		current: getFallbackSequence(),
		next:    getFallbackSequence(),
	}
	c.Step()
	return c
}

func connectToSQLite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("database.sqlite"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

type Sequences struct {
	Id          string `gorm:"primaryKey"`
	Active      bool
	Category    string
	Locations   []string `gorm:"type:text[]; serializer:json"`
	displayJSON string   `gorm:"column:displayJSON"`
	// lastUpdated time.Time `gorm:"column:lastUpdated"`
}

func getFallbackSequence() *Sequences {
	return &Sequences{
		Id:          "fallback",
		Locations:   []string{},
		displayJSON: "{}",
	}
}

func (c *SequenceService) getNextSequence() *Sequences {
	var sequence Sequences
	result := c.db.Where("id <> ?", c.current.Id).Order("RANDOM()").Take(&sequence)
	if result.Error != nil {
		return getFallbackSequence()
	}
	if result.RowsAffected == 0 {
		return c.current
	}
	return &sequence
}

func (c *SequenceService) Step() {
	c.current = c.next
	c.next = c.getNextSequence()
}
