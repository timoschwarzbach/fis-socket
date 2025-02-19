package sequences

import (
	"log"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SequenceService struct {
	db      *gorm.DB
	dbSync  chan bool
	current *Sequences
	next    *Sequences
}

func CreateSequenceService(dbSync chan bool) *SequenceService {
	db := connectToSQLite()
	c := &SequenceService{
		db:      db,
		dbSync:  dbSync,
		current: getFallbackSequence(),
		next:    getFallbackSequence(),
	}
	c.ListenToUpdates()
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

func (c *SequenceService) forceReconnect() {
	// sqlDB, err := c.db.DB()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = sqlDB.Close()
	// sqlCon, err := sqlDB.Conn(context.Background())
	log.Println("SequenceService:\tReconnecting to SQLite")
	c.db = connectToSQLite()
}

type Sequences struct {
	Id        string `gorm:"primaryKey"`
	Active    bool
	Category  string
	Locations []string `gorm:"type:text[]; serializer:json"`
	Slides    string   `gorm:"serializer:json"`
	// lastUpdated time.Time `gorm:"column:lastUpdated"`
}

type Files struct {
	Id           string `gorm:"primaryKey"`
	Bucket       string
	FileName     string `gorm:"column:fileName"`
	FileType     string `gorm:"column:fileType"`
	OriginalName string `gorm:"column:originalName"`
}

func getFallbackSequence() *Sequences {
	return &Sequences{
		Id:        "fallback",
		Locations: []string{},
		Slides:    "[]",
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

// please rework the project to either download to the correct id, or save the remote id in the display json
func (c *SequenceService) getLocalFileReferenceFromId(id string) (string, string) {
	var file Files
	result := c.db.Where("id = ?", id).Take(&file)
	if result.Error != nil || result.RowsAffected == 0 || file.Id == "" {
		return "image", "fallback"
	}
	if strings.HasPrefix(file.FileType, "video/") {
		return "video", file.FileName
	}
	return "image", file.FileName
}

func (c *SequenceService) Step() {
	c.current = c.next
	c.next = c.getNextSequence()
}

func (c *SequenceService) ListenToUpdates() {
	log.Println("SequenceService:\tListening to SQLite update requests")
	go func() {
		for {
			<-c.dbSync
			log.Println("SequenceService:\tDatabase reconnect update request received")
			c.forceReconnect()
		}
	}()
}
