package update

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PerrorOne/GetAnything-Server/logger"
	logger2 "github.com/apsdehal/go-logger"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"os"
	"strings"
	"time"
)

var (
	queryUrl = "https://api.github.com/repos/%s/%s/commits"
	log      = logger.NewLogger("GetAnything", logger2.InfoLevel)
)

type Update struct {
	Pid       int
	GithubUrl string
	ShaNow    string
}

type githubResp struct {
	Sha    string `json:"sha"`
	Commit struct {
		Author struct {
			Date time.Time `json:"date"`
		} `json:"author"`
		Message string `json:"message"`
	} `json:"commit"`
}

func NewUpdate(url string) (*Update, error) {
	u, err := url2.Parse(url)
	if err != nil {
		return nil, err
	}
	p := strings.Split(u.Path, "/")
	if len(p) != 3 {
		return nil, errors.New("GitHub仓库url错误")
	}
	update := &Update{Pid: os.Getpid(), GithubUrl: fmt.Sprintf(queryUrl, p[1], p[2])}
	c, err := ioutil.ReadFile("./git.log")
	if err != nil {
		resp, err := http.Get(fmt.Sprintf(queryUrl, p[1], p[2]))
		if err != nil {
			return nil, err
		}
		g := make([]*githubResp, 0)
		c, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(c, &g); err != nil {
			return nil, errors.New(fmt.Sprintf("%s:%s", err.Error(), string(c)))
		}
		update.ShaNow = g[0].Sha
		update.WriteSha()
	} else {
		update.ShaNow = string(c)
	}
	return update, nil
}

func (u *Update) WriteSha() {
	ioutil.WriteFile("./git.log", []byte(u.ShaNow), 0777)
}

func (u *Update) restart() error {
	p, err := os.FindProcess(u.Pid)
	if err != nil {
		return err
	}
	return p.Signal(os.Kill)
}

func (u *Update) build() error {
	return nil
}

func (u *Update) Restart() {
	for range time.Tick(time.Minute) {
		resp, err := http.Get(u.GithubUrl)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		g := make([]*githubResp, 0)
		c, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(c, &g); err != nil {
			log.Error(err.Error())
			continue
		}
		if len(g) > 0 {
			if g[0].Sha != u.ShaNow && u.ShaNow != "" { // update
				u.ShaNow = g[0].Sha
				u.WriteSha()
				if err := u.build(); err != nil {
					log.Error(err.Error())
					continue
				}
				if err := u.restart(); err != nil {
					log.Error(err.Error())
					continue
				}
			}
		}
	}
}
