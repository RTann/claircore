package test

import (
	"fmt"
	"strconv"

	"github.com/quay/claircore"
)

// GenUniqueRepositories creates an array of unique repositories.
//
// The array is guaranteed not to have any duplicated fields.
func GenUniqueRepositories(n int, opts ...GenRepoOption) []*claircore.Repository {
	repos := []*claircore.Repository{}
	for i := range n {
		r := claircore.Repository{
			ID:   strconv.Itoa(i),
			Name: fmt.Sprintf("repository-%d", i),
			Key:  fmt.Sprintf("key-%d", i),
			URI:  fmt.Sprintf("uri-%d", i),
		}
		for _, f := range opts {
			f(&r)
		}
		repos = append(repos, &r)
	}
	return repos
}

// GenRepoOption is a functional option to [GenUniqueRepositories].
type GenRepoOption func(*claircore.Repository)
