package main

import (
	"app/context"
	"app/domain/entity/master"
	mastervalue "app/domain/value/master"
	"app/domain/value/user"
	masterrepository "app/infrastructure/repository/master"
	"encoding/csv"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// TODO: validation対応
func main() {
	userId := user.NewUserId(uuid.New())
	ctx := &gin.Context{}
	ctx.Set("GameContext", &context.GameContext{
		UserId: &userId,
		Udctx:  &context.UserDbContext{},
		Ucctx:  &context.UserCacheContext{},
		Mdctx:  &context.MasterDbContext{},
		Mcctxt: &context.MasterCacheContext{},
		UtcNow: time.Now(),
	})
	loadEnv()
	d := LoadTsv()
	save(ctx, d)
}

// envファイル読み込み
func loadEnv() {
	err := godotenv.Load("../../../.env")
	if err != nil {
		panic(err)
	}
}

func LoadTsv() map[string][]map[string]string {
	ps := getPath("../../../resource/master")

	d := map[string][]map[string]string{}
	for _, p := range ps {
		mn := strings.Replace(filepath.Base(p), filepath.Ext(p), "", -1)
		d[mn] = loadTsv(p)
	}

	return d
}

func getPath(d string) []string {
	fs, err := os.ReadDir(d)
	if err != nil {
		panic(err)
	}

	var ps []string
	for _, f := range fs {
		if f.IsDir() {
			ps = append(ps, getPath(filepath.Join(d, f.Name()))...)
		}
		ps = append(ps, filepath.Join(d, f.Name()))
	}

	return ps
}

func loadTsv(p string) []map[string]string {
	file, err := os.Open(p)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	c := transform.NewReader(file, japanese.ShiftJIS.NewDecoder())

	r := csv.NewReader(c)
	r.Comma = '\t'
	r.Comment = '#'
	r.LazyQuotes = true

	cnt := 0
	var columns []string
	var ls []map[string]string
	for {
		l, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		if cnt == 0 {
			columns = l
		} else {
			line := map[string]string{}
			for idx, d := range l {
				line[columns[idx]] = d
			}
			ls = append(ls, line)
		}

		cnt++
	}

	return ls
}

func save(ctx *gin.Context, d map[string][]map[string]string) {
	mdc := &context.MasterDbContext{}

	err := mdc.Dc.TransactionScope(func() error {
		for n, m := range d {
			err := saveSwitcher(ctx, n, m)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return
	}
}

func saveSwitcher(
	ctx *gin.Context,
	n string,
	m []map[string]string,
) error {
	for _, l := range m {
		// TODO: 動的にしたい
		switch n {
		case "item_master":
			err := saveItemMaster(ctx, l)
			if err != nil {
				return err
			}
		case "schedule_master":
			err := saveScheduleMaster(ctx, l)
			if err != nil {
				return err
			}
		case "dictionary_master":
			err := saveDictionaryMaster(ctx, l)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func saveItemMaster(
	ctx *gin.Context,
	line map[string]string,
) error {
	d := &master.ItemMaster{}
	id, _ := strconv.ParseUint(line["id"], 10, 64)
	d.ID = mastervalue.NewItemId(id)

	d.Name = line["name"]

	t, _ := strconv.ParseUint(line["type"], 10, 64)
	d.Type = uint(t)

	s, _ := strconv.ParseUint(line["sell_coin"], 10, 64)
	d.SellCoin = s

	et, _ := strconv.ParseUint(line["effect_type"], 10, 64)
	d.EffectType = uint(et)

	ev, _ := strconv.ParseUint(line["effect_value"], 10, 64)
	d.EffectValue = ev

	si, _ := strconv.ParseInt(line["schedule_id"], 10, 64)
	d.ScheduleId = mastervalue.NewScheduleId(si)

	mc, _ := strconv.ParseUint(line["max_count"], 10, 64)
	d.MaxCount = mc

	mvc, _ := strconv.ParseUint(line["max_view_count"], 10, 64)
	d.MaxViewCount = mvc

	r := masterrepository.NewItemMasterRepository()
	return r.Save(ctx, *d)
}

func saveScheduleMaster(
	ctx *gin.Context,
	l map[string]string,
) error {
	d := &master.ScheduleMaster{}
	id, _ := strconv.ParseInt(l["id"], 10, 64)
	d.ID = mastervalue.NewScheduleId(id)

	jst, _ := time.LoadLocation("Asia/Tokyo")
	startAt, err := time.ParseInLocation(time.DateTime, l["start_at"], jst)
	if err != nil {
		log.Printf("failed start_at to time parse. %v", err)
		return err
	}
	d.StartAt = startAt

	endAt, err := time.ParseInLocation(time.DateTime, l["end_at"], jst)
	if err != nil {
		log.Printf("failed start_at to time parse. %v", err)
		return err
	}
	d.EndAt = endAt

	closeAt, err := time.ParseInLocation(time.DateTime, l["close_at"], jst)
	if err != nil {
		log.Printf("failed start_at to time parse. %v", err)
		return err
	}
	d.CloseAt = closeAt

	r := masterrepository.NewScheduleMasterRepository()
	return r.Save(ctx, *d)
}

func saveDictionaryMaster(
	ctx *gin.Context,
	l map[string]string,
) error {
	d := &master.DictionaryMaster{}
	d.Key = mastervalue.NewDictionaryKey(l["key"])
	d.Ja = l["ja"]
	d.En = l["en"]

	r := masterrepository.NewDictionaryMasterRepository()
	return r.Save(ctx, *d)
}
