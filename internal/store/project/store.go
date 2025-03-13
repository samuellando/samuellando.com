package project

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"samuellando.com/data"
	"samuellando.com/internal/cache"
	"samuellando.com/internal/datatypes"
	"samuellando.com/internal/errors"
	"samuellando.com/internal/store"
	"samuellando.com/internal/store/tag"
)

const URL = "https://api.github.com/users/samuellando/repos?per_page=100"
const API_VERSION = "2022-11-28"

type Store struct {
	db           *sql.DB
	options      Options
	materialized *store.MaterializedStore[Project]
}

type Options struct {
	Url string
}

func CreateStore(db *sql.DB, opts ...func(*Options)) Store {
	o := Options{Url: URL}
	for _, opt := range opts {
		opt(&o)
	}
	return Store{
		db:           db,
		options:      o,
		materialized: nil,
	}
}

func (ps Store) GetById(id int64) (Project, error) {
	if ps.materialized != nil {
		return ps.materialized.GetById(id)
	}
	ghProjects, err := loadGitHubProjects(ps.db, ps.options.Url)
	if err != nil {
		return Project{}, err
	}
	found := false
	var external Project
	for _, project := range ghProjects {
		if project.Id() == id {
			external = project
			found = true
			break
		}
	}
	if !found {
		return Project{}, errors.CreateNotFoundError("Project")
	}
	ctx := context.TODO()
	queries := data.New(ps.db)
	internal, err := getInternalProjectData(ctx, queries, id)
	if err != nil {
		return Project{}, err
	}
	project := coallesceProjectData(internal, external)
	project.db = ps.db
	return project, nil
}

func (ps Store) GetAll() ([]Project, error) {
	if ps.materialized != nil {
		return ps.materialized.GetAll()
	}
	ghProjects, err := loadGitHubProjects(ps.db, ps.options.Url)
	if err != nil {
		return nil, err
	}
	ctx := context.TODO()
	queries := data.New(ps.db)
	internals, err := getAllInternalProjectData(ctx, queries)
	if err != nil {
		return nil, err
	}
	projects := make([]Project, len(ghProjects))
	for i, external := range ghProjects {
		if internal, ok := internals[external.Id()]; ok {
			projects[i] = coallesceProjectData(internal, external)
		} else {
			projects[i] = external
		}
		projects[i].db = ps.db
	}
	return projects, nil
}

func (ps Store) Filter(f func(Project) bool) (store.Store[Project], error) {
	var filtered store.Store[Project]
	var err error
	if ps.materialized != nil {
		filtered, err = ps.materialized.Filter(f)
	} else {
		filtered, err = store.Filter(ps, f)
	}
	if err != nil {
		return ps, err
	}
	if ms, ok := filtered.(store.MaterializedStore[Project]); ok {
		return Store{db: ps.db, materialized: &ms}, nil
	} else {
		panic("Could not type cast to MaterializedStore!")
	}
}

func (ps Store) Group(f func(Project) string) (datatypes.OrderedMap[string, store.Store[Project]], error) {
	return store.Group(ps, f)
}

func (ps Store) Sort(f func(Project, Project) bool) (store.Store[Project], error) {
	var sorted store.Store[Project]
	var err error
	if ps.materialized != nil {
		sorted, err = ps.materialized.Sort(f)
	} else {
		sorted, err = store.Sort(ps, f)
	}
	if err != nil {
		return ps, err
	}
	if ms, ok := sorted.(store.MaterializedStore[Project]); ok {
		return Store{db: ps.db, materialized: &ms}, nil
	} else {
		panic("Could not type cast to MaterializedStore!")
	}
}

func (ds Store) AllTags() ([]tag.ProtoTag, error) {
	ctx := context.TODO()
	queries := data.New(ds.db)
	tagRows, err := queries.GetAllProjectTags(ctx)
	if err != nil {
		return nil, err
	}
	tags := make([]tag.ProtoTag, len(tagRows))
	for i, tagRow := range tagRows {
		tags[i] = tag.ProtoTag{
			Value: tagRow.Tag.Value,
			Color: tagRow.Tag.Color,
		}
	}
	return tags, nil
}

func (ds Store) AllSharedTags(tagValue string) ([]tag.ProtoTag, error) {
	ctx := context.TODO()
	queries := data.New(ds.db)
	tagRows, err := queries.GetSharedProjectTags(ctx, tagValue)
	if err != nil {
		return nil, err
	}
	tags := make([]tag.ProtoTag, len(tagRows))
	for i, tagRow := range tagRows {
		tags[i] = tag.ProtoTag{
			Value: tagRow.Tag.Value,
			Color: tagRow.Tag.Color,
		}
	}
	return tags, nil
}

func loadGitHubProjects(db *sql.DB, url string) ([]Project, error) {
	bytes, err := getGithubProjects(url, db)
	if err != nil {
		return nil, err
	}
	return unmarshalResponse(bytes)
}

func getGithubProjects(url string, db *sql.DB) ([]byte, error) {
	cachedFunc := cache.Cached(func() ([]byte, error) {
		req, err := createRequest(url)
		if err != nil {
			return nil, fmt.Errorf("Failed to crete request : %s", err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("Failed to get response : %s", err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			return nil, fmt.Errorf("Bad response code : %d", res.StatusCode)
		}
		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("Failed to read response body : %s", err)
		}
		return bytes, nil
	}, func(o *cache.CacheOptions) {
		o.MaxAge = 5 * time.Minute
		o.Db = db
	})
	return cachedFunc()
}

func createRequest(url string) (*http.Request, error) {
	if req, err := http.NewRequest("GET", url, nil); err != nil {
		return nil, err
	} else {
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", API_VERSION)
		return req, nil
	}
}

func unmarshalResponse(b []byte) ([]Project, error) {
	data := make([]*schema, 0)
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("Failed to Unmarshal Json : %s", err)
	}
	projects := make([]Project, 0, len(data))
	for _, d := range data {
		projects = append(projects, Project{
			id:          int64(d.ID),
			name:        d.Name,
			description: &d.Description,
			created:     d.CreatedAt,
			pushed:      d.PushedAt,
			url:         d.HTMLURL,
		})
	}
	return projects, nil
}

func getInternalProjectData(ctx context.Context, queries *data.Queries, id int64) (Project, error) {
	rows, err := queries.GetProject(ctx, id)
	if err != nil {
		return Project{}, err
	}
	if len(rows) == 0 {
		return Project{id: id, tags: []tag.ProtoTag{}}, nil
	}
	tags := make([]tag.ProtoTag, 0)
	for _, row := range rows {
		if row.TagID.Valid {
			tags = append(tags, tag.ProtoTag{Value: row.TagValue.String, Color: row.TagColor.String})
		}
	}
	var desc *string
	var imageLink *string
	if rows[0].Project.Description.Valid {
		desc = &rows[0].Project.Description.String
	}
	if rows[0].Project.ImageLink.Valid {
		imageLink = &rows[0].Project.ImageLink.String
	}
	return Project{
		id:          rows[0].Project.ID,
		description: desc,
		imageLink:   imageLink,
		hidden:      rows[0].Project.Hidden,
		tags:        tags,
	}, nil
}

func getAllInternalProjectData(ctx context.Context, queries *data.Queries) (map[int64]Project, error) {
	docRows, err := queries.GetProjects(ctx)
	if err != nil {
		return nil, err
	}
	projs := make(map[int64]*Project)
	for _, row := range docRows {
		if _, ok := projs[row.Project.ID]; !ok {
			var desc *string
			var imageLink *string
			if row.Project.Description.Valid {
				desc = &row.Project.Description.String
			}
			if row.Project.ImageLink.Valid {
				imageLink = &row.Project.ImageLink.String
			}
			projs[row.Project.ID] = &Project{
				id:          row.Project.ID,
				description: desc,
				imageLink:   imageLink,
				hidden:      row.Project.Hidden,
				tags:        make([]tag.ProtoTag, 0),
			}
		}
		if row.TagID.Valid {
			tag := tag.ProtoTag{
				Value: row.TagValue.String,
				Color: row.TagColor.String,
			}
			projs[row.Project.ID].tags = append(projs[row.Project.ID].tags, tag)
		}
	}
	res := make(map[int64]Project)
	for k, v := range projs {
		res[k] = *v
	}
	return res, nil
}

func coallesceProjectData(internal, external Project) Project {
	desc := external.description
	if internal.description != nil {
		desc = internal.description
	}
	return Project{
		id:          internal.id,
		description: desc,
		tags:        internal.tags,
		name:        external.name,
		created:     external.created,
		pushed:      external.pushed,
		url:         external.url,
		imageLink:   internal.imageLink,
		hidden:      internal.hidden,
	}
}
