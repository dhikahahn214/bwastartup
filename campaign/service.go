package campaign

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(userID int) ([]Campaign, error) {
	// yang ditangkap bukan dalam bentuk JSON, maka dari itu parameternya berbeda dengan service.go di package user
	if userID != 0 {
		campaigns, err := s.repository.FindByUserID(userID)
		if err != nil {
			return campaigns, err
		}

		return campaigns, nil
	}

	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

///////////////////////////////////
func (s *service) GetCampaignByID(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.repository.FindByID(input.ID)
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	campaign := Campaign{} // object campaign
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.GoalAmount = input.GoalAmount
	campaign.UserID = input.User.ID
	// pembuatan slug
	slugCandidate := fmt.Sprintf("%s %d", input.Name, input.User.ID) // penamaan ke database {%s = string, %d = integer}
	campaign.Slug = slug.Make(slugCandidate)                         // exp : Nama Campaign, user id => nama-campaign-1

	newCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}
	return newCampaign, nil
}

func (s *service) UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error) {
	// menangkap data lama dari database
	campaign, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		return campaign, err
	}
	// data sudah dengan data baru

	if campaign.UserID != inputData.User.ID {
		return campaign, errors.New("Not an owner of the cmapaign")
	}
	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.Perks = inputData.Perks
	campaign.GoalAmount = inputData.GoalAmount

	// update ke database
	updatedCampaign, err := s.repository.Update(campaign)
	if err != nil {
		return updatedCampaign, err
	}
	return updatedCampaign, nil
}
