package can

import (
	"net/http"
	"strings"

	"github.com/go-chi/jwtauth/v5"
)

type Verifier func(string, int64) error

type Guard struct {
	Verifiers []Verifier
}

func NewGuard() *Guard {
	return &Guard{}
}

func (g *Guard) AddVerifier(v Verifier) {
	g.Verifiers = append(g.Verifiers, v)
}

func (g *Guard) Can(ja *jwtauth.JWTAuth, perm string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			token, claims, err := jwtauth.FromContext(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			if token == nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			id, ok := claims["user_id"].(float64)
			if !ok {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			errs := []string{}
			for _, v := range g.Verifiers {
				if err := v(perm, int64(id)); err != nil {
					errs = append(errs, err.Error())
				}
			}
			if len(errs) > 0 {
				http.Error(w, strings.Join(errs, "\n"), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
