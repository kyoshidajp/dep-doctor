package cmd

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const CACHE_FILE = ".dep-doctor.yaml"

type CacheStore struct {
	PackageManagers []CachePackageManager `yaml:"package_managers"`
}

func (r *CacheStore) URLbyPackageManager(packageManager string) map[string]string {
	repos := make(map[string]string)
	for i, p := range r.PackageManagers {
		if p.Name == packageManager {
			for _, repo := range r.PackageManagers[i].Repositories {
				repos[repo.Name] = repo.SourceURL
			}
		}
	}

	return repos
}

func (r *CacheStore) RebuildCacheRoot(diagnoses map[string]Diagnosis, packageManager string) CacheStore {
	var packageManagers []CachePackageManager
	// get from Diagnosis
	repos := map[string]CacheRepository{}
	for _, diagnosis := range diagnoses {
		if !diagnosis.Diagnosed {
			continue
		}

		repos[diagnosis.Name] = CacheRepository{
			Name:      diagnosis.Name,
			SourceURL: diagnosis.URL,
		}
	}

	// get from cache for padding
	cache := r.URLbyPackageManager(packageManager)
	for key := range cache {
		_, ok := repos[key]
		if !ok {
			// only isn't diagnosed
			repos[key] = CacheRepository{
				Name:      key,
				SourceURL: cache[key],
			}
		}
	}
	var crepos []CacheRepository
	for _, v := range repos {
		crepos = append(crepos, v)
	}
	packageManagers = append(packageManagers,
		CachePackageManager{
			Name:         o.packageManager,
			Repositories: crepos,
		},
	)

	for _, pm := range r.PackageManagers {
		if pm.Name != packageManager {
			// get from cache
			packageManagers = append(packageManagers, pm)
		}

	}

	return CacheStore{
		PackageManagers: packageManagers,
	}
}

type CachePackageManager struct {
	Name         string            `yaml:"name"`
	Repositories []CacheRepository `yaml:"repositories"`
}

type CacheRepository struct {
	Name      string `yaml:"name"`
	SourceURL string `yaml:"source_url"`
}

func BuildCacheStore() CacheStore {
	var store CacheStore
	file, err := os.ReadFile(CACHE_FILE)
	if err != nil {
		fmt.Println("could not read cache.")
		return store
	}

	if err := yaml.Unmarshal(file, &store); err != nil {
		fmt.Println("could not read cache.")
		return store
	}

	return store
}

func SaveCache(diagnoses map[string]Diagnosis, cacheStore CacheStore, packageManager string) error {
	root := cacheStore.RebuildCacheRoot(diagnoses, packageManager)
	yamlData, err := yaml.Marshal(&root)
	if err != nil {
		return err
	}

	file, err := os.Create(CACHE_FILE)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(yamlData)
	if err != nil {
		return err
	}

	return nil
}
