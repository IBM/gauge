package sbom

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/IBM/gauge/pkg/common"
)

//ParseOSSPackagesCSV :
func ParseOSSPackagesCSV(sbomfp string) ([]common.PackageProps, error) {
	pkgList := []common.PackageProps{}
	f, err := os.Open(sbomfp)
	if err != nil {
		return pkgList, err
	}

	// remember to close the file at the end of the program
	defer f.Close()
	csvReader := csv.NewReader(f)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return pkgList, err
		}
		if rec[2] == "Python" || rec[2] == "javascript/Node.js" {
			pkg := common.PackageProps{}
			pkg.Ecosystem = getPackageEcosystemFromCSV(rec[2])
			pkg.Name = rec[0]
			pkg.Key = rec[1]
			pkgList = append(pkgList, pkg)
		}
	}
	return pkgList, nil
}

func getPackageEcosystemFromCSV(ecosystem string) string {
	switch ecosystem {
	case "Python":
		return common.PYTHON_ECOSYSTEM
	case "javascript/Node.js":
		return common.NODE_ECOSYSTEM
	default:
		return "unknown"
	}
}
