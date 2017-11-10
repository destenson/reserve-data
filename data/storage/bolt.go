package storage

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/boltdb/bolt"
)

const (
	PRICE_BUCKET            string = "prices"
	BALANCE_BUCKET          string = "balances"
	EXCHANGE_BALANCE_BUCKET string = "ebalances"
	RATE_BUCKET             string = "rates"
)

type BoltStorage struct {
	db *bolt.DB
}

func NewBoltStorage(path string) (*BoltStorage, error) {
	// init instance
	var err error
	var db *bolt.DB
	db, err = bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	// init buckets
	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucket([]byte(PRICE_BUCKET))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket([]byte(BALANCE_BUCKET))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket([]byte(EXCHANGE_BALANCE_BUCKET))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket([]byte(RATE_BUCKET))
		if err != nil {
			return err
		}
		return nil
	})
	return &BoltStorage{db}, nil
}

func uint64ToBytes(u uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, u)
	return b
}

func bytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func reverseSeek(timepoint uint64, c *bolt.Cursor) (uint64, error) {
	version, _ := c.Seek(uint64ToBytes(timepoint))
	if version == nil {
		version, _ = c.Prev()
		if version == nil {
			return 0, errors.New(fmt.Sprintf("There is no data before timepoint %d", timepoint))
		} else {
			return bytesToUint64(version), nil
		}
	} else {
		v := bytesToUint64(version)
		if v == timepoint {
			return v, nil
		} else {
			version, _ = c.Prev()
			if version == nil {
				return 0, errors.New(fmt.Sprintf("There is no data before timepoint %d", timepoint))
			} else {
				return bytesToUint64(version), nil
			}
		}
	}
}

func (self *BoltStorage) CurrentPriceVersion(timepoint uint64) (common.Version, error) {
	var result uint64
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(PRICE_BUCKET)).Cursor()
		result, err = reverseSeek(timepoint, c)
		return nil
	})
	return common.Version(result), err
}

func (self *BoltStorage) GetAllPrices(version common.Version) (map[common.TokenPairID]common.OnePrice, error) {
	result := map[common.TokenPairID]common.OnePrice{}
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PRICE_BUCKET))
		data := b.Get(uint64ToBytes(uint64(version)))
		if data == nil {
			err = errors.New(fmt.Sprintf("version %s doesn't exist", version))
		} else {
			err = json.Unmarshal(data, &result)
		}
		return nil
	})
	return result, err
}

func (self *BoltStorage) GetOnePrice(pair common.TokenPairID, version common.Version) (common.OnePrice, error) {
	result := map[common.TokenPairID]common.OnePrice{}
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PRICE_BUCKET))
		data := b.Get(uint64ToBytes(uint64(version)))
		if data == nil {
			err = errors.New(fmt.Sprintf("version %s doesn't exist", version))
		} else {
			err = json.Unmarshal(data, &result)
		}
		return nil
	})
	if err != nil {
		return common.OnePrice{}, err
	} else {
		pair, exist := result[pair]
		if exist {
			return pair, nil
		} else {
			return common.OnePrice{}, errors.New("Pair of token is not supported")
		}
	}
}

func (self *BoltStorage) CurrentBalanceVersion(timepoint uint64) (common.Version, error) {
	var result uint64
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(BALANCE_BUCKET)).Cursor()
		result, err = reverseSeek(timepoint, c)
		return nil
	})
	return common.Version(result), err
}
func (self *BoltStorage) GetAllBalances(version common.Version) (map[string]common.BalanceEntry, error) {
	result := map[string]common.BalanceEntry{}
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PRICE_BUCKET))
		data := b.Get(uint64ToBytes(uint64(version)))
		if data == nil {
			err = errors.New(fmt.Sprintf("version %s doesn't exist", version))
		} else {
			err = json.Unmarshal(data, &result)
		}
		return nil
	})
	return result, err
}

func (self *BoltStorage) CurrentEBalanceVersion(timepoint uint64) (common.Version, error) {
	var result uint64
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(EXCHANGE_BALANCE_BUCKET)).Cursor()
		result, err = reverseSeek(timepoint, c)
		return nil
	})
	return common.Version(result), err
}

func (self *BoltStorage) GetAllEBalances(version common.Version) (map[common.ExchangeID]common.EBalanceEntry, error) {
	result := map[common.ExchangeID]common.EBalanceEntry{}
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PRICE_BUCKET))
		data := b.Get(uint64ToBytes(uint64(version)))
		if data == nil {
			err = errors.New(fmt.Sprintf("version %s doesn't exist", version))
		} else {
			err = json.Unmarshal(data, &result)
		}
		return nil
	})
	return result, err
}

func (self *BoltStorage) CurrentRateVersion(timepoint uint64) (common.Version, error) {
	var result uint64
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(RATE_BUCKET)).Cursor()
		result, err = reverseSeek(timepoint, c)
		return nil
	})
	return common.Version(result), err
}

func (self *BoltStorage) GetAllRates(version common.Version) (common.AllRateEntry, error) {
	result := common.AllRateEntry{}
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(RATE_BUCKET))
		data := b.Get(uint64ToBytes(uint64(version)))
		if data == nil {
			err = errors.New(fmt.Sprintf("version %s doesn't exist", version))
		} else {
			err = json.Unmarshal(data, &result)
		}
		return nil
	})
	return result, err
}

func (self *BoltStorage) StorePrice(data map[common.TokenPairID]common.OnePrice, timepoint uint64) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		var dataJson []byte
		b := tx.Bucket([]byte(PRICE_BUCKET))
		dataJson, err = json.Marshal(data)
		if err != nil {
			return err
		}
		return b.Put(uint64ToBytes(timepoint), dataJson)
	})
	return err
}

func (self *BoltStorage) StoreBalance(data map[string]common.BalanceEntry, timepoint uint64) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		var dataJson []byte
		b := tx.Bucket([]byte(BALANCE_BUCKET))
		dataJson, err = json.Marshal(data)
		if err != nil {
			return err
		}
		return b.Put(uint64ToBytes(timepoint), dataJson)
	})
	return err
}

func (self *BoltStorage) StoreEBalance(data map[common.ExchangeID]common.EBalanceEntry, timepoint uint64) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		var dataJson []byte
		b := tx.Bucket([]byte(EXCHANGE_BALANCE_BUCKET))
		dataJson, err = json.Marshal(data)
		if err != nil {
			return err
		}
		return b.Put(uint64ToBytes(timepoint), dataJson)
	})
	return err
}

func (self *BoltStorage) StoreRate(data common.AllRateEntry, timepoint uint64) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		var dataJson []byte
		b := tx.Bucket([]byte(RATE_BUCKET))
		dataJson, err = json.Marshal(data)
		if err != nil {
			return err
		}
		return b.Put(uint64ToBytes(timepoint), dataJson)
	})
	return err
}