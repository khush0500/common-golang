package cache

import "github.com/kukkar/common-golang/pkg/utils/depchecker"

//RegisterDependencyChecker register dependency checker
func RegisterDependencyChecker() {
	depchecker.RegisterDependency(func() depchecker.Dependency {
		redisDep := new(CacheChecker)
		return redisDep
	}())
}

// CacheChecker health checker for cache
type CacheChecker struct{}

//GetPinger pinger ping to mysq conn to check conn is alive
func (this *CacheChecker) GetPinger() func() (bool, error) {
	return func() (bool, error) {
		_, err := GetPool(DefaultKey)
		if err != nil {
			return false, err
		}
		return true, nil
	}
}

//GetName get healthchecker service name
func (this *CacheChecker) GetName() string {
	return "cache"
}
