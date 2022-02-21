package service

import (
	"os"
	"path/filepath"

	"github.com/juju/errors"
)

type Trait struct {
	//groupRepository domain.GroupRepository
	//traitRepository domain.TraitRepository
}

func NewTrait() Trait {
	return Trait{}
}

func (s *Trait) Import(root string) error {
	//priority := 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		fileExtension := filepath.Ext(info.Name())

		if fileExtension != ".png" {
			return errors.New("file is not a .png")
		}

		//splitted := strings.Split(path, "/")
		//traitType := splitted[1]
		//traitRarity := splitted[2]
		//traitName := strings.TrimSuffix(info.Name(), fileExtension)

		/*
			foundGroup, _ := s.groupRepository.GetByName(groupName)
			if foundGroup == nil {
				if _, err := s.groupRepository.Create(&domain.GroupWrite{
					Name:     groupName,
					Priotiry: priority,
				}); err != nil {
					return errors.Annotate(err, "creating group failed")
				}

				priority++
				fmt.Printf("%s - %d\n", groupName, priority)
			}
			foundGroup, _ = s.groupRepository.GetByName(groupName)
			if foundGroup == nil {
				return errors.Errorf("group not found: group=%s", groupName)
			}

			b, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.Annotate(err, "reading file failed")
			}

			// Create trait.
			if _, err = s.traitRepository.Create(&domain.TraitWrite{
				Name:         traitName,
				Group:        foundGroup,
				Image:        imgBytes,
				RarenessKind: rarenessKind,
			}); err != nil {
				return errors.Annotate(err, "creating trait failed")
			}*/

		return nil
	})
	if err != nil {
		return errors.Annotate(err, "walking path failed")
	}

	return nil
}
