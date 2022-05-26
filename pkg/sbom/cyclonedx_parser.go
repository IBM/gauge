package sbom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	cdx "github.com/CycloneDX/cyclonedx-go"
	"github.com/IBM/gauge/pkg/common"
)

//ParseOSSPackagesCyclonedx :
func ParseOSSPackagesCyclonedx(sbomfp string) ([]common.PackageProps, error) {
	pkgList := []common.PackageProps{}

	sbomBuf, err := ioutil.ReadFile(sbomfp)
	if err != nil {
		return nil, fmt.Errorf("error reading sbom file `%s`: %v", sbomfp, err)
	}
	bomObj := cdx.NewBOM()
	if err = json.Unmarshal(sbomBuf, &bomObj); err != nil {
		return nil, fmt.Errorf("error parsing sbom file `%s`: %v", sbomfp, err)
	}
	for _, fpc := range *bomObj.Components {
		if fpc.Type == cdx.ComponentTypeFile {
			ecosystem := getPackageEcosystem(fpc.Name)
			for _, pkg := range *fpc.Components {
				p := common.PackageProps{
					Ecosystem: ecosystem,
					Name:      pkg.Name,
					Key:       pkg.Version,
				}
				pkgList = append(pkgList, p)
			}
		} else if fpc.Type == cdx.ComponentTypeLibrary {
			var ecosystem string
			if strings.Contains(fpc.BOMRef, "pkg:npm") {
				ecosystem = common.NODE_ECOSYSTEM
			} else if strings.Contains(fpc.BOMRef, "pkg:pypi") {
				ecosystem = common.PYTHON_ECOSYSTEM
			} else if strings.Contains(fpc.BOMRef, "pkg:maven") {
				// java doesnt exist yet
				// ecosystem = common.JAVA_ECOSYSTEM
				ecosystem = "unknown"
			} else if strings.Contains(fpc.BOMRef, "pkg:cpan") {
				// perl doesnt exist yet
				// ecosystem = common.PERL_ECOSYSTEM
				ecosystem = "unknown"
			} else {
				ecosystem = "unknown"
			}

			p := common.PackageProps{
				Ecosystem: ecosystem,
				Name:      fpc.Name,
				Key:       fpc.Version,
			}
			pkgList = append(pkgList, p)
		}
	}

	return pkgList, nil
}

func getPackageEcosystem(manifestName string) string {
	switch manifestName {
	case "requirements.txt":
		return common.PYTHON_ECOSYSTEM
	case "package-lock.json":
		return common.NODE_ECOSYSTEM
	case "package.json":
		return common.NODE_ECOSYSTEM
	default:
		return "unknown"
	}
}
