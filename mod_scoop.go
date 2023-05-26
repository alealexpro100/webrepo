package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/labstack/echo/v4"
)

type mod_scoop struct {
	cfg *Config
}

func (m *mod_scoop) init(config *Config) {
	// Set config
	m.cfg = config
}

func (m *mod_scoop) list_repos() ([]string, error) {
	entries, err := os.ReadDir(m.cfg.ScoopDir)
	l := []string{}
	if err == nil {
		for _, file := range entries {
			if file.IsDir() {
				l = append(l, file.Name())
			}
		}
	}
	return l, err
}

func (m *mod_scoop) clone_repo(name, repo string) error {
	_, err := git.PlainClone(path.Join(m.cfg.ScoopDir, name), false, &git.CloneOptions{
		URL:      repo,
		Progress: os.Stdout,
	})
	return err
}

func (m *mod_scoop) init_repo(name string) error {
	return m.clone_repo(name, m.cfg.ScoopInitRepo)
}

func (m *mod_scoop) del_repo(name string) error {
	return os.RemoveAll(path.Join(m.cfg.ScoopDir, name))
}

func (m *mod_scoop) list_packages(repo string) ([]string, error) {
	entries, err := os.ReadDir(path.Join(m.cfg.ScoopDir, repo, "bucket"))
	l := []string{}
	if err == nil {
		for _, file := range entries {
			if !file.IsDir() {
				l = append(l, file.Name())
			}
		}
	}
	return l, err
}

// We do not need mod function, as this function will do the same work
func (m *mod_scoop) add_package(repo, package_name string, contents string) error {
	// Write or mod file
	f, err := os.Create(path.Join(m.cfg.ScoopDir, repo, "bucket", package_name+".json"))
	if err != nil {
		return err
	}
	if _, err = f.WriteString(contents); err != nil {
		return err
	}
	f.Close()
	// git add
	repogit, err := git.PlainOpen(path.Join(m.cfg.ScoopDir, repo))
	if err != nil {
		return err
	}
	worktree, err := repogit.Worktree()
	if err != nil {
		return err
	}
	_, err = worktree.Add(path.Join("bucket", package_name+".json"))
	if err != nil {
		return err
	}
	// git commit
	_, err = worktree.Commit(package_name+": auto mod", &git.CommitOptions{
		Author: &object.Signature{
			Name:  m.cfg.AdminUser,
			Email: "worker@test.org",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *mod_scoop) get_package(repo, package_name string) (string, error) {
	f, err := os.Open(path.Join(m.cfg.ScoopDir, repo, "bucket", package_name+".json"))
	if err != nil {
		return "", err
	}
	c, err := ioutil.ReadAll(f)
	f.Close()
	return string(c), err
}

func (m *mod_scoop) del_package(repo, package_name string) error {
	err := os.Remove(path.Join(m.cfg.ScoopDir, repo, "bucket", package_name+".json"))
	if err != nil {
		return err
	}
	// git add
	repogit, err := git.PlainOpen(path.Join(m.cfg.ScoopDir, repo))
	if err != nil {
		return err
	}
	worktree, err := repogit.Worktree()
	if err != nil {
		return err
	}
	_, err = worktree.Remove(path.Join("bucket", package_name+".json"))
	if err != nil {
		return err
	}
	// git commit
	_, err = worktree.Commit(package_name+": auto del", &git.CommitOptions{
		Author: &object.Signature{
			Name:  m.cfg.AdminUser,
			Email: "worker@test.org",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}
	return err
}

func (m *mod_scoop) init_route(loc *echo.Group) {
	loc.GET("/repo/list", func(c echo.Context) error {
		repos, err := m.list_repos()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to receive repo list.")
		}
		type Data struct {
			Repos []string `json:"repos"`
		}
		data := &Data{
			Repos: repos,
		}
		return c.JSON(http.StatusOK, data)
	})
	loc.POST("/repo/init", func(c echo.Context) error {
		name := c.QueryParam("name")
		if name == "" {
			return c.String(http.StatusBadRequest, "No name param.")
		}
		err := m.init_repo(name)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to init repo.")
		}
		return c.String(http.StatusOK, "ok")
	})
	loc.POST("/repo/clone", func(c echo.Context) error {
		name := c.QueryParam("name")
		url_string := c.QueryParam("url")
		if name == "" {
			return c.String(http.StatusBadRequest, "No name param.")
		}
		if url_string == "" {
			return c.String(http.StatusBadRequest, "No url param.")
		}
		err := m.clone_repo(name, url_string)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to init repo.")
		}
		return c.String(http.StatusOK, "ok")
	})
	loc.POST("/repo/del", func(c echo.Context) error {
		name := c.QueryParam("name")
		if name == "" {
			return c.String(http.StatusBadRequest, "No name param.")
		}
		err := m.del_repo(name)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to del repo.")
		}
		return c.String(http.StatusOK, "ok")
	})
	loc.GET("/packages/list", func(c echo.Context) error {
		repo := c.QueryParam("repo")
		if repo == "" {
			return c.String(http.StatusBadRequest, "No repo param.")
		}
		packages, err := m.list_packages(repo)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to receive packages list.")
		}
		type Data struct {
			Packages []string `json:"packages"`
		}
		data := &Data{
			Packages: packages,
		}
		return c.JSON(http.StatusOK, data)
	})
	loc.GET("/packages/get", func(c echo.Context) error {
		repo := c.QueryParam("repo")
		if repo == "" {
			return c.String(http.StatusBadRequest, "No repo param.")
		}
		name := c.QueryParam("name")
		if name == "" {
			return c.String(http.StatusBadRequest, "No name param.")
		}
		packagestring, err := m.get_package(repo, name)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to get package.")
		}
		return c.String(http.StatusOK, packagestring)
	})
	loc.POST("/packages/add", func(c echo.Context) error {
		repo := c.QueryParam("repo")
		if repo == "" {
			return c.String(http.StatusBadRequest, "No repo param.")
		}
		name := c.QueryParam("name")
		if name == "" {
			return c.String(http.StatusBadRequest, "No name param.")
		}
		// Source
		json_file := c.FormValue("document")
		if json_file == "" {
			return c.String(http.StatusBadRequest, "Failed to get package file.")
		}
		err := m.add_package(repo, name, json_file)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to create package file.")
		}
		return c.String(http.StatusOK, "ok")
	})
	loc.POST("/packages/del", func(c echo.Context) error {
		repo := c.QueryParam("repo")
		if repo == "" {
			return c.String(http.StatusBadRequest, "No repo param.")
		}
		name := c.QueryParam("name")
		if name == "" {
			return c.String(http.StatusBadRequest, "No name param.")
		}
		err := m.del_package(repo, name)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to get package.")
		}
		return c.String(http.StatusOK, "ok")
	})
}
