package sql

import (
	"github.com/kukkar/common-golang/pkg/utils/depchecker"
)

//RegisterDependencyChecker register dependency checker
func RegisterDependencyChecker() {
	depchecker.RegisterDependency(func() depchecker.Dependency {
		mysqlDep := new(MysqlChecker)
		return mysqlDep
	}())
}

// MysqlChecker health checker for mysql
type MysqlChecker struct{}

//GetPinger pinger ping to mysq conn to check conn is alive
func (this *MysqlChecker) GetPinger() func() (bool, error) {
	return func() (bool, error) {
		_, err := GetOldDbPool(DefaultKey)
		if err != nil {
			return false, err
		}
		_, err = GetNewDbPool(DefaultKey)
		if err != nil {
			return false, err
		}
		// mdberr := mdb.Insert("test", map[string]interface{}{
		// 	"testedAt": time.Now().Unix(),
		// })
		// if mdberr != nil {
		// 	return false, fmt.Errorf(mdberr.DeveloperMessage)
		// }
		return true, nil
	}
}

//GetName get healthchecker service name
func (this *MysqlChecker) GetName() string {
	return "mysql"
}
