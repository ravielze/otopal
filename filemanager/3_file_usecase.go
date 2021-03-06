package filemanager

import (
	"github.com/ravielze/oculi/common"
	"github.com/ravielze/otopal/auth"
)

type Usecase struct {
	repo IRepo
}

func NewUsecase(repo IRepo) IUsecase {
	return Usecase{repo: repo}
}

func (uc Usecase) AddFile(user auth.User, fileGroup string, item common.FileAttachment) (FileResponse, error) {
	result, err := uc.repo.AddFile(user.ID, fileGroup, item.Attachment)
	if err != nil {
		return FileResponse{}, err
	}
	return result.Convert(), nil
}

func (uc Usecase) DeleteFile(idFile string) error {
	return uc.repo.DeleteFile(idFile)
}

func (uc Usecase) GetFile(idFile string) (FileResponse, error) {
	result, err := uc.repo.GetFile(idFile)
	if err != nil {
		return FileResponse{}, err
	}
	return result.Convert(), nil
}

func (uc Usecase) GetRawFile(idFile string) (File, error) {
	return uc.repo.GetFile(idFile)
}

func (uc Usecase) GetFilesByGroup(fileGroup string) ([]FileResponse, error) {
	rawResult, err := uc.repo.GetFilesByGroup(fileGroup)
	if err != nil {
		return []FileResponse{}, err
	}
	result := make([]FileResponse, len(rawResult))
	for i, x := range rawResult {
		result[i] = x.Convert()
	}
	return result, nil

}
