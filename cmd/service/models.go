package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/liqmix/ebiten-holiday-2024/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

// Store handles BadgerDB operations
type Store struct {
	db *badger.DB
}

// User model
type User struct {
	ID           uint       `json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	Username     string     `json:"username"`
	Password     string     `json:"password"`
	Rank         float64    `json:"rank"`
	RefreshToken string     `json:"refresh_token"`
	LastIP       string     `json:"last_ip"`
	LastLoginAt  time.Time  `json:"last_login_at"`
}

func (u *User) Clean() {
	u.Password = ""
	u.RefreshToken = ""
	u.LastIP = ""
}

// Score model
type Score struct {
	ID         uint      `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	UserID     uint      `json:"user_id"`
	Username   string    `json:"username"`
	SongHash   string    `json:"song_hash"`
	Score      int       `json:"score"`
	Rank       float64   `json:"rank"`
	Accuracy   float64   `json:"accuracy"`
	MaxCombo   int       `json:"max_combo"`
	PlayedAt   time.Time `json:"played_at"`
	Difficulty int       `json:"difficulty"`
}

func NewStore(path string) (*Store, error) {
	opts := badger.DefaultOptions(path)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

// Key prefixes for different types of data
const (
	userPrefix     = "user:"
	scorePrefix    = "score:"
	usernameIndex  = "username_index:"
	userScoreIndex = "user_score:"
)

// CreateUser creates a new user
func (s *Store) CreateUser(user *User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Check if username exists
	exists, err := s.UsernameExists(user.Username)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("username already exists")
	}

	// Get next ID
	id, err := s.getNextID("user_id")
	if err != nil {
		return err
	}
	user.ID = id

	// Create username index
	err = s.db.Update(func(txn *badger.Txn) error {
		// Store username index
		err := txn.Set([]byte(usernameIndex+user.Username), []byte(fmt.Sprintf("%d", user.ID)))
		if err != nil {
			return err
		}

		// Store user data
		userData, err := json.Marshal(user)

		if err != nil {
			return err
		}
		return txn.Set([]byte(userPrefix+fmt.Sprintf("%d", user.ID)), userData)
	})
	return err
}

// UpdateUser updates an existing user
func (s *Store) UpdateUser(user *User) error {
	user.UpdatedAt = time.Now()

	err := s.db.Update(func(txn *badger.Txn) error {
		userData, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return txn.Set([]byte(userPrefix+fmt.Sprintf("%d", user.ID)), userData)
	})
	return err
}

// GetUserByID retrieves a user by ID
func (s *Store) GetUserByID(id uint) (*User, error) {
	var user User
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(userPrefix + fmt.Sprintf("%d", id)))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &user)
		})
	})

	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (s *Store) GetUserByUsername(username string) (*User, error) {
	var userID uint
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(usernameIndex + username))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			var id uint
			_, err := fmt.Sscanf(string(val), "%d", &id)
			if err != nil {
				return err
			}
			userID = id
			return nil
		})
	})

	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return s.GetUserByID(userID)
}

// CreateScore creates a new score entry
func (s *Store) CreateScore(score *Score) error {
	score.CreatedAt = time.Now()
	score.UpdatedAt = time.Now()

	// Get next ID
	id, err := s.getNextID("score_id")
	if err != nil {
		return err
	}
	score.ID = id

	return s.db.Update(func(txn *badger.Txn) error {
		// Store score data
		scoreData, err := json.Marshal(score)
		if err != nil {
			return err
		}
		err = txn.Set([]byte(scorePrefix+fmt.Sprintf("%d", score.ID)), scoreData)
		if err != nil {
			return err
		}

		// Store user-score index
		return txn.Set(
			[]byte(fmt.Sprintf("%s%d:%d", userScoreIndex, score.UserID, score.ID)),
			nil,
		)
	})
}

// In your Store struct, add a new method:
func (s *Store) CreateScoreAndUpdateRating(score *Score, userID uint, ratingIncrease float64) error {
	return s.db.Update(func(txn *badger.Txn) error {
		// First get the user
		user, err := s.getUserInTx(txn, userID)
		if err != nil {
			return err
		}

		// Update user's rating if there's an increase
		if ratingIncrease > 0 {
			user.Rank += ratingIncrease
			user.UpdatedAt = time.Now()

			// Save updated user
			userData, err := json.Marshal(user)
			if err != nil {
				return err
			}
			if err := txn.Set([]byte(userPrefix+fmt.Sprintf("%d", user.ID)), userData); err != nil {
				return err
			}
		}

		// Set score metadata
		score.UserID = userID
		score.Username = user.Username
		score.CreatedAt = time.Now()
		score.UpdatedAt = time.Now()

		// Get next score ID
		id, err := s.getNextIDInTx(txn, "score_id")
		if err != nil {
			return err
		}
		score.ID = id

		// Store score data
		scoreData, err := json.Marshal(score)
		logger.Debug("Storing score %v", score)
		if err != nil {
			return err
		}
		if err := txn.Set([]byte(scorePrefix+fmt.Sprintf("%d", score.ID)), scoreData); err != nil {
			return err
		}

		// Store user-score index
		return txn.Set(
			[]byte(fmt.Sprintf("%s%d:%d", userScoreIndex, score.UserID, score.ID)),
			nil,
		)
	})
}

// Helper function to get user within a transaction
func (s *Store) getUserInTx(txn *badger.Txn, id uint) (*User, error) {
	item, err := txn.Get([]byte(userPrefix + fmt.Sprintf("%d", id)))
	if err != nil {
		return nil, err
	}

	var user User
	err = item.Value(func(val []byte) error {
		return json.Unmarshal(val, &user)
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Helper function to get next ID within a transaction
func (s *Store) getNextIDInTx(txn *badger.Txn, sequence string) (uint, error) {
	var id uint
	item, err := txn.Get([]byte("seq:" + sequence))
	if err == badger.ErrKeyNotFound {
		id = 1
	} else if err != nil {
		return 0, err
	} else {
		err = item.Value(func(val []byte) error {
			_, err := fmt.Sscanf(string(val), "%d", &id)
			id++
			return err
		})
		if err != nil {
			return 0, err
		}
	}
	err = txn.Set([]byte("seq:"+sequence), []byte(fmt.Sprintf("%d", id)))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Store) GetUserScores(userID uint) ([]Score, error) {
	var scores []Score
	prefix := []byte(fmt.Sprintf("%s%d:", userScoreIndex, userID))

	// Map to track highest score per song and difficulty
	highestScores := make(map[string]Score) // key: "songHash:difficulty"

	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = prefix

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.Item().Key()
			var scoreID uint
			_, err := fmt.Sscanf(string(key[len(prefix):]), "%d", &scoreID)
			if err != nil {
				continue
			}

			score, err := s.getScoreByID(txn, scoreID)
			if err != nil {
				continue
			}

			// Create composite key for song and difficulty
			mapKey := fmt.Sprintf("%s:%d", score.SongHash, score.Difficulty)

			// If this is the first score for this song/difficulty or if it's higher than existing score
			if existing, exists := highestScores[mapKey]; !exists || score.Score > existing.Score {
				highestScores[mapKey] = *score
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert map to slice
	for _, score := range highestScores {
		scores = append(scores, score)
	}

	// Sort scores by played time, most recent first
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].PlayedAt.After(scores[j].PlayedAt)
	})

	return scores, nil
}

func (s *Store) GetLeaderboard(song string, difficulty int) ([]Score, error) {
	var scores []Score
	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte(scorePrefix)

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(opts.Prefix); it.Valid(); it.Next() {
			item := it.Item()
			var score Score
			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &score)
			})
			if err != nil {
				continue
			}
			// Only include scores for the specified song and difficulty
			if score.SongHash == song && score.Difficulty == difficulty {
				scores = append(scores, score)
			}
		}
		return nil
	})

	// Sort scores by score value in descending order
	if err == nil {
		sort.Slice(scores, func(i, j int) bool {
			return scores[i].Score > scores[j].Score
		})

		// Keep only the best score per user
		seen := make(map[uint]bool)
		uniqueScores := make([]Score, 0)
		for _, score := range scores {
			if !seen[score.UserID] {
				seen[score.UserID] = true
				uniqueScores = append(uniqueScores, score)
			}
		}

		// Limit to top 10 scores
		if len(uniqueScores) > 10 {
			uniqueScores = uniqueScores[:10]
		}
		scores = uniqueScores
	}

	// for each score found, attach the user's current rank
	for i := range scores {
		user, err := s.GetUserByID(scores[i].UserID)
		if err != nil {
			return nil, err
		}
		scores[i].Rank = user.Rank
	}

	return scores, err
}

// Helper method to get a score by ID within a transaction
func (s *Store) getScoreByID(txn *badger.Txn, id uint) (*Score, error) {
	item, err := txn.Get([]byte(scorePrefix + fmt.Sprintf("%d", id)))
	if err != nil {
		return nil, err
	}

	var score Score
	err = item.Value(func(val []byte) error {
		return json.Unmarshal(val, &score)
	})
	if err != nil {
		return nil, err
	}

	return &score, nil
}

// Helper method to get next ID for a sequence
func (s *Store) getNextID(sequence string) (uint, error) {
	var id uint
	err := s.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("seq:" + sequence))
		if err == badger.ErrKeyNotFound {
			id = 1
		} else if err != nil {
			return err
		} else {
			err = item.Value(func(val []byte) error {
				_, err := fmt.Sscanf(string(val), "%d", &id)
				id++
				return err
			})
			if err != nil {
				return err
			}
		}
		return txn.Set([]byte("seq:"+sequence), []byte(fmt.Sprintf("%d", id)))
	})
	return id, err
}

// UsernameExists checks if a username is already taken
func (s *Store) UsernameExists(username string) (bool, error) {
	err := s.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(usernameIndex + username))
		return err
	})

	if err == badger.ErrKeyNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// User methods
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) SetRefreshToken(token string) error {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.RefreshToken = string(hashedToken)
	return nil
}

func (u *User) CheckRefreshToken(token string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.RefreshToken), []byte(token))
}
